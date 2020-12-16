package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

// Rule represent a role with its name and 2 validation ranges
type Rule struct {
	Name   string
	Range1 Range
	Range2 Range
}

// Contains returns true if i is contained in one of its ranges
func (r Rule) Contains(i int) bool {
	return r.Range1.Contains(i) || r.Range2.Contains(i)
}

var ruleR = regexp.MustCompile(`(.*): (\d+)-(\d+) or (\d+)-(\d+)`)

// NewRule parses the input string to create a new Rule
func NewRule(s string) (Rule, error) {
	var r Rule
	fields := ruleR.FindStringSubmatch(s)
	if len(fields) != 6 {
		return r, fmt.Errorf("Can't parse %s", s)
	}

	r.Name = fields[1]
	var err error
	r.Range1.From, err = strconv.Atoi(fields[2])
	if err != nil {
		return r, err
	}
	r.Range1.To, err = strconv.Atoi(fields[3])
	if err != nil {
		return r, err
	}
	r.Range2.From, err = strconv.Atoi(fields[4])
	if err != nil {
		return r, err
	}
	r.Range2.To, err = strconv.Atoi(fields[5])
	if err != nil {
		return r, err
	}

	return r, nil
}

// Range represents a from-to range
type Range struct {
	From int
	To   int
}

// Contains returns true if i is contained in the range
func (r Range) Contains(i int) bool {
	return i >= r.From && i <= r.To
}

func main() {
	var myTicket []int
	tickets := make([][]int, 0, 250)
	rules := make([]Rule, 0, 20)
	var zone byte
	var part1, part2 int
	err := helpers.ScanLine("input.txt", func(s string) error {
		if s == "your ticket:" {
			zone = 1
			return nil
		} else if s == "nearby tickets:" {
			zone = 2
			return nil
		} else if s == "" {
			return nil
		}

		var erri error

		switch zone {
		case 0:
			r, erri := NewRule(s)
			if erri != nil {
				return erri
			}
			rules = append(rules, r)
		case 1:
			myTicket, erri = NewTicket(s)
			if erri != nil {
				return erri
			}
		case 2:
			t, erri := NewTicket(s)
			if erri != nil {
				return erri
			}
			tickets = append(tickets, t)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	part1 = countErrorRate(rules, tickets)
	part2 = getDepartures(rules, myTicket, tickets)
	fmt.Println(part1, part2)
}

func countErrorRate(rules []Rule, tickets [][]int) int {
	var sum int
	for _, t := range tickets {
	vLoop:
		for _, v := range t {
			for _, r := range rules {
				if r.Contains(v) {
					continue vLoop
				}
			}
			sum += v
		}
	}
	return sum
}

func getDepartures(rules []Rule, myTicket []int, tickets [][]int) int {
	// for each field, we count how many tickets satisfy each rule
	fieldCount := make([]map[string]int, len(tickets[0]))
	for i := range fieldCount {
		fieldCount[i] = make(map[string]int, 20)
	}
	var validTickets int

ticketLoop:
	for _, t := range tickets {
		tmpFields := make([]map[string]bool, 0, len(t))
		for _, v := range t {
			validRules := make(map[string]bool, 20)
			for _, r := range rules {
				if r.Contains(v) {
					validRules[r.Name] = true
				}
			}

			// we keep only the tickets where all fields can be validated by at least one rule
			if len(validRules) == 0 {
				continue ticketLoop
			}
			tmpFields = append(tmpFields, validRules)
		}

		for i, f := range tmpFields {
			for vr := range f {
				fieldCount[i][vr]++
			}
		}
		validTickets++
	}

	// for each field, we keep only the rules that validate all tickets
	for i, f := range fieldCount {
		for rn, c := range f {
			if c != validTickets {
				// removing the entry from the map, but we could create a new one instead
				delete(f, rn)
			}
		}
		// keeping the field position with a special key in the map, as we'll sort after
		f["index"] = i
	}

	// sorting by the lower number of possible rules remaining
	sort.Slice(fieldCount, func(i, j int) bool {
		return len(fieldCount[i]) < len(fieldCount[j])
	})

	// keeping for each rule name the found index
	rulePosition := make(map[string]int, 20)

	// making sure that we continue until we find all rule index,
	// was not needed on my input as with the sort above only one loop on
	// fieldCount was necessary
	for len(rulePosition) != len(rules) {
		// for each field, we try to see if there's an only remaining possible rule,
		// if that's the case we save the result and "remove" the field by setting
		// its map to nil in the slice
	fieldLoop:
		for i, f := range fieldCount {
			if f == nil {
				continue
			}

			var only string
			for fn := range f {
				if fn == "index" {
					continue
				}

				if _, ok := rulePosition[fn]; !ok {
					if only == "" {
						only = fn
					} else {
						continue fieldLoop
					}
				}
			}
			rulePosition[only] = f["index"]
			fieldCount[i] = nil
		}
	}

	// multiplying all departure fields from my ticket
	mul := 1

	for rn, i := range rulePosition {
		if strings.HasPrefix(rn, "departure") {
			mul *= myTicket[i]
		}
	}
	return mul
}

// NewTicket parses the string input to return an int slice representing a ticket
func NewTicket(s string) ([]int, error) {
	split := strings.Split(strings.TrimSpace(string(s)), ",")
	ints := make([]int, 0, len(split))
	for _, c := range split {
		i, err := strconv.Atoi(c)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}
