package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var part1, part2 int

	var ints [][]bool
	err := helpers.ScanLine("input.txt", func(s string) error {
		var bits []bool
		for _, c := range s {
			if c == '0' {
				bits = append(bits, false)
			} else {
				bits = append(bits, true)
			}
		}
		ints = append(ints, bits)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	part1 = power(ints)
	part2 = life(ints)
	fmt.Println(part1, part2)
}

func power(ints [][]bool) int {
	count := make([]int, len(ints[0]))

	for _, i := range ints {
		for k, b := range i {
			if b {
				count[k]++
			}
		}
	}

	gammaB := make([]byte, 0, len(count))

	for _, c := range count {
		if c > len(ints)/2 {
			gammaB = append(gammaB, '1')
		} else {
			gammaB = append(gammaB, '0')
		}
	}
	epsB := make([]byte, 0, len(count))
	for _, b := range gammaB {
		if b == '1' {
			epsB = append(epsB, '0')
		} else {
			epsB = append(epsB, '1')
		}
	}
	gamma, _ := strconv.ParseInt(string(gammaB), 2, 64)
	eps, _ := strconv.ParseInt(string(epsB), 2, 64)

	return int(eps * gamma)
}

func life(ints [][]bool) int {
	ox := lifeFilter(ints, true)
	co2 := lifeFilter(ints, false)
	return ox * co2
}

func lifeFilter(ints [][]bool, ox bool) int {
	for k := range ints[0] {
		if len(ints) == 1 {
			break
		}

		m := most(ints, k) == ox

		intsF := make([][]bool, 0, len(ints)/2)

		for _, i := range ints {
			if i[k] == m {
				intsF = append(intsF, i)
			}
		}
		ints = intsF
	}
	return convert(ints[0])
}

func most(ints [][]bool, k int) bool {
	var c int

	for _, i := range ints {
		if i[k] {
			c++
		}
	}
	return c >= len(ints)-c
}

func convert(i []bool) int {
	resB := make([]byte, 0, len(i))

	for _, b := range i {
		if b {
			resB = append(resB, '1')
		} else {
			resB = append(resB, '0')
		}
	}

	res, _ := strconv.ParseInt(string(resB), 2, 64)
	return int(res)
}
