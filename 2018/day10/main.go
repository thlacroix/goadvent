package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var rPlane = regexp.MustCompile(`position=<\s*(-|\s)(\d+),\s*(-|\s)(\d+)> velocity=<(-|\s)(\d), (-|\s)(\d)>`)

const maxTime = 15000

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if planes, timeToComplete, err := getMessage(fileName); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("The message seen in", timeToComplete, "seconds is")
		buildAndPrintZone(planes)
	}
}

type Plane struct {
	ID     int
	X      int
	Y      int
	XSpeed int
	YSpeed int
}

// parsing a plane line to return a Plane
func NewPlaneFromText(text string) Plane {
	extract := rPlane.FindAllStringSubmatch(text, -1)[0]
	return Plane{
		X:      valueFromPair(extract[1], extract[2]),
		Y:      valueFromPair(extract[3], extract[4]),
		XSpeed: valueFromPair(extract[5], extract[6]),
		YSpeed: valueFromPair(extract[7], extract[8]),
	}
}

// helper to get the int value from a sign + abs pair
func valueFromPair(sign, abs string) int {
	value, _ := strconv.Atoi(abs)
	if sign == "-" {
		return -value
	}
	return value
}

type ColumnMeta struct {
	Count int
	Min   int
	Max   int
}

func getMessage(fileName string) ([]Plane, int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var planes []Plane

	scanner := bufio.NewScanner(file)
	// building plane list
	for scanner.Scan() {
		planeText := scanner.Text()
		planes = append(planes, NewPlaneFromText(planeText))
	}

	var time int
	for time < maxTime {
		time++
		columnCount := make(map[int]ColumnMeta)

		// counting planes in column
		for j := range planes {
			planes[j].X += planes[j].XSpeed
			planes[j].Y += planes[j].YSpeed
			currentCount := columnCount[planes[j].X]
			if currentCount.Count == 0 || planes[j].Y < currentCount.Min {
				currentCount.Min = planes[j].Y
			}
			if currentCount.Count == 0 || planes[j].Y > currentCount.Max {
				currentCount.Max = planes[j].Y
			}
			currentCount.Count++
			columnCount[planes[j].X] = currentCount
		}

		// finding columns with max consecutive planes, corresponding to the vertical
		// lines of the letters
		var messageSize, messageSizeCount int
		for _, m := range columnCount {
			if m.Max-m.Min == m.Count-1 {
				if m.Count == messageSize {
					messageSizeCount++
				} else if m.Count > messageSize {
					messageSize = m.Count
					messageSizeCount = 1
				}
			}
		}
		// if we have enough lines that are long enough, could be a text
		if messageSize > 3 && messageSizeCount > 3 {
			return planes, time, nil
		}
	}
	return nil, 0, errors.New("Can't find the message")
}

// just displaying planes on a limited area
func buildAndPrintZone(planes []Plane) {
	// first building the base zone
	var minx, maxx, miny, maxy int
	for i, plane := range planes {
		if i == 0 || plane.X < minx {
			minx = plane.X
		}
		if i == 0 || plane.X > maxx {
			maxx = plane.X
		}
		if i == 0 || plane.Y < miny {
			miny = plane.Y
		}
		if i == 0 || plane.Y > maxy {
			maxy = plane.Y
		}
	}
	messageMap := make([][]int, maxy-miny+1)
	for i := range messageMap {
		messageMap[i] = make([]int, maxx-minx+1)
	}

	// adding the planes
	for _, plane := range planes {
		messageMap[plane.Y-miny][plane.X-minx] = 1
	}

	for _, row := range messageMap {
		fmt.Println(row)
	}
}
