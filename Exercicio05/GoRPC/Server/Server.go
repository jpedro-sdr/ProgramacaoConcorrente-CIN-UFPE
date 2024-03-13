package main

import (
	"fmt"
	common "module05/Common"
	"net"
	"net/rpc"
)

type WordCountService struct{}

func (w *WordCountService) WordCount(request string, response *string) error {
	timeWordCountWithConcurrency, timeWordCountWithoutConcurrency := common.ExecuteWordCount(request)
	*response = fmt.Sprintf("Com concorrência: %d  || Sem concorrência: %d", timeWordCountWithConcurrency, timeWordCountWithoutConcurrency)
	return nil
}

func main() {
	wordCountService := new(WordCountService)
	rpc.Register(wordCountService)

	listener, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor RPC:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor RPC ouvindo em", listener.Addr())

	rpc.Accept(listener)
}
