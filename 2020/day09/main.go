package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var part1, part2 int
	ints, err := helpers.GetIntsNL("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	part1 = slidingTwoSUm(ints)
	part2 = slidingSum(ints, part1)
	fmt.Println(part1, part2)
}

func slidingTwoSUm(data []int) int {
	for i, c := range data {
		if i >= 25 {
			if !twoSum(data[i-25:i], c) {
				return c
			}
		}
	}
	return -1
}

func twoSum(data []int, target int) bool {
	lookup := make(map[int]bool, len(data))

	for _, v := range data {
		if lookup[v] {
			return true
		}
		lookup[target-v] = true
	}
	return false
}

func slidingSum(data []int, target int) int {
	var i, j, sum = 0, 1, 0
	sum = data[i] + data[j]

	for sum != target {
		if sum > target && i+1 < j {
			sum -= data[i]
			i++
		} else if j+1 < len(data) {
			j++
			sum += data[j]
		} else {
			return -1
		}
	}

	min, max := minMax(data[i : j+1])
	return min + max

}

func minMax(data []int) (int, int) {
	var min, max int

	for _, c := range data {
		if c > max {
			max = c
		}
		if min == 0 || c < min {
			min = c
		}
	}
	return min, max
}
