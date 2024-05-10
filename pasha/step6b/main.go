package main

import (
	"bufio"
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

const numWorkers = 1

func main() {
	//getting book path from terminal
	var path string
	fmt.Println("Hello, please provide a path to your text file(can use './pasha/book/text.txt' as default):")
	_, err := fmt.Scan(&path)
	if err != nil {
		panic(err)
	}

	// creating context for listening os signal and stop concurrently running function
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// create goroutine for listening context signal and stop program if signal received
	go func() {
		<-ctx.Done()
		fmt.Println("Received termination signal. Exiting...")
		os.Exit(0)
	}()

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Please provide a character for prediction:")
		start, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		start = strings.TrimSuffix(start, "\n")

		if len(start) < 5 {
			fmt.Println("Character must be a more than 5 characters long:")
			os.Exit(1)
		}
		now := time.Now()

		// Complete processed logic and return 100 predicted symbols
		model, err := BuildPredictions(path, 5)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		res := Complete(model, start[len(start)-5:], 100)
		fmt.Println("Time taken:", time.Since(now))

		// print result to terminal
		fmt.Println(fmt.Sprintf("%c", res))
	}
}

func BuildPredictions(path string, maxLength int) (map[string]map[rune]float64, error) {
	wg := sync.WaitGroup{}
	resultChan := make(chan map[string]map[rune]float64, numWorkers)

	chunks, err := SplitText(path)
	if err != nil {
		return nil, err
	}

	for c := range chunks {
		wg.Add(1)

		go func(wg *sync.WaitGroup, text string, ch chan<- map[string]map[rune]float64) {
			defer wg.Done()

			nextCharCount := make(map[string]map[rune]float64)
			for i := 0; i < len(text); i++ {
				for j := 1; j <= maxLength && i+j <= len(text); j++ {
					prefix := text[i : i+j]
					if utils.IsValidPrefix(prefix) {
						if nextCharCount[prefix] == nil {
							nextCharCount[prefix] = make(map[rune]float64)
						}
						if i+j < len(text) {
							nextLetter := rune(text[i+j])
							if utils.IsValidLetter(string(nextLetter)) {
								nextCharCount[prefix][nextLetter]++
							}
						}
					}
				}
			}

			ch <- nextCharCount
		}(&wg, chunks[c], resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	resultCount := make(map[string]map[rune]float64)
	for res := range resultChan {
		for k, v := range res {
			if _, ok := resultCount[k]; !ok {
				resultCount[k] = make(map[rune]float64)
			}
			for kk, vv := range v {
				resultCount[k][kk] += vv
			}
		}
	}

	utils.Normalize(resultCount)

	return resultCount, nil
}

func Complete(model map[string]map[rune]float64, start string, maxLength int) []rune {
	res := make([]rune, 0, maxLength)
	res = append(res, []rune(start)...)
	word := start

	for i := 0; i < maxLength; i++ {
		prefix := findPrefix(model, word)
		nextChar := utils.SelectCharacterWithProbabilities(model[prefix])
		word = string(word[1:]) + string(nextChar)
		res = append(res, nextChar)
	}

	return res
}

func findPrefix(model map[string]map[rune]float64, start string) string {
	prefix := ""

	for i := 0; i < len(start); i++ {
		if _, ok := model[start[i:]]; ok {
			prefix = start[i:]
			break
		}
	}

	return prefix
}

func SplitText(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	text := strings.ToLower(string(data))

	chunkSize := len(text) / numWorkers
	chunks := make([]string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize + 1
		if i == numWorkers-1 {
			end = len(text)
		}

		chunks[i] = text[start:end]
	}

	return chunks, nil
}
