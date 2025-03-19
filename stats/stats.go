package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/thenopholo/go_file_manager/config"
	"github.com/thenopholo/go_file_manager/monitor"
)

type Stats struct {
	Timestamp      time.Time           `json:"timestamp"`
	FileCount      int                 `json:"file_count"`
	TotalSize      int64               `json:"total_size"`
	TotalSizeHuman string              `json:"total_size_human"`
	ByExtentions   map[string]ExtStats `json:"by_extention"`
}

type ExtStats struct {
	Count int   `json:"count"`
	Size  int64 `json:"size_bytes"`
}

func GenerateStats(monitor *monitor.FileMonitor, cfg config.Config) (*Stats, error) {
	stats := &Stats{
		Timestamp:    time.Now(),
		FileCount:    monitor.GetFileCount(),
		TotalSize:    monitor.GetTotalSize(),
		ByExtentions: make(map[string]ExtStats),
	}

	stats.TotalSizeHuman = formatSize(stats.TotalSize)

	return stats, nil
}

func SaveStatsToFile(stats *Stats, cfg config.Config) error {
	statsDir := filepath.Join(cfg.LogDir, "stats")

	if _, err := os.Stat(statsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(statsDir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório de estatística: %w", err)
		}
	}

	fileName := fmt.Sprintf("stats_%s.json", time.Now().Format("2006-01-02_15-04-05"))
	filePath := filepath.Join(statsDir, fileName)

	jsonData, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar estatísticas: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo de estatísticas: %w", err)
	}

	return nil
}

func formatSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
