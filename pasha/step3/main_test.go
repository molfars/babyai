package main

import (
	"fmt"
	"reflect"
	"testing"
)

var want = map[rune]map[rune]int{
	'w': {'h': 1},
	't': {'h': 3, 'e': 2, 'a': 1},
	's': {' ': 2, 'e': 1, 't': 1},
	'a': {'n': 2, 't': 1},
	'y': {'o': 1, 'w': 1},
	'd': {' ': 1},
	'b': {'o': 1},
	'k': {' ': 1},
	' ': {'a': 2, 'i': 2, 's': 1, 'e': 1, 'u': 2, 'f': 1, 't': 2, 'o': 1},
	'e': {'r': 1, 'b': 1, ' ': 5, 'd': 1, 's': 1},
	'u': {'s': 1, 'n': 1},
	'r': {' ': 1, 'e': 1},
	'o': {'n': 1, 'o': 1, 'k': 1, 'r': 1, 'f': 1},
	'n': {'y': 2, 'e': 1, ' ': 1, 'i': 1},
	'h': {'i': 1, 'e': 3},
	'i': {'s': 2, 'n': 1, 't': 1},
	'f': {'o': 1, ' ': 1},
}

func TestCountNextCharAppearance(t *testing.T) {
	res, err := CountNextCharAppearance("../book/test.txt")
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range res {
		fmt.Println(fmt.Sprintf("char: %c", k))

		for kk, vv := range v {
			fmt.Println(fmt.Sprintf("%c : %d", kk, vv))

		}
	}

	if !reflect.DeepEqual(want, res) {
		t.Error("not equal")
	}
}

func TestComplete(t *testing.T) {
	res, err := Complete("../book/test.txt", 5, 'e')
	if err != nil {
		t.Fatal(err)
	}
	wantC := "e anyo"

	if !reflect.DeepEqual(wantC, res) {
		t.Error("not equal")
	}
}
