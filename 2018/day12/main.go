package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const generations = 20
const part2Generations = 50000000000
const shift = 20

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if count, err := getPlantCount(fileName); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Plan count after", generations, "generations is", count)
		fmt.Println("Plan count after", part2Generations, "generations is", part2())

	}
}

type Note struct {
	Input  []bool
	Output bool
}

func parseInitial(initial string, shift int) []bool {
	// shifting the array with initial false values to get some room for expansion
	// on the left
	pots := make([]bool, shift)
	for _, c := range initial {
		if c == '#' {
			pots = append(pots, true)
		} else {
			pots = append(pots, false)
		}
	}
	return pots
}

func getPlantCount(fileName string) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// parsing the input to get the pot and note lists
	var lineIndex int
	var pots []bool
	var notes []Note
	for scanner.Scan() {
		line := scanner.Text()
		if lineIndex == 0 {
			// reading the first line
			setup := strings.Split(line, "initial state: ")[1]
			pots = parseInitial(setup, shift)
		} else if lineIndex > 1 {
			// reading the notes
			note := Note{}
			splitted := strings.Split(line, " => ")
			for _, c := range splitted[0] {
				if c == '#' {
					note.Input = append(note.Input, true)
				} else {
					note.Input = append(note.Input, false)
				}
			}
			if splitted[1] == "#" {
				note.Output = true
			} else {
				note.Output = false
			}
			notes = append(notes, note)
		}
		lineIndex++
	}

	// for each generation, we build the new list of pots
	for generation := 0; generation < generations; generation++ {
		newPots := make([]bool, len(pots))
		for i := range pots {
			// looping on each notes to see if one match
			for _, note := range notes {
				allGood := true
				for j, condition := range note.Input {
					var potValue bool
					// if the needed values are not in the list (too short), we know they're
					// false
					if i+j < len(pots) && pots[i+j] {
						potValue = true
					}
					if potValue != condition {
						allGood = false
						break
					}
				}
				if allGood && note.Output {
					// extending the list with 2 false elements to be able to write the new
					// one
					if i+2 >= len(newPots) {
						newPots = append(newPots, false, false)
					}
					newPots[i+2] = note.Output
					break
				}
			}
		}
		pots = newPots
	}

	// counting the sum by removing the index shift
	var sum int
	for i, pot := range pots {
		if pot {
			sum += i - shift
		}
	}
	return sum, nil
}

// Input size (generations) becomes quite high, so using compute for that is tedious
// After around bewteen 150 and 200 generations, the pattern of pots is the same
// (see below), shifted by generations - 36 empty pots, so output can easily be
// calculated in O(1) (considering that input size is number of generations)
func part2() int {
	pots := parseInitial("#.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.##.#...#.#", 0)
	// counting
	var sum int
	for i, pot := range pots {
		if pot {
			sum += 50000000000 - 36 + i
		}
	}
	return sum
}

// simple helper to visualise the pots
func printPots(pots []bool) {
	var s strings.Builder
	for _, pot := range pots {
		if pot {
			s.WriteRune('#')
		} else {
			s.WriteRune('.')
		}
	}
	fmt.Println(s.String())
}
