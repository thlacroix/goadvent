package main

import "fmt"

func main() {
	fmt.Println(countPasswords(146810, 612564, validatePasswordPart1))
	fmt.Println(countPasswords(146810, 612564, validatePasswordPart2))
}

func countPasswords(from, to int, validator func(i int) bool) int {
	var count int
	for i := from; i <= to; i++ {
		if validator(i) {
			count++
		}
	}
	return count
}

func validatePasswordPart1(password int) bool {
	var hasDouble bool
	numbers := getNumbers(password)
	previous := numbers[0]
	for _, n := range numbers[1:] {
		if n < previous {
			return false
		} else if n == previous {
			hasDouble = true
		}
		previous = n
	}
	return hasDouble
}

func validatePasswordPart2(password int) bool {
	numbers := getNumbers(password)
	var hasDouble bool
	groupCount := 1
	previous := numbers[0]
	for _, n := range numbers[1:] {
		if n < previous {
			return false
		}
		if n == previous {
			groupCount++
		} else {
			if groupCount == 2 {
				hasDouble = true
			}
			groupCount = 1
		}
		previous = n
	}
	return hasDouble || groupCount == 2
}

func getNumbers(i int) [6]byte {
	numbers := [6]byte{}
	inc := 1
	for j := 0; j < 6; j++ {
		numbers[5-j] = byte(i / inc % 10)
		inc = inc * 10
	}
	return numbers
}
