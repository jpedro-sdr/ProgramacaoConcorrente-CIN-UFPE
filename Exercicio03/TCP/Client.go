package main

import (
	"encoding/json"
	"fmt"
	"net"
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

func main() {
	// Conectar ao servidor TCP
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	// Construir a mensagem a ser enviada ao servidor
	requestData := WordCountRequest{
		Text:     "Seu texto aqui",
		NumParts: 3,
	}

	// Enviar a mensagem para o servidor
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(requestData)
	if err != nil {
		fmt.Println("Erro ao codificar e enviar a mensagem:", err)
		return
	}

	// Receber a resposta do servidor
	var response WordCountResponse
	decoder := json.NewDecoder(conn)
	err = decoder.Decode(&response)
	if err != nil {
		fmt.Println("Erro ao decodificar a resposta:", err)
		return
	}

	// Processar a resposta conforme necessário
	fmt.Printf("Tempo sem concorrência: %v\n", response.TimeWithoutConcurrency)
	fmt.Printf("Tempo com concorrência: %v\n", response.TimeWithConcurrency)
}
