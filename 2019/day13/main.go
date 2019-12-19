package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/2019/intcode"
	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	ints, err := helpers.GetInts("day13input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(countTiles(ints))
	ints[0] = 2
	fmt.Println(playGame(ints))

}

// countTiles simply counts the number of blocks in the screen
func countTiles(ints []int) int {
	var count int
	m := intcode.NewBufferedMachine(ints, 0, 3)
	go m.Run()

	for {
		_, end := m.GetOutputOrEnd()
		if end {
			return count
		}
		m.GetOutput()

		tile := m.GetOutput()
		if tile == 2 {
			count++
		}
	}
}

// Object represents the object of a tile
type Object int

const (
	Empty Object = iota
	Wall
	Block
	Paddle
	Ball
)

// Game holds all the game information
type Game struct {
	Map    [22][43]Object
	Paddle Point
	Ball   Point
	Score  int
}

// Print pretty print the game
func (g *Game) Print() {
	fmt.Printf("\033[0;0H") // comment the line if not using an ANSI terminal
	for _, l := range g.Map {
		for _, t := range l {
			c := '.'
			switch t {
			case Wall:
				c = 'X'
			case Ball:
				c = 'O'
			case Paddle:
				c = '_'
			case Block:
				c = '□'
			}
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}
	fmt.Println("Score:", g.Score)
}

// Move tells us where to move the joystick
func (g *Game) Move() int {
	if g.Paddle.X > g.Ball.X {
		return -1
	} else if g.Paddle.X < g.Ball.X {
		return 1
	} else {
		return 0
	}
}

// Point simply holds coordinates of an object
type Point struct {
	X int
	Y int
}

// playGame plays the game by moving the joystick where the ball is
// and returns the end score.
// Uncomment the game.Print() lines if you want to visualize the game
func playGame(ints []int) int {
	game := &Game{}
	m := intcode.NewBufferedMachine(ints, 0, 2)
	go m.Run()

	for {
		x, input, end := m.GetOutputOrAddInputOrEnd(game.Move())
		if input {
			continue
		}
		if end {
			//game.Print()
			return game.Score
		}
		y := m.GetOutput()
		tile := m.GetOutput()

		if x == -1 {
			game.Score = tile
		} else {
			o := Object(tile)
			game.Map[y][x] = o
			switch o {
			case Ball:
				game.Ball = Point{x, y}
				//game.Print()
			case Paddle:
				game.Paddle = Point{x, y}
			}
		}
	}
}
