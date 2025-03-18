# Explicação Detalhada: Pacote de Configuração (config.go)

Este documento explica detalhadamente o funcionamento do arquivo `config.go`, que é responsável por gerenciar as configurações da aplicação de monitoramento de arquivos.

## Estrutura do Código

### Importações

```go
import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)
```

- `os`: Pacote que fornece funções para interagir com o sistema operacional, como acessar variáveis de ambiente e verificar arquivos.
- `path/filepath`: Permite manipular caminhos de arquivos de forma compatível com o sistema operacional.
- `strconv`: Usado para converter strings em outros tipos (números).
- `time`: Permite trabalhar com durações e intervalos de tempo.

### Estrutura Config

```go
type Config struct {
	WatchDir      string
	LogDir        string
	CheckInterval time.Duration
	MaxLogSize    int64
	AutoAction    bool
	IgnoreExts    []string
}
```

Esta estrutura é como uma "planta" que define todas as configurações necessárias para nossa aplicação:

- `WatchDir`: O diretório que vamos monitorar (como um vigia observando uma área específica).
- `LogDir`: Onde vamos guardar os registros de eventos (como um diário de bordo).
- `CheckInterval`: De quanto em quanto tempo verificamos por mudanças (como a frequência de rondas de um segurança).
- `MaxLogSize`: Tamanho máximo que permitimos para o arquivo de log (como um limite para não deixar o diário ficar grande demais).
- `AutoAction`: Indica se devemos agir automaticamente ao detectar mudanças (como um robô que toma decisões sozinho).
- `IgnoreExts`: Lista de extensões de arquivo que devemos ignorar (como uma lista de "não perturbe").

### Função LoadConfig

A função `LoadConfig` é como um chef de cozinha preparando um prato:

1. Primeiro, prepara os ingredientes básicos (valores padrão):

```go
config := Config {
    WatchDir:      ".",
    LogDir:        "./logs",
    CheckInterval: 5 * time.Second,
    MaxLogSize:    1024 * 1024 * 10, // 10MB
    AutoAction:    false,
    IgnoreExts:    []string{".temp", ".swp"},
}
```

- Começa vigiando o diretório atual (`.`)
- Armazena logs em uma pasta `./logs`
- Verifica a cada 5 segundos
- Limita os logs a 10MB (calculado como 1024 bytes × 1024 = 1MB, × 10 = 10MB)
- Não executa ações automáticas por padrão
- Ignora arquivos temporários (.temp, .swp)

2. Em seguida, verifica se o "cliente" pediu alterações (variáveis de ambiente):

```go
if dir := os.Getenv("WATCH_DIR"); dir != "" {
  config.WatchDir = dir
}
```

Aqui, a função `os.Getenv` é como perguntar ao sistema "Existe alguma instrução especial para o diretório a ser vigiado?". Se existir (não for vazio), usamos esse valor personalizado.

Esse padrão se repete para todas as configurações, cada uma com suas peculiaridades:

- Para o `CheckInterval`, convertemos a string da variável de ambiente para um número inteiro (segundos) e depois para uma duração:

```go
if interval := os.Getenv("CHECK_INTERVAL"); interval != "" {
  if seconds, err := strconv.Atoi(interval); err == nil {
    config.CheckInterval = time.Duration(seconds) * time.Second
  }
}
```

- Para o `MaxLogSize`, convertemos a string para um número inteiro de 64 bits:

```go
if size := os.Getenv("MAX_LOG_SIZE"); size != "" {
  if bytes, err := strconv.ParseInt(size, 10, 64); err == nil {
    config.MaxLogSize = bytes
  }
}
```

- Para `AutoAction`, verificamos se o valor da variável é exatamente "true":

```go
if autoAction := os.Getenv("AUTO_ACTION"); autoAction == "true" {
  config.AutoAction = true
}
```

- Para `IgnoreExts`, usamos a função `filepath.SplitList` que divide a string em uma lista conforme o separador padrão do sistema operacional:

```go
if ignoreExist := os.Getenv("IGNORE_EXTS"); ignoreExist != "" {
  config.IgnoreExts = filepath.SplitList(ignoreExist)
}
```

3. Por fim, garantimos que os diretórios necessários existam:

```go
ensureDir(config.WatchDir)
ensureDir(config.LogDir)
```

### Função ensureDir

```go
func ensureDir(dir string) {
  if _, err := os.Stat(dir); os.IsNotExist(err) {
    os.MkdirAll(dir, 0755)
  }
}
```

Esta função é como um zelador que verifica se uma sala existe e, se não existir, a constrói:

1. `os.Stat(dir)` verifica se o diretório existe (como tentar abrir a porta da sala)
2. `os.IsNotExist(err)` confirma se o erro é porque o diretório não existe (como confirmar que a sala realmente não está lá)
3. `os.MkdirAll(dir, 0755)` cria o diretório e todos os diretórios pais necessários (como construir não só a sala, mas todo o caminho até ela)
   - `0755` são as permissões do diretório (7=rwx para o dono, 5=r-x para grupo e outros)

## Analogia Geral

O pacote de configuração funciona como o painel de controle de um sistema de vigilância:

- Você pode usar as configurações padrão (apertar os botões pré-configurados)
- Ou pode personalizar cada configuração através de variáveis de ambiente (como ajustar os botões do painel)
- O sistema garante que todas as áreas necessárias (diretórios) estejam prontas para uso

Este design é flexível e permite que você ajuste o comportamento da aplicação sem alterar o código, ideal para diferentes ambientes (desenvolvimento, teste, produção).

````

## Como Usar

Para personalizar as configurações, defina variáveis de ambiente antes de iniciar a aplicação:

```bash
export WATCH_DIR="/meus/arquivos/importantes"
export CHECK_INTERVAL="10"  # 10 segundos
export AUTO_ACTION="true"
````

Quando a aplicação iniciar, ela usará essas configurações personalizadas em vez dos valores padrão.
