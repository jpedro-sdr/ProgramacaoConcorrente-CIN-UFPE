package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var NumRequests = 10000
var TimeSleep time.Duration = 50

func SaveToFile(roundTripTimes []time.Duration, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Erro ao criar arquivo:", err)
		return
	}
	defer file.Close()

	for _, rt := range roundTripTimes {
		_, err := file.WriteString(fmt.Sprintf("%f\n", rt.Seconds()*1000))
		if err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return
		}
	}
}

func ReadBibleText() (string, error) {
	file, err := os.Open("C:/Projetos Faculdade/ProgramacaoConcorrente/Exercicio05/biblia.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()
	var content string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content, nil
}

func wordCount(s string) map[string]int {
	words := strings.Fields(s)
	m := make(map[string]int)
	for _, word := range words {
		m[word]++
	}
	return m
}

func concurrentWordCount(s string, numParts int) map[string]int {
	parts := make([]string, numParts)
	words := make([]map[string]int, numParts)
	var wg sync.WaitGroup
	for i := 0; i < numParts; i++ {
		start := i * len(s) / numParts
		end := (i + 1) * len(s) / numParts
		parts[i] = s[start:end]
		wg.Add(1)
		go func(i int) {
			words[i] = wordCount(parts[i])
			wg.Done()
		}(i)
	}
	wg.Wait()
	m := make(map[string]int)
	for _, wordMap := range words {
		for word, count := range wordMap {
			m[word] += count
		}
	}
	return m
}

func ExecuteWordCount(text string) (time.Duration, time.Duration) {
	numParts := 2
	start := time.Now()
	wordCount(text)
	timeWordCountWithoutConcurrency := time.Since(start)
	start = time.Now()
	concurrentWordCount(text, numParts)
	timeWordCountWithConcurrency := time.Since(start)

	return timeWordCountWithoutConcurrency, timeWordCountWithConcurrency
}
