package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/2019/intcode"
	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	ints, err := helpers.GetInts("day11input.txt")
	if err != nil {
		log.Fatal(err)
	}
	painted, _ := paintAndCount(ints, InitialBlack)
	fmt.Println(painted)
	_, tableau := paintAndCount(ints, White)
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
	Left Direction = iota
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
	var x, y int

	var (
		paintColor    Color
		currentColor  Color
		turnDirection Direction
	)

	robotDirection := Up
	tableau[Point{0, 0}] = initialColor

	m := intcode.NewBufferedMachine(ints, 0, 2)
	go m.Run()

	for {
		currentColor = tableau[Point{x, y}]
		var colorInput int
		if currentColor == White {
			colorInput = 1
		}
		ok := m.AddInput(colorInput)
		if !ok {
			break
		}
		paintColor = Color(m.GetOuput() + 1)
		turnDirection = Direction(m.GetOuput())

		if currentColor == InitialBlack {
			painted++
		}
		tableau[Point{x, y}] = paintColor

		switch turnDirection {
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
