package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

type Recipe struct {
	Ingredients []string
	Allergens   []string
}

func main() {
	var part1 int
	var part2 string
	recipes := make([]Recipe, 0, 50)
	err := helpers.ScanLine("input.txt", func(s string) error {
		var r Recipe
		split := strings.Split(s, " (contains ")
		if len(split) != 2 {
			return fmt.Errorf("can't split %s", s)
		}
		r.Ingredients = strings.Split(split[0], " ")
		r.Allergens = strings.Split(strings.TrimSuffix(split[1], ")"), ", ")
		recipes = append(recipes, r)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	part1, part2 = findNoAllergens(recipes)
	fmt.Println(part1, part2)
}

func findNoAllergens(recipes []Recipe) (int, string) {
	allergenIngredients := make(map[string]map[string]int, len(recipes))
	allIngs := make(map[string]int, len(recipes))
	for _, r := range recipes {
		for _, i := range r.Ingredients {
			allIngs[i]++
		}
	}
	for _, r := range recipes {
		for _, a := range r.Allergens {
			var aIngs map[string]int
			var ok bool

			if aIngs, ok = allergenIngredients[a]; !ok {
				aIngs = map[string]int{a: 0}
				allergenIngredients[a] = aIngs
			}

			for _, i := range r.Ingredients {
				aIngs[i]++
			}
			aIngs[a]++
		}
	}

	possibleIngs := make(map[string]bool, len(allIngs))

	for a, aIngs := range allergenIngredients {
		for i, c := range aIngs {
			if i == a {
				continue
			}
			if c == aIngs[a] {
				possibleIngs[i] = true
			} else {
				delete(aIngs, i)
			}
		}
		delete(aIngs, a)
	}

	var c int

	for i, v := range allIngs {
		if !possibleIngs[i] {
			c += v
		}
	}

	found := make(map[string]string, len(allergenIngredients))

	for len(found) != len(allergenIngredients) {
		// random iteration order, could be improved by length order
	aLoop:
		for a, aIngs := range allergenIngredients {
			if aIngs == nil {
				continue
			}
			var possible string

			for i := range aIngs {
				if found[i] == "" {
					if possible == "" {
						possible = i
					} else {
						continue aLoop
					}
				}
			}
			allergenIngredients[a] = nil
			found[possible] = a
		}
	}

	mapSlice := make([][2]string, 0, len(found))

	for i, a := range found {
		mapSlice = append(mapSlice, [2]string{i, a})
	}

	sort.Slice(mapSlice, func(i, j int) bool {
		return mapSlice[i][1] < mapSlice[j][1]
	})

	var s strings.Builder

	for i, p := range mapSlice {
		if i != 0 {
			s.WriteRune(',')
		}
		s.WriteString(p[0])
	}

	return c, s.String()
}
