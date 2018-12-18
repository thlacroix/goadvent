package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const size = 50
const minutes = 1000000000

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if initialMap, err := getMap(fileName); err != nil {
		log.Fatal(err)
	} else {
		res, _ := processMap(initialMap, 10)
		fmt.Println("Resource value for Part1 is", res)
		res2, frequency := processMap(initialMap, 2008)
		/*
			frequencyCountForTarget := (minutes - 2000) / frequency
			timeOnSameFrequency := minutes - frequency*frequencyCountForTarget
		*/
		fmt.Println("Resource value is for Part2", res2, "and frequency is", frequency)

	}
}

type SquareType int

const (
	OpenGround SquareType = iota
	Tree
	Lumberyard
)

func getMap(fileName string) ([][]SquareType, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// building defaut map
	initialMap := make([][]SquareType, size)
	for i := range initialMap {
		initialMap[i] = make([]SquareType, size)
	}
	var x, y int
	for scanner.Scan() {
		line := scanner.Text()
		x = 0
		for _, c := range line {
			switch c {
			case '.':
				initialMap[y][x] = OpenGround
			case '|':
				initialMap[y][x] = Tree
			case '#':
				initialMap[y][x] = Lumberyard
			}
			x++
		}
		y++
	}
	return initialMap, nil
}

func processMap(initialMap [][]SquareType, minutes int) (int, int) {
	var woods, lumberyards int
	var max, min int
	var frequency int
	currentMap := initialMap
	// for each minute, we build the new map
	for m := 1; m <= minutes; m++ {
		var newMap [][]SquareType
		woods = 0
		lumberyards = 0
		for y, row := range currentMap {
			var newRow []SquareType
			for x, square := range row {
				// we apply the rules for each square of the map
				switch square {
				case OpenGround:
					if countAdjacentType(currentMap, Tree, x, y) >= 3 {
						newRow = append(newRow, Tree)
						woods++
					} else {
						newRow = append(newRow, OpenGround)
					}
				case Tree:
					if countAdjacentType(currentMap, Lumberyard, x, y) >= 3 {
						newRow = append(newRow, Lumberyard)
						lumberyards++
					} else {
						newRow = append(newRow, Tree)
						woods++
					}
				case Lumberyard:
					if countAdjacentType(currentMap, Lumberyard, x, y) >= 1 && countAdjacentType(currentMap, Tree, x, y) >= 1 {
						newRow = append(newRow, Lumberyard)
						lumberyards++
					} else {
						newRow = append(newRow, OpenGround)
					}
				}
			}
			newMap = append(newMap, newRow)
		}
		currentMap = newMap
		res := woods * lumberyards

		// getting min and max to find frequency
		if m > 1000 && m < 1500 {
			if min == 0 || res < min {
				min = res
			}
			if res > max {
				max = res
			}
		}

		// getting frequency
		if m >= 1500 && frequency <= 0 && res == max {
			if frequency == 0 {
				frequency = -m
			} else if frequency < 0 {
				frequency = m + frequency
			}
		}
	}
	return woods * lumberyards, frequency
}

func countAdjacentType(initialMap [][]SquareType, square SquareType, x, y int) int {
	var count int
	// top left
	if x > 0 && y > 0 && initialMap[y-1][x-1] == square {
		count++
	}
	// left
	if x > 0 && initialMap[y][x-1] == square {
		count++
	}
	// bottom left
	if x > 0 && y < len(initialMap)-1 && initialMap[y+1][x-1] == square {
		count++
	}
	// bottom
	if y < len(initialMap)-1 && initialMap[y+1][x] == square {
		count++
	}
	// bottom right
	if x < len(initialMap[0])-1 && y < len(initialMap[0])-1 && initialMap[y+1][x+1] == square {
		count++
	}
	// right
	if x < len(initialMap)-1 && initialMap[y][x+1] == square {
		count++
	}
	// top right
	if y > 0 && x < len(initialMap)-1 && initialMap[y-1][x+1] == square {
		count++
	}
	// top
	if y > 0 && initialMap[y-1][x] == square {
		count++
	}
	return count
}

// helper to print the map
func printMap(initalMap [][]SquareType) {
	for _, row := range initalMap {
		var s strings.Builder
		for _, square := range row {
			switch square {
			case Lumberyard:
				s.WriteRune('#')
			case OpenGround:
				s.WriteRune('.')
			case Tree:
				s.WriteRune('|')
			}
		}
		fmt.Println(s.String())
	}
}
