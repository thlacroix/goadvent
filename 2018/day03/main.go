package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var rClaim = regexp.MustCompile("^#(\\d+) @ (\\d+),(\\d+): (\\d+)x(\\d+)$")

const maxSize = 1000

type Claim struct {
	ID       int
	Position Position
	Size     Size
}

type Position struct {
	Left int
	Top  int
}

type Size struct {
	Length int
	Height int
}

type ClaimStatus int

const (
	NeverSeen ClaimStatus = iota
	SeenAloneOnly
	SeenWithOtherClaims
)

// building claim from input
func NewClaim(claimContent string) (Claim, error) {
	claim := Claim{}
	extracts := rClaim.FindAllStringSubmatch(claimContent, -1)
	if len(extracts) != 1 && len(extracts[0]) != 6 {
		return claim, errors.New("Can't parse claim")
	}
	extract := extracts[0]
	claim.ID = toInt(extract[1])
	claim.Position = Position{Left: toInt(extract[2]), Top: toInt(extract[3])}
	claim.Size = Size{Length: toInt(extract[4]), Height: toInt(extract[5])}
	return claim, nil
}

// herlper to inline string -> int conversion when we know we have an int
func toInt(in string) int {
	res, _ := strconv.Atoi(in)
	return res
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if table, err := getClaimedTable(fileName); err != nil {
		log.Fatal(err)
	} else {
		overlapCount := getOverlapCount(table)
		bestClaim := getBestClaim(table)
		fmt.Println("Overlap count is", overlapCount, "and best claim is", bestClaim)
	}
}

// building the claim table, with a 2D slice where values are the list of claims
// the inches
func getClaimedTable(fileName string) ([][][]int, error) {
	// bootstrapping the base table with defaults
	claimed := make([][][]int, maxSize)
	for i := range claimed {
		claimed[i] = make([][]int, maxSize)
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		claimContent := scanner.Text()
		claim, err := NewClaim(claimContent)
		if err != nil {
			return nil, err
		}
		for i := claim.Position.Top; i < claim.Position.Top+claim.Size.Height; i++ {
			for j := claim.Position.Left; j < claim.Position.Left+claim.Size.Length; j++ {
				claimed[i][j] = append(claimed[i][j], claim.ID)
			}
		}
	}
	return claimed, nil
}

// counting inches with more than one claim
func getOverlapCount(claimedTable [][][]int) int {
	var claimedCount int
	for _, row := range claimedTable {
		for _, inch := range row {
			if len(inch) >= 2 {
				claimedCount++
			}
		}
	}
	return claimedCount
}

// finding claim that doesn't overlap by looking at all inches and setting the
// status of the claims as we go along, and keeping the one that has always been
// seen alone
func getBestClaim(claimedTable [][][]int) int {
	statusMap := make(map[int]ClaimStatus)
	for _, row := range claimedTable {
		for _, inch := range row {
			inchClaimCount := len(inch)
			for _, id := range inch {
				// updating claim status
				if inchClaimCount != 1 {
					statusMap[id] = SeenWithOtherClaims
				} else if statusMap[id] != SeenWithOtherClaims {
					statusMap[id] = SeenAloneOnly
				}
			}
		}
	}

	// finding not overlaping claim
	for id, status := range statusMap {
		if status == 1 {
			return id
		}
	}
	return 0
}
