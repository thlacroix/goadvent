package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/scanner"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if length, moreThan100, err := getDirections(fileName); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Max shortest distance is", length, "with", moreThan100, "rooms at least 1000 doors away")
	}
}

// A room is just a point
type Room struct {
	X int
	Y int
}

// Status to return during recursive parsing
type Status int

const (
	OtherGroups Status = iota
	EndGroup
	EndRegexp
	Error
)

// Getting the file a a char scanner, and recursively parsing
func getDirections(fileName string) (int, int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	scan := bufio.NewScanner(file)
	scan.Scan()
	var textScanner scanner.Scanner
	textScanner.Init(strings.NewReader(scan.Text()))

	// storing room distances in a map
	distances := make(map[Room]int)

	// starting a a 0 distance and central room
	status := getDirectionsR(&textScanner, 0, Room{}, distances)
	var maxDistance, moreThan1000 int
	for _, distance := range distances {
		if distance > maxDistance {
			maxDistance = distance
		}
		if distance >= 1000 {
			moreThan1000++
		}
	}
	if status != EndRegexp {
		return 0, 0, errors.New("Missed the end of the regexp")
	}
	return maxDistance, moreThan1000, nil
}

func getDirectionsR(r *scanner.Scanner, currentDistance int, currentRoom Room, distances map[Room]int) Status {
	for c := r.Next(); c != scanner.EOF; c = r.Next() {
		switch c {
		case '^': // nothing to do
		case '$':
			return EndRegexp
		case 'E', 'W', 'N', 'S':
			currentDistance++
			// moving to other rooms
			switch c {
			case 'E':
				currentRoom = Room{X: currentRoom.X + 1, Y: currentRoom.Y}
			case 'W':
				currentRoom = Room{X: currentRoom.X - 1, Y: currentRoom.Y}
			case 'N':
				currentRoom = Room{X: currentRoom.X, Y: currentRoom.Y - 1}
			case 'S':
				currentRoom = Room{X: currentRoom.X, Y: currentRoom.Y + 1}
			}
			// keeping shortest distance for the room
			if distance, ok := distances[currentRoom]; !ok || currentDistance < distance {
				distances[currentRoom] = currentDistance
			}
		case '(':
			var status Status
			// iterating on each group
			for status != EndGroup {
				status = getDirectionsR(r, currentDistance, currentRoom, distances)
				if status == Error {
					return Error
				}
			}
		case '|':
			return OtherGroups
		case ')':
			return EndGroup
		default:
			return Error
		}
	}
	return Error
}
