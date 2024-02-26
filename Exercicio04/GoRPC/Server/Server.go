package main

import (
	"fmt"
	"net"
	"net/rpc"
	"strings"
	"sync"
	"time"
)

type WordCountService struct{}

func (w *WordCountService) WordCount(request string, response *string) error {
	timeWordCountWithConcurrency, timeWordCountWithoutConcurrency := executeWordCount(request)
	*response = fmt.Sprintf("Com concorrência: %d  || Sem concorrência: %d", timeWordCountWithConcurrency, timeWordCountWithoutConcurrency)
	return nil
}

func main() {
	wordCountService := new(WordCountService)
	rpc.Register(wordCountService)

	listener, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor RPC:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor RPC ouvindo em", listener.Addr())

	rpc.Accept(listener)
}

func wordCount(s string) map[string]int {
	words := strings.Fields(s)
	m := make(map[string]int)
	for _, word := range words {
		m[word]++
	}
	return m
}

func concurrentWordCount(s string, numParts int) map[string]int {
	parts := make([]string, numParts)
	words := make([]map[string]int, numParts)
	var wg sync.WaitGroup
	for i := 0; i < numParts; i++ {
		start := i * len(s) / numParts
		end := (i + 1) * len(s) / numParts
		parts[i] = s[start:end]
		wg.Add(1)
		go func(i int) {
			words[i] = wordCount(parts[i])
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

func executeWordCount(text string) (time.Duration, time.Duration) {
	numParts := 2
	start := time.Now()
	wordCount(text)
	timeWordCountWithoutConcurrency := time.Since(start)
	start = time.Now()
	concurrentWordCount(text, numParts)
	timeWordCountWithConcurrency := time.Since(start)

	return timeWordCountWithoutConcurrency, timeWordCountWithConcurrency
}
