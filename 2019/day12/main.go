package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/thlacroix/goadvent/helpers"
)

var rInstruction = regexp.MustCompile(`<x=(-?\d+), y=(-?\d+), z=(-?\d+)>`)

func main() {
	moons, err := getMoons("day12input.txt")
	if err != nil {
		log.Fatal(err)
	}
	initialMoons := copyMoons(moons)

	fmt.Println(getEnergyAfterXSteps(copyMoons(moons), 2000))

	// getting history of the first 600000 steps
	// less steps is not enough to make sure we're in a loop
	history := getHistoryXSteps(moons, 600000)

	// If you want to see the moon pattern, you can uncomment the call
	// belown which will produce all the images in the CWD
	// But first, reduce the number of points above from 600000
	// to something between 200 and 500

	// plotMoons(history)

	frequencies, err := getFrequencies(history, initialMoons)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(getOverallFreq(frequencies))
}

// Coordinate of a moon
type Coordinate struct {
	X int
	Y int
	Z int
}

// CoordinateHistory holds the history of the coordinates of a moon
type CoordinateHistory struct {
	X []int
	Y []int
	Z []int
}

func (c Coordinate) String() string {
	return fmt.Sprintf("(%d, %d, %d)", c.X, c.Y, c.Z)
}

// AbsSum returns the sum of the absolute values of a coordinate
func (c Coordinate) AbsSum() int {
	return helpers.Abs(c.X) + helpers.Abs(c.Y) + helpers.Abs(c.Z)
}

// A Moon has a position coordinate and a velocity
type Moon struct {
	Position Coordinate
	Velocity Coordinate
}

func (m Moon) String() string {
	return fmt.Sprintf("pos=%v, vel=%v", m.Position, m.Velocity)
}

// Energy returns the energy of a moon
func (m Moon) Energy() int {
	return m.Position.AbsSum() * m.Velocity.AbsSum()
}

// MoonHistory has the history of the positions and velocities of a moon
type MoonHistory struct {
	PositionHistory CoordinateHistory
	VelocityHistory CoordinateHistory
}

// copyMoons creates a copy of the moons
func copyMoons(moons []*Moon) []*Moon {
	moonCopies := make([]*Moon, len(moons))
	for i, m := range moons {
		c := *m
		moonCopies[i] = &c
	}
	return moonCopies
}

// getMoons reads the input and returns the moons
func getMoons(fileName string) ([]*Moon, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var moons []*Moon
	for scanner.Scan() {
		line := scanner.Text()
		extracts := rInstruction.FindAllStringSubmatch(line, -1)
		if len(extracts) != 1 || len(extracts[0]) != 4 {
			return nil, errors.New("Can't parse instruction line " + line)
		}
		moon := &Moon{
			Position: Coordinate{
				X: atoi(extracts[0][1]),
				Y: atoi(extracts[0][2]),
				Z: atoi(extracts[0][3]),
			},
		}
		moons = append(moons, moon)
	}
	return moons, nil
}

// getEnergyAfterXSteps simulates the universe for X steps
func getEnergyAfterXSteps(moons []*Moon, steps int) int {
	for step := 1; step <= steps; step++ {
		for i, m1 := range moons {
			for _, m2 := range moons[i+1:] {
				updateVelocity(m1, m2)
			}
			applyVelocity(m1)
		}
	}
	return getEnergy(moons)
}

// getHistoryXSteps return the history of the universe after X steps
func getHistoryXSteps(moons []*Moon, steps int) []MoonHistory {
	moonHistory := make([]MoonHistory, len(moons))
	for step := 1; step <= steps; step++ {
		for i, m1 := range moons {
			for _, m2 := range moons[i+1:] {
				updateVelocity(m1, m2)
			}
			applyVelocity(m1)

			moonHistory[i].PositionHistory.X = append(moonHistory[i].PositionHistory.X, m1.Position.X)
			moonHistory[i].PositionHistory.Y = append(moonHistory[i].PositionHistory.Y, m1.Position.Y)
			moonHistory[i].PositionHistory.Z = append(moonHistory[i].PositionHistory.Z, m1.Position.Z)

			moonHistory[i].VelocityHistory.X = append(moonHistory[i].VelocityHistory.X, m1.Velocity.X)
			moonHistory[i].VelocityHistory.Y = append(moonHistory[i].VelocityHistory.Y, m1.Velocity.Y)
			moonHistory[i].VelocityHistory.Z = append(moonHistory[i].VelocityHistory.Z, m1.Velocity.Z)
		}
	}
	return moonHistory
}

// getEnergy returns the sum of the energy of the moons
func getEnergy(moons []*Moon) int {
	var energy int

	for _, m := range moons {
		energy += m.Energy()
	}

	return energy
}

