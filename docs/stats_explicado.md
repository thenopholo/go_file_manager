# Explicação Detalhada: Pacote de Estatísticas (stats.go)

Este documento explica detalhadamente o funcionamento do arquivo `stats.go`, que é responsável por gerar e armazenar estatísticas sobre os arquivos monitorados na aplicação.

## Estrutura do Código

### Importações

```go
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/thenopholo/go_file_manager/config"
	"github.com/thenopholo/go_file_manager/monitor"
)
```

- `encoding/json`: Como um tradutor, converte estruturas de dados Go para o formato JSON e vice-versa.
- `fmt`: Uma biblioteca de formatação, como um designer de texto para criar mensagens bem estruturadas.
- `os`: O principal pacote para interagir com o sistema operacional - nossa ponte direta para o sistema de arquivos.
- `path/filepath`: Um assistente especializado em manipular caminhos de arquivos, garantindo que estejam no formato correto para o sistema operacional atual.
- `time`: Como um relógio e calendário digital, permite registrar quando as estatísticas foram geradas.
- `config`: Nosso pacote local que contém as configurações da aplicação, como um manual de instruções personalizado.
- `monitor`: Pacote que contém as informações sobre os arquivos sendo monitorados, como um vigilante que nos reporta o que está acontecendo.

### Estrutura Stats

```go
type Stats struct {
	Timestamp      time.Time           `json:"timestamp"`
	FileCount      int                 `json:"file_count"`
	TotalSize      int64               `json:"total_size"`
	TotalSizeHuman string              `json:"total_size_human"`
	ByExtentions   map[string]ExtStats `json:"by_extention"`
}
```

Esta estrutura é como um relatório completo sobre o estado atual dos arquivos:

- `Timestamp`: Registra o momento exato em que as estatísticas foram coletadas (como o carimbo de data em um documento oficial).
- `FileCount`: Contador simples de quantos arquivos estão sendo monitorados (como o número total de itens em um inventário).
- `TotalSize`: O tamanho total de todos os arquivos em bytes (como o peso total de todos os itens em um armazém).
- `TotalSizeHuman`: O mesmo tamanho total, mas formatado de maneira legível para humanos (como dizer "2 quilos" em vez de "2000 gramas").
- `ByExtentions`: Um mapa que organiza estatísticas por extensão de arquivo (como separar itens de um armazém por categorias).

As tags `json:"..."` são como etiquetas de identificação que determinam como cada campo será nomeado quando convertido para JSON.

### Estrutura ExtStats

```go
type ExtStats struct {
	Count int   `json:"count"`
	Size  int64 `json:"size_bytes"`
}
```

Esta estrutura é como um relatório resumido para cada tipo de arquivo:

- `Count`: Quantos arquivos deste tipo específico existem (como contar quantos produtos de cada categoria).
- `Size`: O espaço total ocupado por arquivos deste tipo (como o espaço que cada categoria de produto ocupa na prateleira).

### Função GenerateStats

```go
func GenerateStats(monitor *monitor.FileMonitor, cfg config.Config) (*Stats, error) {
	stats := &Stats{
		Timestamp:    time.Now(),
		FileCount:    monitor.GetFileCount(),
		TotalSize:    monitor.GetTotalSize(),
		ByExtentions: make(map[string]ExtStats),
	}

	stats.TotalSizeHuman = formatSize(stats.TotalSize)

	return stats, nil
}
```

Esta função é como um contador de inventário que faz um relatório rápido do estado atual:

1. Recebe como parâmetros o monitor (que tem as informações sobre os arquivos) e as configurações.
2. Cria uma nova estrutura de estatísticas (`stats`) com:
   - O momento atual (`time.Now()`) - como colocar a data no cabeçalho do relatório
   - O número de arquivos obtido do monitor (`monitor.GetFileCount()`) - como contar todos os itens no estoque
   - O tamanho total também obtido do monitor (`monitor.GetTotalSize()`) - como somar o peso de todos os produtos
   - Um mapa vazio para estatísticas por extensão - como preparar categorias vazias para classificar os produtos
3. Converte o tamanho total para um formato legível (`formatSize(stats.TotalSize)`) - como transformar "1048576 bytes" em "1 MB"
4. Retorna as estatísticas prontas

Note que esta função ainda não está preenchendo as estatísticas por extensão (o mapa `ByExtentions` permanece vazio). Em uma implementação completa, esta função percorreria todos os arquivos e agruparia as estatísticas por extensão.

### Função SaveStatsToFile

```go
func SaveStatsToFile(stats *Stats, cfg config.Config) error {
	statsDir := filepath.Join(cfg.LogDir, "stats")

	if _, err := os.Stat(statsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(statsDir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório de estatística: %w", err)
		}
	}

	fileName := fmt.Sprintf("stats_%s.json", time.Now().Format("2006-01-02_15-04-05"))
	filePath := filepath.Join(statsDir, fileName)

	jsonData, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar estatísticas: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo de estatísticas: %w", err)
	}

	return nil
}
```

Esta função é como um arquivista que guarda relatórios em arquivos organizados:

1. Primeiro, determina onde guardar o relatório, criando um caminho para um subdiretório "stats" dentro do diretório de logs:

   ```go
   statsDir := filepath.Join(cfg.LogDir, "stats")
   ```

   - Isto é como decidir guardar todos os relatórios em uma pasta específica dentro do arquivo geral.

