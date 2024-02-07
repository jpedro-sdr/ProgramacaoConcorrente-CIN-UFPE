package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	// Leitura da requisição do cliente
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao ler a requisição do cliente:", err)
		return
	}

	requestText := string(buffer[:n])
	timeWordCountWithConcurrency, timeWordCountWithoutConcurrency := executeWordCount(requestText)

	response := fmt.Sprintf("Com concorrência: %d  || Sem concorrência: %d", timeWordCountWithConcurrency, timeWordCountWithoutConcurrency)
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Erro ao enviar resposta para o cliente:", err)
		return
	}
}

func main() {
	address := "localhost:8081"

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor TCP ouvindo em", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			// fmt.Println("Erro ao aceitar conexão do cliente:", err)
			continue
		}

		go handleConnection(conn)
	}
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
	numParts := 6
	start := time.Now()
	wordCount(text)
	timeWordCountWithoutConcurrency := time.Since(start)

	start = time.Now()
	concurrentWordCount(text, numParts)
	timeWordCountWithConcurrency := time.Since(start)

	return timeWordCountWithoutConcurrency, timeWordCountWithConcurrency
}
