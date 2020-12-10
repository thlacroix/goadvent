package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/thlacroix/goadvent/helpers"
)

// cache used in part2 to store already computed combinations.
// The performance without are good enough (not really noticeable) but looks cleaner with,
// and might be needed for other inputs
var cache map[int]int

func main() {
	var part1, part2 int
	ints, err := helpers.GetIntsNL("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	sort.Ints(ints)

	ints = append(ints, ints[len(ints)-1]+3)
	cache = make(map[int]int, len(ints))

	// for part2 we're separating the input by groups before a +3 jump
	currentGroup := []int{0}

	var oneDiff, threeDiff, previous int
	part2 = 1
	for _, c := range ints {
		diff := c - previous
		if diff == 1 {
			oneDiff++
		} else if diff == 3 {
			threeDiff++
			// for each group we count the combinations inside it
			part2 *= countGroupCombinations(currentGroup)
			currentGroup = nil
		}
		currentGroup = append(currentGroup, c)
		previous = c
	}
	part1 = oneDiff * threeDiff

	fmt.Println(part1, part2)
}

// couting all possible combinations for a group recursively
func countGroupCombinations(group []int) int {
	if len(group) == 1 {
		return 1
	}

	// the target is the rightmost element, that might already be in the cache
	target := group[len(group)-1]
	if v, ok := cache[target]; ok {
		return v
	}

	var combinations int

	// we count the combinations that can reach the target
	for i := len(group) - 2; i >= 0 && target-group[i] <= 3; i-- {
		combinations += countGroupCombinations(group[:i+1])
	}
	cache[target] = combinations
	return combinations

}
