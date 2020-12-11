package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var part1, part2 int
	var in [][]rune
	err := helpers.ScanLine("input.txt", func(s string) error {
		in = append(in, []rune(s))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	changed := true
	currentIn := in

	for changed {
		currentIn, changed = processMap(currentIn)
	}
	part1 = countOccupied(currentIn)

	changed = true
	currentIn = in
	for changed {
		currentIn, changed = processMap2(currentIn)
	}
	part2 = countOccupied(currentIn)
	fmt.Println(part1, part2)
}

func processMap(in [][]rune) ([][]rune, bool) {
	var changed bool
	newIn := make([][]rune, 0, len(in))

	for i := range in {
		line := make([]rune, len(in[i]))
		for j := range in[i] {
			v := in[i][j]
			var occupied int

			if j > 0 && in[i][j-1] == '#' {
				occupied++
			}

			if j < len(in[0])-1 && in[i][j+1] == '#' {
				occupied++
			}

			if i > 0 {
				if in[i-1][j] == '#' {
					occupied++
				}

				if j > 0 && in[i-1][j-1] == '#' {
					occupied++
				}

				if j < len(in[0])-1 && in[i-1][j+1] == '#' {
					occupied++
				}
			}

			if i < len(in)-1 {
				if in[i+1][j] == '#' {
					occupied++
				}

				if j > 0 && in[i+1][j-1] == '#' {
					occupied++
				}

				if j < len(in[0])-1 && in[i+1][j+1] == '#' {
					occupied++
				}
			}

			if v == 'L' && occupied == 0 {
				line[j] = '#'
				changed = true
			} else if v == '#' && occupied >= 4 {
				line[j] = 'L'
				changed = true
			} else {
				line[j] = v
			}
		}
		newIn = append(newIn, line)
	}
	return newIn, changed
}

func processMap2(in [][]rune) ([][]rune, bool) {
	var changed bool
	newIn := make([][]rune, 0, len(in))

	for i := range in {
		line := make([]rune, len(in[i]))
		for j := range in[i] {
			v := in[i][j]
			var occupied int

			for jj := j - 1; jj >= 0; jj-- {
				if in[i][jj] == '#' {
					occupied++
					break
				} else if in[i][jj] == 'L' {
					break
				}
			}

			for jj := j + 1; jj < len(in[0]); jj++ {
				if in[i][jj] == '#' {
					occupied++
					break
				} else if in[i][jj] == 'L' {
					break
				}
			}

			for ii := i - 1; ii >= 0; ii-- {
				if in[ii][j] == '#' {
					occupied++
					break
				} else if in[ii][j] == 'L' {
					break
				}
			}

			for ii := i + 1; ii < len(in); ii++ {
				if in[ii][j] == '#' {
					occupied++
					break
				} else if in[ii][j] == 'L' {
					break
				}
			}

			for inc := 1; j-inc >= 0 && i-inc >= 0; inc++ {
				if in[i-inc][j-inc] == '#' {
					occupied++
					break
				} else if in[i-inc][j-inc] == 'L' {
					break
				}
			}

			for inc := 1; j+inc < len(in[0]) && i-inc >= 0; inc++ {
				if in[i-inc][j+inc] == '#' {
					occupied++
					break
				} else if in[i-inc][j+inc] == 'L' {
					break
				}
			}

			for inc := 1; j+inc < len(in[0]) && i+inc < len(in); inc++ {
				if in[i+inc][j+inc] == '#' {
					occupied++
					break
				} else if in[i+inc][j+inc] == 'L' {
					break
				}
			}

			for inc := 1; j-inc >= 0 && i+inc < len(in); inc++ {
				if in[i+inc][j-inc] == '#' {
					occupied++
					break
				} else if in[i+inc][j-inc] == 'L' {
					break
				}
			}

			if v == 'L' && occupied == 0 {
				line[j] = '#'
				changed = true
			} else if v == '#' && occupied >= 5 {
				line[j] = 'L'
				changed = true
			} else {
				line[j] = v
			}
		}
		newIn = append(newIn, line)
	}
	return newIn, changed
}

func countOccupied(in [][]rune) int {
	var count int

	for _, l := range in {
		for _, c := range l {
			if c == '#' {
				count++
			}
		}
	}
	return count
}
