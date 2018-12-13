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

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if tracks, karts, err := getTracks(fileName); err != nil {
		log.Fatal(err)
	} else {
		x, y := moveKarts(tracks, karts, true)
		fmt.Println("Accident at", x, ",", y)
	}
}

type TrackType int

const (
	Horizontal TrackType = iota
	Vertical
	Intersection
	CurveUp
	CurveDown
)

type Track struct {
	X    int
	Y    int
	Type TrackType
	Kart *Kart
}

type KartDirection int

const (
	Top KartDirection = iota
	Bottom
	Left
	Right
)

type KartTurn int

const (
	LeftTurn KartTurn = iota
	RightTurn
	StraightTurn
)

type Kart struct {
	CurrentTrack *Track
	Direction    KartDirection
	NextTurn     KartTurn
}

func (k *Kart) MoveUp(tracks [][]*Track) {
	k.CurrentTrack = tracks[k.CurrentTrack.Y-1][k.CurrentTrack.X]
	k.Direction = Top
}

func (k *Kart) MoveDown(tracks [][]*Track) {
	k.CurrentTrack = tracks[k.CurrentTrack.Y+1][k.CurrentTrack.X]
	k.Direction = Bottom
}

func (k *Kart) MoveLeft(tracks [][]*Track) {
	k.CurrentTrack = tracks[k.CurrentTrack.Y][k.CurrentTrack.X-1]
	k.Direction = Left
}

func (k *Kart) MoveRight(tracks [][]*Track) {
	k.CurrentTrack = tracks[k.CurrentTrack.Y][k.CurrentTrack.X+1]
	k.Direction = Right
}

func (k *Kart) Move(tracks [][]*Track) bool {
	initalTrack := k.CurrentTrack
	// moving kart based on direction and track type
	switch k.Direction {
	case Top:
		switch initalTrack.Type {
		case Vertical:
			k.MoveUp(tracks)
		case CurveUp:
			k.MoveRight(tracks)
		case CurveDown:
			k.MoveLeft(tracks)
		case Intersection:
			switch k.NextTurn {
			case LeftTurn:
				k.MoveLeft(tracks)
				k.NextTurn = StraightTurn
			case StraightTurn:
				k.MoveUp(tracks)
				k.NextTurn = RightTurn
			case RightTurn:
				k.MoveRight(tracks)
				k.NextTurn = LeftTurn
			}
		}
	case Bottom:
		switch initalTrack.Type {
		case Vertical:
			k.MoveDown(tracks)
		case CurveUp:
			k.MoveLeft(tracks)
		case CurveDown:
			k.MoveRight(tracks)
		case Intersection:
			switch k.NextTurn {
			case LeftTurn:
				k.MoveRight(tracks)
				k.NextTurn = StraightTurn
			case StraightTurn:
				k.MoveDown(tracks)
				k.NextTurn = RightTurn
			case RightTurn:
				k.MoveLeft(tracks)
				k.NextTurn = LeftTurn
			}
		}
	case Right:
		switch initalTrack.Type {
		case Horizontal:
			k.MoveRight(tracks)
		case CurveUp:
			k.MoveUp(tracks)
		case CurveDown:
			k.MoveDown(tracks)
		case Intersection:
			switch k.NextTurn {
			case LeftTurn:
				k.MoveUp(tracks)
				k.NextTurn = StraightTurn
			case StraightTurn:
				k.MoveRight(tracks)
				k.NextTurn = RightTurn
			case RightTurn:
				k.MoveDown(tracks)
				k.NextTurn = LeftTurn
			}
		}
	case Left:
		switch initalTrack.Type {
		case Horizontal:
			k.MoveLeft(tracks)
		case CurveUp:
			k.MoveDown(tracks)
		case CurveDown:
			k.MoveUp(tracks)
		case Intersection:
			switch k.NextTurn {
			case LeftTurn:
				k.MoveDown(tracks)
				k.NextTurn = StraightTurn
			case StraightTurn:
				k.MoveLeft(tracks)
				k.NextTurn = RightTurn
			case RightTurn:
				k.MoveUp(tracks)
				k.NextTurn = LeftTurn
			}
		}
	}
	// cleaning previous track
	initalTrack.Kart = nil
	if k.CurrentTrack.Kart != nil {
		// kart crash
		return true
	}
	k.CurrentTrack.Kart = k
	return false
}

