package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		fmt.Println("Received termination signal. Exiting...")
		os.Exit(0)
	}()

	for {
		var char string

		fmt.Println("Please provide a character for prediction:")
		_, err = fmt.Scan(&char)
		if err != nil {
			panic(err)
		}

		if len(char) != 1 {
			fmt.Println("Character must be a single character:")
			os.Exit(1)
		}

		res, err := CountNextCharAppearance(path, []rune(char)[0])
		if err != nil {
			fmt.Println("Error:", err)
		}

		for k, v := range res {
			fmt.Println(fmt.Sprintf("%c : %d", k, v))
		}
	}
}

func CountNextCharAppearance(path string, givenChar rune) (map[rune]int, error) {
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
			var foundChar bool

			for _, char := range text {
				if foundChar {
					if (char >= 'a' && char <= 'z') || char == ' ' {
						charCount[char]++
					}
				}
				if char == givenChar {
					foundChar = true
				} else {
					foundChar = false
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
