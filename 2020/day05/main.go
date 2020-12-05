package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var seats [128 * 8]bool
	var min, max int
	err := helpers.ScanLine("input.txt", func(s string) error {
		id := getSeatIDBinary(s)
		if min == 0 || id < min {
			min = id
		}
		if id > max {
			max = id
		}

		seats[id] = true

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	for id, used := range seats[min+1 : max] {
		if !used {
			fmt.Println(max, min+1+id)
			return
		}
	}
}

// takes a seat string definition and returns the id with bit manipulatiohn
func getSeatIDBinary(s string) int {
	var id int
	for _, c := range s {
		id <<= 1
		if c == 'B' || c == 'R' {
			id++
		}
	}
	return id
}

// takes a seat string definition and returns the id with range calculation
func getSeatID(s string) int {
	rowMin, rowMax, columnMin, columnMax := 0, 127, 0, 7
	for _, c := range s {
		switch c {
		case 'F':
			rowMax = rowMax - (rowMax-rowMin+1)/2
		case 'B':
			rowMin = rowMin + (rowMax-rowMin+1)/2
		case 'L':
			columnMax = columnMax - (columnMax-columnMin+1)/2
		case 'R':
			columnMin = columnMin + (columnMax-columnMin+1)/2
		}
	}
	return rowMin*8 + columnMin
}
