package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

type P struct {
	X, Y, Z, W int
}

func (p P) String() string {
	return fmt.Sprintf("%d/%d/%d", p.X, p.Y, p.Z)
}

func main() {
	var part1, part2 int
	var x int
	activePoints := make(map[P]bool, 40)
	err := helpers.ScanLine("input.txt", func(s string) error {
		for i, c := range s {
			if c == '#' {
				activePoints[P{x, i, 0, 0}] = true
			}
		}
		x++
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	part1 = process(activePoints, 6, false)
	part2 = process(activePoints, 6, true)
	fmt.Println(part1, part2)
}

func process(activePoints map[P]bool, n int, fourD bool) int {
	for i := 0; i < n; i++ {
		activePoints = apply(activePoints, fourD)
	}

	return len(activePoints)
}

func apply(activePoints map[P]bool, fourD bool) map[P]bool {
	l := 26
	if fourD {
		l = 80
	}
	affectedPoints := make(map[P]int, l)

	for p := range activePoints {
		for x := p.X - 1; x <= p.X+1; x++ {
			for y := p.Y - 1; y <= p.Y+1; y++ {
				for z := p.Z - 1; z <= p.Z+1; z++ {
					if fourD {
						for w := p.W - 1; w <= p.W+1; w++ {
							np := P{x, y, z, w}
							if np == p {
								continue
							}
							affectedPoints[np]++
						}
					} else {
						np := P{x, y, z, 0}
						if np == p {
							continue
						}
						affectedPoints[np]++
					}

				}
			}
		}
	}

	newActives := make(map[P]bool, len(activePoints))

	for p, c := range affectedPoints {
		if activePoints[p] {
			if c == 2 || c == 3 {
				newActives[p] = true
			}
		} else {
			if c == 3 {
				newActives[p] = true
			}
		}
	}
	return newActives
}
