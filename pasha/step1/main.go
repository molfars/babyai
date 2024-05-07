package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numWorkers = 50

func main() {
	var path string
	fmt.Println("Hello, please provide a path to your text file(can use './pasha/book/text.txt' as default):")
	_, err := fmt.Scan(&path)
	if err != nil {
		panic(err)
	}
	res, err := ReadWithChunkConcurrently(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for k, v := range res {
		fmt.Println(fmt.Sprintf("%c : %d", k, v))
	}
}

func ReadWithChunkConcurrently(path string) (map[rune]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	text := strings.ToLower(string(data))

	chunkSize := len(text) / numWorkers
	chunk := make([]string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(text)
		}

		chunk[i] = text[start:end]
	}

	wg := sync.WaitGroup{}

	resultChan := make(chan map[rune]int, numWorkers)

	for i := range chunk {
		wg.Add(1)
		go func(wg *sync.WaitGroup, text string, ch chan<- map[rune]int) {
			defer wg.Done()
			charCount := make(map[rune]int, len(alphabet))
			for _, char := range text {
				if char >= 'a' && char <= 'z' {
					charCount[char]++
				}
			}
			resultChan <- charCount
		}(&wg, chunk[i], resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	resultCount := make(map[rune]int, len(alphabet))
	for res := range resultChan {
		for k, v := range res {
			resultCount[k] += v
		}
	}

	return resultCount, nil
}
