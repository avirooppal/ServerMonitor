package metrics

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"
	"os"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"bytes"
	"fmt"
	"github.com/docker/docker/pkg/stdcopy"
)

type Collector struct {
	lastNetIO    []net.IOCountersStat
	lastDiskIO   map[string]disk.IOCountersStat
	lastTime     time.Time
	dockerClient *client.Client

	// Caching for heavy operations
	diskUsageMutex  sync.RWMutex
	cachedDiskUsage []FolderSize
}

func NewCollector() *Collector {
	c := &Collector{
		lastDiskIO: make(map[string]disk.IOCountersStat),
		lastTime:   time.Now(),
	}

	// Init Docker Client
	var err error
	c.dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Failed to create Docker client: %v", err)
	} else {
		// Test connection
		_, err := c.dockerClient.Ping(context.Background())
		if err != nil {
			log.Printf("Docker not available: %v", err)
			c.dockerClient = nil
		} else {
			log.Println("Docker client connected")
		}
	}

	return c
}

func (c *Collector) Collect() SystemMetrics {
	now := time.Now()
	timeDiff := now.Sub(c.lastTime).Seconds()

	var metrics SystemMetrics

	// Host Info
	h, _ := host.Info()
	metrics.HostInfo = h

	// CPU
	// Use 1s interval to get accurate reading, preventing 100% spikes
	cpuPerc, _ := cpu.Percent(1*time.Second, true)
	
	var total float64
	for _, p := range cpuPerc {
		total += p
	}
	if len(cpuPerc) > 0 {
		metrics.CPUTotal = total / float64(len(cpuPerc))
	}
	metrics.CPU = cpuPerc
	
	l, _ := load.Avg()
	metrics.LoadAvg = l

	// Memory
	m, _ := mem.VirtualMemory()
	s, _ := mem.SwapMemory()
	metrics.Memory = &ExtendedMemoryStat{
		VirtualMemoryStat: m,
		Buffers:           m.Buffers,
		Cached:            m.Cached,
	}
	metrics.Swap = s

	// Disks (Usage & I/O)
	parts, _ := disk.Partitions(true)
	var disks []DiskInfo
	ioCounters, _ := disk.IOCounters()

	log.Printf("DEBUG: Found %d partitions", len(parts))
	for _, p := range parts {
		log.Printf("DEBUG: Checking partition: %s (%s)", p.Mountpoint, p.Fstype)
		// Filter out Docker bind mounts and irrelevant system paths
		if p.Mountpoint == "/etc/hostname" || 
		   p.Mountpoint == "/etc/hosts" || 
		   p.Mountpoint == "/etc/resolv.conf" ||
		   strings.HasPrefix(p.Mountpoint, "/dev") ||
		   strings.HasPrefix(p.Mountpoint, "/sys") ||
		   strings.HasPrefix(p.Mountpoint, "/proc") ||
		   strings.HasPrefix(p.Mountpoint, "/run") ||
		   p.Fstype == "tmpfs" ||
		   p.Fstype == "devtmpfs" ||
		   p.Fstype == "squashfs" ||
		   (p.Fstype == "overlay" && p.Mountpoint != "/") {
			continue
		}

		u, err := disk.Usage(p.Mountpoint)
		if err != nil {
			log.Printf("DEBUG: Failed to get usage for %s: %v", p.Mountpoint, err)
			continue
		}

		// Calculate I/O Rates
		var rRate, wRate uint64
		deviceName := p.Device
		if strings.HasPrefix(deviceName, "/dev/") {
			deviceName = strings.TrimPrefix(deviceName, "/dev/")
		}

		if curIO, ok := ioCounters[deviceName]; ok {
			if prevIO, ok := c.lastDiskIO[deviceName]; ok && timeDiff > 0 {
				rRate = uint64(float64(curIO.ReadBytes-prevIO.ReadBytes) / timeDiff)
				wRate = uint64(float64(curIO.WriteBytes-prevIO.WriteBytes) / timeDiff)
			}
			c.lastDiskIO[deviceName] = curIO
		}

		dInfo := DiskInfo{
			Path:        p.Mountpoint,
			Total:       u.Total,
			Used:        u.Used,
			Free:        u.Free,
			UsedPercent: u.UsedPercent,
			ReadRate:    rRate,
			WriteRate:   wRate,
		}
		log.Printf("DEBUG: Added partition: %s", p.Mountpoint)
		disks = append(disks, dInfo)
	}
	metrics.Disks = disks

	// Network
	netIO, _ := net.IOCounters(true)
	var netStats NetworkStats
	var totalRecv, totalSent uint64

	if timeDiff > 0 {
		// Calculate Net Rates
		for _, cur := range netIO {
			var prev net.IOCountersStat
			found := false
			for _, p := range c.lastNetIO {
				if p.Name == cur.Name {
					prev = p
					found = true
					break
				}
			}

			if found {
				rRate := uint64(float64(cur.BytesRecv-prev.BytesRecv) / timeDiff)
				sRate := uint64(float64(cur.BytesSent-prev.BytesSent) / timeDiff)
				netStats.Interfaces = append(netStats.Interfaces, NetInterface{
					Name:     cur.Name,
					RecvRate: rRate,
					SentRate: sRate,
				})
				totalRecv += rRate
				totalSent += sRate
			}
		}
	}
	netStats.TotalRecv = totalRecv
	netStats.TotalSent = totalSent
	c.lastNetIO = netIO
	metrics.Network = netStats

	// Processes (Top 20 by CPU)
	procs, _ := process.Processes()
	var procInfos []ProcessInfo
	for _, p := range procs {
		cpuP, _ := p.CPUPercent()
		if cpuP < 0.1 {
			continue
		} // Optimization

		name, _ := p.Name()
		memP, _ := p.MemoryPercent()
		username, _ := p.Username()

		procInfos = append(procInfos, ProcessInfo{
			PID:      p.Pid,
			Name:     name,
			CPU:      cpuP,
			Mem:      memP,
			Username: username,
		})
	}

	// Sort procInfos by CPU (desc)
	for i := 0; i < len(procInfos); i++ {
		for j := i + 1; j < len(procInfos); j++ {
			if procInfos[i].CPU < procInfos[j].CPU {
				procInfos[i], procInfos[j] = procInfos[j], procInfos[i]
			}
		}
	}
	if len(procInfos) > 20 {
		procInfos = procInfos[:20]
	}
	metrics.Processes = procInfos

	// Docker Containers
	var containerInfos []ContainerInfo
	if c.dockerClient != nil {
		containers, err := c.dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
		if err == nil {
			for _, ctr := range containers {
				name := "unknown"
				if len(ctr.Names) > 0 {
					name = ctr.Names[0]
				}

				var cpuPercent float64
				var memUsage, memLimit uint64

				if ctr.State == "running" {
					stats, err := c.dockerClient.ContainerStats(context.Background(), ctr.ID, false)
					if err == nil {
						var statsData map[string]interface{}
						if err := json.NewDecoder(stats.Body).Decode(&statsData); err == nil {
							if mem, ok := statsData["memory_stats"].(map[string]interface{}); ok {
								if usage, ok := mem["usage"].(float64); ok {
									memUsage = uint64(usage)
								}
								if limit, ok := mem["limit"].(float64); ok {
									memLimit = uint64(limit)
								}
							}
						}
						stats.Body.Close()
					}
				}

				containerInfos = append(containerInfos, ContainerInfo{
					ID:          ctr.ID,
					Name:        name,
					Image:       ctr.Image,
					State:       ctr.State,
					Status:      ctr.Status,
					Created:     ctr.Created,
					CPUPercent:  cpuPercent,
					MemoryUsage: memUsage,
					MemoryLimit: memLimit,
				})
			}
		}
	}
	metrics.Containers = containerInfos

	c.lastTime = now
	metrics.LastUpdate = now

	return metrics
}

