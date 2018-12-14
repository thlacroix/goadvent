package main

import (
	"fmt"
	"strconv"
	"strings"
)

const recipesAfter = 509671
const elfCount = 2

func main() {
	recipes := getRecipes(recipesAfter+10, nil)
	fmt.Println("Recipes after", recipesAfter, "are", getLastTen(recipes))
	pattern := getDigits(recipesAfter)
	patternRecipes := getRecipes(300000000, pattern)
	fmt.Println("There are", len(patternRecipes)-len(pattern), "recipes left to sequence")
}

func getRecipes(max int, pattern []int) []int {
	recipes := []int{3, 7}
	elfCurrentRecipeIndex := make([]int, elfCount)
	// initial indexes for elf recipes
	for i := range recipes {
		elfCurrentRecipeIndex[i] = i
	}

	for len(recipes) < max {
		// geting new recipes
		var currentRecipes []int
		for _, recipeIndex := range elfCurrentRecipeIndex {
			currentRecipes = append(currentRecipes, recipes[recipeIndex])
		}
		newRecipes := getNewRecipes(currentRecipes)
		recipes = append(recipes, newRecipes...)

		// checking pattern for Part2
		if pattern != nil && len(recipes) > len(pattern) {
			// as sometimes we add 1 and sometimes we add 2 recipes, checking for both
			// cases
			if recipesMatch(recipes, pattern) {
				return recipes
			} else if recipesMatch(recipes[:len(recipes)-1], pattern) {
				return recipes[:len(recipes)-1]
			}
		}

		// elves take new recipes
		for i, recipeIndex := range elfCurrentRecipeIndex {
			elfCurrentRecipeIndex[i] = (recipeIndex + recipes[recipeIndex] + 1) % len(recipes)
		}
	}

	return recipes
}

func getNewRecipes(currentRecipes []int) []int {
	var sum int
	for _, recipe := range currentRecipes {
		sum += recipe
	}
	return getDigits(sum)
}

func getDigits(d int) []int {
	if d < 10 {
		return []int{d}
	} else {
		return append(getDigits(d/10), d%10)
	}
}

func getLastTen(recipes []int) string {
	var s strings.Builder
	for _, recipe := range recipes[len(recipes)-10:] {
		s.WriteString(strconv.Itoa(recipe))
	}
	return s.String()
}

func recipesMatch(recipes, pattern []int) bool {
	subRecipes := recipes[len(recipes)-len(pattern):]
	for i, d := range pattern {
		if subRecipes[i] != d {
			return false
		}
	}
	return true
}
