// servidor.go

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type WordCountResponse struct {
	TimeWithoutConcurrency time.Duration `json:"timeWithoutConcurrency"`
	TimeWithConcurrency    time.Duration `json:"timeWithConcurrency"`
}

type WordCountRequest struct {
	Text     string `json:"text"`
	NumParts int    `json:"numParts"`
}

func handleClient(conn *net.UDPConn) {
	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			return
		}

		var request WordCountRequest
		err = json.Unmarshal(buffer[:n], &request)
		if err != nil {
			fmt.Println("Error decoding JSON request:", err)
			return
		}

		timeWithoutConcurrency, timeWithConcurrency := executeWordCount(request.Text, request.NumParts)

		response := WordCountResponse{
			TimeWithoutConcurrency: timeWithoutConcurrency,
			TimeWithConcurrency:    timeWithConcurrency,
		}

		// Enviar a resposta JSON de volta para o cliente
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error encoding JSON response:", err)
			return
		}

		_, err = conn.WriteToUDP(responseBytes, addr)
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
			return
		}
	}
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening on UDP:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Servidor UDP aguardando conexões na porta 8080...")

	handleClient(conn)
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
