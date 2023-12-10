package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/go-echarts/go-echarts/charts"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	m := make(map[string]int)
	for _, word := range words {
		word = strings.ToLower(word) // Converter a palavra para minúsculas antes de adicionar ao mapa.
		m[word]++
	}
	return m
}

func ExpandText(text string, quantidadeRepeticao int) string {
	return strings.Repeat(text, quantidadeRepeticao)
}

func ConcurrentWordCount(s string, numParts int) map[string]int {
	parts := make([]string, numParts)
	words := make([]map[string]int, numParts)
	strLen := len(s) // Otimização o len(s) é chamado apenas uma vez.

	var wg sync.WaitGroup

	for i := 0; i < numParts; i++ {
		start := i * strLen / numParts
		end := (i + 1) * strLen / numParts
		parts[i] = s[start:end]
		wg.Add(1)
		go func(i int) {
			words[i] = WordCount(parts[i])
			wg.Done()
		}(i)
	}

	wg.Wait()
	m := make(map[string]int)
	for _, wordMap := range words {
		for word, count := range wordMap {
			m[word] += count
		}
	}

	return m
}

func PrintCPUUsage() (uint64, int) {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	memUsage := stat.Sys / 1024 / 1024 // Memória alocada em MB
	cpuUsage := int(stat.Sys * 100 / stat.TotalAlloc)

	return memUsage, cpuUsage
}

func main() {
	// Leitura do texto a partir de um arquivo (por exemplo, input.txt)
	file, err := os.Open("files text/input.txt")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Leitura do texto linha a linha
	scanner := bufio.NewScanner(file)
	var text string

	for scanner.Scan() {
		text += scanner.Text() + " "
	}

	// Inicializar o perfil de CPU
	f, err := os.Create("cpu_profile.pprof")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo de perfil de CPU:", err)
		return
	}
	defer f.Close()

	// Iniciar a coleta do perfil de CPU
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("Erro ao iniciar a coleta de perfil de CPU:", err)
		return
	}
	defer pprof.StopCPUProfile()

	var xData []string
	var yDataWithoutConcurrency []time.Duration
	var yDataWithConcurrency []time.Duration

	// Inicializar o temporizador da aplicação
	start := time.Now()
	timeWordCountWithoutConcurrency := time.Since(start)
	timeWordCountWithConcurrency := time.Since(start)

	re := regexp.MustCompile(`[[:punct:]]`) //Regex de remoção de pontuação.
	text = re.ReplaceAllString(text, "")    //Remova a pontuação do texto.
	numParts := 3                           // Que também é o número de threads a serem executados.

	// Número de execuções desejadas
	numExecutions := 10
	var totalTimeWithoutConcurrency time.Duration
	var totalTimeWithConcurrency time.Duration

	//Parametros de CPU e memória
	var totalMemoryAllocatedWithoutConcurreny uint64
	var totalMemoryAllocatedWithConcurrency uint64
	var stat runtime.MemStats // Variável de execução de memória

	for i := 0; i < numExecutions; i++ {
		// Executar o programa não concorrente.
		start = time.Now()
		WordCount(text)
		timeWordCountWithoutConcurrency = time.Since(start)
		totalTimeWithoutConcurrency += timeWordCountWithoutConcurrency

		// Leitura de uso da memória sem concorrência.
		runtime.ReadMemStats(&stat)
		totalMemoryAllocatedWithoutConcurreny += stat.Alloc

		// Executar o programa concorrente.
		start = time.Now()
		ConcurrentWordCount(text, numParts)
		timeWordCountWithConcurrency = time.Since(start)
		totalTimeWithConcurrency += timeWordCountWithConcurrency

		// Leitura de uso da memória com concorrência.
		runtime.ReadMemStats(&stat)
		totalMemoryAllocatedWithConcurrency += stat.Alloc
	}

	// Calcular o tempo médio de execução sem concorrência.
	averageWithoutConcurrency := totalTimeWithoutConcurrency / time.Duration(numExecutions)
	// Calcular o uso médio de memória.
	averageMemory := totalMemoryAllocatedWithoutConcurreny / uint64(numExecutions)
	fmt.Printf("Memória média alocada do programa sem concorrência: %v MB\n", averageMemory/1024/1024)

	// Calcular o tempo médio de execução com concorrência.
	averageWithConcurrency := totalTimeWithConcurrency / time.Duration(numExecutions)
	// Calcular o uso médio de memória.
	averageMemory = totalMemoryAllocatedWithConcurrency / uint64(numExecutions)
	fmt.Printf("Memória média alocada do programa com concorrência: %v MB\n", averageMemory/1024/1024)

	// Criar dados para o gráfico com o tempo de execução.
	yDataWithoutConcurrency = append(yDataWithoutConcurrency, averageWithoutConcurrency)
	yDataWithConcurrency = append(yDataWithConcurrency, averageWithConcurrency)

	// Configurando o gráfico
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "Tempos de execução"}, charts.ToolboxOpts{Show: true})
	bar.AddXAxis(xData).AddYAxis("Sem concorrência", yDataWithoutConcurrency).AddYAxis("Com concorrência", yDataWithConcurrency)

	// Agora, use outra variável para o arquivo HTML
	fHTML, _ := os.Create("bar.html")
	bar.Render(fHTML)
	fmt.Printf("Tamanho do texto em palavras: %v\n", len(text))
	fmt.Printf("Quantidade de repetições: %v\n", numExecutions)
}
