package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	ints, err := getInts("day16input.txt")
	if err != nil {
		log.Fatal(err)
	}
	pattern := []int{0, 1, 0, -1}

	fmt.Println(intsToString(processInts(ints, pattern, 100)[:8]))

	mInts := intsMultipliedFrom(ints, 10000, 5977567)
	fmt.Println(intsToString(processIntsSimplified(mInts, pattern, 100)[:8]))
}

// processInts computes the pattern during N phases on the ints input
// It does the changes in place, as computing a digit doesn't need the
// previous ones.
func processInts(ints []int, pattern []int, N int) []int {
	for phase := 1; phase <= N; phase++ {
		for i := 0; i < len(ints); i++ {
			var sum int
			for j := i; j < len(ints); j++ {
				v := patternValue(pattern, j, i)
				sum += v * ints[j]
			}
			ints[i] = getLastDigit(sum)
		}
	}
	return ints
}

// processIntsSimplified works on a subset of the input, starting from
// where we want, as long as it's after the middle.
// In this case, for each digit, we don't need what comes before,
// and the pattern after in only 1s.
// We just need to compute the sum once, and then just proceed by
// removing the previous
func processIntsSimplified(ints []int, pattern []int, N int) []int {
	for phase := 1; phase <= N; phase++ {
		previous := -1
		var sum int
		for i := 0; i < len(ints); i++ {
			if previous == -1 {
				sum = 0
				for j := i; j < len(ints); j++ {
					sum += ints[j]
				}
			} else {
				sum -= previous
			}
			previous = ints[i]
			ints[i] = getLastDigit(sum)
		}
	}
	return ints
}

// gets the pattern value, based on the index of the
// digit processed (j) and the index of the digit used as input (i)
func patternValue(pattern []int, i, j int) int {
	index := ((i + 1) / (j + 1)) % len(pattern)

	return pattern[index]
}

// gets the last digit of a number
func getLastDigit(i int) int {
	return helpers.Abs(i % 10)
}

// gets a subset of a list repeated N times, starting at from
func intsMultipliedFrom(ints []int, N, from int) []int {
	mInts := make([]int, N*len(ints)-from)
	for i := range mInts {
		index := (from + i) % len(ints)
		mInts[i] = ints[index]
	}
	return mInts
}

// prints a list of ints as a number
func intsToString(ints []int) string {
	var s strings.Builder
	for _, i := range ints {
		s.WriteByte(byte(i + 48))
	}
	return s.String()
}

// reads the input to get a list of ints
func getInts(fileName string) ([]int, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return intsFromString(string(content)), nil
}

// takes a string representing a number and returns the list of ints
func intsFromString(i string) []int {
	var ints []int

	for _, c := range strings.TrimSpace(string(i)) {
		ints = append(ints, atoi(string(c)))
	}
	return ints
}

// unsafe string -> integer parsing
func atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}
