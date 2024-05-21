package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
)

var tokenLength = flag.Int("tokens", 10, "token `length` (how many characters to take into account)")
var filePath = flag.String("file", "books/comb.txt", "`filename` to parse")
var generateLength = flag.Int("length", 100, "number of characters to generate")
var interactive = flag.Bool("i", false, "interactive mode")
var randomize = flag.Bool("random", true, "randomize output")

type Prob struct {
	ch   rune
	freq int
}

func generateNext(inp string, m map[string]int, randomize bool) rune {
	if len(inp) > *tokenLength {
		inp = inp[len(inp)-*tokenLength:]
	}
	res := make([]Prob, 0)
	totalFreq := 0
a:
	for j := 0; j < len(inp); j++ {
		addToken := func(ch rune) {
			fr := m[inp[j:]+string(ch)]
			if fr != 0 {
				res = append(res, Prob{ch, fr})
				totalFreq += fr
			}
		}
		for _, t := range SYMBOLS {
			addToken(rune(t))
		}
		for i := 'a'; i <= 'z'; i++ {
			addToken(i)
		}
		if len(res) > 10 {
			break a
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].freq > res[j].freq
	})
	if !randomize {
		return res[0].ch
	}
	if len(res) == 0 {
		return '?'
	}
	cur := 0
	target := rand.Intn(totalFreq)
	for _, p := range res {
		cur += p.freq
		if target <= cur {
			return p.ch
		}
	}
	return res[0].ch
}

const SYMBOLS = ",.!? "

func main() {

	flag.Parse()

	if !*interactive && flag.NArg() == 0 {
		fmt.Println("In non-interactive mode pass string to complete")
		os.Exit(2)
	}

	b, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Printf("Loaded file. %v bytes\n", len(b))

	result := make(map[string]int)
	token := ""
	percent := 0
	for i, ch := range b {
		npercent := i * 100 / len(b)
		if npercent > percent {
			fmt.Printf("Processing %d%%\r", npercent)
			percent = npercent
		}
		if ch > 'A' && ch < 'Z' {
			ch = ch - 'A' + 'a'
		}
		if (ch < 'a' || ch > 'z') && !strings.ContainsRune(SYMBOLS, rune(ch)) {
			continue
		}
		token += string(ch)
		if len(token) > *tokenLength {
			token = token[1:]
		}
		for i := 0; i < len(token); i++ {
			subtoken := token[i:]
			result[subtoken]++
		}
	}
	fmt.Printf("Done!                         \n")
	if !*interactive {
		s := strings.Join(flag.Args(), " ")
		for i := 0; i < *generateLength || strings.ContainsRune(SYMBOLS, rune(s[len(s)-1])); i++ {
			s += string(generateNext(s, result, *randomize))
		}
		fmt.Printf("%s\n", s)
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter text (press Enter to finish): ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSuffix(text, "\n")
			if len(text) == 0 {
				break
			}
			for i := 0; i < *generateLength || strings.ContainsRune(SYMBOLS, rune(text[len(text)-1])); i++ {
				text += string(generateNext(text, result, *randomize))
			}
			fmt.Println(text)
		}
		fmt.Println(("Done."))
	}
}
