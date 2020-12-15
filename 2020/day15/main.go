package main

import (
	"fmt"
)

var input = []int{0, 1, 4, 13, 15, 12, 16}

const part1Target = 2020

const part2Target = 30000000

func main() {
	var part1, part2 int
	part1 = play(input, part1Target)
	// bruteforcing part2 seems to work well, both run in 1.4s on my machine
	// will try to visualize to see if there's an easy pattern here
	part2 = play(input, part2Target)
	fmt.Println(part1, part2)
}

func play(start []int, target int) int {
	mem := make(map[int]int, target)

	for i, n := range start[:len(start)-1] {
		mem[n] = i
	}

	last := start[len(start)-1]
	var next int

	for i := len(start) - 1; i < target-1; i++ {
		if j, ok := mem[last]; !ok {
			next = 0
		} else {
			next = i - j
		}
		mem[last] = i
		last = next
	}
	return last
}
