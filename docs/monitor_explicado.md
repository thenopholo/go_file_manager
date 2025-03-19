# Explicação Detalhada: Pacote de Monitoramento (monitor.go)

Este documento explica detalhadamente o funcionamento do arquivo `monitor.go`, que é responsável por monitorar alterações em arquivos e diretórios na aplicação.

## Estrutura do Código

### Importações

```go
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
```

- `crypto/md5`: Como uma máquina de impressão digital, cria uma "identidade única" para cada arquivo baseada em seu conteúdo.
- `encoding/hex`: Transforma os dados binários da impressão digital em texto legível (hexadecimal).
- `fmt`: Biblioteca para formatação de texto, como um editor de textos para nossas mensagens.
- `io`: Fornece ferramentas para mover dados entre diferentes partes do sistema, como um sistema de encanamento para informações.
- `os`: Principal pacote para interagir com o sistema operacional - é como nosso "embaixador" para se comunicar com o sistema de arquivos.
- `path/filepath`: Especialista em lidar com caminhos de arquivo, como um guia que conhece todos os atalhos e rotas da cidade.
- `sync`: Fornece mecanismos de controle de tráfego para evitar "colisões" quando múltiplas partes do programa acessam os mesmos dados.
- `time`: Como um relógio e calendário, ajuda a medir passagem de tempo e agendar tarefas.
- Pacotes locais `config` e `logger`: São como o manual de instruções e o diário de bordo da nossa aplicação.

### Estrutura FileInfo

```go
type FileInfo struct {
	Path    string
	Size    int64
	ModTime time.Time
	Hash    string
	IsDir   bool
}
```

Esta estrutura é como uma "ficha cadastral" para cada arquivo monitorado:

- `Path`: O endereço completo do arquivo (como o endereço residencial de uma pessoa)
- `Size`: Tamanho em bytes (como o peso de um objeto)
- `ModTime`: Quando foi modificado pela última vez (como a data de atualização de um documento)
- `Hash`: A "impressão digital" única do conteúdo do arquivo
- `IsDir`: Indica se é uma pasta ou um arquivo (como diferenciar entre uma casa e um apartamento)

### Estrutura FileMonitor

```go
type FileMonitor struct {
	config   config.Config
	logger   *logger.Logger
	files    map[string]FileInfo
	mutex    sync.RWMutex
	isRunnig bool
	stopChan chan struct{}
}
```

O FileMonitor é como um vigilante que supervisiona um bairro inteiro:

- `config`: As regras e instruções de como o monitoramento deve funcionar
- `logger`: O caderno de anotações onde registra tudo que observa
- `files`: Um mapa de todos os arquivos sob vigilância, como uma lista de endereços
- `mutex`: Um sistema de trancas para garantir acesso seguro às informações (como uma porta que só permite uma pessoa por vez)
- `isRunnig`: Uma flag indicando se a vigilância está ativa
- `stopChan`: Um canal de comunicação para sinalizar quando parar (como um rádio para receber ordens)

### Função NewFileMonitor

```go
func NewFileMonitor(cfg config.Config, log *logger.Logger) *FileMonitor {
	return &FileMonitor{
		config:   cfg,
		logger:   log,
		files:    make(map[string]FileInfo),
		stopChan: make(chan struct{}),
	}
}
```

Esta função é como contratar um novo vigilante:

1. Recebe as instruções de trabalho (config) e um diário para anotações (logger)
2. Prepara um mapa em branco para registrar os arquivos
3. Cria um canal de comunicação para futuras ordens de parada
4. Retorna o vigilante pronto para começar o trabalho

### Função Start

```go
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
```

Esta função é como o primeiro dia de trabalho do vigilante:

1. Tranca a porta do escritório (mutex.Lock) para verificar se já não tem alguém trabalhando
2. Se verificar que já está trabalhando, destrava a porta e diz "já estamos monitorando"
3. Caso contrário, marca que começou a trabalhar, destrava a porta
4. Faz uma primeira ronda completa (scanDirectory) para conhecer o território
5. Se algo der errado nessa ronda inicial, marca que não está mais trabalhando e reporta o problema
6. Se tudo correr bem, inicia uma rotina separada (goroutine) que fará rondas contínuas
7. Retorna indicando que começou o trabalho com sucesso

### Função Stop

```go
func (m *FileMonitor) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.isRunnig {
		close(m.stopChan)
		m.isRunnig = false
	}
}
```

Como encerrar o turno do vigilante:

1. Tranca a porta do escritório
2. Programa para destrancar automaticamente ao sair (defer)
3. Se o vigilante estiver trabalhando:
   - Envia um sinal pelo rádio (stopChan) para interromper as rondas
   - Marca que não está mais trabalhando

### Função monitorLoop

```go
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
```

