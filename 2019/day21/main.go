package main

import (
	"fmt"
	"log"

	"github.com/thlacroix/goadvent/2019/intcode"
	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	ints, err := helpers.GetInts("day21input.txt")
	if err != nil {
		log.Fatal(err)
	}
	// we jump if there is a hole in front, of if there is a hole in
	// 3 steps and a platform at 4 steps
	sequences := []string{
		"NOT A J",
		"NOT C T",
		"AND D T",
		"OR T J",
		"WALK",
	}
	fmt.Println(jump(ints, sequences))

	// we jump if we see in hole in the next 3 steps and if there is a
	// platform at 4 steps, and if either there is a platfrom also at 5,
	// or we can jump again on a platform at 8
	runSequences := []string{
		"OR A T",
		"AND B T",
		"AND C T",
		"NOT T J",
		"AND D J",
		"OR E T",
		"OR H T",
		"AND T J",
		"RUN",
	}
	fmt.Println(jump(ints, runSequences))
}

// feeding the machine the input sequence, and returning the last output
func jump(ints []int, sequences []string) int {
	m := intcode.NewMachine(ints)
	go m.Run()

	for _, s := range sequences {
		printAndInput(m, s)
	}
	var last int
	for {
		c, end := m.GetOutputOrEnd()
		if end {
			return last
		}
		fmt.Printf("%c", last)
		last = c
	}
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
