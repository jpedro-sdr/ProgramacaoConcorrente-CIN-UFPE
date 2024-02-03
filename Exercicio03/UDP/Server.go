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

func handleClient(conn *net.UDPConn) {
	buffer := make([]byte, 4096)

	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Erro ao ler dados:", err)
		return
	}

	var matrices Matrices
	decoder := gob.NewDecoder(&gobDecoderWriter{Reader: &buffer, n: n})
	err = decoder.Decode(&matrices)
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

type gobDecoderWriter struct {
	Reader *[]byte
	n      int
}

func (gdw *gobDecoderWriter) Read(p []byte) (n int, err error) {
	copy(p, (*gdw.Reader)[:gdw.n])
	return gdw.n, nil
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5000")
	if err != nil {
		fmt.Println("Erro ao resolver endereço UDP:", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor UDP:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Servidor UDP ouvindo em 127.0.0.1:5000")

	for {
		handleClient(conn)
	}
}
