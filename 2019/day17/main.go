package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/thlacroix/goadvent/2019/intcode"
	"github.com/thlacroix/goadvent/helpers"
)

const scaffoldChar = '#'

func main() {
	ints, err := helpers.GetInts("day17input.txt")
	if err != nil {
		log.Fatal(err)
	}
	scaffold, robot := getScaffold(ints)
	count := countIntersects(scaffold)
	moves := getScaffoldMoves(scaffold, robot)
	moves = moves[2 : len(moves)-1]

	mainRoutine, functions := compress(moves)
	fmt.Println("Moves:", moves)
	fmt.Println("Main routine:", mainRoutine)
	fmt.Println("Functions:", functions)

	// I initially did the compression manually with vs code, and later
	// proceeded to automate the process.
	// Keeping the values manually found below:
	// mainRoutine := "A,B,A,C,B,C,B,A,C,B"
	// functions := [3]string{"L,10,L,6,R,10", "R,6,R,8,R,8,L,6,R,8", "L,10,R,8,R,8,L,10"}

	ints[0] = 2
	dust := moveOnScaffold(ints, mainRoutine, functions)
	fmt.Println(count)
	fmt.Println(dust)
}

// getScaffoldMoves make the robot move on the scaffold and build the
// sequences of moves, by moving forward until it has to turn
func getScaffoldMoves(scaffold [][]int, r Robot) string {
	var s strings.Builder
	var moves int
	for {
		var moved bool
		switch r.Direction {
		case Up:
			if r.X-1 >= 0 && scaffold[r.X-1][r.Y] == scaffoldChar {
				r.X--
				moves++
				moved = true
			}
		case Down:
			if r.X+1 < len(scaffold) && scaffold[r.X+1][r.Y] == scaffoldChar {
				r.X++
				moves++
				moved = true
			}
		case Left:
			if r.Y-1 >= 0 && scaffold[r.X][r.Y-1] == scaffoldChar {
				r.Y--
				moves++
				moved = true
			}
		case Right:
			if r.Y+1 < len(scaffold[0]) && scaffold[r.X][r.Y+1] == scaffoldChar {
				r.Y++
				moves++
				moved = true
			}
		}
		if !moved {
			s.WriteString(fmt.Sprintf("%d,", moves))
			moves = 0
			switch r.Direction {
			case Up:
				if r.Y-1 >= 0 && scaffold[r.X][r.Y-1] == scaffoldChar {
					r.Direction = Left
					r.Y--
					s.WriteString("L,")
					moves++
				} else if r.Y+1 < len(scaffold[0]) && scaffold[r.X][r.Y+1] == scaffoldChar {
					r.Direction = Right
					r.Y++
					s.WriteString("R,")
					moves++
				}
			case Down:
				if r.Y-1 >= 0 && scaffold[r.X][r.Y-1] == scaffoldChar {
					r.Direction = Left
					r.Y--
					s.WriteString("R,")
					moves++
				} else if r.Y+1 < len(scaffold[0]) && scaffold[r.X][r.Y+1] == scaffoldChar {
					r.Direction = Right
					r.Y++
					s.WriteString("L,")
					moves++
				}
			case Left:
				if r.X-1 >= 0 && scaffold[r.X-1][r.Y] == scaffoldChar {
					r.Direction = Up
					r.X--
					s.WriteString("R,")
					moves++
				} else if r.X+1 < len(scaffold) && scaffold[r.X+1][r.Y] == scaffoldChar {
					r.Direction = Down
					r.X++
					s.WriteString("L,")
					moves++
				}
			case Right:
				if r.X-1 >= 0 && scaffold[r.X-1][r.Y] == scaffoldChar {
					r.Direction = Up
					r.X--
					s.WriteString("L,")
					moves++
				} else if r.X+1 < len(scaffold) && scaffold[r.X+1][r.Y] == scaffoldChar {
					r.Direction = Down
					r.X++
					s.WriteString("R,")
					moves++
				}
			}
			if moves == 0 {
				break
			}
		}
	}
	return s.String()
}

type Robot struct {
	X, Y      int
	Direction Direction
}

type Direction byte

const (
	Up Direction = iota
	Down
	Left
	Right
)

// getting the scaffold as a map from the machine
func getScaffold(ints []int) ([][]int, Robot) {
	m := intcode.NewMachine(ints)
	var p Robot
	go m.Run()
	var scaffold [][]int
	var line []int
	var x, y int
	for {
		c, end := m.GetOutputOrEnd()
		if end {
			break
		}
		switch c {
		case '\n':
			if line != nil {
				scaffold = append(scaffold, line)
			}
			line = nil
			x++
			y = 0
		case '^':
			p.X, p.Y = x, y
			p.Direction = Up
		case 'v':
			p.X, p.Y = x, y
			p.Direction = Down
		case '<':
			p.X, p.Y = x, y
			p.Direction = Left
		case '>':
			p.X, p.Y = x, y
			p.Direction = Right
		}

		if c != '\n' {
			line = append(line, c)
			y++
		}
	}
	return scaffold, p
}

