package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	// Consumidor para receber mensagens
	msgs, err := ch.Consume(
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

	go func() {
		for msg := range msgs {
			wordCount := countWords(string(msg.Body))
			fmt.Printf("Número de palavras na mensagem recebida: %d\n", wordCount)
		}
	}()

	// Publicador para enviar mensagens
	message := "Esta é uma mensagem fixa."
	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Fatalf("Erro ao publicar a mensagem: %v", err)
	}

	fmt.Println("Mensagem enviada com sucesso!")

	// Aguarde o sinal de interrupção para encerrar o programa
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
}

func countWords(text string) int {
	// Função para contar palavras
	words := strings.Fields(text)
	return len(words)
}
