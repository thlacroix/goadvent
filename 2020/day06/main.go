package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var count, count2, currentGroupLength int
	currentGroup := make(map[rune]int, 26)
	err := helpers.ScanLine("input.txt", func(s string) error {
		if s == "" {
			count += len(currentGroup)
			for _, c := range currentGroup {
				if c == currentGroupLength {
					count2++
				}
			}
			currentGroup = nil
			currentGroupLength = 0
			return nil
		}
		if currentGroup == nil {
			currentGroup = make(map[rune]int, 26)
		}
		for _, c := range s {
			currentGroup[c]++
		}
		currentGroupLength++
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if currentGroup != nil {
		count += len(currentGroup)
		for _, c := range currentGroup {
			if c == currentGroupLength {
				count2++
			}
		}
	}
	fmt.Println(count, count2)
}
