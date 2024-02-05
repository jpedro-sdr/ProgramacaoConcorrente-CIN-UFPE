package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// WordCountRequest é a estrutura para a solicitação de contagem de palavras
type WordCountRequest struct {
	Text     string `json:"text"`
	NumParts int    `json:"numParts"`
}

// WordCountResponse é a estrutura para a resposta da contagem de palavras
type WordCountResponse struct {
	TimeWithoutConcurrency time.Duration `json:"timeWithoutConcurrency"`
	TimeWithConcurrency    time.Duration `json:"timeWithConcurrency"`
}

func communicateWithServer(requestData WordCountRequest, conn net.Conn, iteration int) {
	// Medir o tempo inicial
	startTime := time.Now()

	// Enviar a mensagem para o servidor
	encoder := json.NewEncoder(conn)
	err := encoder.Encode(requestData)
	if err != nil {
		fmt.Println("Erro ao codificar e enviar a mensagem:", err)
		return
	}

	fmt.Println(requestData)

	// Receber a resposta do servidor
	var response WordCountResponse
	decoder := json.NewDecoder(conn)
	err = decoder.Decode(&response)
	if err != nil {
		fmt.Println("Erro ao decodificar a resposta:", decoder)
		fmt.Println("Erro ao decodificar a resposta:", err)
		return
	}

	// Medir o tempo final e calcular o round-trip time
	endTime := time.Now()
	roundTripTime := endTime.Sub(startTime)

	// Processar a resposta conforme necessário
	fmt.Printf("Iteração %d\n", iteration+1)
	fmt.Printf("Tempo sem concorrência: %v\n", response.TimeWithoutConcurrency)
	fmt.Printf("Tempo com concorrência: %v\n", response.TimeWithConcurrency)
	fmt.Printf("Round-trip time: %v\n", roundTripTime)
	fmt.Println("------------------------------")
}

func main() {
	// Conectar ao servidor TCP
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup

	// Executar o código 10,000 vezes
	for i := 0; i < 10000; i++ {
		// Construir a mensagem a ser enviada ao servidor
		requestData := WordCountRequest{
			Text:     "Seu texto aqui",
			NumParts: 3,
		}
		wg.Add(1)
		// Executar a comunicação com o servidor em uma goroutine
		go func(iteration int) {
			defer wg.Done()
			communicateWithServer(requestData, conn, iteration)
		}(i)

		// Aguardar um curto período entre as iterações
		time.Sleep(time.Millisecond * 1000)
	}
	wg.Wait()

	// Aguardar um tempo suficiente para permitir a conclusão das goroutines antes de encerrar o programa
	time.Sleep(time.Second * 5)
}
