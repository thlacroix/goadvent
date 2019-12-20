package main

import (
	"fmt"
	"log"
	"unicode"

	"github.com/thlacroix/goadvent/helpers"
)

const portal = '®'

func main() {
	maze, err := getMaze("day20input.txt")
	if err != nil {
		log.Fatal(err)
	}
	maze, portals, start, end := simplifyMaze(maze)
	fmt.Println(steps(maze, portals, start, end, false))
	fmt.Println(steps(maze, portals, start, end, true))

}

// getting the raw maze as rune matrix
func getMaze(filename string) ([][]rune, error) {
	var maze [][]rune

	err := helpers.ScanLine(filename, func(s string) error {
		maze = append(maze, []rune(s))
		return nil
	})

	return maze, err
}

// Point hold coordinates
type Point struct {
	X, Y int
}

// Portal is a point with a name, and could be inner or outer
type Portal struct {
	Point
	Name  string
	Inner bool
}

// processes the maze to get the real location of the portals, their types, and where the go to
func simplifyMaze(maze [][]rune) ([][]rune, map[Point]Portal, Point, Point) {
	portals := make(map[Point]Portal)
	portalMap := make(map[string]Portal)

	for y, l := range maze {
		for x, v := range l {
			if unicode.IsLetter(v) {
				// we need to find the adjacent letter, and then find which
				// one is close to a '.'
				portalEntrance := Point{-1, -1}
				var name string
				var inner bool
				if x-1 >= 0 {
					if letter := maze[y][x-1]; letter == '.' {
						portalEntrance.X, portalEntrance.Y = x-1, y
						if x-1 < len(maze[0])-4 {
							inner = true
						}
					} else if unicode.IsLetter(letter) {
						name = fmt.Sprintf("%c%c", letter, v)
					}
				}
				if x+1 < len(maze[0]) {
					if letter := maze[y][x+1]; letter == '.' {
						portalEntrance.X, portalEntrance.Y = x+1, y
						if x+1 > 4 {
							inner = true
						}
					} else if unicode.IsLetter(letter) {
						name = fmt.Sprintf("%c%c", v, letter)
					}
				}
				if y-1 >= 0 {
					if letter := maze[y-1][x]; letter == '.' {
						portalEntrance.X, portalEntrance.Y = x, y-1
						if y-1 < len(maze)-4 {
							inner = true
						}
					} else if unicode.IsLetter(letter) {
						name = fmt.Sprintf("%c%c", letter, v)
					}
				}
				if y+1 < len(maze) {
					if letter := maze[y+1][x]; letter == '.' {
						portalEntrance.X, portalEntrance.Y = x, y+1
						if y+1 > 4 {
							inner = true
						}
					} else if unicode.IsLetter(letter) {
						name = fmt.Sprintf("%c%c", v, letter)
					}
				}

				if portalEntrance.X == -1 {
					continue
				}

				if p, ok := portalMap[name]; ok {
					portals[portalEntrance] = Portal{Point: p.Point, Name: name, Inner: inner}
					portals[p.Point] = Portal{Point: portalEntrance, Name: name, Inner: p.Inner}
				} else {
					portalMap[name] = Portal{Point: portalEntrance, Name: name, Inner: inner}
				}
				maze[portalEntrance.Y][portalEntrance.X] = portal
			}
		}
	}
	return maze, portals, portalMap["AA"].Point, portalMap["ZZ"].Point
}

// State for the BFS
type State struct {
	Point    Point
	Distance int
	Level    int
}

// Seen is a helper for the list of state already seen
type Seen struct {
	Point
	Level int
}

// count the number of steps from start to end with BFS
func steps(maze [][]rune, portals map[Point]Portal, start, end Point, recurse bool) int {
	queue := make(chan State, 10000)
	seen := map[Seen]bool{Seen{start, 0}: true}
	queue <- State{start, 0, 0}
loop:
	for {
		select {
		case s := <-queue:
			if s.Point == end && s.Level == 0 {
				return s.Distance
			}
			if p := (Point{s.Point.X - 1, s.Point.Y}); p.X >= 0 && !seen[Seen{p, s.Level}] && maze[p.Y][p.X] == '.' || maze[p.Y][p.X] == portal {
				queue <- State{Point: p, Distance: s.Distance + 1, Level: s.Level}
				seen[Seen{p, s.Level}] = true
			}
			if p := (Point{s.Point.X + 1, s.Point.Y}); p.X < len(maze[0]) && !seen[Seen{p, s.Level}] && maze[p.Y][p.X] == '.' || maze[p.Y][p.X] == portal {
				queue <- State{Point: p, Distance: s.Distance + 1, Level: s.Level}
				seen[Seen{p, s.Level}] = true
			}
			if p := (Point{s.Point.X, s.Point.Y - 1}); p.Y >= 0 && !seen[Seen{p, s.Level}] && maze[p.Y][p.X] == '.' || maze[p.Y][p.X] == portal {
				queue <- State{Point: p, Distance: s.Distance + 1, Level: s.Level}
				seen[Seen{p, s.Level}] = true
			}
			if p := (Point{s.Point.X, s.Point.Y + 1}); p.Y < len(maze) && !seen[Seen{p, s.Level}] && maze[p.Y][p.X] == '.' || maze[p.Y][p.X] == portal {
				queue <- State{Point: p, Distance: s.Distance + 1, Level: s.Level}
				seen[Seen{p, s.Level}] = true
			}
			if s.Point != start && s.Point != end && maze[s.Point.Y][s.Point.X] == portal {
				portal := portals[s.Point]
				var nextLevel int
				if !recurse {
					nextLevel = s.Level
				} else if portal.Inner {
					nextLevel = s.Level + 1
				} else {
					nextLevel = s.Level - 1
				}
				if (!recurse || portal.Inner || s.Level > 0) && !seen[Seen{portal.Point, nextLevel}] {
					queue <- State{Point: portal.Point, Distance: s.Distance + 1, Level: nextLevel}
					seen[Seen{portal.Point, nextLevel}] = true
				}

			}
		default:
			break loop
		}
	}

	return 0
}

// helper to print the maze
func printMaze(maze [][]rune) {
	for _, l := range maze {
		for _, v := range l {
			fmt.Printf("%c", v)
		}
		fmt.Println()
	}
}
