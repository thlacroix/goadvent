package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const attackPower = 3
const initialHealth = 200

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if squareMap, elves, goblins, err := getMap(fileName); err != nil {
		log.Fatal(err)
	} else {
		turns, totalHealth := fight(squareMap, elves, goblins, attackPower, false)
		fmt.Println("Outcome is", turns*totalHealth)
	}

	// Part 2
	for force := 15; force < 50; force++ {
		// reading each time, as we modify everything as we go
		if squareMap, elves, goblins, err := getMap(fileName); err != nil {
			log.Fatal(err)
		} else {
			turns, totalHealth := fight(squareMap, elves, goblins, force, true)
			if turns != 0 && totalHealth != 0 {
				fmt.Println("Win outcome is", turns*totalHealth)
				break
			}
		}
	}
}

type PersoType int

const (
	Elf PersoType = iota
	Goblin
)

type Perso struct {
	CurrentSquare *Square
	Type          PersoType
	Health        int
}

// lookking for ennemy already adjacent wil lesser health
func (p *Perso) HasEnnemyInRange(squareMap [][]*Square) *Perso {
	var target *Perso
	if perso := squareMap[p.CurrentSquare.Y-1][p.CurrentSquare.X].Perso; perso != nil && perso.Type != p.Type {
		if target == nil || perso.Health < target.Health {
			target = perso
		}
	}
	if perso := squareMap[p.CurrentSquare.Y][p.CurrentSquare.X-1].Perso; perso != nil && perso.Type != p.Type {
		if target == nil || perso.Health < target.Health {
			target = perso
		}
	}
	if perso := squareMap[p.CurrentSquare.Y][p.CurrentSquare.X+1].Perso; perso != nil && perso.Type != p.Type {
		if target == nil || perso.Health < target.Health {
			target = perso
		}
	}
	if perso := squareMap[p.CurrentSquare.Y+1][p.CurrentSquare.X].Perso; perso != nil && perso.Type != p.Type {
		if target == nil || perso.Health < target.Health {
			target = perso
		}
	}
	return target
}

type SquareType int

const (
	Wall SquareType = iota
	Cavern
)

type Square struct {
	X     int
	Y     int
	Type  SquareType
	Perso *Perso
}

type Movement int

const (
	None Movement = iota
	Up
	Down
	Left
	Right
)

type Path struct {
	Distance      int
	Target        *Square
	PreviousSteps map[*Square][]*Square
}

// parsing the input
func getMap(fileName string) ([][]*Square, []*Perso, []*Perso, error) {
	var elves []*Perso
	var goblins []*Perso
	var squareMap [][]*Square

	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var x, y int
	for scanner.Scan() {
		line := scanner.Text()
		var squareLine []*Square
		x = 0
		for _, s := range line {
			square := &Square{X: x, Y: y}
			switch s {
			case '#':
				square.Type = Wall
			case '.':
				square.Type = Cavern
			case 'E':
				square.Type = Cavern
				elf := &Perso{CurrentSquare: square, Type: Elf, Health: initialHealth}
				square.Perso = elf
				elves = append(elves, elf)
			case 'G':
				square.Type = Cavern
				goblin := &Perso{CurrentSquare: square, Type: Goblin, Health: initialHealth}
				square.Perso = goblin
				goblins = append(goblins, goblin)
			default:
				return nil, nil, nil, errors.New("Can't parse the track")
			}
			squareLine = append(squareLine, square)
			x++
		}
		squareMap = append(squareMap, squareLine)
		y++
	}
	return squareMap, elves, goblins, nil
}

