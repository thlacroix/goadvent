package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// defining the size of image and the number of layers
// (precomputed) to be used in array sizes
const (
	wide     = 25
	tall     = 6
	nbLayers = 100
)

// Layers is a helper type representing all the layers of
// an image
type Layers [nbLayers][tall][wide]int

func main() {
	ints, err := getInts("day08input.txt")
	if err != nil {
		log.Fatal(err)
	}
	layers := getLayers(ints)
	fmt.Println(getLowestLayerScore(layers))
	image := getFinalImage(layers)
	printImage(image)
}

func getInts(fileName string) ([]int, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	str := strings.TrimSpace(string(content))
	ints := make([]int, 0, len(str))
	for _, c := range str {
		i, err := strconv.Atoi(string(c))
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

// Builing the image layers from the input
func getLayers(ints []int) Layers {
	var index int
	var layers Layers
	for i := 0; i < nbLayers; i++ {
		for j := 0; j < tall; j++ {
			for k := 0; k < wide; k++ {
				layers[i][j][k] = ints[index]
				index++
			}
		}
	}
	return layers
}

// Score help count the number of 0, 1 and 2 in a layer
type Score struct {
	Zeros int
	Ones  int
	Twos  int
}

// getLowestLayerScore returns the score (num of 1 * num of 2)
// of the layer with the lowest num of 0
func getLowestLayerScore(layers Layers) int {
	var layerScores [nbLayers]Score
	minZeroes := -1
	minIndex := -1

	for i, layer := range layers {
		var s Score
		for _, row := range layer {
			for _, v := range row {
				switch v {
				case 0:
					s.Zeros++
				case 1:
					s.Ones++
				case 2:
					s.Twos++
				}
			}
		}
		if minZeroes == -1 || s.Zeros < minZeroes {
			minZeroes = s.Zeros
			minIndex = i
		}
		layerScores[i] = s
	}
	score := layerScores[minIndex]
	return score.Ones * score.Twos
}

// Getting the final message by superposing the layers
func getFinalImage(layers Layers) [tall][wide]int {
	var image [tall][wide]int
	for i := 0; i < tall; i++ {
		for j := 0; j < wide; j++ {
		layersLoop:
			for l := 0; l < nbLayers; l++ {
				switch layers[l][i][j] {
				case 2:
					continue layersLoop
				case 1:
					image[i][j] = 1
					break layersLoop
				case 0:
					image[i][j] = 0
					break layersLoop
				}
			}
		}
	}
	return image
}

// Printing the image to help read the message
func printImage(image [tall][wide]int) {
	for _, row := range image {
		fmt.Println(row)
	}
}
