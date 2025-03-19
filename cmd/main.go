package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thenopholo/go_file_manager/config"
	"github.com/thenopholo/go_file_manager/logger"
	"github.com/thenopholo/go_file_manager/monitor"
	"github.com/thenopholo/go_file_manager/stats"
)

func main() {
	fmt.Println("||---------Iniciando sistema de monitoramento de arquivos---------||")

	cfg := config.LoadConfig()
	fmt.Printf("Monitorando diretório: %s\n", cfg.WatchDir)
	fmt.Printf("Intervalo de verificação: %s\n", cfg.CheckInterval)

	log, err := logger.NewLogger(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Log("||---------Sistema de monitoramento iniciado---------||")

	fileMonitor := monitor.NewFileMonitor(cfg, log)
	if err := fileMonitor.Start(); err != nil {
		log.Log(fmt.Sprintf("Erro ao iniciar monitoramento: %v", err))
		fmt.Fprintf(os.Stderr, "Erro ao iniciar monitoramento: %v\n", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	statsTicker := time.NewTicker(1 * time.Hour)
	defer statsTicker.Stop()

	fmt.Println("Sistema em execução. Pressione Ctrl+C para encerrar.")

	running := true
	for running {
		select {
		case <-statsTicker.C:
			statsData, err := stats.GenerateStats(fileMonitor, cfg)
			if err != nil {
				log.Log(fmt.Sprintf("Erro ao gerar estatísticas: %v", err))
			} else {
				if err := stats.SaveStatsToFile(statsData, cfg); err != nil {
					log.Log(fmt.Sprintf("Erro ao salvar estatísticas: %v", err))
				} else {
					log.Log(fmt.Sprintf("Estatísticas geradas: %d arquivos, %s",
						statsData.FileCount, statsData.TotalSizeHuman))
				}
			}
		case <-sigCh:
			fmt.Println("\nSinal de interrupção recebido. Encerrando...")
			log.Log("Sistema de monitoramento encerrando")
			fileMonitor.Stop()
			running = false
		}
	}

	finalStats, _ := stats.GenerateStats(fileMonitor, cfg)
	if finalStats != nil {
		if err := stats.SaveStatsToFile(finalStats, cfg); err != nil {
			log.Log(fmt.Sprintf("Erro ao salvar relatório final: %v", err))
		} else {
			fmt.Printf("Relatório final: %d arquivos, %s\n",
				finalStats.FileCount, finalStats.TotalSizeHuman)
		}
	}

	fmt.Println("Sistema encerrado com sucesso.")
}
