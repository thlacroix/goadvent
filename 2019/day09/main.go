package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/2019/intcode"
	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	ints, err := helpers.GetInts("day09input.txt")
	if err != nil {
		log.Fatal(err)
	}
	m := intcode.NewMachine(ints)
	go m.Run()
	m.AddInput(1)
	fmt.Println(m.GetOuput())
	m2 := intcode.NewMachine(ints)
	go m2.Run()
	m2.AddInput(2)
	fmt.Println(m2.GetOuput())
}
