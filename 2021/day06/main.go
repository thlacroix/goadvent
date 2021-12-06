package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var part1, part2 int
	ints, err := helpers.GetInts("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	part1 = simulate(ints, 80)
	part2 = simulate(ints, 256)
	fmt.Println(part1, part2)
}

func simulate(ints []int, days int) int {
	lanterns := make(map[int]int, len(ints))

	for _, i := range ints {
		lanterns[i]++
	}

	for i := 0; i < days; i++ {
		newLanterns := make(map[int]int, len(lanterns))

		for l, c := range lanterns {
			if l == 0 {
				newLanterns[8] += c
				newLanterns[6] += c
			} else {
				newLanterns[l-1] += c
			}
		}
		lanterns = newLanterns
	}

	var sum int

	for _, c := range lanterns {
		sum += c
	}
	return sum
}
