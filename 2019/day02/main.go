package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	ints, err := getInts("day02input.txt")
	if err != nil {
		log.Fatal(err)
	}
	input := make([]int, len(ints))
	copy(input, ints)
	input = initInts(input, 12, 2)
	input = processInts(input)
	fmt.Println(input[0])

	nonce, verb := findNonceVerb(ints, 19690720)
	fmt.Println(100*nonce + verb)
}

func getInts(fileName string) ([]int, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	split := strings.Split(strings.TrimSpace(string(content)), ",")
	ints := make([]int, 0, len(split))
	for _, c := range split {
		i, err := strconv.Atoi(c)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func initInts(ints []int, nonce, verb int) []int {
	ints[1] = nonce
	ints[2] = verb
	return ints
}

func processInts(ints []int) []int {
	var index int

	for index < len(ints) {
		switch ints[index] {
		case 1:
			ints[ints[index+3]] = ints[ints[index+1]] + ints[ints[index+2]]
		case 2:
			ints[ints[index+3]] = ints[ints[index+1]] * ints[ints[index+2]]
		case 99:
			return ints
		}
		index += 4
	}

	return ints
}

func findNonceVerb(ints []int, target int) (int, int) {
	input := make([]int, len(ints))
	for nonce := 0; nonce < 100; nonce++ {
		for verb := 0; verb < 100; verb++ {
			copy(input, ints)
			input = initInts(input, nonce, verb)
			input = processInts(input)
			if input[0] == target {
				return nonce, verb
			}
		}
	}
	return 0, 0
}
