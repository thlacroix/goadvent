package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

type bagCount struct {
	name  string
	count int
}

const shinyGold = "shiny gold"

func main() {
	var part1, part2 int
	isContainedBy := make(map[string][]bagCount, 50)
	contains := make(map[string][]bagCount, 50)

	err := helpers.ScanLine("input.txt", func(s string) error {
		split := strings.Split(s, " contain ")
		if len(split) != 2 {
			return fmt.Errorf("'%s' can't be splitted", s)
		}
		name := strings.TrimSuffix(split[0], " bags")
		var bags []bagCount
		bagsSplit := strings.Split(strings.TrimSuffix(split[1], "."), ", ")
		for _, b := range bagsSplit {
			b = strings.TrimSuffix(b, "s")
			b = strings.TrimSuffix(b, " bag")
			bagSplit := strings.SplitN(b, " ", 2)
			if len(bagSplit) != 2 {
				return fmt.Errorf("'%s' can't be split", b)
			}
			if bagSplit[0] != "no" {
				c, err := strconv.Atoi(bagSplit[0])
				if err != nil {
					return err
				}
				bags = append(bags, bagCount{bagSplit[1], c})
			}
		}
		contains[name] = bags

		for _, b := range bags {
			isContainedBy[b.name] = append(isContainedBy[b.name], bagCount{name: name, count: b.count})
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// part 1 with BFS
	containsShiny := make(map[string]bool, 50)

	containingBags := []bagCount{{shinyGold, 0}}

	for len(containingBags) != 0 {
		b := containingBags[0]
		containingBags = containingBags[1:]
		for _, cb := range isContainedBy[b.name] {
			if !containsShiny[cb.name] {
				containsShiny[cb.name] = true
				containingBags = append(containingBags, cb)
			}
		}
	}

	part1 = len(containsShiny)

	// part 2 recursively
	var countBags func(b string) int
	countBags = func(b string) int {
		var count int
		for _, cb := range contains[b] {
			count += cb.count * (countBags(cb.name) + 1)
		}
		return count
	}

	part2 = countBags(shinyGold)

	fmt.Println(part1, part2)
}