func fight(squareMap [][]*Square, elves []*Perso, goblins []*Perso, elvesAttack int, earlyStop bool) (int, int) {
	var turns int

	// figthing up until when we can't find an ennemy, or early stop if needed
Fight:
	for {
		turns++

		// getting turn order
		persos := make([]*Perso, len(elves))

		// copying needed to avoid append issues corrupting elves slice
		copy(persos, elves)
		persos = append(persos, goblins...)
		sort.Slice(persos, func(i, j int) bool {
			if persos[i].CurrentSquare.Y < persos[j].CurrentSquare.Y {
				return true
			} else if persos[i].CurrentSquare.Y == persos[j].CurrentSquare.Y {
				return persos[i].CurrentSquare.X < persos[j].CurrentSquare.X
			} else {
				return false
			}
		})

		for _, perso := range persos {
			// making sure that a dead perso doesn't act
			if perso.Health > 0 {
				// Search for ennemy already in range
				target := perso.HasEnnemyInRange(squareMap)

				// Move if no ennemy in range
				if target == nil {
					// getting in range squares
					var inRangeSquares []*Square
					switch perso.Type {
					case Elf:
						if len(goblins) == 0 {
							// breaking fight if no ennemies remaining
							break Fight
						}
						inRangeSquares = getInRange(squareMap, goblins)
					case Goblin:
						if len(elves) == 0 {
							// breaking fight if no ennemies remaining
							break Fight
						}
						inRangeSquares = getInRange(squareMap, elves)
					}

					if len(inRangeSquares) > 0 {
						// getting available paths from current position to in range squares
						paths := getReachableDistances(squareMap, perso.CurrentSquare, inRangeSquares)
						if len(paths) > 0 {
							// if some paths exist, we take the shortest, and we move
							nextSquare := getClosestSquare(paths, perso.CurrentSquare)
							perso.CurrentSquare.Perso = nil
							perso.CurrentSquare = nextSquare
							nextSquare.Perso = perso
							// looking again for adjacent ennemy
							target = perso.HasEnnemyInRange(squareMap)
						}
					}
				}

				// attack if target available
				if target != nil {
					if perso.Type == Elf {
						target.Health -= elvesAttack
					} else {
						target.Health -= attackPower
					}
					if target.Health <= 0 {
						if earlyStop && target.Type == Elf {
							// early stop for part 2 if an elf dies
							return 0, 0
						}
						// removing dead perso
						target.CurrentSquare.Perso = nil
						target.CurrentSquare = nil
						switch target.Type {
						case Elf:
							elves = removePersoFromList(elves, target)
						case Goblin:
							goblins = removePersoFromList(goblins, target)
						}
					}
				}
			}
		}
	}
	// computing remaining health once fight is over
	var totalHealth int
	for _, elf := range elves {
		totalHealth += elf.Health
	}
	for _, goblin := range goblins {
		totalHealth += goblin.Health
	}
	return turns - 1, totalHealth
}

func getClosestSquare(paths []Path, current *Square) *Square {
	var minDistance int
	var minPaths []Path

	// getting the shortest paths
	for _, path := range paths {
		if minDistance == 0 || path.Distance < minDistance {
			minDistance = path.Distance
			minPaths = []Path{path}
		} else if path.Distance == minDistance {
			minPaths = append(minPaths, path)
		}
	}

	// keeping only the shortest path with top-leftmost target
	var bestPath Path
	for _, path := range minPaths {
		if bestPath.Distance == 0 || path.Target.Y < bestPath.Target.Y || (path.Target.Y == bestPath.Target.Y && path.Target.X < bestPath.Target.X) {
			bestPath = path
		}
	}

	// returning first step of the path by priority
	return getTopLeftMost(getFirstSquares(bestPath.PreviousSteps, current, bestPath.Target))
}

// helper to get the top-leftmost square in the list
func getTopLeftMost(squares []*Square) *Square {
	var best *Square
	for _, square := range squares {
		if best == nil || square.Y < best.Y || (square.Y == best.Y && square.X < best.X) {
			best = square
		}
	}
	return best
}

// helper to remove a perso from a slice
func removePersoFromList(persos []*Perso, perso *Perso) []*Perso {
	for i, p := range persos {
		if p == perso {
			return append(persos[:i], persos[i+1:]...)
		}
	}
	return persos
}

// looking for squares in range for all persos alive
func getInRange(squareMap [][]*Square, persos []*Perso) []*Square {
	var squares []*Square
	for _, perso := range persos {
		if square := squareMap[perso.CurrentSquare.Y-1][perso.CurrentSquare.X]; square.Type == Cavern && square.Perso == nil {
			squares = append(squares, square)
		}
		if square := squareMap[perso.CurrentSquare.Y][perso.CurrentSquare.X-1]; square.Type == Cavern && square.Perso == nil {
			squares = append(squares, square)
		}
		if square := squareMap[perso.CurrentSquare.Y][perso.CurrentSquare.X+1]; square.Type == Cavern && square.Perso == nil {
			squares = append(squares, square)
		}
		if square := squareMap[perso.CurrentSquare.Y+1][perso.CurrentSquare.X]; square.Type == Cavern && square.Perso == nil {
			squares = append(squares, square)
		}
	}
	return squares
}

func getReachableDistances(squareMap [][]*Square, from *Square, inRangeSquares []*Square) []Path {
	var paths []Path
	for _, inRangeSquare := range inRangeSquares {
		path := getShortestPath(squareMap, from, inRangeSquare)
		if path.Target != nil {
			paths = append(paths, path)
		}
	}
	return paths
}

