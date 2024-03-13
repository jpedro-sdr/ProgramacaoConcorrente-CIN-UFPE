package main

import (
	"fmt"
	"log"
	common "module05/Common"

	"github.com/streadway/amqp"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir o canal: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"wordCountQueue", // Nome da fila
		false,            // Durable
		false,            // Delete when unused
		false,            // Exclusive
		false,            // No-wait
		nil,              // Arguments
	)
	if err != nil {
		log.Fatalf("Erro ao declarar a fila: %v", err)
	}

	msgsFromClient, err := ch.Consume(
		q.Name, // Fila
		"",     // Consumidor
		true,   // Auto-ack
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Args
	)
	if err != nil {
		log.Fatalf("Erro ao registrar o consumidor: %v", err)
	}

	for msg := range msgsFromClient {
		timeWordCountWithoutConcurrency, timeWordCountWithConcurrency :=
			common.ExecuteWordCount(string(msg.Body))
		fmt.Printf("Tempo sem concorrência: %d\n", timeWordCountWithoutConcurrency)
		fmt.Printf("Tempo com concorrência: %d\n", timeWordCountWithConcurrency)
	}
}
