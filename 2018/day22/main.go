package main

import (
	"fmt"
	"strings"
)

const depth = 8103

const xTarget, yTarget = 9, 758

const additionalLinesX, additionalLinesY = 40, 10

type Region struct {
	Type          RegionType
	GeologicIndex int
	ErosionLevel  int
	X             int
	Y             int
}

func (r Region) String() string {
	var regionType string
	switch r.Type {
	case Wet:
		regionType = "wet"
	case Rocky:
		regionType = "rocky"
	case Narrow:
		regionType = "narrow"
	}
	return fmt.Sprintf("%d,%d (%s)", r.X, r.Y, regionType)
}

type RegionType int

const (
	Rocky RegionType = iota
	Wet
	Narrow
)

func main() {
	cave := getCave()
	printCave(cave)
	res := computeCaveRiskLevel(cave)
	fastest := getFastestWay(cave)
	fmt.Println(res, fastest)
}

func getCave() (cave [yTarget + additionalLinesY + 1][xTarget + additionalLinesX + 1]Region) {
	// top line
	for i := range cave[0] {
		region := getRegion(i * 16807)
		region.X = i
		cave[0][i] = region
	}

	// left column
	for i := range cave {
		region := getRegion(i * 48271)
		region.Y = i
		cave[i][0] = region
	}

	// rest
	for i := 1; i <= yTarget+additionalLinesY; i++ {
		for j := 1; j <= xTarget+additionalLinesX; j++ {
			var geologicIndex int
			if i == yTarget && j == xTarget {
				geologicIndex = 0
			} else {
				geologicIndex = cave[i][j-1].ErosionLevel * cave[i-1][j].ErosionLevel
			}

			region := getRegion(geologicIndex)
			region.X = j
			region.Y = i
			cave[i][j] = region
		}
	}
	return
}

func getRegion(geologicIndex int) Region {
	erosionLevel := (geologicIndex + depth) % 20183
	return Region{Type: RegionType(erosionLevel % 3), ErosionLevel: erosionLevel, GeologicIndex: geologicIndex}
}

func computeCaveRiskLevel(cave [yTarget + additionalLinesY + 1][xTarget + additionalLinesX + 1]Region) int {
	var risk int
	for _, row := range cave {
		for _, region := range row {
			risk += int(region.Type)
		}
	}
	return risk
}

type EquipmentType int

const (
	Torch EquipmentType = iota
	ClimbingGear
	Neither
)

type NodeStatus struct {
	Distance  uint64
	Equipment EquipmentType
}

type RegionGear struct {
	Region Region
	Gear   EquipmentType
}

func (r RegionGear) String() string {
	var gearName string
	switch r.Gear {
	case Torch:
		gearName = "torch"
	case ClimbingGear:
		gearName = "climbing gear"
	case Neither:
		gearName = "neither"
	}
	return fmt.Sprintf("%v | %s", r.Region, gearName)
}

func getFastestWay(cave [yTarget + additionalLinesY + 1][xTarget + additionalLinesX + 1]Region) int {
	target := RegionGear{Region: cave[yTarget][xTarget], Gear: Torch}
	regions := make(map[RegionGear]uint64)
	visited := make(map[RegionGear]bool)

	start := RegionGear{Region: cave[0][0], Gear: Torch}

	regions[start] = 0
	var currentRegion RegionGear

	for len(regions) > 0 {
		minDistance := -1
		var minRegion RegionGear
		for region, distance := range regions {
			if minDistance == -1 || int(distance) < minDistance {
				minDistance = int(distance)
				minRegion = region
			}
		}

		currentRegion = minRegion
		if currentRegion == target {
			return minDistance
		} else if currentRegion.Region == target.Region {
			currentTargetDistance, ok := regions[target]
			newTargetDistance := regions[currentRegion] + 1
			if !ok || (ok && newTargetDistance < currentTargetDistance) {
				regions[target] = newTargetDistance
			}
			return minDistance + 7
		} else {
			if currentRegion.Region.X < xTarget+additionalLinesX {
				updateNeighbourDistance(currentRegion, cave[currentRegion.Region.Y][currentRegion.Region.X+1], regions, visited)
			}
			if currentRegion.Region.Y < yTarget+additionalLinesY {
				updateNeighbourDistance(currentRegion, cave[currentRegion.Region.Y+1][currentRegion.Region.X], regions, visited)
			}
			if currentRegion.Region.X-1 >= 0 {
				updateNeighbourDistance(currentRegion, cave[currentRegion.Region.Y][currentRegion.Region.X-1], regions, visited)
			}
			if currentRegion.Region.Y-1 >= 0 {
				updateNeighbourDistance(currentRegion, cave[currentRegion.Region.Y-1][currentRegion.Region.X], regions, visited)
			}
		}

		visited[currentRegion] = true
		delete(regions, currentRegion)
	}
	return -1
}

func updateNeighbourDistance(originRegion RegionGear, neighbour Region, regions map[RegionGear]uint64, visited map[RegionGear]bool) {
	neighbourRegionGear := RegionGear{Region: neighbour}
	newNeighbourDistance := regions[originRegion] + 1
	switch originRegion.Gear {
	case Torch:
		if neighbour.Type == Wet {
			newNeighbourDistance += 7
			switch originRegion.Region.Type {
			case Rocky:
				neighbourRegionGear.Gear = ClimbingGear
			case Narrow:
				neighbourRegionGear.Gear = Neither
			}
		} else {
			neighbourRegionGear.Gear = Torch
		}
	case ClimbingGear:
		if neighbour.Type == Narrow {
			newNeighbourDistance += 7
			switch originRegion.Region.Type {
			case Rocky:
				neighbourRegionGear.Gear = Torch
			case Wet:
				neighbourRegionGear.Gear = Neither
			}
		} else {
			neighbourRegionGear.Gear = ClimbingGear
		}
	case Neither:
		if neighbour.Type == Rocky {
			newNeighbourDistance += 7
			switch originRegion.Region.Type {
			case Wet:
				neighbourRegionGear.Gear = ClimbingGear
			case Narrow:
				neighbourRegionGear.Gear = Torch
			}
		} else {
			neighbourRegionGear.Gear = Neither
		}
	}

	currentNeighbourDistance, ok := regions[neighbourRegionGear]

	if (!ok && !visited[neighbourRegionGear]) || (ok && newNeighbourDistance < currentNeighbourDistance) {
		regions[neighbourRegionGear] = newNeighbourDistance
	}
}

func printCave(cave [yTarget + additionalLinesY + 1][xTarget + additionalLinesX + 1]Region) {
	for _, row := range cave {
		var s strings.Builder
		for _, region := range row {
			if region.X == 0 && region.Y == 0 {
				s.WriteRune('M')
			} else if region.X == xTarget && region.Y == yTarget {
				s.WriteRune('T')
			} else {
				switch region.Type {
				case Wet:
					s.WriteRune('=')
				case Narrow:
					s.WriteRune('|')
				case Rocky:
					s.WriteRune('.')
				}
			}
		}
		fmt.Println(s.String())
	}
}
