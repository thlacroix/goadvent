package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var metadataTotal int

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if root, err := getTree(fileName); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Total metadata is", metadataTotal, "and root value is", getNodeValue(root))
	}
}

type Node struct {
	Metadatas  []int
	Children   []*Node
	StartIndex int
}

func getTree(fileName string) (*Node, error) {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	// building input list
	input := strings.Split(strings.TrimSpace(string(fileContent)), " ")
	ints := make([]int, 0, len(input))
	for _, s := range input {
		d, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		ints = append(ints, d)
	}
	// getting root node and end index recursively, and checking that we're at the end
	root, end := getNextNode(ints, 0)
	if end != len(input) {
		return nil, errors.New(fmt.Sprintln("Last index is", end, " and should be", len(ints)-1))
	}
	return root, nil
}

func getNextNode(input []int, startIndex int) (*Node, int) {
	node := &Node{StartIndex: startIndex}
	childrenCount := input[startIndex]
	metadataCount := input[startIndex+1]
	nextIndex := startIndex + 2
	// for each child, we get it recursively with its end index, where we get the
	// next one
	for i := 0; i < childrenCount; i++ {
		child, end := getNextNode(input, nextIndex)
		node.Children = append(node.Children, child)
		nextIndex = end
	}
	// starting from where we stopped getting the children, we get the metadata
	for i := 0; i < metadataCount; i++ {
		node.Metadatas = append(node.Metadatas, input[nextIndex])
		metadataTotal += input[nextIndex]
		nextIndex++
	}
	// returning the node and the next index
	return node, nextIndex
}

// computing the node value recursively
func getNodeValue(node *Node) int {
	var value int
	for _, metadata := range node.Metadatas {
		if len(node.Children) > 0 {
			if metadata <= len(node.Children) {
				value += getNodeValue(node.Children[metadata-1])
			}
		} else {
			value += metadata
		}
	}
	return value
}
