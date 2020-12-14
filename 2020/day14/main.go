package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

const size = 36

// Group a mask with its assignment instructions
type Group struct {
	Mask Mask
	Ins  []Assignment
}

// Assignment is a key value tuple
type Assignment struct {
	I int
	V int
}

// Mask represent an input mask with raw data in BitArray,
// and precomputed And and Or ints for part 1 and 2
type Mask struct {
	BitArray [size]byte
	And      int
	Or       int
}

// NewMask creates a Mask from its string definition,
// computing the And and Or ints
func NewMask(s string) Mask {
	var m Mask
	for i, c := range []byte(s) {
		m.BitArray[i] = c
		if c == '1' {
			m.Or += (1 << (size - 1 - i))
		}
		if c != '0' {
			m.And += (1 << (size - 1 - i))
		}
	}
	return m
}

// Apply the mask to an int for part 1
func (m Mask) Apply(i int) int {
	i = i | m.Or
	i = i & m.And
	return i
}

// Apply2 the mask to an int for part 2, generating a combination of ints
func (m Mask) Apply2(n int) []int {
	n = n | m.Or

	ns := []int{n}

	for i, b := range m.BitArray {
		if b == 'X' {
			newNs := make([]int, 0, len(ns)*2)
			for _, m := range ns {
				newNs = append(newNs, m&((1<<size-1)-1<<(size-1-i)))
				newNs = append(newNs, m|1<<(size-1-i))
			}
			ns = newNs
		}
	}

	return ns
}

var assignementR = regexp.MustCompile(`mem\[(\d+)\] = (\d+)`)

func main() {
	var part1, part2 int
	var currentGroup Group
	groups := make([]Group, 0, 20)
	err := helpers.ScanLine("input.txt", func(s string) error {
		if strings.HasPrefix(s, "mask") {
			mv := strings.Split(s, " = ")[1]

			m := NewMask(mv)

			if len(currentGroup.Ins) > 0 {
				groups = append(groups, currentGroup)
			}

			currentGroup = Group{Mask: m}
		} else {
			p := assignementR.FindStringSubmatch(s)
			if len(p) != 3 {
				return fmt.Errorf("Can't parse %s", s)
			}
			i, err := strconv.Atoi(p[1])
			if err != nil {
				return err
			}
			v, err := strconv.Atoi(p[2])
			if err != nil {
				return err
			}
			currentGroup.Ins = append(currentGroup.Ins, Assignment{I: i, V: v})
		}
		return nil
	})
	groups = append(groups, currentGroup)
	if err != nil {
		log.Fatal(err)
	}
	part1 = processGroups(groups)
	part2 = processGroups2(groups)
	fmt.Println(part1, part2)
}

func processGroups(groups []Group) int {
	mem := make(map[int]int, 1000)

	for _, g := range groups {
		for _, i := range g.Ins {
			v := g.Mask.Apply(i.V)
			mem[i.I] = v
		}
	}

	var sum int

	for _, v := range mem {
		sum += v
	}
	return sum
}

func processGroups2(groups []Group) int {
	mem := make(map[int]int, 1000)

	for _, g := range groups {
		for _, i := range g.Ins {
			is := g.Mask.Apply2(i.I)
			for _, k := range is {
				mem[k] = i.V
			}
		}
	}

	var sum int

	for _, v := range mem {
		sum += v
	}
	return sum
}

func a2i(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
