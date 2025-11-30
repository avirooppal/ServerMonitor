package main

import (
	"log"
	"github.com/user/server-moni/internal/db"
	"github.com/user/server-moni/internal/metrics"
)

func snapshotDisk(collector *metrics.Collector) {
	// We want overall disk usage. Let's take the root partition "/"
	// In Docker, we might want the host fs if mounted.
	// The collector.Collect() gets all disks.
	m := collector.Collect()
	
	var total, used uint64
	var usedPercent float64
	found := false

	// Try to find the root or hostfs
	for _, d := range m.Disks {
		if d.Path == "/" || d.Path == "/hostfs" {
			total = d.Total
			used = d.Used
			usedPercent = d.UsedPercent
			found = true
			break
		}
	}

	// If not found, just take the first one?
	if !found && len(m.Disks) > 0 {
		total = m.Disks[0].Total
		used = m.Disks[0].Used
		usedPercent = m.Disks[0].UsedPercent
	}

	if total > 0 {
		err := db.AddDiskHistory(usedPercent, total, used)
		if err != nil {
			log.Printf("Failed to save disk history: %v", err)
		} else {
			log.Printf("Saved daily disk snapshot: %.2f%% used", usedPercent)
		}
	}
}
