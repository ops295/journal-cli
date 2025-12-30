package stats

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Stats holds journal statistics
type Stats struct {
	TotalEntries int
	LastMissed   time.Time
}

// GetStats calculates statistics for the given journal directory.
func GetStats(journalDir string) (Stats, error) {
	var stats Stats

	// 1. Count Total Entries
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		if os.IsNotExist(err) {
			return stats, nil
		}
		return stats, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			stats.TotalEntries++
		}
	}

	// 2. Find Last Missed Date
	// Iterate backwards from yesterday up to 30 days.
	// We skip today because the user might just be starting to journal.
	now := time.Now()
	for i := 1; i <= 30; i++ {
		date := now.AddDate(0, 0, -i)
		filename := date.Format("2006-01-02") + ".md"
		path := filepath.Join(journalDir, filename)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			stats.LastMissed = date
			break
		}
	}

	return stats, nil
}
