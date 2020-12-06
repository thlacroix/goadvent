package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var part1, part2 int
	err := helpers.ScanLine("input.txt", func(s string) error {
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(part1, part2)
}
