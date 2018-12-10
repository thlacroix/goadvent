package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filename passed")
	}
	totalFrequency, doubleFrequency, err := getFrequencies(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total frequence is", totalFrequency, "and double frequency is", doubleFrequency)

}

func getFrequencies(fileName string) (int, int, error) {
	var totalFrequency, frequency int
	fileContent, err := ioutil.ReadFile(fileName)
	lines := strings.Split(string(fileContent), "\n")
	if err != nil {
		return 0, 0, err
	}

	// parsing diffs properly and counting total for part 1
	diffs := make([]int, 0, len(lines)-1)
	seenFrequencies := make(map[int]bool)

	for _, line := range lines[:len(lines)-1] {
		// splitting here is probably not the optimal operation here, but gives
		// directly the sign and the value, so quite convenient
		if positiveSplit := strings.Split(line, "+"); len(positiveSplit) == 2 {
			diffInt, err := strconv.Atoi(positiveSplit[1])
			if err != nil {
				return 0, 0, err
			}
			diffs = append(diffs, diffInt)
			frequency = frequency + diffInt
			seenFrequencies[frequency] = true
		} else if negativeSplit := strings.Split(line, "-"); len(negativeSplit) == 2 {
			diffInt, err := strconv.Atoi(negativeSplit[1])
			if err != nil {
				return 0, 0, err
			}
			diffs = append(diffs, -diffInt)
			frequency = frequency - diffInt
			seenFrequencies[frequency] = true
		}
	}
	totalFrequency = frequency

	// now reiterating to find frequency appearing twice
	for {
		for _, diff := range diffs {
			frequency += diff
			if seenFrequencies[frequency] {
				return totalFrequency, frequency, nil
			}
			seenFrequencies[frequency] = true
		}
	}
}
