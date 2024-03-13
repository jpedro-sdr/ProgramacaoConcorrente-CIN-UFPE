package main

import (
	"fmt"
	common "module05/Common"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

func makeRequest(wg *sync.WaitGroup, roundTripTimes *[]time.Duration, totalTime *time.Duration) {
	defer wg.Done()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Erro ao conectar ao RabbitMQ:", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Erro ao abrir o canal:", err)
		return
	}
	defer ch.Close()

	queueName := "bibleQueue"
	_, err = ch.QueueDeclare(
		queueName, // Nome da fila
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)

	if err != nil {
		fmt.Println("Erro ao declarar a fila:", err)
		return
	}

	bibleText, err := common.ReadBibleText()
	if err != nil {
		fmt.Println("Erro ao ler o conteúdo do arquivo:", err)
		return
	}

	startTime := time.Now()
	for _, chunk := range bibleText {
		chunkBytes := []byte(string(chunk))

		err = ch.Publish(
			"",        // Exchange
			queueName, // Routing key
			false,     // Mandatory
			false,     // Immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        chunkBytes, // Enviando slice de bytes
			})
		if err != nil {
			fmt.Println("Erro ao enviar mensagem para o RabbitMQ:", err)
			return
		}
	}

	endTime := time.Now()
	roundTripTime := endTime.Sub(startTime)
	*roundTripTimes = append(*roundTripTimes, roundTripTime)
	*totalTime += roundTripTime
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
