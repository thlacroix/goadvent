package main

import (
	"fmt"
	"log"
	"math"

	"github.com/thlacroix/goadvent/helpers"
)

const size = 5

const (
	bug   = '#'
	space = 0 // using 0 instead of '.' as it's the rune default value
	rec   = '?'
)

func main() {
	eris, err := getEris("day24input.txt")
	//eris, err := getEris("example")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(getBiodiversityRating(eris))
	eris[2][2] = '?'
	fmt.Println(countBugsAfterNMinutes(eris, 200))

}

// getting the raw maze as rune matrix
func getEris(filename string) ([size][size]rune, error) {
	var eris [size][size]rune
	var y int

	err := helpers.ScanLine(filename, func(s string) error {
		x := 0
		for _, c := range s {
			eris[y][x] = c
			x++
		}
		y++
		return nil
	})

	return eris, err
}

// gets the biodiversity rating of eris when we see twice the same situation
func getBiodiversityRating(eris [size][size]rune) int {
	seen := map[[size][size]rune]int{eris: 0}

	minute := 1
	for {
		var newEris [size][size]rune

		for y, l := range eris {
			for x, v := range l {
				// counting adjacent bugs
				var bugs int
				if x-1 >= 0 && eris[y][x-1] == bug {
					bugs++
				}
				if x+1 < size && eris[y][x+1] == bug {
					bugs++
				}
				if y-1 >= 0 && eris[y-1][x] == bug {
					bugs++
				}
				if y+1 < size && eris[y+1][x] == bug {
					bugs++
				}

				var newValue rune
				switch v {
				case bug:
					if bugs != 1 {
						newValue = space
					} else {
						newValue = bug
					}
				default:
					if bugs == 1 || bugs == 2 {
						newValue = bug
					} else {
						newValue = space
					}
				}
				newEris[y][x] = newValue

			}
		}
		if _, ok := seen[newEris]; ok {
			return countBiodiversity(newEris)
		}
		seen[newEris] = minute
		eris = newEris

		minute++
	}
}

// runs the simulation of the infinite recursion levels and returns the number
// of bugs seen in the levels after N minutes
func countBugsAfterNMinutes(eris [size][size]rune, minutes int) int {
	levels := map[int][size][size]rune{0: eris}
	var newLevel [size][size]rune
	newLevel[2][2] = rec

	minute := 1
	for minute <= minutes {
		newLevels := make(map[int][size][size]rune, len(levels))
		impactedLevels := make(map[int]bool)

		f := func(i int, level [size][size]rune) {
			var newEris [size][size]rune

			for y, l := range level {
				for x, v := range l {
					if v == rec {
						newEris[y][x] = rec
						continue
					}
					var bugs int

					// Left
					if x-1 >= 0 {
						switch level[y][x-1] {
						case bug:
							bugs++
						case rec:
							innerLevel, ok := levels[i+1]
							if ok {
								for y := range innerLevel {
									if innerLevel[y][size-1] == bug {
										bugs++
									}
								}
							} else if v == bug {
								impactedLevels[i+1] = true
							}
						}
					} else {
						outerLevel, ok := levels[i-1]
						if ok {
							if outerLevel[size/2][size/2-1] == bug {
								bugs++
							}
						} else if v == bug {
							impactedLevels[i-1] = true
						}
					}

					// Right
					if x+1 < size {
						switch level[y][x+1] {
						case bug:
							bugs++
						case rec:
							innerLevel, ok := levels[i+1]
							if ok {
								for y := range innerLevel {
									if innerLevel[y][0] == bug {
										bugs++
									}
								}
							} else if v == bug {
								impactedLevels[i+1] = true
							}
						}
					} else {
						outerLevel, ok := levels[i-1]
						if ok {
							if outerLevel[size/2][size/2+1] == bug {
								bugs++
							}
						} else if v == bug {
							impactedLevels[i-1] = true
						}
					}

					// Top
					if y-1 >= 0 {
						switch level[y-1][x] {
						case bug:
							bugs++
						case rec:
							innerLevel, ok := levels[i+1]
							if ok {
								for x := range innerLevel {
									if innerLevel[size-1][x] == bug {
										bugs++
									}
								}
							} else if v == bug {
								impactedLevels[i+1] = true
							}
						}
					} else {
						outerLevel, ok := levels[i-1]
						if ok {
							if outerLevel[size/2-1][size/2] == bug {
								bugs++
							}
						} else if v == bug {
							impactedLevels[i-1] = true
						}
					}

					// Right
					if y+1 < size {
						switch level[y+1][x] {
						case bug:
							bugs++
						case rec:
							innerLevel, ok := levels[i+1]
							if ok {
								for x := range innerLevel {
									if innerLevel[0][x] == bug {
										bugs++
									}
								}
							} else if v == bug {
								impactedLevels[i+1] = true
							}
						}
					} else {
						outerLevel, ok := levels[i-1]
						if ok {
							if outerLevel[size/2+1][size/2] == bug {
								bugs++
							}
						} else if v == bug {
							impactedLevels[i-1] = true
						}
					}

					var newValue rune
					switch v {
					case bug:
						if bugs != 1 {
							newValue = space
						} else {
							newValue = bug
						}
					default:
						if bugs == 1 || bugs == 2 {
							newValue = bug
						} else {
							newValue = space
						}
					}
					newEris[y][x] = newValue
				}
			}
			newLevels[i] = newEris
		}

		// running the main function for all levels where we currently have bugs
		for i, level := range levels {
			f(i, level)
		}

		// running the main function on new levels where we now we'll get bugs
		// from the previous loop
		for i := range impactedLevels {
			f(i, newLevel)
		}

		levels = newLevels
		minute++
	}

	// counting the bugs
	var count int
	for _, level := range levels {
		for _, l := range level {
			for _, v := range l {
				if v == bug {
					count++
				}
			}
		}
	}
	return count
}

// getting the biodiversity rating of eris at a specific time
func countBiodiversity(eris [size][size]rune) int {
	var sum int
	for y, l := range eris {
		for x, v := range l {
			if v == bug {
				n := y*size + x
				sum += int(math.Pow(2, float64(n)))
			}

		}
	}
	return sum
}

// helper to print the eris
func printEris(eris [size][size]rune) {
	for _, l := range eris {
		for _, v := range l {
			if v == 0 {
				v = '.'
			}
			fmt.Printf("%c", v)
		}
		fmt.Println()
	}
}
