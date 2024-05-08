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

const numWorkers = 1

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
		// getting first character from terminal
		var char string
		fmt.Println("Please provide a character for prediction:")
		_, err := fmt.Scan(&char)
		if err != nil {
			panic(err)
		}

		if len(char) != 1 {
			fmt.Println("Character must be a single character:")
			os.Exit(1)
		}
		now := time.Now()

		// Complete processed logic and return 100 predicted symbols
		res, err := Complete(path, 100, []rune(char)[0])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Time taken:", time.Since(now))

		// print result to terminal
		fmt.Println(fmt.Sprintf("%c", res))
	}
}

func Complete(path string, length int, givenChar rune) ([]rune, error) {
	chunks, err := SplitText(path)
	if err != nil {
		return nil, err
	}

	// received all prediction for given character (need for the first letter and case if prediction for pair return nothing)
	allCharAppearances, err := CountNextCharAppearance(chunks)
	if err != nil {
		return nil, err
	}

	// create pair from given char and next predicted
	pair := []rune{givenChar, getMostProbabilityChar(allCharAppearances[givenChar])}

	// create slice(array) for results
	var res []rune
	res = append(res, pair...)

	for i := 0; i < length; i++ {
		var nextChar rune
		// predict next char for pair
		pairAppear, err := CountPairNextCharAppearance(chunks, pair)
		if err != nil {
			return nil, err
		}
		// in case if prediction for pair return nothing - take next char for one character prediction
		if len(pairAppear) < 1 {
			nextChar = getMostProbabilityChar(allCharAppearances[pair[0]])
			res = append(res, nextChar)
			pair = []rune{pair[1], nextChar}

			continue
		}
		// if prediction for pair return some result - take next char from this prediction
		nextChar = getMostProbabilityChar(pairAppear)
		res = append(res, getMostProbabilityChar(pairAppear))
		// create nre pair from second character from previous pair and new predicted character
		pair = []rune{pair[1], nextChar}
	}

	return res, nil
}

func CountPairNextCharAppearance(chunks []string, pair []rune) (map[rune]int, error) {
	wg := sync.WaitGroup{}

	resultChan := make(chan map[rune]int, numWorkers)

	for i := range chunks {
		wg.Add(1)

		go func(wg *sync.WaitGroup, text string, ch chan<- map[rune]int) {
			defer wg.Done()

			nextCharCount := make(map[rune]int)

			for j, char := range text {
				// loop for text, if current character = first letter from pair and after this character in text more than 2 character are exist
				// go next and check if next character in text = second character in pair - write next character appearance to hash map
				if char == pair[0] && (j+1) < len(text)-1 {
					if rune(text[j+1]) == pair[1] {
						if utils.IsChar(rune(text[j+2])) {
							nextCharCount[rune(text[j+2])]++
						}
						continue
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

	resultCount := make(map[rune]int)
	for res := range resultChan {
		for k, v := range res {
			resultCount[k] += v
		}
	}

	return resultCount, nil
}

func CountNextCharAppearance(chunks []string) (map[rune]map[rune]int, error) {
	wg := sync.WaitGroup{}

	resultChan := make(chan map[rune]map[rune]int, numWorkers)

	for i := range chunks {
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
		}(&wg, chunks[i], resultChan)
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

func getMostProbabilityChar(data map[rune]int) rune {
	percents := utils.PercentageAppear(data)
	return utils.SelectCharacterWithProbabilities(percents)
}
