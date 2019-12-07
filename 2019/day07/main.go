package main

import (
	"fmt"
	"github.com/thlacroix/goadvent/helpers"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	ints, err := getInts("day07input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(getMaxThrust(ints, [5]int{0, 1, 2, 3, 4}, getThrustPart1))
	fmt.Println(getMaxThrust(ints, [5]int{5, 6, 7, 8, 9}, getThrustPart2))
}

// getMaxThrust gets the permutations of the possible phases,
// computes the thrust for each,and return the max
func getMaxThrust(ints []int, phaseValues [5]int, getThrust func([]int, []int) int) int {
	var max int
	for _, perm := range helpers.Permute(phaseValues[:]) {
		thrust := getThrust(ints, perm)
		if thrust > max {
			max = thrust
		}
	}
	return max
}

// getThrustPart1 takes a list of amplifier and phases, and
// get the final thrust
func getThrustPart1(ints []int, phases []int) int {
	var input int
	intsCopy := make([]int, len(ints))
	for _, phase := range phases {
		copy(intsCopy, ints)
		input, _, _ = processInts(intsCopy, 0, phase, input)
	}
	return input
}

// Amplifier holds the machine and the last index
type Amplifier struct {
	Machine []int
	Index   int
}

// getThrustPart1 takes a list of amplifier and phases, and
// get the final thrust, with feedback loop.
func getThrustPart2(ints []int, phases []int) int {
	var (
		input, newInput, index int
		stop                   bool
	)
	amplifiers := make([]*Amplifier, 5)

	for i := 0; !stop; i++ {
		var phase int
		if i < 5 {
			phase = phases[i]
		} else {
			phase = input
		}
		amplifierID := i % 5
		var amplifier *Amplifier
		if amp := amplifiers[amplifierID]; amp != nil {
			amplifier = amp
		} else {
			machine := make([]int, len(ints))
			copy(machine, ints)
			amplifier = &Amplifier{Machine: machine}
			amplifiers[amplifierID] = amplifier
		}

		newInput, index, stop = processInts(amplifier.Machine, amplifier.Index, phase, input)
		amplifier.Index = index
		if !stop {
			input = newInput
		}
	}
	return input
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

// processInts now return the output, the last index, and tells if it stopped
func processInts(ints []int, index, input1, input2 int) (int, int, bool) {
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
			ints[ints[index+1]] = input1
			input1 = input2
			index += 2
		case 4:
			modes, parameters := getModesParameters(ints[index:], 1)
			a := getValue(parameters[0], modes[0], ints)
			index += 2
			return a, index, false
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
			return 0, 0, true
		}
	}

	log.Fatal("Program shouldn't be there")
	return 0, 0, false
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
