package main

import (
	"fmt"
	"net"
)

func main() {
	// Conecta ao servidor TCP na porta 8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	// Envia uma mensagem para o servidor
	message := "Olá, servidor!"
	conn.Write([]byte(message))

	// Aguarda a resposta do servidor
	buffer := make([]byte, 1024)
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Erro ao ler a resposta do servidor:", err)
		return
	}

	// Exibe a resposta recebida do servidor
	response := string(buffer[:bytesRead])
	fmt.Printf("Resposta do servidor: %s\n", response)
}