Esta função é como uma rotina de rondas programadas:

1. Configura um relógio de alarme (ticker) que tocará em intervalos regulares
2. Garante que o alarme será desligado quando terminar o turno
3. Entra em um loop infinito onde:
   - Se o alarme tocar, faz uma nova ronda (scanDirectory)
   - Se receber um sinal no rádio (stopChan), encerra as rondas e volta para casa

### Função shouldIgnore

```go
func (m *FileMonitor) shouldIgnore(path string) bool {
	ext := filepath.Ext(path)
	for _, ignoreExt := range m.config.IgnoreExts {
		if ext == ignoreExt {
			return true
		}
	}
	return false
}
```

Como um filtro de "pessoas de interesse":

1. Examina a "característica" (extensão) do arquivo
2. Compara com uma lista de características para ignorar
3. Se for para ignorar, responde "sim, ignore este"
4. Caso contrário, responde "não, este precisa ser monitorado"

### Função calculateHash

```go
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
```

Como tirar a impressão digital de um documento:

1. Abre o documento (arquivo) para leitura
2. Garante que o documento será fechado quando terminar
3. Prepara a máquina de impressão digital (md5)
4. Passa o documento inteiro pela máquina
5. Coleta o resultado binário da impressão digital
6. Converte para um formato legível em texto e retorna

### Função scanDirectory

```go
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
```

Como fazer uma ronda completa pelo bairro:

1. Tranca o escritório para tirar uma foto da lista atual de arquivos
2. Faz uma cópia da lista atual (como tirar uma foto do estado atual)
3. Destrava o escritório para continuar o trabalho
4. Prepara uma nova lista vazia para a ronda atual
5. Caminha por todos os lugares do bairro (filepath.Walk) com estas regras:
   - Se encontrar problemas ao examinar um local, reporta o problema
   - Se encontrar o próprio escritório de logs, pula ele (não queremos vigiar a nós mesmos)
   - Se encontrar um tipo de arquivo que deve ser ignorado, continua a ronda
   - Para cada local válido, anota todas as características importantes
   - Adiciona estas informações à nova lista
6. Se tiver problemas durante a ronda, anota no diário e reporta
7. Compara a foto antiga com a nova situação para detectar mudanças
8. Tranca o escritório para atualizar a lista oficial
9. Substitui a lista antiga pela nova
10. Destrava o escritório e conclui a ronda

### Função detectChanges

```go
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
```

Como um jogo de "encontre as diferenças" entre duas fotos:

1. Examina cada item na foto nova:
   - Se não existia na foto antiga, anota "CRIADO" no diário
   - Se existia mas mudou de tamanho, data ou impressão digital, anota "MODIFICADO"
2. Examina cada item na foto antiga:
   - Se não existe mais na foto nova, anota "EXCLUÍDO" no diário

### Função GetFileCount

```go
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
```

Como contar quantos arquivos (não pastas) estão sendo monitorados:

1. Tranca o escritório em modo de leitura (permitindo que outros também leiam)
2. Garante que a porta será destrancada quando terminar
3. Prepara um contador começando do zero
4. Examina cada item na lista
5. Se for um arquivo (não uma pasta), incrementa o contador
6. Retorna o total de arquivos encontrados

### Função GetTotalSize

```go
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
```

Como calcular o tamanho total de todos os arquivos:

1. Tranca o escritório em modo de leitura
2. Garante que a porta será destrancada quando terminar
3. Prepara uma variável para somar o tamanho total
4. Examina cada item na lista
5. Se for um arquivo (não uma pasta), adiciona seu tamanho ao total
6. Retorna o tamanho total ocupado por todos os arquivos

## Analogia Geral do Sistema

O sistema de monitoramento funciona como um vigilante de segurança em um bairro:

1. Começa fazendo um mapeamento completo do território (arquivos e pastas)
2. Faz rondas regulares para verificar se algo mudou
3. Quando encontra mudanças (novos arquivos, modificações ou exclusões), registra no diário
4. Pode fornecer relatórios como "quantas casas estão no bairro" ou "qual é a área total ocupada"
5. Trabalha de forma organizada para não interferir com outros vigilantes (usando o sistema de trancas)
6. Sabe quais áreas pode ignorar (como o próprio escritório ou tipos específicos de construções)

O uso do pacote `os` é central neste sistema, pois ele permite:

- Abrir e ler arquivos (`os.Open`, `os.OpenFile`)
- Verificar informações sobre arquivos (`os.Stat`)
- Manipular caminhos de arquivos (`filepath.Walk`, `filepath.Ext`)
- Criar diretórios quando necessário

Esta implementação demonstra como construir um sistema de monitoramento de arquivos robusto usando Go, com ênfase na eficiência e na segurança de concorrência.
