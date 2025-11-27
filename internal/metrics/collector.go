package metrics

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type SystemMetrics struct {
	CPU          []float64              `json:"cpu"` // Per core
	CPUTotal     float64                `json:"cpu_total"`
	LoadAvg      *load.AvgStat          `json:"load_avg"`
	Memory       *ExtendedMemoryStat    `json:"memory"`
	Swap         *mem.SwapMemoryStat    `json:"swap"`
	Disks        []DiskInfo             `json:"disks"`
	Network      NetworkStats           `json:"network"`
	Processes    []ProcessInfo          `json:"processes"`
	Containers   []ContainerInfo        `json:"containers"`
	HostInfo     *host.InfoStat         `json:"host_info"`
	LastUpdate   time.Time              `json:"last_update"`
}

type ExtendedMemoryStat struct {
	*mem.VirtualMemoryStat
	Buffers uint64 `json:"buffers"`
	Cached  uint64 `json:"cached"`
}

type DiskInfo struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
	ReadRate    uint64  `json:"read_rate"`  // Bytes per second
	WriteRate   uint64  `json:"write_rate"` // Bytes per second
}

type NetworkStats struct {
	Interfaces []NetInterface `json:"interfaces"`
	TotalRecv  uint64         `json:"total_recv"` // Rate B/s
	TotalSent  uint64         `json:"total_sent"` // Rate B/s
}

type NetInterface struct {
	Name      string `json:"name"`
	RecvRate  uint64 `json:"recv_rate"`
	SentRate  uint64 `json:"sent_rate"`
}

type ProcessInfo struct {
	PID      int32   `json:"pid"`
	Name     string  `json:"name"`
	CPU      float64 `json:"cpu"`
	Mem      float32 `json:"mem"` // Percent
	Username string  `json:"username"`
}

type ContainerInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	State       string  `json:"state"`
	Status      string  `json:"status"`
	Created     int64   `json:"created"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryUsage uint64  `json:"memory_usage"`
	MemoryLimit uint64  `json:"memory_limit"`
}

var (
	currentMetrics SystemMetrics
	mutex          sync.RWMutex
	lastNetIO      []net.IOCountersStat
	lastDiskIO     map[string]disk.IOCountersStat
	lastTime       time.Time
	dockerClient   *client.Client
)

func StartCollector(interval time.Duration) {
	// Initial Host Info
	h, _ := host.Info()
	mutex.Lock()
	currentMetrics.HostInfo = h
	mutex.Unlock()

	// Init Docker Client
	var err error
	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Failed to create Docker client: %v", err)
	} else {
		// Test connection
		_, err := dockerClient.Ping(context.Background())
		if err != nil {
			log.Printf("Docker not available: %v", err)
			dockerClient = nil
		} else {
			log.Println("Docker client connected")
		}
	}

	lastDiskIO = make(map[string]disk.IOCountersStat)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		collect()
	}
}

func GetMetrics() SystemMetrics {
	mutex.RLock()
	defer mutex.RUnlock()
	return currentMetrics
}

func collect() {
	now := time.Now()
	timeDiff := now.Sub(lastTime).Seconds()
	
	// CPU
	c, _ := cpu.Percent(0, true)
	cTotal, _ := cpu.Percent(0, false)
	l, _ := load.Avg()

	// Memory
	m, _ := mem.VirtualMemory()
	s, _ := mem.SwapMemory()
	extMem := &ExtendedMemoryStat{
		VirtualMemoryStat: m,
		Buffers:           m.Buffers,
		Cached:            m.Cached,
	}

	// Disks (Usage & I/O)
	parts, _ := disk.Partitions(false)
	var disks []DiskInfo
	ioCounters, _ := disk.IOCounters()

	for _, p := range parts {
		u, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		
		// Calculate I/O Rates
		var rRate, wRate uint64
		// Find matching IO counter (heuristic matching by device name)
		// p.Device is like "/dev/sda1", ioCounters keys are like "sda", "sda1"
		deviceName := p.Device
		if strings.HasPrefix(deviceName, "/dev/") {
			deviceName = strings.TrimPrefix(deviceName, "/dev/")
		}
		
		if curIO, ok := ioCounters[deviceName]; ok {
			if prevIO, ok := lastDiskIO[deviceName]; ok && timeDiff > 0 {
				rRate = uint64(float64(curIO.ReadBytes-prevIO.ReadBytes) / timeDiff)
				wRate = uint64(float64(curIO.WriteBytes-prevIO.WriteBytes) / timeDiff)
			}
			lastDiskIO[deviceName] = curIO
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
		disks = append(disks, dInfo)
	}
	
	// Network
	netIO, _ := net.IOCounters(true)
	var netStats NetworkStats
	var totalRecv, totalSent uint64
	
	if timeDiff > 0 {
		// Calculate Net Rates
		for _, cur := range netIO {
			var prev net.IOCountersStat
			found := false
			for _, p := range lastNetIO {
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
	lastNetIO = netIO
	lastTime = now

	// Processes (Top 20 by CPU)
	procs, _ := process.Processes()
	var procInfos []ProcessInfo
	for _, p := range procs {
		cpuP, _ := p.CPUPercent()
		if cpuP < 0.1 { continue } // Optimization
		
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

	// Docker Containers with Stats
	var containerInfos []ContainerInfo
	if dockerClient != nil {
		containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
		if err == nil {
			for _, ctr := range containers {
				name := "unknown"
				if len(ctr.Names) > 0 {
					name = ctr.Names[0]
				}
				
				var cpuPercent float64
				var memUsage, memLimit uint64

				// Only fetch stats for running containers to save resources
				if ctr.State == "running" {
					stats, err := dockerClient.ContainerStats(context.Background(), ctr.ID, false)
					if err == nil {
						// Use a simple map for decoding to avoid importing full types
						var statsData map[string]interface{}
						if err := json.NewDecoder(stats.Body).Decode(&statsData); err == nil {
							// Calculate CPU %
							// This is complex, simplified version:
							// (cpu_delta / system_cpu_delta) * number_cpus * 100.0
							// For now, let's just try to get basic memory usage.
							if mem, ok := statsData["memory_stats"].(map[string]interface{}); ok {
								if usage, ok := mem["usage"].(float64); ok {
									memUsage = uint64(usage)
								}
								if limit, ok := mem["limit"].(float64); ok {
									memLimit = uint64(limit)
								}
							}
							
							// CPU calculation (very simplified placeholder)
							// Proper calculation requires previous stats.
							// For this iteration, we might skip complex CPU calc or implement it fully later.
							// Let's try to get cpu_stats.cpu_usage.total_usage
							if cpuObj, ok := statsData["cpu_stats"].(map[string]interface{}); ok {
								if cpuUsage, ok := cpuObj["cpu_usage"].(map[string]interface{}); ok {
									if total, ok := cpuUsage["total_usage"].(float64); ok {
										// This is cumulative, need delta. 
										// Storing state per container is complex for this single function.
										// We will leave CPU as 0 for now or implement a stateful collector later.
										_ = total
									}
								}
							}
						}
						stats.Body.Close()
					}
				}

				containerInfos = append(containerInfos, ContainerInfo{
					ID:          ctr.ID[:12],
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
	
	mutex.Lock()
	currentMetrics.CPU = c
	if len(cTotal) > 0 {
		currentMetrics.CPUTotal = cTotal[0]
	}
	currentMetrics.LoadAvg = l
	currentMetrics.Memory = extMem
	currentMetrics.Swap = s
	currentMetrics.Disks = disks
	currentMetrics.Network = netStats
	currentMetrics.Processes = procInfos
	currentMetrics.Containers = containerInfos
	currentMetrics.LastUpdate = now
	mutex.Unlock()
}
