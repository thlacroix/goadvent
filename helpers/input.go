package helpers

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// GetInts reads a file containing a comma separated list of ints
// a return a slice of ints
func GetInts(fileName string) ([]int, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	split := strings.Split(strings.TrimSpace(string(content)), ",")
	ints := make([]int, 0, len(split))
	for _, c := range split {
		i, err := strconv.Atoi(c)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

// GetIntsNL reads a file containing a newline separated list of ints
// a return a slice of ints
func GetIntsNL(fileName string) ([]int, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	split := strings.Split(strings.TrimSpace(string(content)), "\n")
	ints := make([]int, 0, len(split))
	for _, c := range split {
		i, err := strconv.Atoi(c)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}
