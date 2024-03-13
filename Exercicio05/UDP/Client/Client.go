package main

import (
	"fmt"
	common "module05/Common"
	"net"
	"strconv"
	"sync"
	"time"
)

func makeRequest(wg *sync.WaitGroup, roundTripTimes *[]time.Duration, totalTime *time.Duration) {
	defer wg.Done()

	serverAddress := "localhost"
	serverPort := 8080
	server, err := net.ResolveUDPAddr("udp", serverAddress+":"+strconv.Itoa(serverPort))
	if err != nil {
		fmt.Println("Erro ao resolver o endereço UDP:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, server)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor UDP:", err)
		return
	}
	defer conn.Close()

	// Envio da requisição
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
	buffer := make([]byte, 65536)
	_, _, err = conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Erro ao ler resposta do servidor:", err)
		return
	}

	endTime := time.Now()
	roundTripTime := endTime.Sub(startTime)
	*roundTripTimes = append(*roundTripTimes, roundTripTime)
	*totalTime += roundTripTime

	// fmt.Printf("Resposta do servidor: %s\n", buffer[:n])
}

func main() {
	var wg sync.WaitGroup

	var roundTripTimesTCP []time.Duration
	var totalTime time.Duration

	for i := 0; i < common.NumRequests; i++ {
		wg.Add(1)
		go makeRequest(&wg, &roundTripTimesTCP, &totalTime)
		time.Sleep(common.TimeSleep * time.Millisecond)
	}

	wg.Wait()

	common.SaveToFile(roundTripTimesTCP, "../udp.txt")
}
