package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/PavelDonchenko/mentorship/babyai/pasha/utils"
)

const numWorkers = 10

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
		now := time.Now()

		res, err := Complete(path, 200, []rune(char)[0])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Time taken:", time.Since(now))

		fmt.Println(res)
	}
}

func Complete(path string, length int, givenChar rune) (string, error) {
	builder := strings.Builder{}
	builder.WriteRune(givenChar)
	nextChar := givenChar

	appear, err := CountNextCharAppearance(path)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		if len(appear[nextChar]) < 1 {
			builder.WriteRune('?')
			return builder.String(), nil
		}

		nextChar = utils.MaxChar(appear[nextChar])
		builder.WriteRune(nextChar)
	}

	return builder.String(), nil
}

func CountNextCharAppearance(path string) (map[rune]map[rune]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	text := strings.ToLower(string(data))

	chunkSize := len(text) / numWorkers
	chunk := make([]string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize + 1
		if i == numWorkers-1 {
			end = len(text)
		}

		chunk[i] = text[start:end]
	}

	wg := sync.WaitGroup{}

	resultChan := make(chan map[rune]map[rune]int, numWorkers)

	for i := range chunk {
		wg.Add(1)

		go func(wg *sync.WaitGroup, text string, ch chan<- map[rune]map[rune]int) {
			defer wg.Done()

			charCount := make(map[rune]map[rune]int)

			for j, char := range text {
				if utils.IsChar(char) && j < len(text)-1 {
					nextChar := rune(text[j+1])
					if utils.IsChar(nextChar) {
						if _, ok := charCount[char]; !ok {
							charCount[char] = make(map[rune]int)
						}
						charCount[char][nextChar]++
					}
				}
			}

			ch <- charCount
		}(&wg, chunk[i], resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	resultCount := make(map[rune]map[rune]int)
	for res := range resultChan {
		for k, v := range res {
			if _, ok := resultCount[k]; !ok {
				resultCount[k] = make(map[rune]int)
			}
			for kk, vv := range v {
				resultCount[k][kk] += vv
			}
		}
	}

	return resultCount, nil
}
