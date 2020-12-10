package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/thlacroix/goadvent/2019/intcode"
	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	ints, err := helpers.GetInts("day25input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(findPassword(ints))
}

func findPassword(ints []int) string {
	m := intcode.NewBufferedMachine(ints, 0, 100)
	go m.Run()

	end := make(chan bool)

	go func() {
		for {
			select {
			case <-m.Done:
				end <- true
			case v := <-m.Output:
				fmt.Printf("%c", v)
			}
		}
	}()

	instructions := []string{
		"north",
		"take festive hat",
		"west",
		"take sand",
		"east",
		"east",
		"take prime number",
		"west",
		"south",
		"east",
		"north",
		"take weather machine",
		"north",
		"take mug",
		"south",
		"south",
		"east",
		"north",
		"east",
		"east",
		"take astronaut ice cream",
		"west",
		"west",
		"south",
		"west",
		"west",
		"south",
		"south",
		"take mutex",
		"south",
		"take boulder",
		"east",
		"south",
		"east",
		"inv",
	}

	objects := []string{"boulder", "sand", "astronaut ice cream", "prime number", "festive hat", "mutex", "mug", "weather machine"}
	fmt.Println(objects)
	for _, i := range instructions {
		for _, c := range i {
			m.AddInput(int(c))
		}
		m.AddInput('\n')
	}

	// var f func([]string)

	// f = func(obj []string) {
	// 	for _, o := range obj {
	// 		f(obj[1:])
	// 		m.AddInput("take "+obj[0])
	// 	}
	// }

	// for _, o := range objects {
	// 	for _, b := range []bool{true, false} {
	// 		fmt.Println(o)
	// 	}
	// }

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		instruction := s.Text()
		switch instruction {
		case "n":
			instruction = "north"
		case "s":
			instruction = "south"
		case "e":
			instruction = "east"
		case "w":
			instruction = "west"
		}
		for _, c := range instruction {
			m.AddInput(int(c))
		}
		m.AddInput('\n')
	}

	return ""
}

// Helper that prints the output and the send an input
// Returns true if program ends
func printAndInput(m *intcode.Machine, in string) bool {
	for {
		c, input, end := m.GetOutputOrAddInputOrEnd(int(in[0]))
		if input {
			break
		} else if end {
			return true
		}
		fmt.Printf("%c", c)
	}
	for _, c := range in[1:] {
		m.AddInput(int(c))
	}
	m.AddInput('\n')
	return false
}
