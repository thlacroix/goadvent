package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/helpers"
)

type command struct {
	a string
	v int
}

func main() {
	var commands []command

	err := helpers.ScanLine("input.txt", func(s string) error {
		var c command
		_, err := fmt.Sscanf(s, "%s %d", &c.a, &c.v)
		if err != nil {
			return err
		}

		commands = append(commands, c)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	hor, depth := move(commands)
	fmt.Println(hor * depth)
	hor2, depth2 := move2(commands)
	fmt.Println(hor2 * depth2)
}

func move(commands []command) (hor int, depth int) {
	for _, c := range commands {
		switch c.a {
		case "forward":
			hor += c.v
		case "down":
			depth += c.v
		case "up":
			depth -= c.v
		}
	}
	return
}

func move2(commands []command) (hor int, depth int) {
	var aim int
	for _, c := range commands {
		switch c.a {
		case "forward":
			hor += c.v
			depth += c.v * aim
		case "down":
			aim += c.v
		case "up":
			aim -= c.v
		}
	}
	return
}
