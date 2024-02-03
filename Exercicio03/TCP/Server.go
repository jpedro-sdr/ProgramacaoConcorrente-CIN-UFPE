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

func multiplicarMatrizes(m1, m2 [][]int) [][]int {
	// Implementação da multiplicação de matrizes
	// ...
	return nil
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var matrices Matrices
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&matrices)
	if err != nil {
		fmt.Println("Erro ao decodificar dados:", err)
		return
	}

	// Processar dados (por exemplo, multiplicar matrizes)
	resultMatrix := multiplicarMatrizes(matrices.Matrix1, matrices.Matrix2)

	// Enviar resultado de volta ao cliente
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(resultMatrix)
	if err != nil {
		fmt.Println("Erro ao codificar resultado:", err)
		return
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:5000")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor TCP ouvindo em 127.0.0.1:5000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleClient(conn)
	}
}
