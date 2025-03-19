package monitor

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/thenopholo/go_file_manager/config"
	"github.com/thenopholo/go_file_manager/logger"
)

type FileInfo struct {
	Path    string
	Size    int64
	ModTime time.Time
	Hash    string
	IsDir   bool
}

type FileMonitor struct {
	config   config.Config
	logger   *logger.Logger
	files    map[string]FileInfo
	mutex    sync.RWMutex
	isRunnig bool
	stopChan chan struct{}
}

func NewFileMonitor(cfg config.Config, log *logger.Logger) *FileMonitor {
	return &FileMonitor{
		config:   cfg,
		logger:   log,
		files:    make(map[string]FileInfo),
		stopChan: make(chan struct{}),
	}
}

func (m *FileMonitor) Start() error {
	m.mutex.Lock()
	if m.isRunnig {
		m.mutex.Unlock()
		return fmt.Errorf("o arquivo já está sendo monitorado")
	}

	m.isRunnig = true
	m.mutex.Unlock()

	if err := m.scanDirectory(); err != nil {
		m.mutex.Lock()
		m.isRunnig = false
		m.mutex.Unlock()
		return err
	}

	go m.monitorLoop()

	return nil
}

func (m *FileMonitor) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.isRunnig {
		close(m.stopChan)
		m.isRunnig = false
	}
}

func (m *FileMonitor) monitorLoop() {
	ticker := time.NewTicker(m.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.scanDirectory()
		case <-m.stopChan:
			return
		}
	}
}

func (m *FileMonitor) shouldIgnore(path string) bool {
	ext := filepath.Ext(path)
	for _, ignoreExt := range m.config.IgnoreExts {
		if ext == ignoreExt {
			return true
		}
	}
	return false
}

func calculateHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes), nil
}

func (m *FileMonitor) scanDirectory() error {
	m.mutex.Lock()
	oldFiles := make(map[string]FileInfo, len(m.files))
	for k, v := range m.files {
		oldFiles[k] = v
	}
	m.mutex.Unlock()

	currentFiles := make(map[string]FileInfo)

	err := filepath.Walk(m.config.WatchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path == m.config.LogDir {
			return filepath.SkipDir
		}

		if !info.IsDir() && m.shouldIgnore(path) {
			return nil
		}

		fileInfo := FileInfo{
			Path:    path,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}

		currentFiles[path] = fileInfo
		return nil
	})

	if err != nil {
		m.logger.Log(fmt.Sprintf("Erro ao escanear arquivo: %v", err))
		return err
	}

	m.detectChanges(oldFiles, currentFiles)

	m.mutex.Lock()
	m.files = currentFiles
	m.mutex.Unlock()

	return nil
}

func (m *FileMonitor) detectChanges(oldFiles, newFiles map[string]FileInfo) {
  for path, newInfo := range newFiles{
    oldInfo, exists := oldFiles[path]

    if !exists {
      m.logger.LogEvent("CRIADO", path)
    } else if !newInfo.IsDir && (newInfo.Size != oldInfo.Size || !newInfo.ModTime.Equal(oldInfo.ModTime) || (newInfo.Hash != "" && oldInfo.Hash != "" && newInfo.Hash != oldInfo.Hash)){
      m.logger.LogEvent("MODIFICADO", path)
    }

    for path := range oldFiles {
      if _, exists := newFiles[path]; !exists {
        m.logger.LogEvent("EXCLUÍDO", path)
      }
    }
  }
}

func (m *FileMonitor) GetFileCount() int {
  m.mutex.RLock()
  defer m.mutex.RUnlock()

  count := 0
  for _, info := range m.files{
    if !info.IsDir {
      count ++
    }
  }

  return count
}

func (m *FileMonitor) GetTotalSize() int64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var total int64
	for _, info := range m.files {
		if !info.IsDir {
			total += info.Size
		}
	}

	return total
}