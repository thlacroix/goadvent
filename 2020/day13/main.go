package main

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var part1, part2 int
	var ts int
	var buses []int
	var erri error
	err := helpers.ScanLine("input.txt", func(s string) error {
		if ts == 0 {
			ts, erri = strconv.Atoi(s)
			if erri != nil {
				return erri
			}
		} else {
			split := strings.Split(s, ",")
			for _, b := range split {
				if b == "x" {
					buses = append(buses, -1)
				} else {
					n, erri := strconv.Atoi(b)
					if erri != nil {
						return erri
					}
					buses = append(buses, n)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	part1 = findBus(ts, buses)
	part2 = int(findStart(buses))
	fmt.Println(part1, part2)
}

func findBus(ts int, buses []int) int {
	min, minBus := -1, -1

	for _, b := range buses {
		if b == -1 {
			continue
		}

		afterLast := ts % b

		if afterLast == 0 {
			return 0
		}

		next := b - afterLast

		if min == -1 || next < min {
			min = next
			minBus = b
		}
	}
	return min * minBus
}

func findStart(buses []int) int64 {
	var entries []CRTInput

	for i, b := range buses {
		if b == -1 {
			continue
		}

		entries = append(entries, CRTInput{A: big.NewInt(-int64(i)), N: big.NewInt(int64(b))})
	}
	x := CRT(entries)
	return x.Int64()
}

// CRTInput holds a pair of input value for the equation: x ≡ A (mod N)
type CRTInput struct {
	A *big.Int
	N *big.Int
}

// CRT solved the uses the Chinese Remainder Theorem (https://en.wikipedia.org/wiki/Chinese_remainder_theorem)
// to solve the part 2 problem.
// Using big ints as part2 requires it
// Implementation based on https://fr.wikipedia.org/wiki/Th%C3%A9or%C3%A8me_des_restes_chinois#Algorithme
func CRT(inputs []CRTInput) *big.Int {
	n := big.NewInt(1)
	es := make([]*big.Int, 0, len(inputs))

	for _, in := range inputs {
		n.Mul(n, in.N)
	}

	for _, in := range inputs {
		ni := new(big.Int).Div(n, in.N)
		vi := new(big.Int).ModInverse(ni, in.N)
		ei := new(big.Int).Mul(vi, ni)
		es = append(es, ei)
	}

	x := big.NewInt(0)

	for i, in := range inputs {
		x.Add(x, new(big.Int).Mul(in.A, es[i]))
	}

	return x.Mod(x, n)
}
