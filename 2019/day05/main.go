package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	ints, err := getInts("day05input.txt")
	if err != nil {
		log.Fatal(err)
	}
	intsCopy := make([]int, len(ints))
	copy(intsCopy, ints)
	processInts(intsCopy, 1)
	copy(intsCopy, ints)
	processInts(intsCopy, 5)
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

func processInts(ints []int, input int) []int {
	var index int

	for index < len(ints) {
		operation := ints[index] % 100
		switch operation {
		case 1:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints), getValue(parameters[1], modes[1], ints)
			ints[parameters[2]] = a + b
			index += 4
		case 2:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints), getValue(parameters[1], modes[1], ints)
			ints[parameters[2]] = a * b
			index += 4
		case 3:
			ints[ints[index+1]] = input
			index += 2
		case 4:
			modes, parameters := getModesParameters(ints[index:], 1)
			a := getValue(parameters[0], modes[0], ints)
			fmt.Println(a)
			index += 2
		case 5:
			modes, parameters := getModesParameters(ints[index:], 2)
			a, b := getValue(parameters[0], modes[0], ints), getValue(parameters[1], modes[1], ints)
			if a != 0 {
				index = b
			} else {
				index += 3
			}
		case 6:
			modes, parameters := getModesParameters(ints[index:], 2)
			a, b := getValue(parameters[0], modes[0], ints), getValue(parameters[1], modes[1], ints)
			if a == 0 {
				index = b
			} else {
				index += 3
			}
		case 7:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints), getValue(parameters[1], modes[1], ints)
			if a < b {
				ints[parameters[2]] = 1
			} else {
				ints[parameters[2]] = 0
			}
			index += 4
		case 8:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints), getValue(parameters[1], modes[1], ints)
			if a == b {
				ints[parameters[2]] = 1
			} else {
				ints[parameters[2]] = 0
			}
			index += 4
		case 99:
			return ints
		}
	}

	return ints
}

// Takes a param, its mode and the list of ints, and return the
// values to use
func getValue(a int, mode bool, ints []int) int {
	if mode {
		return a
	}
	return ints[a]
}

// Takes the operation and its paramers, and a number of parameters
// to process, to return the list of modes and parameters.
// A mode to false means by position, true means immediate
func getModesParameters(ints []int, count int) ([]bool, []int) {
	ope := ints[0]
	modes := make([]bool, count)
	parameters := make([]int, count)
	div := 100
	for i := 0; i < count; i++ {
		mode := ope / div % 10
		modes[i] = mode == 1
		div = div * 10
		parameters[i] = ints[i+1]
	}
	return modes, parameters
}
