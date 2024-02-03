//Exercício 01
//  Implementar um programa capaz de contar as ocorrências de cada palavra de um texto.
//  Utilize como exemplo de texto grandes textos, e.g., a bíblia. É preciso implementar
//  duas versões do algoritmo: uma versão sem concorrência e uma versão concorrente.
//  Em seguida, é preciso fazer uma avaliação comparativa de desempenho entre as versões.

package main

import (
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
	var wg sync.WaitGroup
	for i := 0; i < numParts; i++ {
		start := i * len(s) / numParts
		end := (i + 1) * len(s) / numParts
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

func PrintCPUUsage() {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	fmt.Printf("Memória Alocada: %v MB\n", stat.Alloc/1024/1024)
	fmt.Printf("Uso da CPU: %d%%\n", stat.Sys*100/stat.TotalAlloc)
}

func main() {

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
	textSize := 1
	start := time.Now()
	timeWordCountWithoutConcurrency := time.Since(start)
	timeWordCountWithConcurrency := time.Since(start)

	for textSize <= 30 {
		xData = append(xData, fmt.Sprintf("%v", textSize))
		text := strings.Repeat("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum sit amet tellus ut nibh eleifend ornare posuere ut sapien. Sed dignissim sollicitudin velit sit amet tristique. Etiam condimentum elit id lectus dapibus lobortis. Vivamus pellentesque feugiat ultricies. Sed id nunc tempus, egestas dui eu, efficitur ipsum. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Etiam lacus lectus, mollis id mattis posuere, facilisis vitae lacus. Interdum et malesuada fames ac ante ipsum primis in faucibus. Curabitur ut aliquet leo, eget aliquet purus. Nulla facilisi. Morbi ante libero, finibus sit amet interdum id, vehicula in urna. Vivamus eu nisl vitae metus aliquet viverra.	Donec varius turpis ac risus sollicitudin convallis. In urna ligula, mattis quis placerat ac, tristique a purus. Fusce a nunc ut neque mattis lacinia at sit amet purus. Nunc mi orci, interdum a metus at, laoreet finibus nulla. Sed neque enim, mattis et rutrum ac, convallis eget sem. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Aenean purus diam, pulvinar in nulla ut, suscipit vehicula justo. Nulla nisi ipsum, ornare fringilla urna ac, egestas laoreet est. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Nam vitae accumsan magna, sed sollicitudin leo. Cras in egestas nibh, vel ornare mauris. Nullam blandit, justo nec mattis volutpat, purus augue sollicitudin tortor, tristique maximus nibh enim non neque. Etiam dignissim sem non neque posuere, non interdum turpis mattis. Fusce porttitor orci nec tincidunt convallis. Donec vitae mauris at justo maximus viverra vitae id augue. Sed mollis ex sed dui placerat hendrerit et sed ante. Vivamus consequat interdum sodales. Pellentesque libero ligula, ultricies eget est a, volutpat rutrum mi. Pellentesque tincidunt mi vitae faucibus tincidunt. Cras fermentum orci sed enim efficitur cursus. Phasellus eu lorem leo. Sed convallis sagittis finibus. Nunc eget eros sapien. Donec ac massa iaculis, convallis ligula eu, convallis magna. Etiam quam est, semper ut venenatis eget, rhoncus non tellus. Mauris vitae quam vel magna blandit egestas eget ac est. Fusce laoreet vitae risus id auctor. Curabitur dapibus bibendum euismod. Sed mollis metus eu mi vulputate ornare. Donec eu posuere dolor. Vestibulum orci ex, condimentum ac eros nec, bibendum pharetra massa. Nullam efficitur nulla nisl, sit amet condimentum diam facilisis vitae. Duis pharetra, felis in lobortis egestas, est quam consectetur ante, id imperdiet nunc ligula eget magna. Pellentesque in lacus id tortor pretium sollicitudin a vitae felis. Cras euismod mauris ut nisi viverra, at dignissim augue volutpat. Aenean scelerisque aliquam purus vel consectetur. Donec vel aliquam nunc, eget pellentesque lacus. Proin pellentesque euismod lobortis. Cras eros elit, malesuada eget nibh non, lobortis viverra mauris. Nulla in augue ut lorem viverra bibendum suscipit vitae magna. Nullam tempus nulla hendrerit dui sodales, pellentesque aliquet libero hendrerit. In nec nibh ac nulla luctus euismod. Nam bibendum ipsum quis dictum dignissim. Morbi commodo consequat orci. Etiam efficitur eros ac laoreet commodo. Nullam id purus at nisi ullamcorper placerat.", textSize) // insira aqui o texto que você deseja contar as palavras
		re := regexp.MustCompile(`[[:punct:]]`)
		text = re.ReplaceAllString(text, "")
		numParts := 3 // insira aqui o número de partes em que você deseja dividir o texto
		start = time.Now()
		WordCount(text)
		timeWordCountWithoutConcurrency = time.Since(start)
		yDataWithoutConcurrency = append(yDataWithoutConcurrency, timeWordCountWithoutConcurrency)
		fmt.Printf("Tempo de execução sem concorrência: %v\n", timeWordCountWithoutConcurrency)
		PrintCPUUsage()

		start = time.Now()
		ConcurrentWordCount(text, numParts)
		timeWordCountWithConcurrency = time.Since(start)
		yDataWithConcurrency = append(yDataWithConcurrency, timeWordCountWithConcurrency)
		fmt.Printf("Tempo de execução com concorrência: %v\n", timeWordCountWithConcurrency)
		PrintCPUUsage()

		textSize = textSize + 1
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "Tempos de execução"}, charts.ToolboxOpts{Show: true})
	bar.AddXAxis(xData).AddYAxis("Sem concorrência", yDataWithoutConcurrency).AddYAxis("Com concorrência", yDataWithConcurrency)

	// Agora, use outra variável para o arquivo HTML
	fHTML, _ := os.Create("bar.html")
	bar.Render(fHTML)
	fmt.Printf("Tamanho do texto: %v\n", textSize)
}
