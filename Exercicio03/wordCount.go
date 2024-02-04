package wordcount

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	m := make(map[string]int)
	for _, word := range words {
		m[word]++
	}
	return m
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

func ExecuteWordCount(text string, numParts int) (time.Duration, time.Duration) {
	start := time.Now()
	WordCount(text)
	timeWordCountWithoutConcurrency := time.Since(start)

	start = time.Now()
	ConcurrentWordCount(text, numParts)
	timeWordCountWithConcurrency := time.Since(start)

	return timeWordCountWithoutConcurrency, timeWordCountWithConcurrency
}

func PrintCPUUsage() {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)
	fmt.Printf("Memória Alocada: %v MB\n", stat.Alloc/1024/1024)
	fmt.Printf("Uso da CPU: %d%%\n", stat.Sys*100/stat.TotalAlloc)
}
