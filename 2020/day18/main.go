package main

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var part1, part2 int
	equations := make([]string, 0, 400)
	err := helpers.ScanLine("input.txt", func(s string) error {
		equations = append(equations, strings.ReplaceAll(s, " ", ""))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	part1 = sumResults(equations, solve)
	part2 = sumResults(equations, solve2)
	fmt.Println(part1, part2)
}

func sumResults(equations []string, f func([]rune, int) (int, int)) int {
	var sum int
	for _, e := range equations {
		v, _ := f([]rune(e), 0)
		sum += v
	}
	return sum
}

func solve(s []rune, i int) (int, int) {
	var value int

	var op rune

	for ; i < len(s); i++ {
		c := s[i]
		var newV int
		if unicode.IsDigit(c) {
			newV = int(c - '0')
		} else if c == ')' {
			break
		} else if c == '(' {
			newV, i = solve(s, i+1)
		} else {
			op = c
			continue
		}

		switch op {
		case 0:
			value = newV
		case '+':
			value += newV
		case '*':
			value *= newV
		}

	}
	return value, i
}

func solve2(s []rune, i int) (int, int) {
	var value int
	var sum int
	var mul bool

	var op rune

	doPreviousMul := func() {
		if sum > 0 {
			if mul {
				value *= sum
			} else {
				value = sum
			}
			sum = 0
		}
	}

	for ; i < len(s); i++ {
		c := s[i]
		var newV int
		if unicode.IsDigit(c) {
			newV = int(c - '0')
		} else if c == ')' {
			break
		} else if c == '(' {
			newV, i = solve2(s, i+1)
		} else {
			op = c
			continue
		}

		switch op {
		case 0:
			sum += newV
		case '+':
			sum += newV
		case '*':
			doPreviousMul()
			mul = true
			sum = newV
		}
	}
	doPreviousMul()
	return value, i
}
