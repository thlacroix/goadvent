package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

type coord struct {
	x, y int
}

type vent struct {
	from, to coord
}

func main() {
	var part1, part2 int

	var vents []vent
	err := helpers.ScanLine("input.txt", func(s string) error {
		var v vent

		_, err := fmt.Sscanf(s, "%d,%d -> %d,%d", &v.from.x, &v.from.y, &v.to.x, &v.to.y)
		if err != nil {
			return err
		}

		vents = append(vents, v)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	part1 = overlap(vents, true)
	part2 = overlap(vents, false)
	fmt.Println(part1, part2)
}

func overlap(vents []vent, part1 bool) int {
	m := make(map[coord]int)
	for _, v := range vents {
		if v.from.x == v.to.x {
			from, to := v.from.y, v.to.y
			if to < from {
				from, to = to, from
			}

			for i := from; i <= to; i++ {
				m[coord{x: v.from.x, y: i}]++
			}
		} else if v.from.y == v.to.y {
			from, to := v.from.x, v.to.x
			if to < from {
				from, to = to, from
			}

			for i := from; i <= to; i++ {
				m[coord{x: i, y: v.from.y}]++
			}
		} else {
			if part1 {
				continue
			}

			fromx, fromy, tox, toy := v.from.x, v.from.y, v.to.x, v.to.y

			if fromx > tox {
				fromx, tox, fromy, toy = tox, fromx, toy, fromy
			}
			yMod := 1
			if fromy > toy {
				yMod = -1
			}

			cury := fromy

			for i := fromx; i <= tox; i++ {
				m[coord{x: i, y: cury}]++
				cury += yMod
			}
		}
	}

	var c int

	for _, v := range m {
		if v > 1 {
			c++
		}
	}
	return c
}
