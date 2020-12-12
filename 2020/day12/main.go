package main

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/thlacroix/goadvent/helpers"
)

// Cmd is a parsed command from the input, with its Action and its Value
type Cmd struct {
	A rune
	V int
}

func main() {
	var part1, part2 int
	cmds := make([]Cmd, 0, 1000)
	err := helpers.ScanLine("input.txt", func(s string) error {
		a := rune(s[0])
		v, err := strconv.Atoi(string(s[1:]))
		if err != nil {
			return err
		}
		cmds = append(cmds, Cmd{a, v})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	s1 := State{C: C{0, 0}, Dir: C{1, 0}}
	s2 := State{C: C{0, 0}, Dir: C{10, 1}}

	for _, c := range cmds {
		s1 = s1.Move(c)
		s2 = s2.MoveW(c)
	}
	part1 = Abs(s1.X) + Abs(s1.Y)
	part2 = Abs(s2.X) + Abs(s2.Y)
	fmt.Println(part1, part2)
}

// C represents a coordinate
type C struct {
	X, Y int
}

// Turn rotates a coordinate around the center
func (c C) Turn(d int) C {
	dp := float64(d) * math.Pi / 180
	return C{
		X: int(math.Round(math.Cos(dp)*float64(c.X) - math.Sin(dp)*float64(c.Y))),
		Y: int(math.Round(math.Sin(dp)*float64(c.X) + math.Cos(dp)*float64(c.Y))),
	}
}

// State of the board and its direction (which is the relative waypoint in part 2)
type State struct {
	C
	Dir C
}

// Move moves a state to another state from a command for part 1
func (s State) Move(c Cmd) State {
	switch c.A {
	case 'N':
		return State{C: C{s.X, s.Y + c.V}, Dir: s.Dir}
	case 'S':
		return State{C: C{s.X, s.Y - c.V}, Dir: s.Dir}
	case 'E':
		return State{C: C{s.X + c.V, s.Y}, Dir: s.Dir}
	case 'W':
		return State{C: C{s.X - c.V, s.Y}, Dir: s.Dir}
	case 'R':
		return State{C: s.C, Dir: s.Dir.Turn(-c.V)}
	case 'L':
		return State{C: s.C, Dir: s.Dir.Turn(c.V)}
	case 'F':
		return State{C: C{s.X + c.V*s.Dir.X, s.Y + c.V*s.Dir.Y}, Dir: s.Dir}
	}
	return State{}
}

// MoveW moves a state to another state from a command for part 2
func (s State) MoveW(c Cmd) State {
	switch c.A {
	case 'N':
		return State{C: s.C, Dir: C{s.Dir.X, s.Dir.Y + c.V}}
	case 'S':
		return State{C: s.C, Dir: C{s.Dir.X, s.Dir.Y - c.V}}
	case 'E':
		return State{C: s.C, Dir: C{s.Dir.X + c.V, s.Dir.Y}}
	case 'W':
		return State{C: s.C, Dir: C{s.Dir.X - c.V, s.Dir.Y}}
	case 'R':
		return State{C: s.C, Dir: s.Dir.Turn(-c.V)}
	case 'L':
		return State{C: s.C, Dir: s.Dir.Turn(c.V)}
	case 'F':
		return State{C: C{s.X + c.V*s.Dir.X, s.Y + c.V*s.Dir.Y}, Dir: s.Dir}
	}
	return State{}
}

// Abs returns the absolute value of an int
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
