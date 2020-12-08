package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

// Instruction represents a machine instruction: a command and a value
type Instruction struct {
	Cmd string
	V   int
}

// Machine is a list of instructions that can be run
type Machine []Instruction

// RunEndBeforeRepeat will run the instructions, and either stop before repeating itself,
// or when the instruction after last will be run (program terminated)
func (m Machine) RunEndBeforeRepeat() (int, bool, error) {
	var acc, i int
	seen := make(map[int]bool, len(m))

	for {
		if i == len(m) {
			return acc, true, nil
		}
		if i < 0 || i >= len(m) {
			return 0, false, fmt.Errorf("index %d outside of machine length (%d)", i, len(m))
		}
		if seen[i] {
			return acc, false, nil
		}
		seen[i] = true
		ins := m[i]

		switch ins.Cmd {
		case "acc":
			acc += ins.V
			i++
		case "jmp":
			i += ins.V
		case "nop":
			i++
		}
	}
}

func main() {
	var part1, part2 int
	instructions := make([]Instruction, 0, 1000)
	err := helpers.ScanLine("input.txt", func(s string) error {
		split := strings.Split(s, " ")
		if len(split) != 2 {
			return fmt.Errorf("can't split %s", s)
		}
		cmd := split[0]
		v, err := strconv.Atoi(split[1])
		if err != nil {
			return err
		}
		ins := Instruction{Cmd: cmd, V: v}
		instructions = append(instructions, ins)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	var done bool
	m := Machine(instructions)
	part1, done, err = m.RunEndBeforeRepeat()
	if err != nil {
		log.Fatal(err)
	}
	if done {
		log.Fatal("Part1 shouldn't terminate")
	}

	part2, done = searchBug(m)
	if !done {
		log.Fatal("Bug not found")
	}
	fmt.Println(part1, part2)
}

func searchBug(m Machine) (int, bool) {
	for i, ins := range m {
		var mc Machine
		if ins.Cmd == "nop" {
			mc = copyMachine(m)
			mc[i].Cmd = "jmp"

		} else if ins.Cmd == "jmp" {
			mc = copyMachine(m)
			mc[i].Cmd = "nop"
		}

		if mc != nil {
			acc, done, err := mc.RunEndBeforeRepeat()
			if done {
				return acc, true
			}
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return 0, false
}

func copyMachine(m Machine) Machine {
	m2 := make(Machine, len(m))
	copy(m2, m)
	return m2
}
