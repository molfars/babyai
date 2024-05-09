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

const numWorkers = 5

func main() {
	// getting book path from terminal
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
		res, err := Complete(path, 100, []rune(start[len(start)-6:]))
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Time taken:", time.Since(now))

		// print result to terminal
		fmt.Println(fmt.Sprintf("%c", res))
	}
}

func Complete(path string, length int, start []rune) ([]rune, error) {
	chunks, err := SplitText(path)
	if err != nil {
		return nil, err
	}

	res := start

	for i := 0; i < length; i++ {
		appearance, err := CountNextCharAppearance(chunks, start)
		if err != nil {
			return nil, err
		}
		if len(appearance) == 0 {
			res = append(res, '?')
			return res, nil
		}
		next := utils.GetNextFromProbabilities(appearance)

		res = append(res, next)

		start = res[len(res)-6:]
	}

	return res, nil
}

func CountNextCharAppearance(chunks []string, start []rune) (map[string]map[rune]float64, error) {
	wg := sync.WaitGroup{}
	resultChan := make(chan map[string]map[rune]float64, numWorkers)

	for i := range chunks {
		wg.Add(1)

		go func(wg *sync.WaitGroup, text string, ch chan<- map[string]map[rune]float64) {
			defer wg.Done()

			nextCharCount := make(map[string]map[rune]float64)

			for j, char := range text {
				if char == start[len(start)-1] && j > len(start) && j < len([]rune(text)) && utils.IsWanted(rune(text[j+1])) {
					startLength := 0
					next := rune(text[j+1])
					for startLength < len(start) {
						prevText := rune(text[j-startLength])
						prevStart := start[len(start)-(1+startLength)]

						if prevText == prevStart {
							toSave := start[len(start)-(1+startLength):]
							if _, ok := nextCharCount[string(toSave)]; !ok {
								nextCharCount[string(toSave)] = make(map[rune]float64)
							}
							nextCharCount[string(toSave)][next]++
						} else {
							break
						}
						startLength++
					}
				}
			}

			ch <- nextCharCount
		}(&wg, chunks[i], resultChan)
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
