# Explicação Detalhada: Pacote de Logger (logger.go)

Este documento explica detalhadamente o funcionamento do arquivo `logger.go`, que é responsável por registrar eventos e mensagens em arquivos de log da aplicação de monitoramento de arquivos.

## Estrutura do Código

### Importações

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/thenopholo/go_file_manager/config"
)
```

- `fmt`: Pacote para formatação de strings e saída.
- `os`: Fornece funções para interagir com arquivos e o sistema operacional.
- `path/filepath`: Manipula caminhos de arquivos.
- `sync`: Fornece primitivas de sincronização como mutex.
- `time`: Permite trabalhar com datas e horas.
- `config`: Pacote local que contém as configurações da aplicação.

### Estrutura Logger

```go
type Logger struct {
	config      config.Config
	file        *os.File
	mu          sync.Mutex
	currentSize int64
}
```

Esta estrutura é como um "escritor de diário" especializado:

- `config`: Contém as configurações, como onde guardar os logs e qual o tamanho máximo (como as regras para escrever no diário).
- `file`: O arquivo atual onde estamos escrevendo (o diário em si).
- `mu`: Um mutex para evitar que múltiplas partes do código escrevam no mesmo arquivo simultaneamente (como uma trava na porta do escritório).
- `currentSize`: Controla quanto já escrevemos no arquivo atual (como contar quantas páginas já usamos no diário).

### Função NewLogger

```go
func NewLogger(cfg config.Config) (*Logger, error) {
	logger := &Logger{
		config: cfg,
	}

	if err := logger.openLogFile(); err != nil {
		return nil, err
	}

	return logger, nil
}
```

Esta função é como "contratar um novo escritor de diário":

1. Cria uma nova instância de Logger com as configurações fornecidas
2. Tenta abrir um arquivo de log (como dar um novo diário ao escritor)
3. Se conseguir abrir o arquivo, retorna o logger pronto para uso
4. Se falhar, retorna um erro (como dizer "não conseguimos preparar seu escritor")

### Função openLogFile

```go
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
```

Esta função é como "preparar um novo diário para escrever":

1. `l.mu.Lock()`: Tranca a porta do escritório para que ninguém mais interfira
2. `defer l.mu.Unlock()`: Garante que a porta será destrancada quando terminarmos, mesmo se ocorrer um erro
3. `if l.file != nil { l.file.Close() }`: Se já estiver com um diário aberto, fecha-o primeiro
4. `fileName := ...`: Cria um nome para o arquivo baseado na data atual (como "diário de 2023-10-15")
   - A data `2006-01-02` é um formato especial em Go que representa "YYYY-MM-DD"
5. `filePath := ...`: Monta o caminho completo do arquivo
6. `file, err := os.OpenFile(...)`: Abre ou cria o arquivo com estas opções:
   - `O_APPEND`: Adiciona ao final do arquivo se ele existir (continua escrevendo no final do diário)
   - `O_CREATE`: Cria o arquivo se não existir (compra um novo diário se não tiver um)
   - `O_WRONLY`: Abre apenas para escrita (o diário é só para registrar, não para ler)
   - `0644`: Permissões do arquivo (o proprietário pode ler e escrever, outros só podem ler)
7. Verifica informações sobre o arquivo para saber seu tamanho atual
8. Atualiza o estado do logger com o novo arquivo e seu tamanho

### Função Log

```go
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
```

Esta função é como "escrever uma entrada no diário":

1. `l.mu.Lock()`: Tranca o escritório para uso exclusivo
2. `defer l.mu.Unlock()`: Garante que destrancará quando terminar
3. `formattedMsg := ...`: Formata a mensagem com data e hora (como escrever "15 de outubro de 2023, 14:30 - Aconteceu algo")
4. `bytesWritten, err := l.file.WriteString(...)`: Escreve a mensagem no arquivo e conta quantos bytes foram escritos
5. `l.currentSize += int64(bytesWritten)`: Atualiza o contador de tamanho do arquivo
6. `if l.currentSize > l.config.MaxLogSize { ... }`: Se o arquivo ficou muito grande, abre um novo
   - É como dizer "este diário está cheio, preciso de um novo"

### Função LogEvent

```go
func (l *Logger) LogEvent(event, path string) error {
	message := fmt.Sprintf("EVENTO: %s | Arquivo: %s", event, path)
	return l.Log(message)
}
```

Esta função é uma especialização da função `Log` para registrar eventos relacionados a arquivos:

1. Formata uma mensagem específica para eventos (como "EVENTO: CRIADO | Arquivo: /documentos/relatorio.txt")
2. Chama a função `Log` regular para registrar esta mensagem formatada

### Função Close

```go
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}

	return nil
}
```

Esta função é como "encerrar o trabalho do escritor de diário":

1. `l.mu.Lock()`: Tranca o escritório uma última vez
2. `defer l.mu.Unlock()`: Garante que destrancará mesmo se ocorrer um erro
3. `if l.file != nil { return l.file.Close() }`: Se tiver um arquivo aberto, fecha-o corretamente
   - Como guardar o diário na gaveta antes de ir embora

## Funcionalidades Importantes

### Rotação de Logs

O sistema implementa rotação de logs de duas maneiras:

1. **Por Tamanho**: Quando um arquivo de log atinge o tamanho máximo configurado, um novo arquivo é aberto
2. **Por Data**: Os nomes dos arquivos incluem a data, então cada dia naturalmente começa com um novo arquivo

### Thread Safety

O uso de mutex (`sync.Mutex`) garante que o logger possa ser utilizado com segurança em um ambiente com múltiplas goroutines tentando escrever no mesmo arquivo de log simultaneamente.

## Analogia Geral

O logger funciona como um assistente dedicado que mantém um registro detalhado de todas as atividades:

- Ele data e hora cada entrada (como um registro oficial)
- Organiza os registros em arquivos diários
- Quando um arquivo fica muito grande, ele começa um novo
- Mantém tudo organizado na pasta designada
- Garante que apenas uma pessoa escreva no registro por vez

Este design é robusto e evita problemas comuns em logs, como arquivos que crescem infinitamente ou corrupção devido a escritas simultâneas.

## Como Usar

```go
// Inicializar o logger
cfg := config.LoadConfig()
logger, err := logger.NewLogger(cfg)
if err != nil {
    panic(err)
}
defer logger.Close()  // Garante que o arquivo será fechado ao final

// Registrar uma mensagem geral
logger.Log("Aplicação iniciada")

// Registrar um evento específico
logger.LogEvent("CRIADO", "/caminho/para/arquivo.txt")
```

Os logs serão criados na pasta configurada em `LogDir`, com nomes como `file_monitor_2023-10-15.log`.