func (c *Collector) GetContainerLogs(containerID string, tail string) (string, error) {
	if c.dockerClient == nil {
		return "", fmt.Errorf("docker client not available")
	}

	// Inspect to check if TTY is enabled
	inspect, err := c.dockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %v", err)
	}

	opts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	}

	out, err := c.dockerClient.ContainerLogs(context.Background(), containerID, opts)
	if err != nil {
		return "", err
	}
	defer out.Close()

	var buf bytes.Buffer
	if inspect.Config.Tty {
		_, err = io.Copy(&buf, out)
	} else {
		_, err = stdcopy.StdCopy(&buf, &buf, out)
	}

	if err != nil {
		log.Printf("Error copying logs for container %s: %v", containerID, err)
		return "", fmt.Errorf("failed to copy logs: %v", err)
	}

	return buf.String(), nil
}

// StartBackgroundTasks starts periodic heavy tasks
func (c *Collector) StartBackgroundTasks() {
	go func() {
		// Initial run
		c.updateDiskUsage()

		// Run every 15 minutes
		ticker := time.NewTicker(15 * time.Minute)
		for range ticker.C {
			c.updateDiskUsage()
		}
	}()
}

func (c *Collector) updateDiskUsage() {
	// Default to root, or env var
	path := "/" 
	if p := os.Getenv("DISK_USAGE_PATH"); p != "" {
		path = p
	}

	usage, err := c.GetDiskUsage(path)
	if err != nil {
		log.Printf("Error updating disk usage: %v", err)
		return
	}

	c.diskUsageMutex.Lock()
	c.cachedDiskUsage = usage
	c.diskUsageMutex.Unlock()
}

// GetCachedDiskUsage returns the cached disk usage data
func (c *Collector) GetCachedDiskUsage() []FolderSize {
	c.diskUsageMutex.RLock()
	defer c.diskUsageMutex.RUnlock()
	return c.cachedDiskUsage
}
