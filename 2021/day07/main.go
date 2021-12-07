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
	part1 = moveCrabs(ints, func(i int) int {
		return i
	})
	part2 = moveCrabs(ints, func(i int) int {
		return i * (i + 1) / 2
	})
	fmt.Println(part1, part2)
}

func moveCrabs(ints []int, f func(int) int) int {
	minP, maxP := ints[0], ints[0]
	for _, i := range ints {
		minP, maxP = min(i, minP), max(i, maxP)
	}

	var minFuel int

	for t := minP; t <= maxP; t++ {
		var c int

		for _, i := range ints {
			c += f(abs(i - t))
		}

		if minFuel == 0 || c < minFuel {
			minFuel = c
		}
	}
	return minFuel
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}