// Dijkstra to get the shortest paths (a Path contains the length and previous
// nodes
func getShortestPath(squareMap [][]*Square, from, to *Square) Path {
	// building node set
	dijkstraNodesSet := make(map[*Square]int)
	dijkstraPreviousNodes := make(map[*Square][]*Square)
	for _, row := range squareMap {
		for _, square := range row {
			if square.Type == Cavern && square.Perso == nil {
				dijkstraNodesSet[square] = -1 // -1 for infinity
			}
		}
	}

	// Using source node as first node
	currentSquare := from
	dijkstraNodesSet[from] = 0

	var shortestDistance int

	for len(dijkstraNodesSet) > 0 {
		// finding node with shortest tentative distance
		minDistance := -1
		var minSquare *Square
		for square, tentativeDistance := range dijkstraNodesSet {
			if tentativeDistance != -1 && (minDistance == -1 || tentativeDistance < minDistance) {
				minDistance = tentativeDistance
				minSquare = square
			}
		}

		if minSquare == nil {
			// only unreachable nodes left, no path to be found
			return Path{}
		}

		// if we find the target node, we've found all shortest paths, exiting
		currentSquare = minSquare
		if currentSquare == to {
			shortestDistance = dijkstraNodesSet[currentSquare]
			break
		}

		// looking at neighbours, increasing distance. Could be put in a func as it's
		// ùostly the same code, but copy paste is enough here
		if square := squareMap[currentSquare.Y-1][currentSquare.X]; square.Type == Cavern && square.Perso == nil {
			if currentDistance, ok := dijkstraNodesSet[square]; ok {
				if currentDistance == -1 || dijkstraNodesSet[currentSquare]+1 < currentDistance { // if unvisited
					dijkstraNodesSet[square] = dijkstraNodesSet[currentSquare] + 1
					dijkstraPreviousNodes[square] = []*Square{currentSquare}
				} else if dijkstraNodesSet[currentSquare]+1 == currentDistance {
					dijkstraPreviousNodes[square] = append(dijkstraPreviousNodes[square], currentSquare)
				}
			}
		}
		if square := squareMap[currentSquare.Y][currentSquare.X-1]; square.Type == Cavern && square.Perso == nil {
			if currentDistance, ok := dijkstraNodesSet[square]; ok {
				if currentDistance == -1 || dijkstraNodesSet[currentSquare]+1 < currentDistance { // if unvisited
					dijkstraNodesSet[square] = dijkstraNodesSet[currentSquare] + 1
					dijkstraPreviousNodes[square] = []*Square{currentSquare}
				} else if dijkstraNodesSet[currentSquare]+1 == currentDistance {
					dijkstraPreviousNodes[square] = append(dijkstraPreviousNodes[square], currentSquare)
				}
			}
		}
		if square := squareMap[currentSquare.Y][currentSquare.X+1]; square.Type == Cavern && square.Perso == nil {
			if currentDistance, ok := dijkstraNodesSet[square]; ok {
				if currentDistance == -1 || dijkstraNodesSet[currentSquare]+1 < currentDistance { // if unvisited
					dijkstraNodesSet[square] = dijkstraNodesSet[currentSquare] + 1
					dijkstraPreviousNodes[square] = []*Square{currentSquare}
				} else if dijkstraNodesSet[currentSquare]+1 == currentDistance {
					dijkstraPreviousNodes[square] = append(dijkstraPreviousNodes[square], currentSquare)
				}
			}
		}
		if square := squareMap[currentSquare.Y+1][currentSquare.X]; square.Type == Cavern && square.Perso == nil {
			if currentDistance, ok := dijkstraNodesSet[square]; ok {
				if currentDistance == -1 || dijkstraNodesSet[currentSquare]+1 < currentDistance { // if unvisited
					dijkstraNodesSet[square] = dijkstraNodesSet[currentSquare] + 1
					dijkstraPreviousNodes[square] = []*Square{currentSquare}
				} else if dijkstraNodesSet[currentSquare]+1 == currentDistance {
					dijkstraPreviousNodes[square] = append(dijkstraPreviousNodes[square], currentSquare)
				}
			}
		}
		// marking node visited by removing it from the map / set
		delete(dijkstraNodesSet, currentSquare)
	}

	// returning the necessary data to compute first step
	return Path{Distance: shortestDistance, Target: to, PreviousSteps: dijkstraPreviousNodes}
}

// getting the first square from the available paths, according to reading order
// we could early return if we find the top node to improve performance
func getFirstSquares(previousSquares map[*Square][]*Square, initialSquare *Square, currentSquare *Square) []*Square {
	var firstSquares []*Square
	for _, child := range previousSquares[currentSquare] {
		if child == initialSquare {
			return []*Square{currentSquare}
		} else {
			firstSquares = append(firstSquares, getFirstSquares(previousSquares, initialSquare, child)...)
		}
	}
	return firstSquares
}

// helper to print the map
func printMap(squareMap [][]*Square) {
	for _, row := range squareMap {
		var s strings.Builder
		for _, square := range row {
			if square == nil {
				s.WriteRune('X')
			} else if square.Perso != nil {
				switch square.Perso.Type {
				case Elf:
					s.WriteRune('E')
				case Goblin:
					s.WriteRune('G')
				}
			} else {
				switch square.Type {
				case Wall:
					s.WriteRune('#')
				case Cavern:
					s.WriteRune('.')
				}
			}
		}
		fmt.Println(s.String())
	}
}

// helper to print the persos
func printPersos(persos []*Perso) {
	for _, perso := range persos {
		fmt.Printf("[%d] %d,%d : %d\n", perso.Type, perso.CurrentSquare.X, perso.CurrentSquare.Y, perso.Health)
	}
}
