package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"sync"
	"time"
)

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

func makeRequest(wg *sync.WaitGroup, roundTripTimes *[]time.Duration, totalTime *time.Duration) {
	defer wg.Done()

	client, err := rpc.Dial("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor RPC:", err)
		return
	}
	defer client.Close()

	bibleText, err := readBibleText("../../biblia.txt")
	if err != nil {
		fmt.Println("Erro ao ler o conteúdo do arquivo:", err)
		return
	}

	startTime := time.Now()
	var response string
	err = client.Call("WordCountService.WordCount", bibleText, &response)
	if err != nil {
		fmt.Println("Erro ao chamar o método WordCount:", err)
		return
	}

	endTime := time.Now()
	roundTripTime := endTime.Sub(startTime)
	*roundTripTimes = append(*roundTripTimes, roundTripTime)
	*totalTime += roundTripTime
}

func main() {
	var wg sync.WaitGroup
	numRequests := 10000

	var roundTripTimes []time.Duration
	var totalTime time.Duration

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeRequest(&wg, &roundTripTimes, &totalTime)
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()

	saveToFile(roundTripTimes, "../rpc.txt")
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
