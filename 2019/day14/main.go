package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	chemicals, err := getChemicals("day14input.txt")
	if err != nil {
		log.Fatal(err)
	}
	ore := getOre(chemicals, 1)
	fmt.Println(ore)
	fuel := sort.Search(1000000000000, func(i int) bool {
		resetChemical(chemicals)
		return getOre(chemicals, i) > 1000000000000
	})
	fmt.Println(fuel - 1)
}

// Chemical represents how a chemical is produced, and holds
// the count of unused produced
type Chemical struct {
	Name      string
	Inputs    []ChemicalInput
	Output    int
	Remainder int
}

// ChemicalInput pairs a chemical and the number required
type ChemicalInput struct {
	Chemical *Chemical
	Input    int
}

// Getting ore count for get n fuels
func getOre(chemicals map[string]*Chemical, n int) int {
	return getOreForChemicalInput(ChemicalInput{chemicals["FUEL"], n}, 1)
}

// getting ore count to get the wanted input, multiplied by the number of
// chemicals from parents
func getOreForChemicalInput(chemicalInput ChemicalInput, mult int) int {
	// multiplying with the current input asked
	mult *= chemicalInput.Input

	// if we are on ore, we return
	if chemicalInput.Chemical.Name == "ORE" {
		return mult
	}
	var sum int

	// first we get how maybe remainders we can use in this reaction
	remainderToUse := min(chemicalInput.Chemical.Remainder, mult)
	// the we compute how much we need to roduce
	chemicalsToProduce := mult - remainderToUse

	// we deduce the number of reactions we need to do
	numberOfReactions := chemicalsToProduce / chemicalInput.Chemical.Output
	if chemicalsToProduce%chemicalInput.Chemical.Output != 0 {
		numberOfReactions++
	}
	// we adjust the remainder
	chemicalInput.Chemical.Remainder += numberOfReactions*chemicalInput.Chemical.Output - chemicalsToProduce - remainderToUse

	// and we recurse with the ingredients to get the ore count
	for _, in := range chemicalInput.Chemical.Inputs {
		sum += getOreForChemicalInput(in, numberOfReactions)
	}

	return sum
}

// reseting the chemicals remainders
func resetChemical(chemicals map[string]*Chemical) {
	for _, c := range chemicals {
		c.Remainder = 0
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

type ingredient struct {
	Name  string
	Count int
}

func parseChemical(in string) ([]ingredient, ingredient) {
	inOutSplit := strings.Split(in, " => ")
	var inIngredients []ingredient
	for _, i := range strings.Split(inOutSplit[0], ",") {
		inIngredients = append(inIngredients, parseIngredient(i))
	}
	return inIngredients, parseIngredient(inOutSplit[1])

}

func parseIngredient(in string) ingredient {
	split := strings.Split(strings.TrimSpace(in), " ")
	return ingredient{Name: split[1], Count: atoi(split[0])}
}

// unsafe string -> integer parsing
func atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}

// getChemicals reads the input and returns the chemicals
func getChemicals(fileName string) (map[string]*Chemical, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	chemicals := make(map[string]*Chemical)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ingredients, output := parseChemical(line)
		if len(ingredients) < 1 {
			return nil, errors.New("Can't parse instruction line " + line)
		}
		var chemical *Chemical
		if c, ok := chemicals[output.Name]; ok {
			chemical = c
		} else {
			chemical = &Chemical{Name: output.Name}
		}
		chemical.Output = output.Count

		for _, ingredient := range ingredients {
			if in, ok := chemicals[ingredient.Name]; ok {
				chemical.Inputs = append(chemical.Inputs, ChemicalInput{in, ingredient.Count})
			} else {
				in := &Chemical{Name: ingredient.Name}
				if in.Name == "ORE" {
					in.Output = 1
				}
				chemical.Inputs = append(chemical.Inputs, ChemicalInput{in, ingredient.Count})
				chemicals[in.Name] = in
			}
		}

		chemicals[chemical.Name] = chemical
	}
	return chemicals, nil
}
