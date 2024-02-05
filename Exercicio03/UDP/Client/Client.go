// cliente.go

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
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

func main() {
	defaultMessage := "Olá, este é um exemplo de mensagem."

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Erro ao resolver o endereço UDP:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Erro ao discar para o servidor UDP:", err)
		os.Exit(1)
	}
	defer conn.Close()

	request := WordCountRequest{
		Text:     defaultMessage,
		NumParts: 3,
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Erro ao codificar o pedido JSON:", err)
		os.Exit(1)
	}

	_, err = conn.Write(requestBytes)
	if err != nil {
		fmt.Println("Erro ao enviar dados para o servidor:", err)
		os.Exit(1)
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Erro ao receber dados do servidor:", err)
		os.Exit(1)
	}

	var response WordCountResponse
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		os.Exit(1)
	}

	fmt.Println("Resposta do servidor:", response)
}
