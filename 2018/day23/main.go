package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const registerCount = 6

const target = 10551428

var rInstruction = regexp.MustCompile(`pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(\d+)`)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if nanobots, err := getNanobots(fileName); err != nil {
		log.Fatal(err)
	} else {
		res := processNanobots(nanobots)
		// res2 := getMostInRange(nanobots)
		res2 := checkOverlaps(nanobots)
		fmt.Println("Part1 result is", res, res2)
	}
}

type Coordinate struct {
	X int
	Y int
	Z int
}

func (c Coordinate) String() string {
	return fmt.Sprintf("%d,%d,%d", c.X, c.Y, c.Z)
}

type Nanobot struct {
	Coordinate
	Radius int
}

func (n Nanobot) String() string {
	return fmt.Sprintf("%v(%d)", n.Coordinate, n.Radius)
}

func (n Nanobot) intersects(n2 Nanobot) bool {
	return distance(n.Coordinate, n2.Coordinate) <= n.Radius+n2.Radius
}

func getNanobots(fileName string) ([]Nanobot, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var nanobots []Nanobot
	for scanner.Scan() {
		line := scanner.Text()
		extracts := rInstruction.FindAllStringSubmatch(line, -1)
		if len(extracts) != 1 || len(extracts[0]) != 5 {
			return nil, errors.New("Can't parse instruction line " + line)
		}
		nanobot := Nanobot{
			Coordinate: Coordinate{
				X: atoi(extracts[0][1]),
				Y: atoi(extracts[0][2]),
				Z: atoi(extracts[0][3]),
			},
			Radius: atoi(extracts[0][4]),
		}
		nanobots = append(nanobots, nanobot)
	}
	return nanobots, nil
}

func processNanobots(nanobots []Nanobot) int {
	var maxRadius int
	var maxNanobot Nanobot
	for _, nanobot := range nanobots {
		if nanobot.Radius > maxRadius {
			maxRadius = nanobot.Radius
			maxNanobot = nanobot
		}
	}
	var count int
	for _, nanobot := range nanobots {
		if distance(maxNanobot.Coordinate, nanobot.Coordinate) <= maxNanobot.Radius {
			count++
		}
	}
	return count
}

func checkOverlaps(nanobots []Nanobot) int {
	insersections := make(map[Nanobot][]Nanobot)
	for i, n1 := range nanobots {
		for _, n2 := range nanobots[i+1:] {
			if n1.intersects(n2) {
				insersections[n1] = append(insersections[n1], n2)
				insersections[n2] = append(insersections[n2], n1)
			}
		}
	}
	clique := getMaxClique(insersections)

	var maxClosestEdge int
	var origin Coordinate
	for n := range clique {
		if d := distance(origin, n.Coordinate) - n.Radius; d > maxClosestEdge {
			maxClosestEdge = d
		}
	}
	return maxClosestEdge
}

func getMaxClique(intersections map[Nanobot][]Nanobot) map[Nanobot]bool {
	P := make(map[Nanobot]bool, len(intersections))
	for n := range intersections {
		P[n] = true
	}
	return getMaxCliqueR(map[Nanobot]bool{}, P, map[Nanobot]bool{}, intersections)
}

func getMaxCliqueR(R, P, X map[Nanobot]bool, intersections map[Nanobot][]Nanobot) map[Nanobot]bool {
	if len(P) == 0 && len(X) == 0 {
		return R
	}

	var maxClique map[Nanobot]bool

	var u Nanobot

	if len(P) > 0 {
		for k := range P {
			u = k
			break
		}
	} else {
		for k := range X {
			u = k
			break
		}
	}

	for n := range removeFromSet(P, intersections[u]) {
		clique := getMaxCliqueR(
			singleUnion(R, n),
			intersectionIntersection(P, intersections[n]),
			intersectionIntersection(X, intersections[n]),
			intersections,
		)
		if len(clique) > len(maxClique) {
			maxClique = clique
		}
		delete(P, n)
		X[n] = true
	}
	return maxClique
}

func singleUnion(S map[Nanobot]bool, n Nanobot) map[Nanobot]bool {
	SB := make(map[Nanobot]bool, len(S)+1)
	for k, v := range S {
		SB[k] = v
	}
	SB[n] = true
	return SB
}

func intersectionIntersection(S map[Nanobot]bool, intersections []Nanobot) map[Nanobot]bool {
	res := make(map[Nanobot]bool)
	for _, n := range intersections {
		if v, ok := S[n]; v && ok {
			res[n] = v
		}
	}
	return res
}

func removeFromSet(S map[Nanobot]bool, intersections []Nanobot) map[Nanobot]bool {
	lookup := make(map[Nanobot]bool, len(S)-len(intersections))
	for _, n := range intersections {
		lookup[n] = true
	}
	SB := make(map[Nanobot]bool)
	for k := range S {
		if !lookup[k] {
			SB[k] = true
		}
	}
	return SB
}

func distance(c1, c2 Coordinate) int {
	return Abs(c1.X-c2.X) + Abs(c1.Y-c2.Y) + Abs(c1.Z-c2.Z)
}

// unsafe string -> integer parsing
func atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Useless code, too slow
func getMostInRange(nanobots []Nanobot) int {
	inRanges := make(map[Coordinate]int)
	for _, nanobot := range nanobots {
		for _, c := range nanobot.inRange() {
			inRanges[c]++
		}
	}

	var maxInRange int
	var maxInRangeCoordinate Coordinate
	for c, d := range inRanges {
		if d > maxInRange {
			maxInRange = d
			maxInRangeCoordinate = c
		}
	}
	return distance(Coordinate{}, maxInRangeCoordinate)
}

func (n Nanobot) inRange() []Coordinate {
	var coordinates []Coordinate
	for i := -n.Radius; i <= n.Radius; i++ {
		for j := -n.Radius + Abs(i); j <= n.Radius-Abs(i); j++ {
			for k := -n.Radius + Abs(i) + Abs(j); k <= n.Radius-Abs(i)-Abs(j); k++ {
				coordinates = append(coordinates, Coordinate{X: n.X + i, Y: n.Y + j, Z: n.Z + k})
			}
		}
	}
	return coordinates
}

func (n Nanobot) toBox() Box {
	return Box{
		Top:    Coordinate{X: n.X, Y: n.Y, Z: n.Z + n.Radius},
		Bottom: Coordinate{X: n.X, Y: n.Y, Z: n.Z - n.Radius},
		Left:   Coordinate{X: n.X, Y: n.Y - n.Radius, Z: n.Z},
		Right:  Coordinate{X: n.X, Y: n.Y + n.Radius, Z: n.Z},
		Front:  Coordinate{X: n.X - n.Radius, Y: n.Y, Z: n.Z},
		Back:   Coordinate{X: n.X + n.Radius, Y: n.Y, Z: n.Z},
	}
}

type Box struct {
	Top    Coordinate
	Bottom Coordinate
	Left   Coordinate
	Right  Coordinate
	Front  Coordinate
	Back   Coordinate
}
