package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filename passed")
	}
	checksum, codes, err := getChecksum(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	closest := getClosest(codes)
	fmt.Println("Checksum is", checksum, "and closest is", closest)
}

// calculate checksum, and also returns codes for next part
func getChecksum(fileName string) (int, []string, error) {
	var exactly2, exactly3 int
	file, err := os.Open(fileName)
	if err != nil {
		return 0, nil, err
	}
	defer file.Close()
	var codes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		countMap := make(map[rune]int)
		code := scanner.Text()
		codes = append(codes, code)

		// counting similar chars in code
		for _, c := range code {
			countMap[c]++
		}
		var exactly2Done, exactly3Done bool
		// increasing counts of codes with exactly 2 similar chars and exactly 3
		// similar chars
		for _, count := range countMap {
			if count == 2 && !exactly2Done {
				exactly2++
				exactly2Done = true
			}
			if count == 3 && !exactly3Done {
				exactly3++
				exactly3Done = true
			}
		}
	}
	return exactly2 * exactly3, codes, nil
}

// comparing each string with each others, gives a O(N^2) solution. Not ideal,
// but it works fast enough in this case
func getClosest(codes []string) string {
	var min int
	var closest string
	for i, code := range codes {
		for _, otherCode := range codes[i+1:] {
			diff, same := compareCodes(code, otherCode)
			if min == 0 || diff < min {
				min = diff
				closest = same
			}
		}
	}
	return closest
}

func compareCodes(code1, code2 string) (int, string) {
	var diff int
	var same []rune
	// we know that codes have the same length, and are simple letter, so we can
	// just compare the bytes
	for i, c := range code1 {
		if code1[i] == code2[i] {
			same = append(same, c)
		} else {
			diff++
		}
	}
	return diff, string(same)
}
