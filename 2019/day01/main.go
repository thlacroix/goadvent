package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/thlacroix/goadvent/helpers"
)

// initially done on a python interpreter on a phone in a bus,
// rewriting it quickly in go
func main() {
	ints, err := getInts("day01input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(getTotalFuel(ints, getFuel))
	fmt.Println(getTotalFuel(ints, getFuelRec))
}

func getInts(filename string) ([]int, error) {
	var ints []int
	return ints, helpers.ScanLine(filename, func(s string) error {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		ints = append(ints, i)
		return nil
	})
}

func getTotalFuel(ints []int, f func(int) int) int {
	var totalFuel int
	for _, i := range ints {
		totalFuel += f(i)
	}
	return totalFuel
}

func getFuel(i int) int {
	return i/3 - 2
}

func getFuelRec(i int) int {
	if i < 9 {
		return 0
	}
	fuel := getFuel(i)
	return fuel + getFuelRec(fuel)
}
