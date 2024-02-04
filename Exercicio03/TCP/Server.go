package main

import (
	"Exercicio03/Functions"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type WordCountResponse struct {
	TimeWithoutConcurrency time.Duration `json:"timeWithoutConcurrency"`
	TimeWithConcurrency    time.Duration `json:"timeWithConcurrency"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var requestData struct {
		Text     string `json:"text"`
		NumParts int    `json:"numParts"`
	}
	err := decoder.Decode(&requestData)
	if err != nil {
		fmt.Println("Erro ao decodificar a solicitação:", err)
		return
	}

	timeWithoutConcurrency, timeWithConcurrency := Functions.ExecuteWordCount(requestData.Text, requestData.NumParts)

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
