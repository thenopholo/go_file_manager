# Explicação Detalhada: Pacote de Ações (actions.go)

Este documento explica detalhadamente o funcionamento do arquivo `actions.go`, que é responsável por executar ações em arquivos quando determinados eventos ocorrem na aplicação de monitoramento de arquivos.

## Estrutura do Código

### Importações

```go
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/thenopholo/go_file_manager/logger"
)
```

- `fmt`: Pacote para formatação de texto e strings, como um editor de textos para nossa mensagens.
- `os`: Pacote fundamental para interagir com o sistema operacional - é por ele que manipulamos arquivos e diretórios.
- `os/exec`: Permite executar comandos externos do sistema operacional, como se estivéssemos digitando no terminal.
- `path/filepath`: Especialista em manipular caminhos de arquivos de forma compatível com o sistema operacional.
- `strings`: Oferece funções para manipular textos, como um processador de palavras para nossas strings.
- `time`: Fornece funcionalidades relacionadas a tempo e datas, como um relógio e calendário para o código.
- `logger`: Nosso pacote local para registrar eventos e mensagens em arquivos de log.

### Tipos e Constantes

```go
type ActionType string

const (
	ActionBackup  ActionType = "backup"
	ActionArchive ActionType = "archive"
	ActionExecute ActionType = "execute"
)
```

Aqui definimos tipos de ações que nosso sistema pode executar:

- `ActionType`: Um tipo personalizado baseado em string para restringir os valores possíveis (como uma "etiqueta" que só pode ter certos valores).
- `ActionBackup`: Representa a ação de fazer uma cópia de segurança de um arquivo.
- `ActionArchive`: Representa a ação de mover um arquivo para um local de arquivamento (como guardar documentos antigos em um arquivo morto).
- `ActionExecute`: Representa a ação de executar um comando do sistema usando o arquivo.

### Estrutura Action

```go
type Action struct {
	Type    ActionType
	Target  string
	Command string
	Args    []string
	Logger  *logger.Logger
}
```

Esta estrutura é como um "formulário de instruções" que contém:

- `Type`: O tipo de ação a ser executada (como o título principal do formulário).
- `Target`: O diretório alvo para operações (como backup ou arquivamento).
- `Command`: O comando a ser executado (quando Type é ActionExecute).
- `Args`: Os argumentos para o comando (como opções adicionais para o comando).
- `Logger`: Um registrador para documentar o que acontece (como um escriba que anota tudo).

### Funções Construtoras

```go
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
```

Estas funções são como "fábricas especializadas" que criam ações pré-configuradas:

1. `NewBackupAction`: Cria uma ação de backup que copiará arquivos para o diretório alvo.
2. `NewArchiveAction`: Cria uma ação de arquivamento que moverá arquivos antigos para o diretório alvo.
3. `NewExecuteAction`: Cria uma ação que executará um comando externo com os argumentos especificados.

Cada função recebe apenas os parâmetros necessários para seu tipo específico de ação e configura os valores corretos para os outros campos.

### Função Execute

```go
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
```

Esta função é como um "supervisor" que delega o trabalho para o especialista correto:

1. Primeiro, registra no log que vai executar uma ação (como anunciar "vou começar a trabalhar em X").
2. Usa uma estrutura `switch` para determinar qual tipo de ação executar:
   - Para backup, chama `executeBackup`
   - Para arquivamento, chama `executeArchive`
   - Para executar comando, chama `executeCommand`
3. Se o tipo de ação for desconhecido, retorna um erro (como dizer "não sei como fazer isso").

O padrão utilizado aqui é semelhante ao "Strategy Pattern", onde diferentes implementações (estratégias) são selecionadas em tempo de execução.

### Função executeBackup

```go
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
```

Esta função é como um "arquivista digital" cuidadoso:

1. Define onde guardar a cópia de segurança:

   ```go
   backupDir := a.Target
   if backupDir == "" {
       backupDir = "backups"
   }
   ```

   - Usa o diretório alvo especificado ou, se não houver, usa "backups".

2. Cria o diretório de backup se não existir:

   ```go
   if err := os.MkdirAll(backupDir, 0755); err != nil {
       return fmt.Errorf("erro ao criar diretório de backup: %w", err)
   }
   ```

   - `os.MkdirAll` cria o diretório e todos os diretórios pais necessários.
   - `0755` são as permissões: proprietário pode ler/escrever/executar (7), outros podem ler/executar (5).
   - Se houver erro, retorna uma mensagem explicativa usando `fmt.Errorf`.

