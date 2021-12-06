package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

type grid [5][5]int

type position struct {
	grid int
	i    int
	j    int
}

func (g grid) String() string {
	var s strings.Builder
	for _, l := range g {
		s.WriteString(fmt.Sprintln(l))
	}
	return s.String()
}

func main() {
	var part1, part2 int

	var numbers []int
	var grids []grid
	err := helpers.ScanGroup("input.txt", func(s []string) error {
		// first line
		if numbers == nil {
			split := strings.Split(s[0], ",")
			numbers = make([]int, 0, len(split))
			for _, ss := range split {
				numbers = append(numbers, atoi(ss))
			}
			return nil
		}

		// grid
		var g grid

		for i, l := range s {
			split := strings.Fields(l)
			for j, v := range split {
				g[i][j] = atoi(v)
			}
		}

		grids = append(grids, g)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	m := mapGrid(numbers, grids)
	part1 = bingo(numbers, m, grids)
	part2 = bingo2(numbers, m, grids)
	fmt.Println(part1, part2)
}

func mapGrid(numbers []int, grids []grid) map[int][]position {
	m := make(map[int][]position, len(numbers))

	for gi, g := range grids {
		for i, l := range g {
			for j, v := range l {
				m[v] = append(m[v], position{grid: gi, i: i, j: j})
			}
		}
	}
	return m
}

func bingo(numbers []int, m map[int][]position, grids []grid) int {
	lines, columns := make(map[position]int), make(map[position]int)

	var winner, lastNumberIndex int

bingoloop:
	for ni, n := range numbers {
		for _, p := range m[n] {
			lineP, columnP := position{grid: p.grid, i: p.i, j: -1}, position{grid: p.grid, j: p.j, i: -1}
			lines[lineP]++
			if lines[lineP] == 5 {
				winner = p.grid
				lastNumberIndex = ni
				break bingoloop
			}
			columns[columnP]++
			if columns[columnP] == 5 {
				winner = p.grid
				lastNumberIndex = ni
				break bingoloop
			}
		}
	}

	numberMap := make(map[int]bool, len(numbers[:lastNumberIndex+1]))
	for _, n := range numbers[:lastNumberIndex+1] {
		numberMap[n] = true
	}

	g := grids[winner]

	var c int

	for _, l := range g {
		for _, v := range l {
			if !numberMap[v] {
				c += v
			}
		}
	}

	return c * numbers[lastNumberIndex]
}

func bingo2(numbers []int, m map[int][]position, grids []grid) int {
	lines, columns := make(map[position]int), make(map[position]int)

	winners := make(map[int]bool, len(grids))

	var looser, lastNumberIndex int

bingoloop:
	for ni, n := range numbers {
		for _, p := range m[n] {
			lineP, columnP := position{grid: p.grid, i: p.i, j: -1}, position{grid: p.grid, j: p.j, i: -1}
			lines[lineP]++
			if lines[lineP] == 5 {
				winners[p.grid] = true
			}
			columns[columnP]++
			if columns[columnP] == 5 {
				winners[p.grid] = true
			}

			if len(winners) == len(grids) {
				looser = p.grid
				lastNumberIndex = ni
				break bingoloop
			}
		}
	}

	numberMap := make(map[int]bool, len(numbers[:lastNumberIndex+1]))
	for _, n := range numbers[:lastNumberIndex+1] {
		numberMap[n] = true
	}

	g := grids[looser]

	var c int

	for _, l := range g {
		for _, v := range l {
			if !numberMap[v] {
				c += v
			}
		}
	}

	return c * numbers[lastNumberIndex]
}

// unsafe string -> integer parsing
func atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}
