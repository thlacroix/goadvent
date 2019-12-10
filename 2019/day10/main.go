package main

import (
	"fmt"
	"log"
	"math"
	"sort"

	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	amap, asteroids, err := getMapAsteroid("day10input.txt")
	if err != nil {
		log.Fatal(err)
	}
	asteroid := getBestAsteroid(amap, asteroids)
	fmt.Println(asteroid.InSight)
	fmt.Println(get200thAsteroidShooted(asteroid, asteroids, amap))

}

// Asteroid represents an asteroid from the map, keeping track of the number
// of other asteroids in sight
// X and Y are reversed from the puzzle text
// Once they are exploded on part 2, we mark them
type Asteroid struct {
	InSight  int
	X        int
	Y        int
	Exploded bool
}

// we get the map, and keep track of the asteroids on the map
func getMapAsteroid(filename string) ([][]*Asteroid, []*Asteroid, error) {
	var amap [][]*Asteroid
	var asteroids []*Asteroid
	var x int

	err := helpers.ScanLine(filename, func(l string) error {
		var line []*Asteroid
		for y, c := range l {
			if c == '#' {
				asteorid := &Asteroid{X: x, Y: y}
				line = append(line, asteorid)
				asteroids = append(asteroids, asteorid)
			} else {
				line = append(line, nil)
			}
		}
		amap = append(amap, line)
		x++
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return amap, asteroids, nil
}

// For each asteroid, we search all the asteroids on sight
// If A is in sight of B, B is in sight of A, so no need to compute twice
// We then return the asteroid with the highest asteroids on sight
func getBestAsteroid(amap [][]*Asteroid, asteroids []*Asteroid) *Asteroid {
	for i, a1 := range asteroids {
		for _, a2 := range asteroids[i+1:] {
			if !hasPointInMiddle(a1, a2, amap) {
				a1.InSight++
				a2.InSight++
			}
		}
	}

	var maxInSight int
	var maxAsteroid *Asteroid

	for _, a := range asteroids {
		if a.InSight > maxInSight {
			maxInSight = a.InSight
			maxAsteroid = a
		}
	}
	return maxAsteroid
}

// Helper to find is an asteroid is in sight of another
func hasPointInMiddle(a1, a2 *Asteroid, amap [][]*Asteroid) bool {
	dx, dy := directionBetweenTwoPoints(a1, a2)
	x, y := a1.X+dx, a1.Y+dy
	for x != a2.X || y != a2.Y {
		if p := amap[x][y]; p != nil && !p.Exploded {
			return true
		}
		x += dx
		y += dy
	}
	return false
}

// Getting the best ratio of a vector to look for other points
func getSmallestDiagonale(a, b int) (int, int) {
	gcd := helpers.Abs(helpers.GCD(a, b))
	return a / gcd, b / gcd

}

// Getting the best vector ratio between two points
func directionBetweenTwoPoints(a, b *Asteroid) (int, int) {
	return getSmallestDiagonale(b.X-a.X, b.Y-a.Y)
}

// Gets all asteroids in sight of an asteroid, sorted by angle
func getInSightByRotatingOrder(a *Asteroid, asteroids []*Asteroid, amap [][]*Asteroid) []*Asteroid {
	var asteroidsInSight []*Asteroid
	for _, a2 := range asteroids {
		if a2 == a || a2.Exploded {
			continue
		}
		if !hasPointInMiddle(a, a2, amap) {
			asteroidsInSight = append(asteroidsInSight, a2)
		}
	}
	sort.Slice(asteroidsInSight, func(i, j int) bool {
		a1, a2 := asteroidsInSight[i], asteroidsInSight[j]
		return angle(a, a1) < angle(a, a2)
	})
	return asteroidsInSight
}

// Gets the angle from a vector between two points, and the up vector
func angle(center, a *Asteroid) float64 {
	res := math.Atan2(1, 0) - math.Atan2(float64(center.X-a.X), float64(a.Y-center.Y))
	if res < 0 {
		res += 2 * math.Pi
	}
	return res
}

// gets the 200th asteroid being exploded
func get200thAsteroidShooted(center *Asteroid, asteroids []*Asteroid, amap [][]*Asteroid) int {
	inSight := getInSightByRotatingOrder(center, asteroids, amap)
	var count int

	for count < 200 {
		for _, a := range inSight {
			a.Exploded = true
			count++
			if count == 200 {
				return a.Y*100 + a.X
			}
		}
		inSight = getInSightByRotatingOrder(center, asteroids, amap)
	}

	return 0
}
