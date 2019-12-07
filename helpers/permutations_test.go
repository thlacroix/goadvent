package helpers_test

import (
	"fmt"
	"testing"

	"github.com/thlacroix/goadvent/helpers"
)

type validation struct {
	valid   bool
	message string
}

func newValidation(status bool) validation {
	return validation{valid: status}
}

func newValidationWithMessage(status bool, message string) validation {
	return validation{valid: status, message: message}
}

func TestPermute(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	out := helpers.Permute(in)
	if v := validatePermutation(in, out); !v.valid {
		t.Error(v.message)
	}
}

func validatePermutation(in []int, out [][]int) validation {
	// we first check that we get n! permutations
	if len(out) != helpers.Factorial(len(in)) {
		return newValidationWithMessage(false, "Wrong number of permutations")
	}

	// for each permutation, we check that it gets the same elements as
	// the input, and that it doesn't appear twice
	for i, p := range out {
		if len(p) != len(in) {
			return newValidationWithMessage(false, fmt.Sprintf("%v does not has the same length as %v", in, p))
		}
		if !hasSameElements(in, p) {
			return newValidationWithMessage(false, fmt.Sprintf("%v does not has the same elements as %v", in, p))
		}
		for _, p2 := range out[i+1:] {
			if sliceEqual(p, p2) {
				return newValidationWithMessage(false, fmt.Sprintf("%v appears twice", p))
			}
		}
	}
	return newValidation(true)
}

func hasSameElements(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	aValues := make(map[int]int)
	for _, v := range a {
		aValues[v]++
	}
	for _, v := range b {
		aValues[v]--
	}
	for _, v := range aValues {
		if v != 0 {
			return false
		}
	}
	return true
}

func sliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

func TestFactorial(t *testing.T) {
	if res := helpers.Factorial(5); res != 120 {
		t.Errorf("Expected 120, go %d", res)
	}
}
