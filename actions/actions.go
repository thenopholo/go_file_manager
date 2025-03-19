package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/thenopholo/go_file_manager/logger"
)

type ActionType string

const (
	ActionBackup  ActionType = "backup"
	ActionArchive ActionType = "archive"
	ActionExecute ActionType = "execute"
)

type Action struct {
	Type    ActionType
	Target  string
	Command string
	Args    []string
	Logger  *logger.Logger
}

func NewBackupAction(target string, log *logger.Logger) *Action {
	return &Action{
		Type:   ActionBackup,
		Target: target,
		Logger: log,
	}
}

func NewArchiveAction(target string, log *logger.Logger) *Action {
	return &Action{
		Type:   ActionArchive,
		Target: target,
		Logger: log,
	}
}

func NewExecuteAction(command string, args []string, log *logger.Logger) *Action {
	return &Action{
		Type:    ActionExecute,
		Command: command,
		Args:    args,
		Logger:  log,
	}
}

func (a *Action) Execute(filePath string) error {
	a.Logger.Log(fmt.Sprintf("Executando ação %s para %s", a.Type, filePath))

	switch a.Type {
	case ActionBackup:
		return a.executeBackup(filePath)
	case ActionArchive:
		return a.executeArchive(filePath)
	case ActionExecute:
		return a.executeCommand(filePath)
	default:
		return fmt.Errorf("tipo de ação desconhecido: %s", a.Type)
	}
}

func (a *Action) executeBackup(filePath string) error {
	backupDir := a.Target
	if backupDir == "" {
		backupDir = "backups"
	}

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de backup: %w", err)
	}

	fileName := filepath.Base(filePath)
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s_%s", timestamp, fileName)
	backupPath := filepath.Join(backupDir, backupName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo de backup: %w", err)
	}

	a.Logger.Log(fmt.Sprintf("Backup criado: %s", backupPath))
	return nil
}

func (a *Action) executeArchive(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("erro ao obter informações do arquivo: %w", err)
	}

	if time.Since(info.ModTime()).Hours() < 24*30 {
		return nil
	}

	archiveDir := a.Target
	if archiveDir == "" {
		archiveDir = "archive"
	}

	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de arquivo: %w", err)
	}

	fileName := filepath.Base(filePath)
	archivePath := filepath.Join(archiveDir, fileName)

	if err := os.Rename(filePath, archivePath); err != nil {
		return fmt.Errorf("erro ao mover arquivo para o arquivo: %w", err)
	}

	a.Logger.Log(fmt.Sprintf("Arquivo movido para: %s", archivePath))
	return nil
}

func (a *Action) executeCommand(filePath string) error {
	command := strings.Replace(a.Command, "{file}", filePath, -1)

	args := make([]string, len(a.Args))
	for i, arg := range a.Args {
		args[i] = strings.Replace(arg, "{file}", filePath, -1)
	}

	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		a.Logger.Log(fmt.Sprintf("Erro ao executar comando: %v - Saída: %s", err, output))
		return fmt.Errorf("erro ao executar comando: %w", err)
	}

	a.Logger.Log(fmt.Sprintf("Comando executado com sucesso. Saída: %s", output))
	return nil
}
