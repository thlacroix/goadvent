package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// Mode represent the mode of a parameter defined in the operatioh=n
type Mode byte

const (
	Position Mode = iota
	Immediate
	Relative
)

func main() {
	ints, err := getInts("day11input.txt")
	if err != nil {
		log.Fatal(err)
	}
	intsCopy := make([]int, len(ints))
	copy(intsCopy, ints)
	painted, _ := paintAndCount(intsCopy, InitialBlack)
	fmt.Println(painted)
	copy(intsCopy, ints)
	_, tableau := paintAndCount(intsCopy, White)
	paintTableau(tableau)
}

// Point holds the coordinates of a point in the tableau
type Point struct {
	X int
	Y int
}

// Color of a panel
// Using different colors for the initial black (not painted)
// and the painted black
type Color byte

const (
	InitialBlack Color = iota
	Black
	White
)

// Direction can both represent the turn direction of the robot
// and the current direction it looks toward
type Direction byte

const (
	Nowhere Direction = iota
	Left
	Right
	Up
	Down
)

// paintAndCount uses the IntCode program to paint the tableau and move the robot
// It uses an initial color that is different for part 1 and 2
// It returns the number of panels painted, and the tableau
func paintAndCount(ints []int, initialColor Color) (int, map[Point]Color) {
	tableau := make(map[Point]Color)
	var painted int
	var index, base, x, y int

	var (
		paintColor    Color
		currentColor  Color
		turnDirection Direction
	)

	robotDirection := Up
	tableau[Point{0, 0}] = initialColor

	for robotDirection != Nowhere {
		currentColor = tableau[Point{x, y}]
		paintColor, turnDirection, ints, index, base = processInts(ints, index, base, currentColor)
		if turnDirection == Nowhere {
			break
		}
		if currentColor == InitialBlack {
			painted++
		}
		tableau[Point{x, y}] = paintColor

		switch turnDirection {
		case Nowhere:
			robotDirection = Nowhere
		case Left:
			switch robotDirection {
			case Up:
				robotDirection = Left
			case Left:
				robotDirection = Down
			case Down:
				robotDirection = Right
			case Right:
				robotDirection = Up
			}
		case Right:
			switch robotDirection {
			case Up:
				robotDirection = Right
			case Left:
				robotDirection = Up
			case Down:
				robotDirection = Left
			case Right:
				robotDirection = Down
			}
		}

		switch robotDirection {
		case Up:
			y++
		case Down:
			y--
		case Left:
			x--
		case Right:
			x++
		}
	}
	return painted, tableau
}

func getInts(fileName string) ([]int, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	split := strings.Split(strings.TrimSpace(string(content)), ",")
	ints := make([]int, 0, len(split))
	for _, c := range split {
		i, err := strconv.Atoi(c)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func processInts(ints []int, index, base int, color Color) (Color, Direction, []int, int, int) {
	var (
		colorOutput     Color
		directionOutput Direction
	)

	for index < len(ints) {
		operation := ints[index] % 100
		switch operation {
		case 1:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			ints = writeInt(ints, parameters[2], a+b, modes[2], base)
			index += 4
		case 2:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			ints = writeInt(ints, parameters[2], a*b, modes[2], base)
			index += 4
		case 3:
			modes, parameters := getModesParameters(ints[index:], 1)
			var input int
			if color == White {
				input = 1
			}
			ints = writeInt(ints, parameters[0], input, modes[0], base)
			index += 2
		case 4:
			modes, parameters := getModesParameters(ints[index:], 1)
			a := getValue(parameters[0], modes[0], ints, base)
			index += 2
			if colorOutput == InitialBlack {
				colorOutput = Color(a + 1)
			} else {
				directionOutput = Direction(a + 1)
				return colorOutput, directionOutput, ints, index, base
			}
		case 5:
			modes, parameters := getModesParameters(ints[index:], 2)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			if a != 0 {
				index = b
			} else {
				index += 3
			}
		case 6:
			modes, parameters := getModesParameters(ints[index:], 2)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			if a == 0 {
				index = b
			} else {
				index += 3
			}
		case 7:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			if a < b {
				ints = writeInt(ints, parameters[2], 1, modes[2], base)
			} else {
				ints = writeInt(ints, parameters[2], 0, modes[2], base)
			}
			index += 4
		case 8:
			modes, parameters := getModesParameters(ints[index:], 3)
			a, b := getValue(parameters[0], modes[0], ints, base), getValue(parameters[1], modes[1], ints, base)
			if a == b {
				ints = writeInt(ints, parameters[2], 1, modes[2], base)
			} else {
				ints = writeInt(ints, parameters[2], 0, modes[2], base)
			}
			index += 4
		case 9:
			modes, parameters := getModesParameters(ints[index:], 1)
			a := getValue(parameters[0], modes[0], ints, base)
			base += a
			index += 2
		case 99:
			return InitialBlack, Nowhere, ints, index, base
		}
	}

	return InitialBlack, Nowhere, ints, index, base
}

// Using a helper to write to the list, depending on the mode, and if the
// list is long enough
func writeInt(ints []int, index, value int, mode Mode, base int) []int {
	if mode == Relative {
		index += base
	}

	if index < len(ints) {
		ints[index] = value
		return ints
	}

	intsCopy := make([]int, index+1)
	copy(intsCopy, ints)
	intsCopy[index] = value

	return intsCopy
}

// Takes a param, its mode, the list of ints and the base, and return the
// values to use
func getValue(a int, mode Mode, ints []int, base int) int {
	switch mode {
	case Position:
		if a >= len(ints) {
			return 0
		}
		return ints[a]
	case Immediate:
		return a
	case Relative:
		if base+a >= len(ints) {
			return 0
		}
		return ints[base+a]
	}
	return -1
}

// Takes the operation and its paramers, and a number of parameters
// to process, to return the list of modes and parameters.
// A mode to false means by position, true means immediate
func getModesParameters(ints []int, count int) ([]Mode, []int) {
	ope := ints[0]
	modes := make([]Mode, count)
	parameters := make([]int, count)
	div := 100
	for i := 0; i < count; i++ {
		mode := ope / div % 10
		modes[i] = Mode(mode)
		div = div * 10
		parameters[i] = ints[i+1]
	}
	return modes, parameters
}

// Painting the tableau to read the registration ID
func paintTableau(tableau map[Point]Color) {
	var maxx, maxy int

	for p := range tableau {
		if p.X > maxx {
			maxx = p.X
		}
		if -p.Y > maxy {
			maxy = -p.Y
		}
	}

	regID := make([][]rune, maxy+1)
	for i := range regID {
		regID[i] = make([]rune, maxx+1)
		for j := range regID[i] {
			regID[i][j] = '.'
		}
	}

	for p, c := range tableau {
		if c == White {
			regID[-p.Y][p.X] = 'X'
		}
	}

	for _, l := range regID {
		for _, p := range l {
			fmt.Printf("%c", p)
		}
		fmt.Println()
	}
}
