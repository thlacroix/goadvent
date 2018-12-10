package main

import (
	"fmt"
	"strings"
)

var metadataTotal int

func main() {
	fmt.Println("First result is", getHighestScore(477, 70851))
	fmt.Println("Second result is", getHighestScore(477, 70851*100))
}

// linked list that makes a circle
type Marble struct {
	Left  *Marble
	Right *Marble
	Value int
}

// just an helper to print the circle like in the example, not needed for the
// solution, but helps debugging
func (m Marble) String() string {
	var circle strings.Builder
	circle.WriteString(fmt.Sprintf("%d ", m.Value))
	currentMarble := m.Right
	for currentMarble.Value != m.Value {
		circle.WriteString(fmt.Sprintf("%d ", currentMarble.Value))
		currentMarble = currentMarble.Right
	}
	return circle.String()
}

func getHighestScore(playerCount, lastMarble int) int {
	// creating the initial marble
	initialMarble := &Marble{Value: 0}
	initialMarble.Left = initialMarble
	initialMarble.Right = initialMarble
	currentMarble := initialMarble
	// keeping the score for each player
	scores := make(map[int]int)
	for i := 1; i <= lastMarble; i++ {
		// getting current player id
		currentPlayer := (i-1)%playerCount + 1
		// special case, where we remove a marble and increase current player score
		if i != 0 && i%23 == 0 {
			for j := 0; j < 7; j++ {
				currentMarble = currentMarble.Left
			}
			currentMarble.Left.Right = currentMarble.Right
			currentMarble.Right.Left = currentMarble.Left
			scores[currentPlayer] += i + currentMarble.Value
			currentMarble = currentMarble.Right
		} else { // otherwise we just insert the new marble
			beforeNewMarble := currentMarble.Right
			afterNewMarble := beforeNewMarble.Right
			currentMarble = &Marble{Value: i, Left: beforeNewMarble, Right: afterNewMarble}
			beforeNewMarble.Right = currentMarble
			afterNewMarble.Left = currentMarble
		}
	}
	// returning max score
	var max int
	for _, score := range scores {
		if score > max {
			max = score
		}
	}
	return max
}
