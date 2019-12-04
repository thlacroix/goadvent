package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath is passed")
	}
	fileName := os.Args[1]
	if points, err := getPoints(fileName); err != nil {
		log.Fatal(err)
	} else {
		res := processPoints(points)
		fmt.Println(res)
	}
}

type Point struct {
	X int
	Y int
	Z int
	T int
}

func (p Point) distance(pp Point) int {
	return Abs(p.X-pp.X) + Abs(p.Y-pp.Y) + Abs(p.Z-pp.Z) + Abs(p.T-pp.T)
}

func getPoints(fileName string) ([]Point, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var points []Point
	for scanner.Scan() {
		line := scanner.Text()
		coords := strings.Split(line, ",")
		points = append(points, Point{
			X: atoi(coords[0]),
			Y: atoi(coords[1]),
			Z: atoi(coords[2]),
			T: atoi(coords[3]),
		})
	}
	return points, nil
}

func processPoints(points []Point) int {
	inDistance := make(map[Point][]Point)
	for i, point := range points {
		for _, otherPoint := range points[i+1:] {
			if point.distance(otherPoint) <= 3 {
				inDistance[point] = append(inDistance[point], otherPoint)
				inDistance[otherPoint] = append(inDistance[otherPoint], point)
			}
		}
		// making sure the point has an entry in the map
		if _, ok := inDistance[point]; !ok {
			inDistance[point] = nil
		}
	}
	constellations := make(map[Point]int)
	var currentConstellation int

	for len(constellations) != len(points) {
		var currentPoint Point
		// first we take a point that's not currently part of a constellation
		for p := range inDistance {
			currentPoint = p
			break
		}

		// then we run BFS starting from this point
		toVisit := map[Point]bool{currentPoint: true}
		for len(toVisit) > 0 {
			for p := range toVisit {
				currentPoint = p
				break
			}

			// adding neighbours to visiting list if not already visited
			for _, p := range inDistance[currentPoint] {
				if _, ok := constellations[p]; !ok {
					toVisit[p] = true
				}
			}

			//fmt.Println("Adding", currentPoint, "to constellation", currentConstellation, len(constellations))
			constellations[currentPoint] = currentConstellation
			delete(toVisit, currentPoint)
			delete(inDistance, currentPoint)
		}

		// when we're done, we switch to a new constellation, from a new point
		currentConstellation++
	}
	return currentConstellation
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
