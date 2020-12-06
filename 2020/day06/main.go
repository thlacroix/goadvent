package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var count, count2 int
	err := helpers.ScanGroup("input.txt", func(ss []string) error {
		yesCount := make(map[rune]int, 26)

		for _, s := range ss {
			for _, c := range s {
				yesCount[c]++
			}
		}

		count += len(yesCount)
		for _, c := range yesCount {
			if c == len(ss) {
				count2++
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count, count2)
}