3. Prepara o nome do arquivo de backup:

   ```go
   fileName := filepath.Base(filePath)
   timestamp := time.Now().Format("20060102_150405")
   backupName := fmt.Sprintf("%s_%s", timestamp, fileName)
   backupPath := filepath.Join(backupDir, backupName)
   ```

   - `filepath.Base` extrai apenas o nome do arquivo sem o caminho.
   - Gera um timestamp no formato AAAAMMDD_HHMMSS.
   - Combina o timestamp com o nome original para criar um nome único.
   - `filepath.Join` cria um caminho completo de forma compatível com o sistema operacional.

4. Lê o conteúdo do arquivo original:

   ```go
   data, err := os.ReadFile(filePath)
   if err != nil {
       return fmt.Errorf("erro ao ler arquivo: %w", err)
   }
   ```

   - `os.ReadFile` lê todo o conteúdo do arquivo em memória.

5. Escreve o conteúdo no arquivo de backup:

   ```go
   if err := os.WriteFile(backupPath, data, 0644); err != nil {
       return fmt.Errorf("erro ao escrever arquivo de backup: %w", err)
   }
   ```

   - `os.WriteFile` cria um novo arquivo e escreve os dados nele.
   - `0644` são as permissões: proprietário pode ler/escrever (6), outros só podem ler (4).

6. Registra o sucesso e retorna:
   ```go
   a.Logger.Log(fmt.Sprintf("Backup criado: %s", backupPath))
   return nil
   ```

### Função executeArchive

```go
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
```

Esta função é como um "bibliotecário de arquivos antigos":

1. Verifica as informações do arquivo:

   ```go
   info, err := os.Stat(filePath)
   if err != nil {
       return fmt.Errorf("erro ao obter informações do arquivo: %w", err)
   }
   ```

   - `os.Stat` obtém metadados do arquivo, como tamanho e data de modificação.
   - É como examinar um documento para verificar quando foi escrito pela última vez.

2. Verifica se o arquivo é antigo o suficiente para ser arquivado:

   ```go
   if time.Since(info.ModTime()).Hours() < 24*30 {
       return nil
   }
   ```

   - Calcula quanto tempo passou desde a última modificação.
   - Se for menos de 30 dias (24 horas × 30), não faz nada.
   - É como verificar se um documento está na "quarentena" antes de enviá-lo para o arquivo.

3. Define o diretório de arquivamento:

   ```go
   archiveDir := a.Target
   if archiveDir == "" {
       archiveDir = "archive"
   }
   ```

   - Usa o diretório alvo especificado ou, se não houver, usa "archive".

4. Cria o diretório de arquivamento se necessário:

   ```go
   if err := os.MkdirAll(archiveDir, 0755); err != nil {
       return fmt.Errorf("erro ao criar diretório de arquivo: %w", err)
   }
   ```

   - Similar ao que fazemos no backup, garante que o destino exista.

5. Prepara o caminho do arquivo de destino:

   ```go
   fileName := filepath.Base(filePath)
   archivePath := filepath.Join(archiveDir, fileName)
   ```

6. Move o arquivo para o arquivamento:

   ```go
   if err := os.Rename(filePath, archivePath); err != nil {
       return fmt.Errorf("erro ao mover arquivo para o arquivo: %w", err)
   }
   ```

   - `os.Rename` move o arquivo de um local para outro.
   - É como mover fisicamente um documento de uma gaveta ativa para uma gaveta de arquivo morto.
   - A diferença importante em relação ao backup é que aqui o arquivo original é removido do local original.

7. Registra a ação completada:
   ```go
   a.Logger.Log(fmt.Sprintf("Arquivo movido para: %s", archivePath))
   return nil
   ```

### Função executeCommand

```go
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
```

Esta função é como um "assistente de linha de comando":

1. Prepara o comando substituindo placeholders:

   ```go
   command := strings.Replace(a.Command, "{file}", filePath, -1)
   ```

   - `strings.Replace` substitui todas as ocorrências (devido ao `-1`) de "{file}" pelo caminho real do arquivo.
   - É como preencher um template com os valores reais.

2. Prepara os argumentos do comando da mesma forma:

   ```go
   args := make([]string, len(a.Args))
   for i, arg := range a.Args {
       args[i] = strings.Replace(arg, "{file}", filePath, -1)
   }
   ```

   - Cria um novo slice do mesmo tamanho que os argumentos originais.
   - Para cada argumento, substitui "{file}" pelo caminho real do arquivo.

