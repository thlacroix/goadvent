package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	ints, err := helpers.GetIntsNL("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(increase(ints))
	fmt.Println(increase(sliding(ints)))
}

func increase(ints []int) int {
	var count, prev int

	for _, i := range ints {
		if prev != 0 && i > prev {
			count++
		}
		prev = i
	}
	return count
}

func sliding(ints []int) []int {
	out := make([]int, 0, len(ints)-2)

	for k, i := range ints[:len(ints)-2] {
		out = append(out, i+ints[k+1]+ints[k+2])
	}

	return out
}
