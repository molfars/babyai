package utils

import (
	"math/rand"
	"strings"
)

const wanted = "abcdefghijklmnopqrstuvwxyz.,!? "

func MaxChar(m map[rune]int) rune {
	var res rune
	i := 0
	for k, v := range m {
		if v > i {
			res = k
			i = v
		}
	}

	return res
}

func IsChar(char rune) bool {
	if (char >= 'a' && char <= 'z') || char == ' ' {
		return true
	}

	return false
}

func IsWanted(char rune) bool {
	if strings.Contains(wanted, string(char)) {
		return true
	}
	return false
}

func PercentageAppear(data map[rune]int) map[rune]float64 {
	sum := 0
	res := make(map[rune]float64)
	for _, v := range data {
		sum += v
	}

	for k, j := range data {
		res[k] = float64(j) / float64(sum) * 100
	}

	return res
}

func SelectCharacterWithProbabilities(probabilities map[rune]float64) rune {
	randomNum := rand.Float64()

	cumulativeProb := 0.0
	for char, prob := range probabilities {
		cumulativeProb += prob
		if randomNum <= cumulativeProb {
			return char
		}
	}

	return 0
}

func GetNextFromProbabilities(probs map[string]map[rune]float64) rune {
	maxKey := ""
	for key, _ := range probs {
		if len(key) > len(maxKey) {
			maxKey = key
		}
	}

	return SelectCharacterWithProbabilities(probs[maxKey])
}

func Normalize(probs map[string]map[rune]float64) {
	for prefix, nextMap := range probs {
		total := 0.0
		for _, count := range nextMap {
			total += count
		}
		for next := range nextMap {
			probs[prefix][next] /= total
		}
	}
}