// updateVelocity updates the velocity of two moons based on their relative position
func updateVelocity(m1, m2 *Moon) {
	if m1.Position.X < m2.Position.X {
		m1.Velocity.X++
		m2.Velocity.X--
	} else if m1.Position.X > m2.Position.X {
		m1.Velocity.X--
		m2.Velocity.X++
	}

	if m1.Position.Y < m2.Position.Y {
		m1.Velocity.Y++
		m2.Velocity.Y--
	} else if m1.Position.Y > m2.Position.Y {
		m1.Velocity.Y--
		m2.Velocity.Y++
	}

	if m1.Position.Z < m2.Position.Z {
		m1.Velocity.Z++
		m2.Velocity.Z--
	} else if m1.Position.Z > m2.Position.Z {
		m1.Velocity.Z--
		m2.Velocity.Z++
	}
}

// applyVelocity makes a moon move based on its velocity
func applyVelocity(m *Moon) {
	m.Position.X += m.Velocity.X
	m.Position.Y += m.Velocity.Y
	m.Position.Z += m.Velocity.Z
}

// getCycleToInitial search for the loop size of a moon returning to its initial position
func getCycleToInitial(positions, velocities []int, initialPosition, initialVelocity int) int {
	for i := range positions {
		if positions[i] == initialPosition && velocities[i] == initialVelocity {
			// once we've found a return to the initial position,
			// we make sure that we're really in a loop

			// in this case, we don't have enough data
			if 2*i+1 >= len(positions) {
				return -1
			}
			if positions[2*i+1] == initialPosition && velocities[2*i+1] == initialVelocity {
				return i + 1
			}
		}
	}
	return -1
}

// we get the cycle frequencies of each pair of position / velocity of each dimension for all moons
func getFrequencies(histories []MoonHistory, initialMoons []*Moon) ([]int, error) {
	frequencies := make([]int, 0, len(initialMoons)*3)

	for i, h := range histories {
		m := initialMoons[i]
		xFreq := getCycleToInitial(h.PositionHistory.X, h.VelocityHistory.X, m.Position.X, m.Velocity.X)
		if xFreq == -1 {
			return nil, fmt.Errorf("Can't find frequency for moon %d on axe X", i+1)
		}
		yFreq := getCycleToInitial(h.PositionHistory.Y, h.VelocityHistory.Y, m.Position.Y, m.Velocity.Y)
		if yFreq == -1 {
			return nil, fmt.Errorf("Can't find frequency for moon %d on axe Y", i+1)
		}
		zFreq := getCycleToInitial(h.PositionHistory.Z, h.VelocityHistory.Z, m.Position.Z, m.Velocity.Z)
		if zFreq == -1 {
			return nil, fmt.Errorf("Can't find frequency for moon %d on axe Z", i+1)
		}
		frequencies = append(frequencies, xFreq, yFreq, zFreq)
	}
	return frequencies, nil
}

// computing the LCM of all frequencies
func getOverallFreq(frequencies []int) int {
	total := 1
	for _, f := range frequencies {
		gcd := helpers.GCD(f, total)
		total = total * (f / gcd)
	}
	return total
}

// unsafe string -> integer parsing
func atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}

// print the moon histories for all moon, 3 files per moon (for each axe)
func plotMoons(history []MoonHistory) error {
	for i, m := range history {
		if err := plotMoon(m, i); err != nil {
			return err
		}
	}
	return nil
}

func plotMoon(moonH MoonHistory, index int) error {
	if err := plotMoonAxe(moonH.PositionHistory.X, moonH.VelocityHistory.X, index, "X"); err != nil {
		return err
	}
	if err := plotMoonAxe(moonH.PositionHistory.Y, moonH.VelocityHistory.Y, index, "Y"); err != nil {
		return err
	}
	return plotMoonAxe(moonH.PositionHistory.Z, moonH.VelocityHistory.Z, index, "Z")
}

func plotMoonAxe(positions, velocities []int, index int, axe string) error {
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = fmt.Sprintf("Moon %d on axe %s", index, axe)
	p.X.Label.Text = "Index"
	p.Y.Label.Text = "Value"

	err = plotutil.AddLinePoints(p,
		"Pos"+axe, pointPlots(positions),
		"Vel"+axe, pointPlots(velocities),
	)
	if err != nil {
		return err
	}

	return p.Save(20*vg.Inch, 10*vg.Inch, fmt.Sprintf("moon-%d-%s.png", index, axe))
}

func pointPlots(ints []int) plotter.XYs {
	pts := make(plotter.XYs, len(ints))
	for i, v := range ints {
		pts[i] = plotter.XY{X: float64(i), Y: float64(v)}
	}
	return pts
}
