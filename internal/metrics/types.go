package metrics

import (
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
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
