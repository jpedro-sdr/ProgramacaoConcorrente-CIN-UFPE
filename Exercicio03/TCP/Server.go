package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var requestData WordCountRequest

	err := decoder.Decode(&requestData)
	if err != nil {
		fmt.Println("Erro ao decodificar a solicitação:", err)
		return
	}

	timeWithoutConcurrency, timeWithConcurrency := executeWordCount(requestData.Text, requestData.NumParts)
	fmt.Println(requestData.Text)

	response := WordCountResponse{
		TimeWithoutConcurrency: timeWithoutConcurrency,
		TimeWithConcurrency:    timeWithConcurrency,
	}

	encoder := json.NewEncoder(conn)
	err = encoder.Encode(response)
	if err != nil {
		fmt.Println("Erro ao enviar a resposta:", err)
		return
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao criar o servidor:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Servidor aguardando conexões na porta 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar a conexão:", err)
			continue
		}

		go handleConnection(conn)
	}
}

type WordCountResponse struct {
	TimeWithoutConcurrency time.Duration `json:"timeWithoutConcurrency"`
	TimeWithConcurrency    time.Duration `json:"timeWithConcurrency"`
}

type WordCountRequest struct {
	Text     string `json:"text"`
	NumParts int    `json:"numParts"`
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

func executeWordCount(text string, numParts int) (time.Duration, time.Duration) {
	start := time.Now()
	wordCount(text)
	timeWordCountWithoutConcurrency := time.Since(start)

	start = time.Now()
	concurrentWordCount(text, numParts)
	timeWordCountWithConcurrency := time.Since(start)

	return timeWordCountWithoutConcurrency, timeWordCountWithConcurrency
}
