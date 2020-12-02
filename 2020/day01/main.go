package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

const target = 2020

func main() {
	ints, err := helpers.GetIntsNL("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(twoSum(ints, target), threeSum(ints, target))
}

func twoSum(data []int, target int) int {
	lookup := make(map[int]bool, len(data))

	for _, v := range data {
		if lookup[v] {
			return v * (target - v)
		}
		lookup[target-v] = true
	}
	return 0
}

func threeSum(data []int, target int) int {
	lookup := make(map[int]bool, len(data))

	for i, x := range data {
		for _, y := range data[i+1:] {
			z := target - x - y
			if lookup[z] {
				return x * y * z
			}
		}
		lookup[x] = true
	}
	return 0
}
