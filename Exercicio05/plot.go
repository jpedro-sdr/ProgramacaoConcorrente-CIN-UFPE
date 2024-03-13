package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func main() {
	// Create a new bar chart
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "RTT Mediana (ms)",
		}),
	)

	tcpValues, err := readValues("5ms/tcp.txt", true)
	if err != nil {
		log.Fatalf("Error reading tcp.txt: %v", err)
	}
	tcpStdDev := average(tcpValues)

	udpValues, err := readValues("5ms/udp.txt", true)
	if err != nil {
		log.Fatalf("Error reading udp.txt: %v", err)
	}
	udpStdDev := average(udpValues)

	rpcValues, err := readValues("5ms/rpc.txt", true)
	if err != nil {
		log.Fatalf("Error reading rpc.txt: %v", err)
	}
	rpcStdDev := average(rpcValues)

	rabbitMQValues, err := readValues("5ms/rabbitmq.txt", true)
	if err != nil {
		log.Fatalf("Error reading rpc.txt: %v", err)
	}
	rabbitMQStdDev := average(rabbitMQValues)

	// Add the standard deviations to the chart
	bar.SetXAxis([]string{"TCP", "UDP", "RPC", "RabbitMQ"}).
		AddSeries("TCP", generateBarItems([]float64{tcpStdDev})).
		AddSeries("UDP", generateBarItems([]float64{udpStdDev})).
		AddSeries("RPC", generateBarItems([]float64{rpcStdDev})).
		AddSeries("RabbitMQ", generateBarItems([]float64{rabbitMQStdDev}))

	// Save the chart in an HTML file
	page := components.NewPage()
	page.AddCharts(bar)

	f, err := os.Create("rtt.html")
	if err != nil {
		log.Fatalf("Error creating the file: %v", err)
	}
	defer f.Close()

	err = page.Render(f)
	if err != nil {
		log.Fatalf("Error rendering the chart: %v", err)
	}
}

func readValues(filename string, removeZeros bool) ([]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var values []float64
	var value float64

	for {
		_, err := fmt.Fscanf(file, "%f\n", &value)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if !removeZeros || value != 0 {
			values = append(values, value)
		}
	}

	return values, nil
}

func median(values []float64) float64 {
	sort.Float64s(values)
	middle := len(values) / 2
	if len(values)%2 == 0 {
		return (values[middle-1] + values[middle]) / 2
	} else {
		return values[middle]
	}
}

// Function to calculate the average of a slice of float64
func average(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

// Function to generate BarData for a given slice of float64 values
func generateBarItems(values []float64) []opts.BarData {
	items := make([]opts.BarData, 0, len(values))
	for _, value := range values {
		items = append(items, opts.BarData{Value: value})
	}
	return items
}

// Function to calculate the standard deviation of a slice of float64
func stdDev(values []float64) float64 {
	mean := average(values)
	total := 0.0
	for _, value := range values {
		total += math.Pow(value-mean, 2)
	}
	variance := total / float64(len(values)-1)
	return math.Sqrt(variance)
}
