package utils

import "math/rand"

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
