package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	var passports []map[string]string
	var currentPassport map[string]string
	var valid, valid2 int
	err := helpers.ScanLine("input.txt", func(s string) error {
		if s == "" {
			passports = append(passports, currentPassport)
			if isValid(currentPassport) {
				valid++
				if isValid2(currentPassport) {
					valid2++
				}
			}
			currentPassport = nil
			return nil
		}
		if currentPassport == nil {
			currentPassport = make(map[string]string, 8)
		}
		split := strings.Split(s, " ")

		for _, ss := range split {
			field := strings.Split(ss, ":")
			if len(field) != 2 {
				return fmt.Errorf("wrong field %s", ss)
			}
			currentPassport[field[0]] = field[1]
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if currentPassport != nil {
		if isValid(currentPassport) {
			valid++
			if isValid2(currentPassport) {
				valid2++
			}
		}
	}
	fmt.Println(valid, valid2)
}

func isValid(passport map[string]string) bool {
	return (passport["cid"] == "" && len(passport) == 7) || (passport["cid"] != "" && len(passport) == 8)
}

func isValid2(passport map[string]string) bool {
	byr, err := strconv.Atoi(passport["byr"])
	if err != nil || byr < 1920 || byr > 2002 {
		return false
	}

	iyr, err := strconv.Atoi(passport["iyr"])
	if err != nil || iyr < 2010 || iyr > 2020 {
		return false
	}

	eyr, err := strconv.Atoi(passport["eyr"])
	if err != nil || eyr < 2020 || eyr > 2030 {
		return false
	}

	if !validateHeight(passport["hgt"]) {
		return false
	}

	if !validateHairColor(passport["hcl"]) {
		return false
	}

	if !validateEyeColor(passport["ecl"]) {
		return false
	}

	if !validatePID(passport["pid"]) {
		return false
	}

	return true
}

func validateHeight(v string) bool {
	if strings.HasSuffix(v, "cm") {
		v = strings.TrimSuffix(v, "cm")
		h, err := strconv.Atoi(v)
		return !(err != nil || h < 150 || h > 193)
	}
	if strings.HasSuffix(v, "in") {
		v = strings.TrimSuffix(v, "in")
		h, err := strconv.Atoi(v)
		return !(err != nil || h < 59 || h > 76)
	}
	return false
}

var colorValidator = regexp.MustCompile(`^#[a-f0-9]{6}$`)

func validateHairColor(s string) bool {
	return colorValidator.MatchString(s)
}

var colors = map[string]bool{"amb": true, "blu": true, "brn": true, "gry": true, "grn": true, "hzl": true, "oth": true}

func validateEyeColor(s string) bool {
	return colors[s]
}

func validatePID(s string) bool {
	if len(s) != 9 {
		return false
	}

	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
