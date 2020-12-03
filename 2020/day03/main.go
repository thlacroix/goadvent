package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var m [][]bool
	err := helpers.ScanLine("input.txt", func(s string) error {
		line := make([]bool, 0, len(s))

		for _, c := range s {
			line = append(line, c == '#')
		}
		m = append(m, line)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	part1 := processMap(m, 1, 3)
	part2 := part1 * processMap(m, 1, 1) * processMap(m, 1, 5) * processMap(m, 1, 7) * processMap(m, 2, 1)
	fmt.Println(part1, part2)
}

func processMap(m [][]bool, down, right int) int {
	var x, y, count int
	for y < len(m) {
		if m[y][x] {
			count++
		}
		x = (x + right) % len(m[0])
		y = y + down
	}
	return count
}
