package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

// Rule is an structured rule from the input
type Rule struct {
	ID       int
	RL1, RL2 []int
	C        rune
}

// NewRule builds a rule from its string definition
func NewRule(s string) (Rule, error) {
	var r Rule
	var err error
	splitID := strings.Split(s, ": ")
	if len(splitID) != 2 {
		return r, fmt.Errorf("can't split id from s")
	}
	r.ID, err = strconv.Atoi(splitID[0])
	if err != nil {
		return r, err
	}

	if strings.Contains(splitID[1], `"`) {
		r.C = rune(splitID[1][1])
		return r, nil
	}

	splitOR := strings.Split(splitID[1], " | ")
	if len(splitOR) > 2 {
		return r, fmt.Errorf("tpp many or in %s", s)
	}

	r.RL1, err = NewRuleList(splitOR[0])
	if err != nil {
		return r, err
	}

	if len(splitOR) == 2 {
		r.RL2, err = NewRuleList(splitOR[1])
		if err != nil {
			return r, err
		}
	}

	return r, nil
}

// NewRuleList builds a list of space separated ints
func NewRuleList(s string) ([]int, error) {
	rp := make([]int, 0, 2)
	split := strings.Split(s, " ")

	for _, ss := range split {
		v, err := strconv.Atoi(ss)
		if err != nil {
			return nil, err
		}
		rp = append(rp, v)
	}

	return rp, nil
}

func a2i(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

var rules = make(map[int]Rule, 150)

func main() {
	var part1, part2 int

	messages := make([]string, 0, 450)

	var messagePart bool
	err := helpers.ScanLine("input.txt", func(s string) error {
		if s == "" {
			messagePart = true
			return nil
		}

		if !messagePart {
			r, err := NewRule(s)
			if err != nil {
				return err
			}
			rules[r.ID] = r
		} else {
			messages = append(messages, s)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	part1 = validateAll(messages, rules)
	rules[8] = Rule{ID: 8, RL1: []int{42}, RL2: []int{42, 8}}
	rules[11] = Rule{ID: 11, RL1: []int{42, 31}, RL2: []int{42, 11, 31}}
	// luckily part2 for me works without changing anything in the code
	// might not be the case for all inputs (recursive depth to be limited)
	part2 = validateAll(messages, rules)
	fmt.Println(part1, part2)
}

func validateAll(messages []string, rules map[int]Rule) int {
	var c int
	for _, m := range messages {
		if validate([]rune(m), 0, 0, nil, rules) {
			c++
		}
	}
	return c
}

func validate(s []rune, ruleID int, index int, next []int, rules map[int]Rule) bool {
	r := rules[ruleID]
	if r.C != 0 {
		if s[index] != r.C {
			return false
		}
		if len(next) == 0 {
			return index == len(s)-1
		} else if index+1 >= len(s) {
			return false
		} else {
			return validate(s, next[0], index+1, next[1:], rules)
		}
	}

	toAdd := r.RL1[1:]
	newNext := make([]int, len(next)+len(toAdd))
	copy(newNext, toAdd)
	copy(newNext[len(toAdd):], next)
	if validate(s, r.RL1[0], index, newNext, rules) {
		return true
	}

	if len(r.RL2) != 0 {
		toAdd := r.RL2[1:]
		newNext := make([]int, len(next)+len(toAdd))
		copy(newNext, toAdd)
		copy(newNext[len(toAdd):], next)
		if validate(s, r.RL2[0], index, newNext, rules) {
			return true
		}
	}

	return false
}
