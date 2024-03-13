package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func main() {
	// Create a new bar chart
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "RTT Desvio padrão (ms)",
		}),
	)

	// Read the values from the files and calculate the standard deviation, replacing zeros with the mean
	tcpValues, err := readValues("tcp.txt", true)
	if err != nil {
		log.Fatalf("Error reading tcp.txt: %v", err)
	}
	tcpStdDev := stdDev(tcpValues)

	udpValues, err := readValues("udp.txt", true)
	if err != nil {
		log.Fatalf("Error reading udp.txt: %v", err)
	}
	udpStdDev := stdDev(udpValues)

	rpcValues, err := readValues("rpc.txt", true)
	if err != nil {
		log.Fatalf("Error reading rpc.txt: %v", err)
	}
	rpcStdDev := stdDev(rpcValues)

	rabbitMQValues, err := readValues("rpc.txt", true)
	if err != nil {
		log.Fatalf("Error reading rpc.txt: %v", err)
	}
	rabbitMQStdDev := stdDev(rabbitMQValues)

	// Add the standard deviations to the chart
	bar.SetXAxis([]string{"TCP", "UDP", "RPC"}).
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

func readValues(filename string, replaceZeros bool) ([]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var values []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			return nil, err
		}
		if replaceZeros && value == 0 {
			// Substituir zero pela média de todos os valores
			value = calculateMean(values)
		}
		values = append(values, value)
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return values, nil
}

func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
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
