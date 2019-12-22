package main

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	actions, err := getActions("day22input.txt")
	if err != nil {
		log.Fatal(err)
	}
	// getting part 1 two ways:
	// * first by applying the full shuffle on the whole array
	// * then by only tracking tne index of the current card
	// It could also be solved in the same way as part 2, without the mod inverse
	cards := getCards(10007)
	fmt.Println(shuffle(cards, actions, 2019), "==", findCardPositionAfterShuffle(actions, 10007, 2019))

	// part 2
	fmt.Println(findBigValueAfterShuffle(actions, big.NewInt(119315717514047), big.NewInt(101741582076661), big.NewInt(2020)))
}

type ActionType byte

const (
	Reverse ActionType = iota
	Increment
	Cut
)

type Action struct {
	Type ActionType
	N    int
}

// part 1 by shuffling all the deck
func shuffle(cards []int, actions []Action, N int) int {
	for _, action := range actions {
		switch action.Type {
		case Reverse:
			cards = reverse(cards)
		case Cut:
			cards = cut(cards, action.N)
		case Increment:
			cards = increment(cards, action.N)
		}
	}
	return searchCard(cards, N)
}

// parsing input
func getActions(filename string) ([]Action, error) {
	var actions []Action
	err := helpers.ScanLine(filename, func(s string) error {
		if s == "deal into new stack" {
			actions = append(actions, Action{Type: Reverse})
		} else if strings.HasPrefix(s, "deal with increment") {
			actions = append(actions, Action{Type: Increment, N: atoi(strings.Split(s, "deal with increment ")[1])})
		} else if strings.HasPrefix(s, "cut") {
			actions = append(actions, Action{Type: Cut, N: atoi(strings.Split(s, "cut ")[1])})
		} else {
			return errors.New("Can parse action " + s)
		}
		return nil
	})
	return actions, err
}

// part 1 by tracking the index of the only card we want
func findCardPositionAfterShuffle(actions []Action, C, P int) int {
	for _, action := range actions {
		switch action.Type {
		case Reverse:
			P = (C - P - 1) % C
		case Cut:
			cn := action.N
			if cn < 0 {
				cn = C + action.N
			}
			P = (P - cn) % C
		case Increment:
			P = (P * action.N) % C
		}
	}
	if P < 0 {
		P = P + C
	}
	return P
}

// apllying a * b [mod] on big ints
func modularMultiplication(a, b, mod *big.Int) *big.Int {
	c := new(big.Int)
	c = c.Mul(a, b)
	return c.Mod(c, mod)
}

// composing (a**x+b) with (aa*x + bb) with mod
func compose(a, b, aa, bb, mod *big.Int) (*big.Int, *big.Int) {
	r := modularMultiplication(a, bb, mod)
	r = r.Add(r, b)
	r = r.Mod(r, mod)
	return modularMultiplication(a, aa, mod), r
}

// Part 2, by tracking the current index of a card, with big ints
// We're looking for a linear equation of the form a*X+b [C] = P [C]
// C is the number of cards (used as modulo)
// N is the number of times we apply the suffle
// P is the index in the deck after the suffle
// X is the value we're looking for
// First, we do the same thing as findCardPositionAfterShuffle, but only track a and b,
// by composing the linear equations
// Then we apply the equation after one suffle N times, by exponentiation by squaring
// At this point we have a*X+b [C] = P [C], to get X get solve the equation by multiplying
// the modulo inverse of a to (P - b)
func findBigValueAfterShuffle(actions []Action, C, N, P *big.Int) *big.Int {
	a := big.NewInt(1)
	b := big.NewInt(0)
	Cminus1 := new(big.Int).Sub(C, big.NewInt(1))
	for _, action := range actions {
		switch action.Type {
		case Reverse:
			a, b = compose(big.NewInt(-1), Cminus1, a, b, C)
		case Cut:
			cn := big.NewInt(int64(action.N))
			if cn.Cmp(big.NewInt(0)) == -1 {
				cn = cn.Add(cn, C)
			}
			cn = cn.Neg(cn)
			a, b = compose(big.NewInt(1), cn, a, b, C)
		case Increment:
			a, b = compose(big.NewInt(int64(action.N)), big.NewInt(0), a, b, C)
		}
	}
	a, b = applyNTimes(a, b, N, C)
	P = P.Sub(P, b)
	i := new(big.Int).ModInverse(a, C)
	P = modularMultiplication(P, i, C)

	return P
}

// exponentiation by squaring, iterative version
func applyNTimes(a, b, N, C *big.Int) (*big.Int, *big.Int) {
	aa, bb := big.NewInt(1), big.NewInt(0)
	for N.Cmp(big.NewInt(1)) > 0 {
		if mod := new(big.Int).Mod(N, big.NewInt(2)); mod.Cmp(big.NewInt(0)) == 0 {
			a, b = compose(a, b, a, b, C)
			N = N.Div(N, big.NewInt(2))
		} else {
			aa, bb = compose(a, b, aa, bb, C)
			a, b = compose(a, b, a, b, C)
			N = N.Sub(N, big.NewInt(1))
			N = N.Div(N, big.NewInt(2))
		}
	}
	return compose(aa, bb, a, b, C)
}

// generation a deck of N cards
func getCards(N int) []int {
	cards := make([]int, N)
	for i := range cards {
		cards[i] = i
	}
	return cards
}

// reversing the deck
func reverse(cards []int) []int {
	for i := 0; i < len(cards)/2; i++ {
		cards[i], cards[len(cards)-1-i] = cards[len(cards)-1-i], cards[i]
	}
	return cards
}

// cutting the deck with N cards (could be positive or negative)
func cut(cards []int, N int) []int {
	newCards := make([]int, len(cards))
	if N < 0 {
		N = len(cards) + N
	}
	copy(newCards, cards[N:])
	copy(newCards[len(cards)-N:], cards[:N])
	return newCards
}

// shuffling by increment
func increment(cards []int, N int) []int {
	newCards := make([]int, len(cards))
	for i, c := range cards {
		newCards[i*N%len(cards)] = c
	}
	return newCards
}

// searching the position of a card in the deck
func searchCard(cards []int, card int) int {
	for i, c := range cards {
		if c == card {
			return i
		}
	}
	return -1
}

// unsafe string -> integer parsing
func atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}
