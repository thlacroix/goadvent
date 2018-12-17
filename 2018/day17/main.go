package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var rClay = regexp.MustCompile(`(x|y)=(\d+), (x|y)=(\d+)..(\d+)`)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if initialMap, err := getMap(fileName); err != nil {
		log.Fatal(err)
	} else {
		count := getWaterFlow(initialMap, 500, 0)
		fmt.Println("Part1 result is", count)
		fmt.Println("Part2 result is", countResting(initialMap))
	}
}

type SquareType int

const (
	Sand SquareType = iota
	Clay
	WaterSpring
	WaterFalling
	WaterResting
)

func getMap(fileName string) ([][]SquareType, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// building defaut map
	initialMap := make([][]SquareType, 2000)
	for i := range initialMap {
		initialMap[i] = make([]SquareType, 600)
	}
	initialMap[0][500] = WaterSpring
	var minLeft, maxRight, minTop, maxBottom int
	for scanner.Scan() {
		line := scanner.Text()
		extracts := rClay.FindAllStringSubmatch(line, -1)
		if len(extracts) != 1 || len(extracts[0]) != 6 {
			return nil, errors.New("Can't parse line " + line)
		}
		extract := extracts[0]
		singleDirection := extract[1]
		singleValue := atoi(extract[2])

		lineDirection := extract[3]
		lineStart := atoi(extract[4])
		lineEnd := atoi(extract[5])
		if lineDirection == "y" && singleDirection == "x" {
			// seting mins and maxs
			if minLeft == 0 || singleValue < minLeft {
				minLeft = singleValue
			}
			if singleValue > maxRight {
				maxRight = singleValue
			}
			if minTop == 0 || lineStart < minTop {
				minTop = lineStart
			}
			if lineEnd > maxBottom {
				maxBottom = lineEnd
			}

			for i := lineStart; i <= lineEnd; i++ {
				initialMap[i][singleValue] = Clay
			}
		} else if lineDirection == "x" && singleDirection == "y" {
			// seting mins and maxs
			if minTop == 0 || singleValue < minTop {
				minTop = singleValue
			}
			if singleValue > maxBottom {
				maxBottom = singleValue
			}
			if minLeft == 0 || lineStart < minLeft {
				minLeft = lineStart
			}
			if lineEnd > maxRight {
				maxRight = lineEnd
			}

			for i := lineStart; i <= lineEnd; i++ {
				initialMap[singleValue][i] = Clay
			}
		}
	}
	// Cutting the map to keep only from toppest to bottomest, and removing after
	// rightest. Before leftest could also be removed easily, start source index
	// should just be moved in this case
	initialMap = initialMap[minTop : maxBottom+1]
	for i := range initialMap {
		initialMap[i] = initialMap[i][:maxRight+2]
	}
	return initialMap, nil
}

func printMap(initalMap [][]SquareType, left, right, top, bottom int) {
	for _, row := range initalMap[top : bottom+1] {
		var s strings.Builder
		for _, square := range row[left : right+1] {
			switch square {
			case Clay:
				s.WriteRune('#')
			case Sand:
				s.WriteRune('.')
			case WaterSpring:
				s.WriteRune('+')
			case WaterFalling:
				s.WriteRune('|')
			case WaterResting:
				s.WriteRune('~')
			}
		}
		fmt.Println(s.String())
	}
}

type Square struct {
	X int
	Y int
}

func getWaterFlow(initialMap [][]SquareType, waterSourceX, waterSourceY int) int {
	var waterCount int
	// going down
	var x, y int
	for y = waterSourceY; y < len(initialMap); y++ {
		if initialMap[y][waterSourceX] == Clay || initialMap[y][waterSourceX] == WaterResting {
			// if we hit clay or water, we'll fill with water above
			break
		} else if initialMap[y][waterSourceX] == WaterFalling {
			// if we hit water falling, no need to recompute
			return waterCount
		} else {
			// otherwise we go down
			initialMap[y][waterSourceX] = WaterFalling
			waterCount++
			if y == len(initialMap)-1 {
				// reaching bottom
				return waterCount
			}
		}
	}
	var overflowing bool
	var overflows []Square
	for !overflowing {
		// going up
		y--
		// going right
		var leftOverflow, rightOverflow int
		for x = waterSourceX; x < len(initialMap[0]); x++ {
			if initialMap[y][x] == Clay {
				// if we hit a wall, we stop filling
				rightOverflow = x - 1
				break
			} else {
				if initialMap[y+1][x] == WaterResting || initialMap[y+1][x] == Clay {
					// if we have clay or water below
					if initialMap[y][x] != WaterFalling && initialMap[y][x] != WaterResting {
						// increasing only if not already water
						waterCount++
					}
					initialMap[y][x] = WaterResting
				} else {
					// otherwise we overflow
					if initialMap[y][x] != WaterFalling {
						overflows = append(overflows, Square{X: x, Y: y})
					}
					overflowing = true
					rightOverflow = x - 1
					break
				}
			}
		}
		// going left
		for x = waterSourceX - 1; x > 0; x-- {
			if initialMap[y][x] == Clay {
				// if we hit a wall, we stop filling
				leftOverflow = x + 1
				break
			} else {
				if initialMap[y+1][x] == WaterResting || initialMap[y+1][x] == Clay {
					// if we have clay or water below
					if initialMap[y][x] != WaterFalling && initialMap[y][x] != WaterResting {
						// increasing only if not already water
						waterCount++
					}
					initialMap[y][x] = WaterResting
				} else {
					// otherwise we overflow
					if initialMap[y][x] != WaterFalling {
						overflows = append(overflows, Square{X: x, Y: y})
					}
					overflowing = true
					leftOverflow = x + 1
					break
				}
			}
		}
		// two things to do if overflowing: marking the top row as overflowed,
		// and recursively getting the overflow water count
		if overflowing {
			for x = leftOverflow; x <= rightOverflow; x++ {
				initialMap[y][x] = WaterFalling
			}
			for _, square := range overflows {
				waterCount += getWaterFlow(initialMap, square.X, square.Y)
			}
		}
	}
	return waterCount
}

// simply counting water resting
func countResting(initialMap [][]SquareType) int {
	var count int
	for _, row := range initialMap {
		for _, square := range row {
			if square == WaterResting {
				count++
			}
		}
	}
	return count
}

// unsafe string -> integer parsing
func atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}
