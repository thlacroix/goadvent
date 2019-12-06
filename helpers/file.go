package helpers

import (
	"bufio"
	"os"
)

// ScanLine takes a file name and calls f for each line of the file
func ScanLine(filename string, f func(string) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := f(scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}
