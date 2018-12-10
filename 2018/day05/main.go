package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"unicode"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if seq, err := getSequence(fileName); err != nil {
		log.Fatal(err)
	} else {
		seqLength := getSequenceLength(seq, 0)
		minLength := getShortestSequence(seq)
		fmt.Println(
			"Basic sequence length is", seqLength,
			"and minimum sequence length after element removals is", minLength,
		)
	}
}

// just reading the file
func getSequence(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

// iterating on the input sequence with forward lookup to avoid removing too
// often, and backtracks by removing only last if needed.
// `removeChar` is used for part 2 to discard elements
func getSequenceLength(initialSequence []byte, removeChar rune) int {
	explodedSequence := make([]byte, 0, len(initialSequence))

	for i := 0; i < len(initialSequence)-1; i++ {
		if unicode.ToUpper(rune(initialSequence[i])) == unicode.ToUpper(rune(removeChar)) {
			// ignoring the element
			continue
		} else if compareBytes(initialSequence[i], initialSequence[i+1]) {
			// ignoring next 2 elements if they match
			// fmt.Println("Not writing", string(initialSequence[i]), string(initialSequence[i+1]))
			i++
		} else if len(explodedSequence) > 0 && compareBytes(explodedSequence[len(explodedSequence)-1], initialSequence[i]) {
			// removing last element if it matches the next one
			// fmt.Println("Backtracking due to", string(explodedSequence[len(explodedSequence)-1]), string(initialSequence[i]))
			explodedSequence = explodedSequence[:len(explodedSequence)-1]
		} else {
			// otherwise we add the element in the sequence
			// fmt.Println("Writing", string(initialSequence[i]))
			explodedSequence = append(explodedSequence, initialSequence[i])
		}
	}
	return len(explodedSequence)
}

// for each element type, we calculate the length of the resulting sequence
// after removing the element, and we keep the shortest one
func getShortestSequence(initialSequence []byte) int {
	min := -1
	for _, c := range "abcdefghijklmnopqrstuvwxyz" {
		length := getSequenceLength(initialSequence, c)
		if min < 0 || length < min {
			min = length
		}
	}
	return min
}

func compareBytes(b1, b2 byte) bool {
	return b1 != b2 && unicode.ToUpper(rune(b1)) == unicode.ToUpper(rune(b2))
}
