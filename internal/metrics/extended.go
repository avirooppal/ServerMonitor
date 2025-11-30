package metrics

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// GetFail2BanStats parses the fail2ban log file
func (c *Collector) GetFail2BanStats(logPath string) (Fail2BanStats, error) {
	stats := Fail2BanStats{
		BansByIP: make(map[string]int),
	}

	file, err := os.Open(logPath)
	if err != nil {
		return stats, err
	}
	defer file.Close()

	// Regex to find "Ban <IP>"
	// Example: 2023-10-27 10:00:00,000 fail2ban.actions [123]: NOTICE [sshd] Ban 192.168.1.1
	banRegex := regexp.MustCompile(`Ban\s+(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	jailRegex := regexp.MustCompile(`\[(.*?)\]\s+Ban`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Ban") {
			stats.TotalBans++
			
			// Extract IP
			matches := banRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				ip := matches[1]
				stats.BansByIP[ip]++
			}

			// Extract Jail
			jailMatches := jailRegex.FindStringSubmatch(line)
			if len(jailMatches) > 1 {
				jail := jailMatches[1]
				found := false
				for _, j := range stats.Jails {
					if j == jail {
						found = true
						break
					}
				}
				if !found {
					stats.Jails = append(stats.Jails, jail)
				}
			}
		}
	}

	return stats, scanner.Err()
}

// GetDiskUsage calculates size of top-level folders in the given path
func (c *Collector) GetDiskUsage(path string) ([]FolderSize, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var folders []FolderSize

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		size, err := getDirSize(fullPath)
		if err != nil {
			// Skip permission errors etc
			continue
		}

		folders = append(folders, FolderSize{
			Path: entry.Name(),
			Size: size,
		})
	}

	// Sort by size desc
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Size > folders[j].Size
	})

	if len(folders) > 10 {
		folders = folders[:10]
	}

	return folders, nil
}

func getDirSize(path string) (uint64, error) {
	var size uint64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !info.IsDir() {
			size += uint64(info.Size())
		}
		return nil
	})
	return size, err
}

// GetAuthLogs parses auth.log for login attempts
func (c *Collector) GetAuthLogs(logPath string) ([]AuthLog, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var logs []AuthLog
	// Regex for standard linux auth.log
	// Example: Oct 27 10:00:00 host sshd[123]: Accepted password for user from 192.168.1.1 port 123
	// Example: Oct 27 10:00:00 host sshd[123]: Failed password for user from 192.168.1.1 port 123
	
	// Simplified regex for demo purposes
	userRegex := regexp.MustCompile(`for\s+(.*?)\s+from`)
	ipRegex := regexp.MustCompile(`from\s+(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

	scanner := bufio.NewScanner(file)
	// Read from end would be better, but for now read all and take last 50
	
	var allLogs []AuthLog

	for scanner.Scan() {
		line := scanner.Text()
		isAccepted := strings.Contains(line, "Accepted password")
		isFailed := strings.Contains(line, "Failed password")

		if isAccepted || isFailed {
			logEntry := AuthLog{
				Message: line,
				Success: isAccepted,
			}

			// Extract User
			userMatches := userRegex.FindStringSubmatch(line)
			if len(userMatches) > 1 {
				logEntry.User = userMatches[1]
			}

			// Extract IP
			ipMatches := ipRegex.FindStringSubmatch(line)
			if len(ipMatches) > 1 {
				logEntry.IP = ipMatches[1]
			}

			// Parse Time (assuming current year as syslog doesn't have year)
			// Format: Mon Jan 2 15:04:05
			if len(line) > 15 {
				t, err := time.Parse("Jan 2 15:04:05", line[:15])
				if err == nil {
					now := time.Now()
					logEntry.Time = t.AddDate(now.Year(), 0, 0)
				}
			}

			allLogs = append(allLogs, logEntry)
		}
	}

	// Return last 50
	if len(allLogs) > 50 {
		logs = allLogs[len(allLogs)-50:]
	} else {
		logs = allLogs
	}
	
	// Reverse to show newest first
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}

	return logs, scanner.Err()
}
