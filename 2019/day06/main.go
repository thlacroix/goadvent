package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	objects, err := getObjects("day06input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(countOrbits(objects))
	fmt.Println(getTransfers(objects["YOU"], objects["SAN"]))
}

type Object struct {
	Name     string
	Orbit    *Object
	Distance int
}

// Getting the list of objects with their orbit in a map, to facilitate
// lookup
func getObjects(filename string) (map[string]*Object, error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	objects := make(map[string]*Object, len(lines))
	for _, l := range lines {
		split := strings.Split(strings.TrimSpace(l), ")")
		var origin *Object
		if oInMap, ok := objects[split[1]]; ok {
			origin = oInMap
		} else {
			origin = &Object{Name: split[1]}
		}
		if target, ok := objects[split[0]]; ok {
			origin.Orbit = target
		} else {
			target = &Object{Name: split[0]}
			objects[split[0]] = target
			origin.Orbit = target
		}
		objects[origin.Name] = origin

	}
	return objects, nil
}

// Summing the distance of each orbits
func countOrbits(objects map[string]*Object) int {
	var total int
	for _, o := range objects {
		total += getDistance(o)
	}
	return total
}

// Getting the distance of an object to COM, and store the result on the object
// we have it to avoid recomputing each time
func getDistance(object *Object) int {
	if object.Distance != 0 {
		return object.Distance
	}
	if object.Orbit == nil {
		return 0
	}
	distance := getDistance(object.Orbit) + 1
	object.Distance = distance
	return distance
}

// Simpler helper to hel visualize
func (o *Object) String() string {
	targetName := "nil"
	if o.Orbit != nil {
		targetName = fmt.Sprintf("%s (%d)", o.Orbit.Name, o.Orbit.Distance)
	}
	return fmt.Sprintf("%s (%d) -> %s", o.Name, o.Distance, targetName)
}

func getTransfers(you, santa *Object) int {
	// keeping a copy of the origin objects
	youC, santaC := you, santa

	// first we find the intersection between you and santa

	// starting by moving you to same distance as santa, as looking
	// at the data, you is further away
	for you.Distance != santa.Distance {
		you = you.Orbit
	}

	// Moving both at the same time until we find the same orbit
	for you.Orbit != santa.Orbit {
		you = you.Orbit
		santa = santa.Orbit
	}

	// if K is the distance at the intersection
	// if SAN distance = K + S
	// if YOU distance = K + Y
	// then SAN + YOU = 2K + S + Y
	// As we want S + Y, we sum YOU and SAM, and remove 2K
	return santaC.Distance + youC.Distance - 2*you.Distance
}
