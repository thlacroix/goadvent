package main

import (
	"container/list"
	"fmt"
	"log"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

type Crabs interface {
	Play()
	Stars() int
	fmt.Stringer
}

func main() {
	var part1 string
	var part2 int
	var c1, c2 Crabs
	err := helpers.ScanLine("input.txt", func(s string) error {
		c1 = NewCrabsL(s)
		c2 = NewCrabsLTo(s, 1000000)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	part1 = play(c1, 100)
	part2 = play2(c2, 10000000)
	fmt.Println(part1, part2)
}

func play(c Crabs, n int) string {
	for i := 0; i < n; i++ {
		c.Play()
	}

	return c.String()
}

func play2(c Crabs, n int) int {
	for i := 0; i < n; i++ {
		c.Play()
	}

	return c.Stars()
}

type CrabsL struct {
	list      *list.List
	current   *list.Element
	positions []*list.Element
}

func (c *CrabsL) Play() {
	// pick up 3 clockwise
	var nextCups [3]int
	for i := 0; i < len(nextCups); i++ {
		next := c.nextClockwise(c.current)
		nextCups[i] = next.Value.(int)
		c.list.Remove(next)
	}

	// find destination
	currentV := c.current.Value.(int)
	var destinationV int
destLoop:
	for destinationV == 0 {
		currentV = prev(currentV, c.list.Len()+3)

		for _, v := range nextCups {
			if v == currentV {
				continue destLoop
			}
		}
		destinationV = currentV
	}

	destination := c.positions[destinationV]

	// place cups
	for _, v := range nextCups {
		destination = c.list.InsertAfter(v, destination)
		c.positions[v] = destination
	}

	// select new current
	c.current = c.nextClockwise(c.current)
}

func (c CrabsL) String() string {
	one := c.positions[1]

	var s strings.Builder

	for e := c.nextClockwise(one); e.Value.(int) != 1; e = c.nextClockwise(e) {
		s.WriteString(fmt.Sprint(e.Value))
	}
	return s.String()
}

func (c CrabsL) Stars() int {
	one := c.positions[1]

	return c.nextClockwise(one).Value.(int) * c.nextClockwise(c.nextClockwise(one)).Value.(int)
}

func (c CrabsL) nextClockwise(e *list.Element) *list.Element {
	if e.Next() != nil {
		return e.Next()
	}
	return c.list.Front()
}

func NewCrabsL(s string) *CrabsL {
	var crabsL CrabsL
	crabsL.list = list.New()
	crabsL.positions = make([]*list.Element, len(s)+1)
	for _, c := range s {
		v := int(c - '0')
		crabsL.positions[v] = crabsL.list.PushBack(v)
	}
	crabsL.current = crabsL.list.Front()
	return &crabsL
}

func NewCrabsLTo(s string, n int) *CrabsL {
	var crabsL CrabsL
	crabsL.list = list.New()
	crabsL.positions = make([]*list.Element, n+1)
	for _, c := range s {
		v := int(c - '0')
		crabsL.positions[v] = crabsL.list.PushBack(v)
	}

	for i := len(s) + 1; i <= n; i++ {
		crabsL.positions[i] = crabsL.list.PushBack(i)
	}
	crabsL.current = crabsL.list.Front()
	return &crabsL
}

func prev(n, m int) int {
	n = n - 1
	if n == 0 {
		return m
	}
	return n
}
