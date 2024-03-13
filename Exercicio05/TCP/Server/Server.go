package main

import (
	"fmt"
	common "module05/Common"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	// Leitura da requisição do cliente
	n, err := conn.Read(buffer)
	if err != nil {
		// fmt.Println("Erro ao ler a requisição do cliente:", err)
		return
	}

	requestText := string(buffer[:n])
	timeWordCountWithConcurrency, timeWordCountWithoutConcurrency := common.ExecuteWordCount(requestText)

	response := fmt.Sprintf("Com concorrência: %d  || Sem concorrência: %d", timeWordCountWithConcurrency, timeWordCountWithoutConcurrency)
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Erro ao enviar resposta para o cliente:", err)
		return
	}
}

func main() {
	address := "localhost:8081"

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor TCP ouvindo em", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}
