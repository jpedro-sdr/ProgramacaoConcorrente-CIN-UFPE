package main

import (
	"fmt"
	common "module05/Common"
	"net"
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

	bibleText, err := common.ReadBibleText()
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
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao ler resposta do servidor:", err)
		return
	}

	// fmt.Printf("%s\n", buffer[:n])

	endTime := time.Now()
	roundTripTime := endTime.Sub(startTime)
	*roundTripTimes = append(*roundTripTimes, roundTripTime)
	*totalTime += roundTripTime
}

func main() {
	var wg sync.WaitGroup

	var roundTripTimesTCP []time.Duration
	var totalTime time.Duration

	for i := 0; i < common.NumRequests; i++ {
		wg.Add(1)
		go makeRequest(&wg, &roundTripTimesTCP, &totalTime)
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()

	common.SaveToFile(roundTripTimesTCP, "../tcp.txt")
	fmt.Println("RTT", totalTime)

}
