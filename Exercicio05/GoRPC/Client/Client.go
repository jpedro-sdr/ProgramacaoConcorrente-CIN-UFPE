package main

import (
	"fmt"
	common "module05/Common"
	"net/rpc"
	"sync"
	"time"
)

func makeRequest(wg *sync.WaitGroup, roundTripTimes *[]time.Duration, totalTime *time.Duration) {
	defer wg.Done()

	client, err := rpc.Dial("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor RPC:", err)
		return
	}
	defer client.Close()

	bibleText, err := common.ReadBibleText()
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

	var roundTripTimes []time.Duration
	var totalTime time.Duration

	for i := 0; i < common.NumRequests; i++ {
		wg.Add(1)
		go makeRequest(&wg, &roundTripTimes, &totalTime)
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()

	common.SaveToFile(roundTripTimes, "../rpc.txt")
	fmt.Println("RTT", totalTime)
}
