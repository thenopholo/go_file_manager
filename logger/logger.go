package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/thenopholo/go_file_manager/config"
)

type Logger struct {
	config      config.Config
	file        *os.File
	mu          sync.Mutex
	currentSize int64
}

func NewLogger(cfg config.Config) (*Logger, error) {
	logger := &Logger{
		config: cfg,
	}

	if err := logger.openLogFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *Logger) openLogFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		l.file.Close()
	}

	fileName := fmt.Sprintf("file_monitor_%s.log", time.Now().Format("2006-01-02"))
	filePath := filepath.Join(l.config.LogDir, fileName)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("falha ao abrir o arquivo de log: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("falha ao obter informações do arquivo de log: %w", err)
	}

	l.file = file
	l.currentSize = info.Size()
	return nil
}

func (l *Logger) Log(message string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	formattedMsg := fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)

	bytesWritten, err := l.file.WriteString(formattedMsg)
	if err != nil {
		return fmt.Errorf("falha ao escrever no arquivo de log: %w", err)
	}

	l.currentSize += int64(bytesWritten)

	if l.currentSize > l.config.MaxLogSize {
		return l.openLogFile()
	}

	return nil
}

func (l *Logger) LogEvent(event, path string) error {
	message := fmt.Sprintf("EVENTO: %s | Arquivo: %s", event, path)
	return l.Log(message)
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}

	return nil
}
