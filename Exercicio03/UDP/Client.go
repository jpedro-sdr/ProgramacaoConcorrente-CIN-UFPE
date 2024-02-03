package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

type Matrices struct {
	Matrix1 [][]int
	Matrix2 [][]int
}

func main() {
	matriz1 := [][]int{...}
	matriz2 := [][]int{...}
	matrices := Matrices{Matrix1: matriz1, Matrix2: matriz2}

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5000")
	if err != nil {
		fmt.Println("Erro ao resolver endereço UDP:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor UDP:", err)
		return
	}
	defer conn.Close()

	// Enviar dados para o servidor
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(matrices)
	if err != nil {
		fmt.Println("Erro ao enviar dados:", err)
		return
	}

	// Receber resultado do servidor
	buffer := make([]byte, 4096)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Erro ao ler resultado:", err)
		return
	}

	var resultMatrix [][]int
	decoder := gob.NewDecoder(&gobDecoderWriter{Reader: &buffer, n: n})
	err = decoder.Decode(&resultMatrix)
	if err != nil {
		fmt.Println("Erro ao decodificar resultado:", err)
		return
	}

	// Processar resultado (se necessário)
	// ...
}
