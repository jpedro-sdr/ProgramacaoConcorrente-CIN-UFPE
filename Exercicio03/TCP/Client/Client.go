package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func makeRequest(wg *sync.WaitGroup, roundTripTimes *[]time.Duration, totalTime *time.Duration) {
	defer wg.Done()

	server := "localhost:8081"
	conn, err := net.DialTimeout("tcp", server, 2*time.Second)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	bibleText, err := readBibleText("../../biblia.txt")
	if err != nil {
		fmt.Println("Erro ao ler o conteúdo do arquivo:", err)
		return
	}

	chunkSize := 1024
	startTime := time.Now()

	for i := 0; i < len(bibleText); i += chunkSize {
		end := i + chunkSize
		if end > len(bibleText) {
			end = len(bibleText)
		}

		chunk := bibleText[i:end]
		_, err = conn.Write([]byte(chunk))
		if err != nil {
			fmt.Println("Erro ao enviar requisição para o servidor:", err)
			return
		}
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao ler resposta do servidor:", err)
		return
	}

	fmt.Printf("%s\n", buffer[:n])

	endTime := time.Now()
	roundTripTime := endTime.Sub(startTime)
	*roundTripTimes = append(*roundTripTimes, roundTripTime)
	*totalTime += roundTripTime
}

func readBibleText(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content, nil
}

func main() {
	var wg sync.WaitGroup
	numRequests := 1000

	var roundTripTimesTCP []time.Duration
	var totalTime time.Duration

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeRequest(&wg, &roundTripTimesTCP, &totalTime)
		time.Sleep(5 * time.Millisecond)
	}

	wg.Wait()

	saveToFile(roundTripTimesTCP, "../tcp.txt")
	fmt.Println("RTT", totalTime)

}
func saveToFile(roundTripTimes []time.Duration, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Erro ao criar arquivo:", err)
		return
	}
	defer file.Close()

	for _, rt := range roundTripTimes {
		_, err := file.WriteString(fmt.Sprintf("%f\n", rt.Seconds()*1000)) // converter para milissegundos
		if err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return
		}
	}
}