3. Cria o comando para execução:

   ```go
   cmd := exec.Command(command, args...)
   ```

   - `exec.Command` prepara um comando para execução, similar a digitar no terminal.
   - É como preparar a instrução que um assistente vai executar no sistema.

4. Executa o comando e captura a saída:

   ```go
   output, err := cmd.CombinedOutput()
   ```

   - `CombinedOutput()` executa o comando e captura tanto a saída padrão quanto os erros.
   - É como executar um comando no terminal e capturar tudo que é exibido na tela.

5. Verifica se houve erro na execução:

   ```go
   if err != nil {
       a.Logger.Log(fmt.Sprintf("Erro ao executar comando: %v - Saída: %s", err, output))
       return fmt.Errorf("erro ao executar comando: %w", err)
   }
   ```

   - Se o comando falhar, registra o erro e a saída no log.
   - Retorna um erro formatado com o erro original encapsulado (`%w`).

6. Se tudo correu bem, registra o sucesso:
   ```go
   a.Logger.Log(fmt.Sprintf("Comando executado com sucesso. Saída: %s", output))
   return nil
   ```

## Analogia Geral do Sistema de Ações

O pacote de ações funciona como um conjunto de ferramentas especializadas em uma oficina:

1. **Backup (executeBackup)**: Como um fotógrafo que tira fotos dos documentos para preservar seu estado atual, sem alterar os originais.

2. **Archive (executeArchive)**: Como um arquivista que move documentos antigos para uma sala de arquivos, liberando espaço na área de trabalho principal, mas mantendo os documentos acessíveis se necessário.

3. **Execute (executeCommand)**: Como um assistente que segue instruções específicas para trabalhar com os documentos, como "por favor, converta este arquivo para PDF" ou "comprima este arquivo".

O uso do pacote `os` é fundamental neste sistema para:

- Verificar existência de arquivos e diretórios (`os.Stat`)
- Criar diretórios quando necessário (`os.MkdirAll`)
- Ler o conteúdo de arquivos (`os.ReadFile`)
- Escrever em novos arquivos (`os.WriteFile`)
- Mover arquivos (`os.Rename`)
- Executar comandos externos (`os/exec.Command`)

Esta implementação demonstra a flexibilidade do Go para manipulação de arquivos e execução de comandos do sistema operacional, permitindo automatizar diferentes tipos de tarefas relacionadas a arquivos de forma segura e eficiente.

## Como Usar este Módulo

```go
// Inicializar configuração e logger
cfg := config.LoadConfig()
logger, _ := logger.NewLogger(cfg)

// Criar ações
backupAction := actions.NewBackupAction("./backups", logger)
archiveAction := actions.NewArchiveAction("./archive", logger)
execAction := actions.NewExecuteAction("zip", []string{"-r", "arquivo.zip", "{file}"}, logger)

// Executar ações para um arquivo
err := backupAction.Execute("/caminho/para/arquivo.txt")
if err != nil {
    log.Fatalf("Erro ao fazer backup: %v", err)
}

// Tentar arquivar se for antigo
err = archiveAction.Execute("/caminho/para/arquivo_antigo.txt")
if err != nil {
    log.Fatalf("Erro ao arquivar: %v", err)
}

// Executar comando zip no arquivo
err = execAction.Execute("/caminho/para/diretorio")
if err != nil {
    log.Fatalf("Erro ao executar comando: %v", err)
}
```

Este exemplo mostra como o módulo de ações pode ser utilizado para automatizar diferentes tipos de operações em arquivos.

```

## Uso do Pacote os

O arquivo `actions.go` é um excelente exemplo de uso prático do pacote `os` do Go para diversas operações de manipulação de arquivos:

1. **Verificação de arquivos**: Usando `os.Stat()` para obter informações de arquivos existentes
2. **Criação de diretórios**: Utilizando `os.MkdirAll()` para criar estruturas de diretórios
3. **Leitura de arquivos**: Com `os.ReadFile()` para carregar o conteúdo completo de um arquivo
4. **Escrita de arquivos**: Usando `os.WriteFile()` para salvar dados em novos arquivos
5. **Movimentação de arquivos**: Através de `os.Rename()` para mover arquivos entre diretórios

Essas operações formam a base de um sistema de gerenciamento de arquivos eficiente.
```
