package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// Mode represent the mode of a parameter defined in the operatioh=n
type Mode byte

const (
	Position Mode = iota
	Immediate
	Relative
)

func main() {
	ints, err := getInts("day09input.txt")
	if err != nil {
		log.Fatal(err)
	}
	intsCopy := make([]int, len(ints))
	copy(intsCopy, ints)
	processInts(intsCopy, 1)
	copy(intsCopy, ints)
	processInts(intsCopy, 2)
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

func processInts(ints []int, input int) {
	var index, base int

	for index < len(ints) {
		operation := ints[index] % 100
		switch operation {
		case 1:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			ints = writeInt(ints, parameters[2], a+b, modes[2], base)
			index += 4
		case 2:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			ints = writeInt(ints, parameters[2], a*b, modes[2], base)
			index += 4
		case 3:
			modes, parameters := getModesParameters(ints[index:], 1)
			ints = writeInt(ints, parameters[0], input, modes[0], base)
			index += 2
		case 4:
			modes, parameters := getModesParameters(ints[index:], 1)
			a := getValue(parameters[0], modes[0], ints, base)
			fmt.Println(a)
			index += 2
		case 5:
			modes, parameters := getModesParameters(ints[index:], 2)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			if a != 0 {
				index = b
			} else {
				index += 3
			}
		case 6:
			modes, parameters := getModesParameters(ints[index:], 2)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			if a == 0 {
				index = b
			} else {
				index += 3
			}
		case 7:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			if a < b {
				ints = writeInt(ints, parameters[2], 1, modes[2], base)
			} else {
				ints = writeInt(ints, parameters[2], 0, modes[2], base)
			}
			index += 4
		case 8:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			if a == b {
				ints = writeInt(ints, parameters[2], 1, modes[2], base)
			} else {
				ints = writeInt(ints, parameters[2], 0, modes[2], base)
			}
			index += 4
		case 9:
			modes, parameters := getModesParameters(ints[index:], 1)
			a := getValue(parameters[0], modes[0], ints, base)
			base += a
			index += 2
		case 99:
			return
		}
	}

	return
}

// Using a helper to write to the list, depending on the mode, and if the
// list is long enough
func writeInt(ints []int, index, value int, mode Mode, base int) []int {
	if mode == Relative {
		index += base
	}

	if index < len(ints) {
		ints[index] = value
		return ints
	}

	intsCopy := make([]int, index+1)
	copy(intsCopy, ints)
	intsCopy[index] = value

	return intsCopy
}

// Takes a param, its mode, the list of ints and the base, and return the
// values to use
func getValue(a int, mode Mode, ints []int, base int) int {
	switch mode {
	case Position:
		if a >= len(ints) {
			return 0
		}
		return ints[a]
	case Immediate:
		return a
	case Relative:
		if base+a >= len(ints) {
			return 0
		}
		return ints[base+a]
	}
	return -1
}

// Takes the operation and its paramers, and a number of parameters
// to process, to return the list of modes and parameters.
// A mode to false means by position, true means immediate
func getModesParameters(ints []int, count int) ([]Mode, []int) {
	ope := ints[0]
	modes := make([]Mode, count)
	parameters := make([]int, count)
	div := 100
	for i := 0; i < count; i++ {
		mode := ope / div % 10
		modes[i] = Mode(mode)
		div = div * 10
		parameters[i] = ints[i+1]
	}
	return modes, parameters
}
