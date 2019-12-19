package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/2019/intcode"
	"github.com/thlacroix/goadvent/helpers"
)

const squareSize = 100

func main() {
	ints, err := helpers.GetInts("day19input.txt")
	if err != nil {
		log.Fatal(err)
	}

	scans := getScans(ints, 50, 0, 0)
	fmt.Println(countScans(scans))

	x, y := getCoordFromScans(ints)

	fmt.Println(x*10000 + y)
}

// building the scan map from the machine
func getScans(ints []int, N int, fromx, fromy int) [][]bool {
	scans := make([][]bool, N)

	for y := 0; y < N; y++ {
		scans[y] = make([]bool, N)
		for x := 0; x < N; x++ {
			scans[y][x] = getValue(ints, x+fromx, y+fromy)
		}
	}
	return scans
}

func getCoordFromScans(ints []int) (int, int) {
	// first we find the first and last X on the 50 line
	currentY := 50
	var currentMinX, currentMaxX int
	for x := 0; x < currentY; x++ {
		v := getValue(ints, x, currentY)
		if v && currentMinX == 0 {
			currentMinX = x
		} else if currentMinX != 0 && !v {
			currentMaxX = x
			break
		}
	}
	// we keep an history of the max
	maxxs := []int{currentMaxX}

	// we're looking for a line where the min x is euqal to the max x of 99 lines before
	for len(maxxs)-squareSize < 0 || (maxxs[len(maxxs)-squareSize]-currentMinX) != squareSize-1 {
		currentY++
		var v bool
		// from the min and max, we move one line below, then move on the right
		// until we find the new min and max
		for {
			v = getValue(ints, currentMinX, currentY)
			if v {
				break
			}
			currentMinX++
		}

		for {
			v = getValue(ints, currentMaxX, currentY)
			if !v {
				currentMaxX--
				break
			}
			currentMaxX++
		}
		maxxs = append(maxxs, currentMaxX)

	}
	return maxxs[len(maxxs)-squareSize] - squareSize + 1, currentY - squareSize + 1
}

// creating a machine and calling it on a point to get if it's pulled
func getValue(ints []int, x, y int) bool {
	m := intcode.NewBufferedMachine(ints, 0, 2)
	go m.Run()
	m.AddInput(x)
	m.AddInput(y)
	pulled := m.GetOutput()
	_, end := m.GetOutputOrEnd()
	if !end {
		panic("Should have ended")
	}
	return pulled == 1
}

// helper to visiualize the map
func printScans(scans [][]bool) {
	for _, l := range scans {
		for _, v := range l {
			if v {
				fmt.Print(("#"))
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

// counts how many points are pulled
func countScans(scans [][]bool) int {
	var count int
	for _, l := range scans {
		for _, v := range l {
			if v {
				count++
			}
		}
	}
	return count
}
