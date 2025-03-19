# Go File Manager

Um sistema de monitoramento de arquivos robusto, desenvolvido em Go, que observa diretórios, detecta mudanças em arquivos, registra eventos, gera estatísticas e pode executar ações automatizadas.

## Índice

- [Visão Geral](#visão-geral)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Conceitos Go Utilizados](#conceitos-go-utilizados)
- [Instalação](#instalação)
- [Como Usar](#como-usar)
- [Configuração](#configuração)
- [Documentação Detalhada](#documentação-detalhada)
- [Exemplos Práticos](#exemplos-práticos)
- [Conceitos Avançados](#conceitos-avançados)

## Visão Geral

O Go File Manager é um sistema que monitora continuamente um diretório e seus subdiretórios para detectar mudanças em arquivos:

- Detecta arquivos criados, modificados ou excluídos
- Gera logs detalhados de todas as alterações
- Produz estatísticas sobre os arquivos monitorados
- Permite executar ações automáticas como backup, arquivamento e execução de comandos

Este projeto demonstra vários conceitos importantes do Go como goroutines para concorrência, uso eficiente do sistema de arquivos, controle de fluxo, manipulação de erros e estruturação de projetos.

## Estrutura do Projeto

O projeto está organizado em vários pacotes:

- **cmd**: Contém o ponto de entrada principal (`main.go`)
- **config**: Gerencia as configurações do sistema
- **logger**: Implementa o sistema de registro de eventos
- **monitor**: Monitora diretórios e detecta mudanças nos arquivos
- **stats**: Coleta e gera estatísticas sobre os arquivos monitorados
- **actions**: Executa ações nos arquivos (backup, arquivamento, comandos)

Cada componente foi projetado para ser modular e com responsabilidades bem definidas, seguindo o princípio da responsabilidade única.

## Conceitos Go Utilizados

### Concorrência com Goroutines

O sistema utiliza goroutines para executar o monitoramento em segundo plano, sem bloquear o programa principal:

```go
// Exemplo do monitor que executa em uma goroutine separada
go m.monitorLoop()
```

### Canais para Comunicação

Canais são usados para comunicação entre goroutines e para sincronização:

```go
// Criação de canal para sinalização
stopChan := make(chan struct{})

// Uso do canal para parar a execução
select {
case <-ticker.C:
    // Executa algo periodicamente
case <-stopChan:
    // Termina a execução quando o canal é fechado
}
```

### Mutex para Controle de Acesso

O sistema usa mutex para garantir acesso seguro a dados compartilhados em ambiente concorrente:

```go
// Adquire lock para leitura
m.mutex.RLock()
defer m.mutex.RUnlock()

// Código que lê dados compartilhados
```

### Interfaces e Composição

O projeto utiliza interfaces para desacoplamento e composição para reusabilidade:

```go
// A estrutura Logger é composta por outros componentes
type Logger struct {
    config      config.Config
    file        *os.File
    mu          sync.Mutex
    currentSize int64
}
```

### Tratamento de Erros

O projeto demonstra práticas robustas de tratamento de erros:

```go
if err := something(); err != nil {
    return fmt.Errorf("erro ao fazer algo: %w", err)
}
```

### Manipulação de Arquivos e Diretórios

O sistema faz uso extensivo do pacote `os` para interagir com o sistema de arquivos:

```go
// Lendo um arquivo
data, err := os.ReadFile(filePath)

// Verificando existência de um diretório
if _, err := os.Stat(dir); os.IsNotExist(err) {
    // O diretório não existe
}

// Percorrendo recursivamente um diretório
filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
    // Processar cada arquivo/diretório
    return nil
})
```

## Instalação

Para instalar o Go File Manager:

```bash
# Clone o repositório
git clone https://github.com/thenopholo/go_file_manager.git

# Entre no diretório do projeto
cd go_file_manager

# Compile o programa
go build -o filemanager ./cmd/main.go
```

## Como Usar

### Uso Básico

```bash
# Execute o programa para monitorar o diretório atual
./filemanager

# Execute especificando um diretório para monitorar
WATCH_DIR="/caminho/para/diretório" ./filemanager
```

### Testando em uma Pasta Local

Para testar o programa em uma pasta local específica:

```bash
# Crie um diretório de teste
mkdir ~/test_monitor

# Configure o programa para monitorar esse diretório
export WATCH_DIR=~/test_monitor
export CHECK_INTERVAL=3  # Verificar a cada 3 segundos para testes

# Execute o programa
./filemanager

# Em outro terminal, crie, modifique e exclua arquivos no diretório monitorado
touch ~/test_monitor/arquivo1.txt
echo "Conteúdo de teste" > ~/test_monitor/arquivo1.txt
rm ~/test_monitor/arquivo1.txt
```

O programa registrará todos esses eventos e gerará estatísticas sobre os arquivos.

## Configuração

O Go File Manager pode ser configurado através de variáveis de ambiente:

| Variável       | Descrição                                   | Padrão              |
| -------------- | ------------------------------------------- | ------------------- |
| WATCH_DIR      | Diretório a ser monitorado                  | Diretório atual (.) |
| LOG_DIR        | Diretório para armazenar logs               | ./logs              |
| CHECK_INTERVAL | Intervalo de verificação em segundos        | 5                   |
| MAX_LOG_SIZE   | Tamanho máximo dos arquivos de log em bytes | 10485760 (10MB)     |
| AUTO_ACTION    | Ativar ações automáticas                    | false               |
| IGNORE_EXTS    | Lista de extensões para ignorar             | .temp,.swp          |

Exemplo de configuração:

```bash
export WATCH_DIR="/dados/importantes"
export LOG_DIR="/var/log/filemanager"
export CHECK_INTERVAL="10"
export MAX_LOG_SIZE="20971520"  # 20MB
export AUTO_ACTION="true"
export IGNORE_EXTS=".tmp,.bak,.swp"

./filemanager
```

Para mais detalhes sobre configuração, consulte a [documentação do pacote config](./docs/config_explicado.md).

## Documentação Detalhada

O projeto inclui documentação detalhada para cada componente:

- [Pacote de Configuração](./docs/config_explicado.md): Como o sistema é configurado
- [Pacote de Logger](./docs/logger_explicado.md): Como são registrados os eventos
- [Pacote de Monitor](./docs/monitor_explicado.md): Como funciona o monitoramento de arquivos
- [Pacote de Estatísticas](./docs/stats_explicado.md): Como são geradas as estatísticas
- [Pacote de Ações](./docs/actions_explicado.md): Como são executadas ações nos arquivos

## Exemplos Práticos

### Monitorando Múltiplos Diretórios

Para monitorar múltiplos diretórios, você pode executar várias instâncias do programa:

```bash
# Terminal 1 - Monitorando diretório de documentos
WATCH_DIR=~/Documentos ./filemanager

# Terminal 2 - Monitorando diretório de projetos
WATCH_DIR=~/Projetos ./filemanager
```

### Configurando Ações Automáticas

O sistema pode ser estendido para executar ações automáticas quando detecta mudanças nos arquivos. Aqui está um exemplo de como você poderia implementar isso:

```go
// Exemplo de como criar e usar ações
backupAction := actions.NewBackupAction("./backups", logger)

// Executar um backup quando um arquivo for modificado
if event == "MODIFICADO" {
    err := backupAction.Execute(filePath)
    if err != nil {
        log.Printf("Erro ao fazer backup: %v", err)
    }
}
```

## Conceitos Avançados

### Rotação de Logs

O sistema implementa rotação de logs para evitar que os arquivos de log cresçam indefinidamente:

```go
// Verifica se o arquivo de log atingiu o tamanho máximo
if l.currentSize > l.config.MaxLogSize {
    return l.openLogFile()  // Abre um novo arquivo de log
}
```

### Hash de Arquivos

O sistema usa hashing MD5 para identificar alterações no conteúdo dos arquivos:

```go
hash := md5.New()
if _, err := io.Copy(hash, file); err != nil {
    return "", err
}
hashInBytes := hash.Sum(nil)
return hex.EncodeToString(hashInBytes), nil
```

### Formatação Humanizada

O sistema converte tamanhos de bytes para formatos legíveis por humanos:

```go
func formatSize(size int64) string {
    const (
        B  = 1
        KB = 1024 * B
        MB = 1024 * KB
        GB = 1024 * MB
    )

    switch {
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
```

### Manipulação Segura de Sinais

O programa captura sinais do sistema operacional para encerrar graciosamente:

```go
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

// Espera pelo sinal de interrupção
select {
case <-sigCh:
    // Código de encerramento
}
```

---

Este projeto demonstra várias práticas recomendadas de desenvolvimento em Go, incluindo concorrência segura, manipulação robusta de erros, interação eficiente com o sistema de arquivos e organização modular de código.

Para dúvidas, sugestões ou contribuições, sinta-se à vontade para abrir uma issue ou enviar um pull request.
