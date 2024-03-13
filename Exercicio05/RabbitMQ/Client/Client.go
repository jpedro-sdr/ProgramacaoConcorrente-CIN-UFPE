package main

import (
	"fmt"
	"log"
	common "module05/Common"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

func makeRequest(wg *sync.WaitGroup, roundTripTimes *[]time.Duration, totalTime *time.Duration) {
	// Conectar-se ao servidor RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Abrir um canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir o canal: %v", err)
	}
	defer ch.Close()

	// Declarar uma fila
	q, err := ch.QueueDeclare(
		"wordcount_queue", // Nome da fila
		false,             // Durable
		false,             // Delete when unused
		false,             // Exclusive
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		log.Fatalf("Erro ao declarar a fila: %v", err)
	}

	// Ler o conteúdo do arquivo da Bíblia
	bibleText, err := common.ReadBibleText()
	if err != nil {
		log.Fatalf("Erro ao ler o conteúdo do arquivo: %v", err)
	}

	// Enviar a mensagem para a fila
	start := time.Now()
	err = ch.Publish(
		"",     // Exchange
		q.Name, // Key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(bibleText),
		})
	if err != nil {
		log.Fatalf("Erro ao publicar a mensagem: %v", err)
	}

	// Calcular o round trip time
	end := time.Now()
	rtt := end.Sub(start)
	fmt.Printf("Round Trip Time: %v\n", rtt)
}

func main() {
	var wg sync.WaitGroup

	var roundTripTimesRabbitMQ []time.Duration
	var totalTime time.Duration

	for i := 0; i < common.NumRequests; i++ {
		wg.Add(1)
		go makeRequest(&wg, &roundTripTimesRabbitMQ, &totalTime)
		time.Sleep(common.TimeSleep * time.Millisecond)
	}

	wg.Wait()

	common.SaveToFile(roundTripTimesRabbitMQ, "../../rabbitmq.txt")
	fmt.Println("RTT", totalTime)
}