2. Verifica se este diretório existe, e se não, o cria:

   ```go
   if _, err := os.Stat(statsDir); os.IsNotExist(err) {
       if err := os.MkdirAll(statsDir, 0755); err != nil {
           return fmt.Errorf("erro ao criar diretório de estatística: %w", err)
       }
   }
   ```

   - `os.Stat(statsDir)` tenta obter informações sobre o diretório, como verificar se uma gaveta existe.
   - `os.IsNotExist(err)` verifica se o erro retornado indica que o diretório não existe.
   - `os.MkdirAll(statsDir, 0755)` cria o diretório (e quaisquer diretórios pais necessários):
     - `0755` são as permissões: o proprietário pode ler/escrever/executar (7), enquanto grupo e outros podem ler/executar (5).
   - Se houver um erro na criação, retorna uma mensagem descritiva usando `fmt.Errorf`.

3. Cria um nome para o arquivo de estatísticas baseado na data e hora atual:

   ```go
   fileName := fmt.Sprintf("stats_%s.json", time.Now().Format("2006-01-02_15-04-05"))
   filePath := filepath.Join(statsDir, fileName)
   ```

   - O formato `2006-01-02_15-04-05` é um padrão único do Go para representar "YYYY-MM-DD_HH-MM-SS".
   - Isso resultará em nomes como "stats_2023-10-15_14-30-22.json".

4. Converte a estrutura de estatísticas para JSON formatado:

   ```go
   jsonData, err := json.MarshalIndent(stats, "", "  ")
   if err != nil {
       return fmt.Errorf("erro ao serializar estatísticas: %w", err)
   }
   ```

   - `json.MarshalIndent` converte a estrutura para texto JSON com indentação para facilitar a leitura.
   - O primeiro argumento (`stats`) é a estrutura a ser convertida.
   - O segundo (`""`) é um prefixo para cada linha (neste caso, vazio).
   - O terceiro (`"  "`) é a indentação - dois espaços neste caso.
   - É como transformar uma lista organizada mentalmente em um documento formatado.

5. Finalmente, escreve os dados JSON em um arquivo:
   ```go
   if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
       return fmt.Errorf("erro ao escrever arquivo de estatísticas: %w", err)
   }
   ```
   - `os.WriteFile` escreve dados em um arquivo com um único comando.
   - `filePath` é o caminho onde o arquivo será salvo.
   - `jsonData` são os dados a serem escritos.
   - `0644` são as permissões do arquivo: o proprietário pode ler/escrever (6), enquanto grupo e outros podem apenas ler (4).
   - Se houver um erro ao escrever, uma mensagem descritiva é retornada.

### Função formatSize

```go
func formatSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
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

Esta função é como um tradutor especializado de medidas técnicas para linguagem comum:

1. Primeiro, define constantes para cada unidade de medida:

   ```go
   const (
       B  = 1
       KB = 1024 * B
       MB = 1024 * KB
       GB = 1024 * MB
       TB = 1024 * GB
   )
   ```

   - Define uma hierarquia de unidades, onde cada uma é 1024 vezes maior que a anterior.
   - É como definir a relação entre milímetros, centímetros, metros e quilômetros.

2. Em seguida, usa uma estrutura `switch` sem expressão (que avalia cada condição em ordem):
   ```go
   switch {
   case size >= TB:
       return fmt.Sprintf("%.2f TB", float64(size)/TB)
   case size >= GB:
       return fmt.Sprintf("%.2f GB", float64(size)/GB)
   case size >= MB:
       return fmt.Sprintf("%.2f MB", float64(size)/MB)
   case size >= KB:
       return fmt.Sprintf("%.2f KB", float64(size)/KB)
   default:
       return fmt.Sprintf("%d B", size)
   }
   ```
   - Começa verificando a maior unidade (TB) e vai descendo.
   - Quando encontra uma unidade apropriada, divide o tamanho por essa unidade.
   - Formata o resultado com duas casas decimais, exceto para bytes que são inteiros.
   - É como escolher a unidade mais apropriada para expressar uma distância - não faz sentido dizer "1000000 mm" quando podemos dizer "1 km".

## Analogia Geral do Sistema de Estatísticas

O pacote de estatísticas funciona como um departamento de análise de dados em uma empresa:

1. Recebe informações brutas do departamento de monitoramento (módulo monitor).
2. Processa e organiza esses dados em categorias significativas.
3. Formata os números grandes em unidades mais compreensíveis para humanos.
4. Salva relatórios periódicos em arquivos JSON, como relatórios corporativos arquivados com data e hora.
5. Garante que todos os relatórios sejam armazenados em um local apropriado, criando esse local se necessário.

O uso do pacote `os` é fundamental aqui para:

- Verificar se diretórios existem (`os.Stat`)
- Criar diretórios quando necessário (`os.MkdirAll`)
- Escrever dados em arquivos (`os.WriteFile`)

## Como Usar este Módulo

```go
// Inicializar configuração e monitor
cfg := config.LoadConfig()
logger, _ := logger.NewLogger(cfg)
fileMonitor := monitor.NewFileMonitor(cfg, logger)
fileMonitor.Start()

// Gerar estatísticas periodicamente
stats, err := stats.GenerateStats(fileMonitor, cfg)
if err != nil {
    log.Fatalf("Erro ao gerar estatísticas: %v", err)
}

// Salvar em arquivo
if err := stats.SaveStatsToFile(stats, cfg); err != nil {
    log.Fatalf("Erro ao salvar estatísticas: %v", err)
}
```

Este exemplo mostra como o módulo de estatísticas se integra com o resto da aplicação para criar relatórios periódicos sobre os arquivos monitorados.
