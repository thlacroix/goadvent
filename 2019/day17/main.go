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
	fmt.Println(getScaffoldMoves(scaffold, robot))

	// generating the compression sequences manually, using vs code,
	// from the move sequence generated above
	mainRoutine := "A,B,A,C,B,C,B,A,C,B"
	A := "L,10,L,6,R,10"
	B := "R,6,R,8,R,8,L,6,R,8"
	C := "L,10,R,8,R,8,L,10"
	ints[0] = 2
	dust := moveOnScaffold(ints, mainRoutine, A, B, C)
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

func moveOnScaffold(ints []int, mainRoutine, A, B, C string) int {
	m := intcode.NewMachine(ints)
	go m.Run()
	show := "n"
	printAndInput(m, mainRoutine)
	printAndInput(m, A)
	printAndInput(m, B)
	printAndInput(m, C)
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
