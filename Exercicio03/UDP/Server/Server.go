package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func handleConnection(conn *net.UDPConn) {
	buffer := make([]byte, 65536)
	// Leitura da requisição do cliente
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Erro ao ler a requisição do cliente:", err)
		return
	}

	// Processamento da requisição
	requestText := string(buffer[:n])
	timeWordCountWithConcurrency, timeWordCountWithoutConcurrency := executeWordCount(requestText)

	timeWordCountWithConcurrencySeconds := float64(timeWordCountWithConcurrency) / float64(time.Millisecond)
	timeWordCountWithoutConcurrencySeconds := float64(timeWordCountWithoutConcurrency) / float64(time.Millisecond)

	response := fmt.Sprintf("Com concorrência: %.2f ms Sem concorrência %.2f ms", timeWordCountWithConcurrencySeconds, timeWordCountWithoutConcurrencySeconds)

	_, err = conn.WriteToUDP([]byte(response), addr)
	if err != nil {
		fmt.Println("Erro ao enviar resposta para o cliente:", err)
		return
	}
}

func main() {
	address, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao resolver o endereço UDP:", err)
		return
	}

	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor UDP:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Servidor UDP ouvindo em", address)

	for {
		handleConnection(conn)
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
