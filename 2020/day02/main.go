package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var count1, count2 int
	err := helpers.ScanLine("input.txt", func(s string) error {
		v1, v2, errv := valid(s)
		if errv != nil {
			return errv
		}
		if v1 {
			count1++
		}
		if v2 {
			count2++
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count1, count2)
}

// Range holds a min max range
type Range struct {
	Min int
	Max int
}

// NewRange splits a string like "2-7" into a min max range
func NewRange(r string) (Range, error) {
	split := strings.Split(r, "-")
	var res Range
	var err error
	if len(split) != 2 {
		return res, fmt.Errorf("range split of length %d", len(split))
	}

	res.Min, err = strconv.Atoi(split[0])
	if err != nil {
		return res, err
	}

	res.Max, err = strconv.Atoi(split[1])
	if err != nil {
		return res, err
	}

	return res, nil

}

// takes an input line and validate the password for part 1 and part 2
func valid(s string) (bool, bool, error) {
	split := strings.Split(s, " ")
	if len(split) != 3 {
		return false, false, errors.New("can't split input")
	}

	r, err := NewRange(split[0])
	if err != nil {
		return false, false, err
	}
	if len(split[1]) != 2 {
		return false, false, errors.New("too many letter bytes")
	}
	letter := split[1][0]
	letterRune := rune(letter)

	pass := split[2]

	if r.Min < 1 || r.Max > len(pass) {
		return false, false, fmt.Errorf("range %v does not fit password length %d", r, len(pass))
	}

	var count int

	for _, c := range pass {
		if c == letterRune {
			count++
		}
	}
	valid1 := r.Min <= count && count <= r.Max
	valid2 := (pass[r.Min-1] == letter && pass[r.Max-1] != letter) || (pass[r.Min-1] != letter && pass[r.Max-1] == letter)
	return valid1, valid2, nil

}
