package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Direction uint8

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Move struct {
	Direction Direction
	Length    int
}

type Point struct {
	X int
	Y int
}

type Line struct {
	Start Point
	End   Point
}

// For this problem I chose to represent the ropes as a set of lines
// that are used to find the intersections between the horizontal lines
// of a rope and the vertical lines of the other rope.
// The other alternative was a brute force approach, drawing all points
// and adding them in data structure (map, set, ...) to dinc those where
// we've been twice.
// Looking at the input, the line approach looked better in terms of time
// and space complexity, but needed some refacto / duplication for part 2.
func main() {
	moves1, moves2, err := getMoves("day03input.txt")
	if err != nil {
		log.Fatal(err)
	}
	// First we get all intersect points
	intersects := getAllIntersects(moves1, moves2)
	// And we find the closest
	fmt.Println(minDistance(intersects))

	// For part 2 we reuse the intersects to find the closest one in terms of steps
	fmt.Println(getLowestSteps(moves1, moves2, intersects))
}

// Using a line intersection based approach between horizontal and vertical to find the points
func getAllIntersects(moves1, moves2 []Move) []Point {
	horLines1, verLines1 := getLinesFromMoves(moves1)
	horLines2, verLines2 := getLinesFromMoves(moves2)
	intersects := getIntersects(horLines1, verLines2)
	intersects = append(intersects, getIntersects(horLines2, verLines1)...)
	return intersects
}

// Parsing the file to have the moves in structs
func getMoves(filename string) ([]Move, []Move, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 2 {
		return nil, nil, fmt.Errorf("Got %d lines", len(lines))
	}
	moves1, err := getMovesFromString(lines[0])
	if err != nil {
		return nil, nil, err
	}
	moves2, err := getMovesFromString(lines[1])
	if err != nil {
		return nil, nil, err
	}
	return moves1, moves2, nil
}

// Using assumption that input is correct, we take a string list of
// moves and return the structs
func getMovesFromString(s string) ([]Move, error) {
	moves := strings.Split(s, ",")
	res := make([]Move, 0, len(moves))

	for _, m := range moves {
		// First character will always be on 1 byte
		length, err := strconv.Atoi(m[1:])
		if err != nil {
			return nil, err
		}
		var dir Direction
		switch m[0] {
		case 'U':
			dir = Up
		case 'D':
			dir = Down
		case 'L':
			dir = Left
		case 'R':
			dir = Right
		default:
			return nil, fmt.Errorf("Unknown direction %c", m[0])
		}
		res = append(res, Move{dir, length})
	}

	return res, nil
}

// We look at the moves, and extract the horizontal and vertical lines.
// In this case, Start will always be before End on the same direction
func getLinesFromMoves(moves []Move) ([]Line, []Line) {
	var horLines, verLines []Line
	var from Point
	for _, m := range moves {
		to := from
		switch m.Direction {
		case Up:
			to.Y += m.Length
			verLines = append(verLines, Line{from, to})
		case Down:
			to.Y -= m.Length
			verLines = append(verLines, Line{to, from})
		case Left:
			to.X -= m.Length
			horLines = append(horLines, Line{to, from})
		case Right:
			to.X += m.Length
			horLines = append(horLines, Line{from, to})
		}
		from = to
	}
	return horLines, verLines
}

// We look at the moves and extract the lines.
// Using real start and end in this case
func getLinesFromMovesInOrder(moves []Move) []Line {
	var lines []Line
	var from Point
	for _, m := range moves {
		to := from
		switch m.Direction {
		case Up:
			to.Y += m.Length
		case Down:
			to.Y -= m.Length
		case Left:
			to.X -= m.Length
		case Right:
			to.X += m.Length
		}
		lines = append(lines, Line{from, to})
		from = to
	}
	return lines
}

// We get all intersect points between a set of horizontal lines and a set
// of vertical lines
func getIntersects(horLines, verLines []Line) []Point {
	var intersects []Point

	for _, horLine := range horLines {
		for _, verLine := range verLines {
			if p, doIntersect := intersect(horLine, verLine); doIntersect {
				intersects = append(intersects, p)
			}
		}
	}

	return intersects
}

// Helper to check if two lines intersect, assuming that Start will always
// be below End in all directions
func intersect(hL, vL Line) (Point, bool) {
	if (hL.Start.Y >= vL.Start.Y && hL.Start.Y <= vL.End.Y) &&
		(vL.Start.X >= hL.Start.X && vL.Start.X <= hL.End.X) {
		return Point{vL.Start.X, hL.Start.Y}, true
	}
	return Point{}, false
}

// Returns the smallest Manhattan distance of a set of points
func minDistance(points []Point) int {
	var min int
	for _, p := range points {
		d := Abs(p.X) + Abs(p.Y)
		if d == 0 {
			continue
		}
		if min == 0 || d < min {
			min = d
		}
	}
	return min
}

// Abs for ints
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// pointInLine returns true if a point is in a line
func pointInLine(point Point, line Line) bool {
	if line.Start.X == line.End.X && point.X == line.Start.X {
		miny, maxy := getMinMax(line.Start.Y, line.End.Y)
		return point.Y >= miny && point.Y <= maxy
	}

	if line.Start.Y == line.End.Y && point.Y == line.Start.Y {
		minx, maxx := getMinMax(line.Start.X, line.End.X)
		return point.X >= minx && point.X <= maxx
	}

	return false
}

// Orders two ints
func getMinMax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

// We get the lines, and then for each intersection point compute the total
// number of steps, then return the smallest
func getLowestSteps(moves1, moves2 []Move, intersects []Point) int {
	intersectSteps := make(map[Point]int, len(intersects))
	lines1 := getLinesFromMovesInOrder(moves1)
	lines2 := getLinesFromMovesInOrder(moves2)

	for _, point := range intersects {
		intersectSteps[point] += stepsToPoint(lines1, point)
		intersectSteps[point] += stepsToPoint(lines2, point)
	}

	return minSteps(intersectSteps)
}

// Computes how many steps on a rope are necessary to reach a specific point
func stepsToPoint(lines []Line, point Point) int {
	var steps int
	for _, l := range lines {
		if pointInLine(point, l) {
			return steps + distance(l.Start, point)
		}
		steps += distance(l.Start, l.End)
	}
	return 0
}

// computes a min for map values
func minSteps(steps map[Point]int) int {
	var min int
	for _, s := range steps {
		if min == 0 || s < min {
			min = s
		}
	}
	return min
}

// Manhattan distance between two points
func distance(p1, p2 Point) int {
	return Abs(p1.X-p2.X) + Abs(p1.Y-p2.Y)
}