// Counts the sum of alignements
// Initially done manually, see getAlignmentParameter
func countIntersects(scaffold [][]int) int {
	var sum int
	for y, l := range scaffold {
		for x, v := range l {
			if v == scaffoldChar {
				if x-1 > 0 && scaffold[y][x-1] == scaffoldChar &&
					x+1 < len(scaffold[0]) && scaffold[y][x+1] == scaffoldChar &&
					y-1 > 0 && scaffold[y-1][x] == scaffoldChar &&
					y+1 < len(scaffold) && scaffold[y+1][x] == scaffoldChar {
					sum += x * y
				}
			}
		}
	}
	return sum
}

func printScaffold(scaffold [][]int) {
	for _, l := range scaffold {
		for _, v := range l {
			fmt.Printf("%c", v)
		}
		fmt.Println()
	}
}

// Algorithm to compress the sequence of moves, which could
// probably done in a more simple / elegant way, but this works fast.
// The main idea behind it is that it will take a substring at the
// beginning of the input, then will check if this pattern is repeated
// directly after itself.
// It then takes a substring after the first one, does the same process,
// and takes a third substring. It takes max 20 characters for each substring.
// Once we have the three substrings, we proceed on the string by searching if
// the next sequence of characters could be replaced by one of the substrings.
// One additional difficulty is the role of the ',', as we need to offset it.
// !!! There are probably some edge cases not handled for some inputs, there's
// a lot of index manipulation, that could go out of bound. I tried to catch
// some, but it worked directly with my input, so it would need more testing !!!
func compress(moves string) (string, [3]string) {
	var subs [3]string
	for len(subs[0]) <= 20 {
		startIndex := len(subs[0]) + 1
		for startIndex < len(moves) && moves[startIndex] != ',' {
			startIndex++
		}
		subs[0] = moves[:startIndex]
		startIndex++
		subs[1] = ""
		pattern := "A"

		for startIndex < len(moves) && strings.HasPrefix(moves[startIndex:], subs[0]) {
			startIndex += len(subs[0]) + 1
			pattern += ",A"
		}
		if startIndex >= len(moves) {
			continue
		}

		for len(subs[1]) <= 20 {
			startIndexLast := startIndex + len(subs[1]) + 1
			for startIndexLast < len(moves) && moves[startIndexLast] != ',' {
				startIndexLast++
			}
			if startIndexLast >= len(moves) {
				break
			}
			subs[1] = moves[startIndex:startIndexLast]
			startIndexLast++
			pattern := pattern + ",B"
			subs[2] = ""
			for startIndexLast < len(moves) {
				if strings.HasPrefix(moves[startIndexLast:], subs[0]) {
					startIndexLast += len(subs[0]) + 1
					pattern += ",A"
				} else if strings.HasPrefix(moves[startIndexLast:], subs[1]) {
					pattern += ",B"
					startIndexLast += len(subs[1]) + 1
				} else {
					break
				}
			}
			if startIndexLast >= len(moves) {
				continue
			}

			for len(subs[2]) <= 20 {
				pattern := pattern + ",C"
				index := startIndexLast + len(subs[2]) + 1
				for index < len(moves) && moves[index] != ',' {
					index++
				}

				subs[2] = moves[startIndexLast:index]
				index++

				for index < len(moves) {
					if strings.HasPrefix(moves[index:], subs[0]) {
						index += len(subs[0]) + 1
						pattern += ",A"
					} else if strings.HasPrefix(moves[index:], subs[1]) {
						pattern += ",B"
						index += len(subs[1]) + 1
					} else if strings.HasPrefix(moves[index:], subs[2]) {
						pattern += ",C"
						index += len(subs[2]) + 1
					} else {
						break
					}
				}
				if index >= len(moves) {
					return pattern, subs
				}
			}

		}
	}

	return "", [3]string{}
}

func moveOnScaffold(ints []int, mainRoutine string, functions [3]string) int {
	m := intcode.NewMachine(ints)
	go m.Run()
	show := "n"
	printAndInput(m, mainRoutine)
	printAndInput(m, functions[0])
	printAndInput(m, functions[1])
	printAndInput(m, functions[2])
	printAndInput(m, show)
	var last int
	for {
		c, end := m.GetOutputOrEnd()
		if end {
			return last
		}
		fmt.Printf("%c", last)
		last = c
	}
}

// Helper that prints the output and the send an input
// Returns true if program ends
func printAndInput(m *intcode.Machine, in string) bool {
	for {
		c, input, end := m.GetOutputOrAddInputOrEnd(int(in[0]))
		if input {
			break
		} else if end {
			return true
		}
		fmt.Printf("%c", c)
	}
	for _, c := range in[1:] {
		m.AddInput(int(c))
	}
	m.AddInput('\n')
	return false
}

// manual solution for part 1
func getAlignmentParameter() int {
	return 8*36 + 10*34 + 14*28 + 20*38 + 26*8 + 26*42 + 28*44 + 30*32 + 38*42 + 42*36 + 44*34
}
