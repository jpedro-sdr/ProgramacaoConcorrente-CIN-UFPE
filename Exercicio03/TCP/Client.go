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

	conn, err := net.Dial("tcp", "127.0.0.1:5000")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
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
	var resultMatrix [][]int
	decoder := gob.NewDecoder(conn)
	err = decoder.Decode(&resultMatrix)
	if err != nil {
		fmt.Println("Erro ao receber resultado:", err)
		return
	}

	// Processar resultado (se necessário)
	// ...
}
