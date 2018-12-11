package main

import (
	"fmt"
)

const serialNumber = 2187
const gridSize = 300

func main() {
	x3, y3, _ := getCoordinates(3)
	fmt.Println("The coordinates are", x3, y3, "for a square size of 3")
	x, y, maxSize := getCoordinates(-1)
	fmt.Println("The coordinates are", x, y, "and max size is", maxSize)

}

func getCoordinates(size int) (int, int, int) {
	// initializing the grid with default values
	grid := make([][]int, gridSize)
	for i := range grid {
		grid[i] = make([]int, gridSize)
	}

	// building the grid
	for i := 1; i <= gridSize; i++ {
		for j := 1; j <= gridSize; j++ {
			grid[i-1][j-1] = cellValue(j, i)
		}
	}

	// calculating the sliding sum
	var max, maxx, maxy, maxSize int
	if size > 0 {
		_, maxx, maxy = getMaxCoordinatesForSize(grid, size)
		maxSize = size
	} else {
		// instead of recalculating from scratch each time, we could rework how
		// the compute is done to make a sliding sum based on the square size
		for i := 1; i <= gridSize; i++ {
			maxForSize, maxxForSize, maxyForSize := getMaxCoordinatesForSize(grid, i)
			if i == 1 || maxForSize > max {
				max = maxForSize
				maxx = maxxForSize
				maxy = maxyForSize
				maxSize = i
			}
		}
	}
	return maxx, maxy, maxSize
}

func getMaxCoordinatesForSize(grid [][]int, size int) (int, int, int) {
	// calculating the sliding sum
	var max, maxx, maxy int
	var squareValue int

	// getting the initial value
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			squareValue += grid[i][j]
		}
	}
	max = squareValue
	// storing the sliding sum value of the first element of the row, that we'll
	// use to compute the first element of the next row
	lastRowInitialvalue := squareValue
	var previousValue int

	for i := 0; i < gridSize+1-size; i++ {
		for j := 0; j < gridSize+1-size; j++ {
			// first of the column, compute from above
			if j == 0 {
				if i == 0 {
					// special case, keep value already computed
					previousValue = lastRowInitialvalue
				} else {
					// we start from the value above, we remove old values and add the new ones
					previousValue = lastRowInitialvalue
					for k := 0; k < size; k++ {
						previousValue -= grid[i-1][j+k]
						previousValue += grid[i+size-1][j+k]
					}
					lastRowInitialvalue = previousValue
				}
			} else { // moving to right, removing left values, adding right values
				for k := 0; k < size; k++ {
					previousValue -= grid[i+k][j-1]
					previousValue += grid[i+k][j+size-1]
				}
			}

			// checking max
			if previousValue > max {
				max = previousValue
				maxx = j + 1
				maxy = i + 1
			}
		}
	}
	return max, maxx, maxy
}

// helper function to print the grid
func printGrid(grid [][]int) {
	for _, row := range grid {
		fmt.Println(row)
	}
}

// Computing value according to the rules
func cellValue(x, y int) int {
	rackID := x + 10
	powerLevel := rackID * y
	plusSerial := powerLevel + serialNumber
	timesRack := plusSerial * rackID
	var hundreds int
	if timesRack < 100 {
		hundreds = 0
	} else {
		hundreds = timesRack / 100 % 10
	}
	return hundreds - 5
}