func getTracks(fileName string) ([][]*Track, []*Kart, error) {
	var karts []*Kart
	var tracks [][]*Track

	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var x, y int

	for scanner.Scan() {
		line := scanner.Text()
		if tracks == nil { // initilazing the grid from size
			tracks = make([][]*Track, len(line)+1)
			for i := range tracks {
				tracks[i] = make([]*Track, len(line)+1)
			}
		}
		x = 0
		for _, trackC := range line {
			var track *Track
			var kart *Kart
			switch trackC {
			case ' ':
			case '|':
				track = &Track{X: x, Y: y, Type: Vertical}
			case '-':
				track = &Track{X: x, Y: y, Type: Horizontal}
			case '\\':
				track = &Track{X: x, Y: y, Type: CurveDown}
			case '/':
				track = &Track{X: x, Y: y, Type: CurveUp}
			case '+':
				track = &Track{X: x, Y: y, Type: Intersection}
			case 'v':
				track = &Track{X: x, Y: y, Type: Vertical}
				kart = &Kart{CurrentTrack: track, Direction: Bottom}
				track.Kart = kart
			case '^':
				track = &Track{X: x, Y: y, Type: Vertical}
				kart = &Kart{CurrentTrack: track, Direction: Top}
				track.Kart = kart
			case '>':
				track = &Track{X: x, Y: y, Type: Horizontal}
				kart = &Kart{CurrentTrack: track, Direction: Right}
				track.Kart = kart
			case '<':
				track = &Track{X: x, Y: y, Type: Horizontal}
				kart = &Kart{CurrentTrack: track, Direction: Left}
				track.Kart = kart
			default:
				return nil, nil, errors.New("Can't parse the track")
			}

			tracks[y][x] = track
			if kart != nil {
				karts = append(karts, kart)
			}
			x++
		}
		y++
	}
	return tracks, karts, nil
}

func moveKarts(tracks [][]*Track, karts []*Kart, remove bool) (int, int) {
	kartCount := len(karts)
	for kartCount > 1 {
		// sorting the karts for move order, might be optimized as sorting is heavy
		sort.Slice(karts, func(i int, j int) bool {
			if karts[i].CurrentTrack == nil {
				return false
			} else if karts[j].CurrentTrack == nil {
				return true
			}

			if karts[i].CurrentTrack.Y < karts[j].CurrentTrack.Y {
				return true
			} else if karts[i].CurrentTrack.Y == karts[j].CurrentTrack.Y {
				return karts[i].CurrentTrack.X < karts[j].CurrentTrack.X
			} else {
				return false
			}
		})
		for _, kart := range karts {
			if kart.CurrentTrack != nil {
				if boom := kart.Move(tracks); boom {
					if !remove {
						return kart.CurrentTrack.X, kart.CurrentTrack.Y
					} else {
						// removing the karts
						otherKart := kart.CurrentTrack.Kart
						otherKart.CurrentTrack.Kart = nil
						otherKart.CurrentTrack = nil
						kart.CurrentTrack.Kart = nil
						kart.CurrentTrack = nil
						kartCount -= 2
					}
				}
			}
		}
	}
	if kartCount == 1 {
		for _, lastKart := range karts {
			if lastKart.CurrentTrack != nil {
				return lastKart.CurrentTrack.X, lastKart.CurrentTrack.Y
			}
		}
	}
	return -1, -1
}

func printTracks(tracks [][]*Track) {
	for _, row := range tracks {
		var s strings.Builder
		for _, track := range row {
			if track == nil {
				s.WriteRune(' ')
			} else if track.Kart != nil {
				switch track.Kart.Direction {
				case Top:
					s.WriteRune('^')
				case Bottom:
					s.WriteRune('v')
				case Right:
					s.WriteRune('>')
				case Left:
					s.WriteRune('<')
				}
			} else {
				switch track.Type {
				case Intersection:
					s.WriteRune('+')
				case Horizontal:
					s.WriteRune('-')
				case Vertical:
					s.WriteRune('|')
				case CurveUp:
					s.WriteRune('/')
				case CurveDown:
					s.WriteRune('\\')
				}
			}
		}
		fmt.Println(s.String())
	}
}

func printKarts(karts []*Kart) {
	for _, kart := range karts {
		if kart.CurrentTrack != nil {
			fmt.Println(kart.CurrentTrack.X, kart.CurrentTrack.Y)
		}
	}
}
