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
	// Cria um novo gráfico de barras
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "RTT Desvio padrão (ms)",
		}),
	)

	// Lê os valores dos arquivos e calcula a média
	tcpValues, err := readValues("tcp.txt")
	if err != nil {
		log.Fatalf("Erro ao ler tcp.txt: %v", err)
	}
	tcpAvg := stdDev(tcpValues)

	udpValues, err := readValues("udp.txt")
	if err != nil {
		log.Fatalf("Erro ao ler udp.txt: %v", err)
	}
	udpAvg := stdDev(udpValues)

	// Adiciona as médias ao gráfico
	bar.SetXAxis([]string{"TCP", "UDP"}).
		AddSeries("TCP", generateBarItems([]float64{tcpAvg})).
		AddSeries("UDP", generateBarItems([]float64{udpAvg}))

	// Salva o gráfico em um arquivo HTML
	page := components.NewPage()
	page.AddCharts(bar)

	f, err := os.Create("rtt.html")
	if err != nil {
		log.Fatalf("Erro ao criar o arquivo: %v", err)
	}
	defer f.Close()

	err = page.Render(f)
	if err != nil {
		log.Fatalf("Erro ao renderizar o gráfico: %v", err)
	}
}

func readValues(filename string) ([]float64, error) {
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
		values = append(values, value)
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return values, nil
}

func average(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func generateBarItems(values []float64) []opts.BarData {
	items := make([]opts.BarData, 0, len(values))
	for _, value := range values {
		items = append(items, opts.BarData{Value: value})
	}
	return items
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

func stdDev(values []float64) float64 {
	mean := average(values)
	total := 0.0
	for _, value := range values {
		total += math.Pow(value-mean, 2)
	}
	variance := total / float64(len(values)-1)
	return math.Sqrt(variance)
}
