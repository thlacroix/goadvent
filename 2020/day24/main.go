package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

type Dir byte

const (
	E Dir = iota
	SE
	SW
	W
	NW
	NE
)

func NewDirs(s string) ([]Dir, error) {
	var prev, curr rune
	dirs := make([]Dir, 0, len(s)*2)

	for _, curr = range s {
		switch curr {
		case 's', 'n':
			prev = curr
		case 'e':
			switch prev {
			case 'n':
				dirs = append(dirs, NE)
			case 's':
				dirs = append(dirs, SE)
			case 0:
				dirs = append(dirs, E)
			default:
				return nil, fmt.Errorf("Not expected rune %c", prev)
			}
			prev = 0
		case 'w':
			switch prev {
			case 'n':
				dirs = append(dirs, NW)
			case 's':
				dirs = append(dirs, SW)
			case 0:
				dirs = append(dirs, W)
			default:
				return nil, fmt.Errorf("Not expected rune %c", prev)
			}
			prev = 0
		default:
			return nil, fmt.Errorf("Not expected rune %c", prev)
		}
	}

	return dirs, nil
}

func main() {
	var part1, part2 int
	tilesToFlip := make([][]Dir, 0, 500)

	err := helpers.ScanLine("input.txt", func(s string) error {
		dirs, err := NewDirs(s)

		if err != nil {
			return err
		}
		tilesToFlip = append(tilesToFlip, dirs)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	firstBlackTiles := flipTiles(tilesToFlip)
	part1 = countBlackTiles(firstBlackTiles)
	blackTilesFlippedDaily := flipXDays(firstBlackTiles, 100)
	part2 = countBlackTiles(blackTilesFlippedDaily)
	fmt.Println(part1, part2)
}

func flipTiles(tilesToFlip [][]Dir) map[P]bool {
	blackTiles := make(map[P]bool, len(tilesToFlip))

	for _, dirs := range tilesToFlip {
		p := positionFromCenter(dirs)
		blackTiles[p] = !blackTiles[p]
	}

	return blackTiles
}

func countBlackTiles(blackTiles map[P]bool) int {
	var count int

	for _, colour := range blackTiles {
		if colour {
			count++
		}
	}

	return count
}

func flipXDays(blackTiles map[P]bool, x int) map[P]bool {
	for i := 0; i < x; i++ {
		blackTiles = dailyFlip(blackTiles)
	}
	return blackTiles
}

func dailyFlip(blackTiles map[P]bool) map[P]bool {
	blackNeighbours := make(map[P]int, len(blackTiles)*4)

	dirs := [6]Dir{E, SE, SW, W, NW, NE}
	for p, b := range blackTiles {
		if !b {
			continue
		}

		for _, d := range dirs {
			var xd, yd int
			switch d {
			case E:
				xd, yd = 2, 0
			case SE:
				xd, yd = 1, -1
			case SW:
				xd, yd = -1, -1
			case W:
				xd, yd = -2, 0
			case NW:
				xd, yd = -1, 1
			case NE:
				xd, yd = 1, 1
			}
			vp := P{p.X + xd, p.Y + yd}
			blackNeighbours[vp]++
		}
	}

	newBlackTiles := make(map[P]bool, len(blackTiles))

	for p, count := range blackNeighbours {
		previousColour := blackTiles[p]

		if previousColour {
			if !(count == 0 || count > 2) {
				newBlackTiles[p] = true
			}
		} else if count == 2 {
			newBlackTiles[p] = true
		}
	}

	return newBlackTiles
}

func positionFromCenter(dirs []Dir) P {
	var x, y int

	for _, d := range dirs {
		switch d {
		case E:
			x += 2
		case W:
			x -= 2
		case NE:
			x++
			y++
		case NW:
			x--
			y++
		case SE:
			x++
			y--
		case SW:
			x--
			y--
		}
	}
	return P{x, y}
}

type P struct {
	X, Y int
}
