package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"unicode"

	"github.com/thlacroix/goadvent/helpers"
)

// Comments to be added later
func main() {
	tunnelsPart1, err := getTunnels("day18input.txt")
	if err != nil {
		log.Fatal(err)
	}
	objects := getObjects(tunnelsPart1)
	fmt.Println(getShortestPath(tunnelsPart1, objects, []rune{'@'}))
	tunnelsPart2, err := getTunnels("day18inputbis.txt")
	if err != nil {
		log.Fatal(err)
	}
	objectsPart2 := getObjects(tunnelsPart2)
	fmt.Println(getShortestPath(tunnelsPart2, objectsPart2, []rune{'@', '%', '^', '$'}))
}

type State struct {
	Point Point
	Path  Path
}

func getShortestPath(tunnels [][]rune, objects map[rune]Point, starts []rune) int {
	paths := getDistancesToKey(tunnels, objects)
	return getShortestToGetAllkeys(paths, starts)
}

type ObjectState struct {
	Objects  []rune
	Distance int
	Seen     ObjectSet
}

func (o ObjectState) Index() string {
	var s strings.Builder
	for _, o := range o.Objects {
		s.WriteRune(o)
	}
	s.WriteString("->")
	seen := make([]rune, len(o.Seen))
	for v := range o.Seen {
		seen = append(seen, v)
	}
	sort.Slice(seen, func(i, j int) bool {
		return seen[i] < seen[j]
	})
	for _, v := range seen {
		s.WriteRune(v)
	}
	return s.String()
}

func getShortestToGetAllkeys(paths map[rune]map[rune]Path, starts []rune) int {
	startSeen := make(ObjectSet, len(starts))
	for _, s := range starts {
		startSeen[s] = true
	}
	states := []ObjectState{{Objects: starts, Distance: 0, Seen: startSeen}}
	cache := make(map[string]int)

	for {
		// getting closest available in state
		var (
			currentState      ObjectState
			currentStateIndex int
		)
		var minDistance, maxSeen int
		for i, s := range states {
			if minDistance == 0 || s.Distance < minDistance || (s.Distance == minDistance && len(s.Seen) > maxSeen) {
				minDistance = s.Distance
				currentState = s
				currentStateIndex = i
				maxSeen = len(s.Seen)
			}
		}

		// finished if the closest one has found them all
		if len(currentState.Seen) == len(paths) {
			return currentState.Distance
		}
		for i, object := range currentState.Objects {
			for o, p := range paths[object] {
				if currentState.Seen[o] {
					continue
				}
				// searching if the doors are opened to go to another object
				var blocked bool
				for d := range p.Blocked {
					if !currentState.Seen[unicode.ToLower(d)] {
						blocked = true
						break
					}
				}
				if !blocked {
					seen := copySet(currentState.Seen)
					seen[o] = true
					objects := copyObjects(currentState.Objects)
					objects[i] = o
					newState := ObjectState{Objects: objects, Distance: currentState.Distance + p.Distance, Seen: seen}
					index := newState.Index()
					if cache[index] == 0 || newState.Distance < cache[index] {
						cache[index] = newState.Distance
						states = append(states, newState)
					}
				}
			}
		}

		states[currentStateIndex], states[len(states)-1] = states[len(states)-1], states[currentStateIndex]
		states = states[:len(states)-1]
	}
}

func copyObjects(s []rune) []rune {
	c := make([]rune, len(s))
	copy(c, s)
	return c
}

type ObjectSet map[rune]bool

type Path struct {
	Distance int
	Blocked  ObjectSet
}

func getDistancesToKey(tunnels [][]rune, objects map[rune]Point) map[rune]map[rune]Path {
	paths := make(map[rune]map[rune]Path)
	for object, p := range objects {
		if isDoor(object) {
			continue
		}
		queue := make(chan State, 1000)
		seen := make(map[Point]bool)
		paths[object] = make(map[rune]Path)
		queue <- State{Point: p, Path: Path{Distance: 0, Blocked: nil}}
	loop:
		for {
			select {
			case s := <-queue:
				nextPath := Path{Distance: s.Path.Distance + 1, Blocked: s.Path.Blocked}
				if r := tunnels[s.Point.Y][s.Point.X]; r != object && isKey(r) {
					paths[object][r] = s.Path
				} else if isDoor(r) {
					nextPath.Blocked = copySet(nextPath.Blocked)
					nextPath.Blocked[r] = true
				}

				if p := (Point{s.Point.X - 1, s.Point.Y}); p.X > 0 && !seen[p] && !isWall(tunnels[p.Y][p.X]) {
					queue <- State{Point: p, Path: nextPath}
				}
				if p := (Point{s.Point.X + 1, s.Point.Y}); p.X < len(tunnels[0]) && !seen[p] && !isWall(tunnels[p.Y][p.X]) {
					queue <- State{Point: p, Path: nextPath}
				}
				if p := (Point{s.Point.X, s.Point.Y - 1}); p.Y > 0 && !seen[p] && !isWall(tunnels[p.Y][p.X]) {
					queue <- State{Point: p, Path: nextPath}
				}
				if p := (Point{s.Point.X, s.Point.Y + 1}); p.Y < len(tunnels) && !seen[p] && !isWall(tunnels[p.Y][p.X]) {
					queue <- State{Point: p, Path: nextPath}
				}
				seen[s.Point] = true
			default:
				break loop
			}
		}
	}
	return paths
}

func printPaths(paths map[rune]map[rune]Path) {
	for object, p := range paths {
		fmt.Printf("%c ->", object)
		for oo, pp := range p {
			var s strings.Builder
			for ooo := range pp.Blocked {
				s.WriteRune(ooo)
			}
			fmt.Printf(" (%c %d | %v) ", oo, pp.Distance, s.String())
		}
		fmt.Println()
	}
}

func isDoor(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsUpper(r)
}

func isKey(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}

func isWall(r rune) bool {
	return r == '#'
}

func copySet(os ObjectSet) ObjectSet {
	c := make(ObjectSet, len(os))
	for k, v := range os {
		c[k] = v
	}
	return c
}

type Point struct {
	X, Y int
}

func getTunnels(filename string) ([][]rune, error) {
	var tunnels [][]rune

	return tunnels, helpers.ScanLine(filename, func(s string) error {
		tunnels = append(tunnels, []rune(s))
		return nil
	})
}

func getObjects(tunnels [][]rune) map[rune]Point {
	objects := make(map[rune]Point)
	for i, l := range tunnels {
		for j, v := range l {
			if v != '.' && v != '#' {
				objects[v] = Point{j, i}
			}
		}
	}
	return objects
}
