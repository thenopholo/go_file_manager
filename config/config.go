package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	WatchDir      string
	LogDir        string
	CheckInterval time.Duration
	MaxLogSize    int64
	AutoAction    bool
	IgnoreExts    []string
}

func LoadConfig() Config {
  config := Config {
  	WatchDir:      ".",
  	LogDir:        "./logs",
  	CheckInterval: 5 * time.Second,
  	MaxLogSize:    1024 * 1024 * 10, // 10MB
  	AutoAction:    false,
  	IgnoreExts:    []string{".temp", ".swp"},
  }

  if dir := os.Getenv("WATCH_DIR"); dir != "" {
    config.WatchDir = dir
  }

  if dir := os.Getenv("LOG_DIR"); dir != "" {
    config.LogDir = dir
  }

  if interval := os.Getenv("CHECK_INTERVAL"); interval != "" {
    if seconds, err := strconv.Atoi(interval); err == nil{
      config.CheckInterval = time.Duration(seconds) * time.Second
    }
  }

  if size := os.Getenv("MAX_LOG_SIZE"); size != "" {
    if bytes, err := strconv.ParseInt(size, 10, 64); err == nil {
      config.MaxLogSize = bytes
    }
  }

  if autoAction := os.Getenv("AUTO_ACTION"); autoAction == "true" {
    config.AutoAction = true
  }

  if ignoreExist := os.Getenv("IGNORE_EXTS"); ignoreExist != "" {
    config.IgnoreExts = filepath.SplitList(ignoreExist)
  }

  ensureDir(config.WatchDir)
  ensureDir(config.LogDir)

  return config
}

func ensureDir(dir string) {
  if _, err := os.Stat(dir); os.IsNotExist(err) {
    os.MkdirAll(dir, 0755)
  }
}