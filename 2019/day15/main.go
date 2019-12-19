package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/2019/intcode"
	"github.com/thlacroix/goadvent/helpers"
)

// Direction where the robot can go
type Direction int

const (
	Nowhere Direction = iota
	North
	South
	West
	East
)

// Status after a move
type Status int

const (
	HitWall Status = iota
	Moved
	FoundOxygen
)

// Place the robot sees
type Place int

const (
	WhoKnows Place = iota
	Wall
	Empty
	Oxygen
)

// Point coordinates
type Point struct {
	X int
	Y int
}

func main() {
	ints, err := helpers.GetInts("day15input.txt")
	if err != nil {
		log.Fatal(err)
	}

	moves, o := getMovesToOxygen(ints)
	fmt.Println(moves)
	fmt.Println(fillOxygen(ints, o))
}

// getMovesToOxygen uses a backtracking algorithm to find the oxygen place
// Return the moves to oxygen, and the coordinates of the oxygen
func getMovesToOxygen(ints []int) (int, Point) {
	m := intcode.NewMachine(ints)
	go m.Run()
	world := make(map[Point]Place)
	return searchOxygen(m, 0, Point{}, world, false)
}

// fillOxygen counts how many minutes we need to fully fill the room with oxygen
func fillOxygen(ints []int, o Point) int {
	m := intcode.NewMachine(ints)
	go m.Run()
	world := make(map[Point]Place)
	searchOxygen(m, 0, Point{}, world, true)
	oxygens := []Point{o}
	newOxygens := oxygens
	var time int

	for len(newOxygens) != 0 {
		time++
		newOxygens = nil

		for _, oxygen := range oxygens {
			if left := (Point{oxygen.X - 1, oxygen.Y}); world[left] == Empty {
				world[left] = Oxygen
				newOxygens = append(newOxygens, left)
			}
			if right := (Point{oxygen.X + 1, oxygen.Y}); world[right] == Empty {
				world[right] = Oxygen
				newOxygens = append(newOxygens, right)
			}
			if up := (Point{oxygen.X, oxygen.Y + 1}); world[up] == Empty {
				world[up] = Oxygen
				newOxygens = append(newOxygens, up)
			}
			if down := (Point{oxygen.X, oxygen.Y - 1}); world[down] == Empty {
				world[down] = Oxygen
				newOxygens = append(newOxygens, down)
			}
		}
		oxygens = append(oxygens, newOxygens...)
	}

	return time
}

// helper to print the world
func printWorld(world map[Point]Place) {
	var minx, miny, maxx, maxy int
	first := true

	for p := range world {
		if first {
			minx, maxx = p.X, p.X
			miny, maxy = p.Y, p.Y
			first = false
		} else {
			if p.X < minx {
				minx = p.X
			}
			if p.Y < miny {
				miny = p.Y
			}
			if p.X > maxx {
				maxx = p.X
			}
			if p.Y > maxy {
				maxy = p.Y
			}
		}
	}

	matrix := make([][]Place, maxy-miny+1)

	for i := range matrix {
		matrix[i] = make([]Place, maxx-minx+1)
	}

	for point, place := range world {
		matrix[maxy-point.Y][point.X-minx] = place
	}
	fmt.Printf("\033[0;0H") // comment the line if not using an ANSI terminal
	for _, l := range matrix {
		for _, p := range l {
			switch p {
			case Wall:
				fmt.Print("#")
			case Empty:
				fmt.Print(".")
			case WhoKnows:
				fmt.Print("?")
			case Oxygen:
				fmt.Print("O")
			}
		}
		fmt.Println()
	}
}

// searchOxygen is the main part of the backtracking algo
// When exploreAll is false, we explore until we find the oxygen
// Otherwise we explore the all map
func searchOxygen(m *intcode.Machine, moves int, p Point, world map[Point]Place, exploreAll bool) (int, Point) {
	// East
	if d := (Point{p.X + 1, p.Y}); world[d] == WhoKnows {
		m.AddInput(int(East))
		s := Status(m.GetOutput())
		switch s {
		case HitWall:
			world[d] = Wall
		case FoundOxygen:
			if !exploreAll {
				return moves + 1, d
			}
			fallthrough
		case Moved:
			world[d] = Place(s + 1)
			res, o := searchOxygen(m, moves+1, d, world, exploreAll)
			if res != 0 {
				return res, o
			}
			m.AddInput(int(West))
			m.GetOutput()
		}
	}

	// West
	if d := (Point{p.X - 1, p.Y}); world[d] == WhoKnows {
		m.AddInput(int(West))
		s := Status(m.GetOutput())
		switch s {
		case HitWall:
			world[d] = Wall
		case FoundOxygen:
			if !exploreAll {
				return moves + 1, d
			}
			fallthrough
		case Moved:
			world[d] = Place(s + 1)
			res, o := searchOxygen(m, moves+1, d, world, exploreAll)
			if res != 0 {
				return res, o
			}
			m.AddInput(int(East))
			m.GetOutput()
		}
	}

	// North
	if d := (Point{p.X, p.Y + 1}); world[d] == WhoKnows {
		m.AddInput(int(North))
		s := Status(m.GetOutput())
		switch s {
		case HitWall:
			world[d] = Wall
		case FoundOxygen:
			if !exploreAll {
				return moves + 1, d
			}
			fallthrough
		case Moved:
			world[d] = Place(s + 1)
			res, o := searchOxygen(m, moves+1, d, world, exploreAll)
			if res != 0 {
				return res, o
			}
			m.AddInput(int(South))
			m.GetOutput()
		}
	}

	// South
	if d := (Point{p.X, p.Y - 1}); world[d] == WhoKnows {
		m.AddInput(int(South))
		s := Status(m.GetOutput())
		switch s {
		case HitWall:
			world[d] = Wall
		case FoundOxygen:
			if !exploreAll {
				return moves + 1, d
			}
			fallthrough
		case Moved:
			world[d] = Place(s + 1)
			res, o := searchOxygen(m, moves+1, d, world, exploreAll)
			if res != 0 {
				return res, o
			}
			m.AddInput(int(North))
			m.GetOutput()
		}
	}

	return 0, Point{}
}
