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
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if largestSize, busySize, err := getSafestAreaSizes(fileName); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Largest area is", largestSize, "and busiest area size is", busySize)
	}
}

type Point struct {
	ID int
	X  int
	Y  int
}

type SpacePoint struct {
	Point   *Point
	Closest *Point
	Tie     bool
	Total   int
}

// returns the area sizes for Part1 and Part2
func getSafestAreaSizes(fileName string) (int, int, error) {
	var points []Point
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var i int
	var leftest, rightest, topest, bottomest int

	// first we build the list of points, and store the leftest, rightest, topest
	// and bottomest coordinates
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, ",")
		x, err := strconv.Atoi(strings.TrimSpace(split[0]))
		if err != nil {
			return 0, 0, err
		}
		if leftest == 0 || x < leftest {
			leftest = x
		}
		if x > rightest {
			rightest = x
		}
		y, err := strconv.Atoi(strings.TrimSpace(split[1]))
		if err != nil {
			return 0, 0, err
		}

		if topest == 0 || y < topest {
			topest = y
		}
		if y > bottomest {
			bottomest = y
		}

		point := Point{ID: i, X: x, Y: y}
		points = append(points, point)
		i++
	}

	// we get the size of areas for the points, the list of infinite points,
	// and the size of the region for Part2
	countMap, infinitePoints, regionSize := getMapCountForAreaSpace(points, leftest, rightest, topest, bottomest)

	// getting the maximum area that doesn't correspond to an infinite point
	var max int
	for point, count := range countMap {
		if count > max && !infinitePoints[point] {
			max = count
		}
	}

	return max, regionSize, nil
}

// builds the map of points in the reduced area (extended by one on each size),
// and returns the size of the closest area for each main point, and detects
// the infinite points by looking at the edges of the extended area
// Also returns solution for part 2
func getMapCountForAreaSpace(points []Point, leftest, rightest, topest, bottomest int) (map[int]int, map[int]bool, int) {
	const increase = 1 // extending the area to get the edge
	space := make([][]SpacePoint, bottomest+1+increase)
	for i := range space {
		space[i] = make([]SpacePoint, rightest+1+increase)
	}

	// setting the initial points on the map
	for i := range points {
		point := points[i]
		space[point.Y][point.X].Point = &point
		space[point.Y][point.X].Closest = &point
	}

	// for each points of the space, setting closest if only one
	for i := topest - increase; i <= bottomest+increase; i++ {
		for j := leftest - increase; j <= rightest+increase; j++ {
			spacePoint := space[i][j]
			min := -1
			var closest Point

			// for each space point, we calculate the distance with each initial points
			for _, point := range points {
				// getting the distance between point and space point
				dist := getDistance(point, i, j)
				spacePoint.Total += dist
				if dist == min {
					spacePoint.Tie = true // setting ties
				} else if min == -1 || dist < min {
					// unsetting tie and setting closest for the moment
					spacePoint.Tie = false
					min = dist
					closest = point
				}
			}

			if !spacePoint.Tie {
				spacePoint.Closest = &closest
			}
			space[i][j] = spacePoint
		}
	}

	countMap := make(map[int]int)
	var regionSize int
	infinitePoints := make(map[int]bool)

	// for each space point, we increment the count of the initial point areas,
	// but also mark them as infinites if they are on the edge of the extended
	// area
	for i := topest - increase; i <= bottomest+increase; i++ {
		for j := leftest - increase; j <= rightest+increase; j++ {
			spacePoint := space[i][j]
			if (i == topest-increase || i == bottomest+increase) &&
				(j == leftest-increase || i == rightest+increase) {
				infinitePoints[spacePoint.Closest.ID] = true
			} else {
				if !spacePoint.Tie && spacePoint.Closest != nil {
					countMap[spacePoint.Closest.ID]++
				}
				// incrementing size for Part2
				if spacePoint.Total < 10000 {
					regionSize++
				}
			}
		}
	}

	return countMap, infinitePoints, regionSize
}

// gives the distance between a Point and a map point
func getDistance(point Point, row, column int) int {
	return Abs(point.X-column) + Abs(point.Y-row)
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
