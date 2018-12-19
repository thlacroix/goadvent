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

var rBefore = regexp.MustCompile(`Before: \[(\d), (\d), (\d), (\d)\]`)
var rInstruction = regexp.MustCompile(`(\d+) (\d+) (\d+) (\d+)`)
var rAfter = regexp.MustCompile(`After:  \[(\d), (\d), (\d), (\d)\]`)

var opscodes = []Opscode{
	Opscode{
		Name: "addr",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			result[output] = registers[inputa] + registers[inputb]
			return result
		},
	},
	Opscode{
		Name: "addi",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			result[output] = registers[inputa] + inputb
			return result
		},
	},
	Opscode{
		Name: "mulr",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			result[output] = registers[inputa] * registers[inputb]
			return result
		},
	},
	Opscode{
		Name: "muli",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			result[output] = registers[inputa] * inputb
			return result
		},
	},
	Opscode{
		Name: "banr",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			result[output] = registers[inputa] & registers[inputb]
			return result
		},
	},
	Opscode{
		Name: "bani",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			result[output] = registers[inputa] & inputb
			return result
		},
	},
	Opscode{
		Name: "borr",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			result[output] = registers[inputa] | registers[inputb]
			return result
		},
	},
	Opscode{
		Name: "bori",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			result[output] = registers[inputa] | inputb
			return result
		},
	},
	Opscode{
		Name: "setr",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			result[output] = registers[inputa]
			return result
		},
	},
	Opscode{
		Name: "seti",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			result[output] = inputa
			return result
		},
	},
	Opscode{
		Name: "gtir",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			var value int
			if inputa > registers[inputb] {
				value = 1
			} else {
				value = 0
			}
			result[output] = value
			return result
		},
	},
	Opscode{
		Name: "gtri",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			var value int
			if registers[inputa] > inputb {
				value = 1
			} else {
				value = 0
			}
			result[output] = value
			return result
		},
	},
	Opscode{
		Name: "gtrr",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			var value int
			if registers[inputa] > registers[inputb] {
				value = 1
			} else {
				value = 0
			}
			result[output] = value
			return result
		},
	},
	Opscode{
		Name: "eqir",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			var value int
			if inputa == registers[inputb] {
				value = 1
			} else {
				value = 0
			}
			result[output] = value
			return result
		},
	},
	Opscode{
		Name: "eqri",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			var value int
			if registers[inputa] == inputb {
				value = 1
			} else {
				value = 0
			}
			result[output] = value
			return result
		},
	},
	Opscode{
		Name: "eqrr",
		Func: func(registers [4]int, inputa int, inputb int, output int) [4]int {
			result := registers
			var value int
			if registers[inputa] == registers[inputb] {
				value = 1
			} else {
				value = 0
			}
			result[output] = value
			return result
		},
	},
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No filepath passed")
	}
	fileName := os.Args[1]
	if samples, instructions, err := getSamples(fileName); err != nil {
		log.Fatal(err)
	} else {
		res := processSamples(samples)
		fmt.Println("Part1 result is", res)
		if !checkSamples(samples) {
			log.Fatal("Opscode mapping is wrong")
		}
		fmt.Println("Part2 result is", computeInstructions(instructions)[0])
		//fmt.Println(instructions)
	}
}

type Sample struct {
	Before      [4]int
	Instruction Instruction
	After       [4]int
}

type Instruction struct {
	OpscodeID int
	InputA    int
	InputB    int
	Output    int
}

type OpscodeFunc func([4]int, int, int, int) [4]int

type Opscode struct {
	ID   int
	Name string
	Func OpscodeFunc
}

func getSamples(fileName string) ([]Sample, []Instruction, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var i int
	var currentSample Sample
	var samples []Sample
Scan:
	for scanner.Scan() {
		line := scanner.Text()
		switch i % 4 {
		case 0: // Before
			if line == "" { // end of first input
				break Scan
			}
			extracts := rBefore.FindAllStringSubmatch(line, -1)
			if len(extracts) != 1 || len(extracts[0]) != 5 {
				return nil, nil, errors.New("Can't parse before line " + line)
			}
			currentSample.Before = [4]int{atoi(extracts[0][1]), atoi(extracts[0][2]), atoi(extracts[0][3]), atoi(extracts[0][4])}
		case 1: // Instruction
			extracts := rInstruction.FindAllStringSubmatch(line, -1)
			if len(extracts) != 1 || len(extracts[0]) != 5 {
				return nil, nil, errors.New("Can't parse instruction line " + line)
			}
			currentSample.Instruction = Instruction{
				OpscodeID: atoi(extracts[0][1]),
				InputA:    atoi(extracts[0][2]),
				InputB:    atoi(extracts[0][3]),
				Output:    atoi(extracts[0][4]),
			}
		case 2: // Ater
			extracts := rAfter.FindAllStringSubmatch(line, -1)
			if len(extracts) != 1 || len(extracts[0]) != 5 {
				return nil, nil, errors.New("Can't parse after line " + line)
			}
			currentSample.After = [4]int{atoi(extracts[0][1]), atoi(extracts[0][2]), atoi(extracts[0][3]), atoi(extracts[0][4])}
		case 3:
			if line != "" {
				return nil, nil, errors.New("Expected newline no found")
			}
			samples = append(samples, currentSample)
			currentSample = Sample{}
		}
		i++
	}
	scanner.Scan() // other empty row

	// reading escond part
	var instructions []Instruction
	for scanner.Scan() {
		line := scanner.Text()
		extracts := rInstruction.FindAllStringSubmatch(line, -1)
		if len(extracts) != 1 || len(extracts[0]) != 5 {
			return nil, nil, errors.New("Can't parse instruction line " + line)
		}
		instruction := Instruction{
			OpscodeID: atoi(extracts[0][1]),
			InputA:    atoi(extracts[0][2]),
			InputB:    atoi(extracts[0][3]),
			Output:    atoi(extracts[0][4]),
		}
		instructions = append(instructions, instruction)
	}
	return samples, instructions, nil
}

func processSamples(samples []Sample) int {
	var matchMoreThan3Opscode int
	idToNames := make(map[int][]string) // mapping opscodeId to possible opscode names
	for _, sample := range samples {
		var matchingOpscodes []string
		// runnnging all opscodes on the sample, checking the matching ones
		for _, opscode := range opscodes {
			output := opscode.Func(sample.Before, sample.Instruction.InputA, sample.Instruction.InputB, sample.Instruction.Output)
			if compareOutputs(output, sample.After) {
				matchingOpscodes = append(matchingOpscodes, opscode.Name)
			}
		}
		// increasing count if more that 3 opscodes match
		if len(matchingOpscodes) >= 3 {
			matchMoreThan3Opscode++
		}

		// intersecting matching ones with previous matching ones to find opscodes
		// that always match
		if len(idToNames[sample.Instruction.OpscodeID]) == 0 {
			idToNames[sample.Instruction.OpscodeID] = matchingOpscodes
		} else {
			idToNames[sample.Instruction.OpscodeID] = intersect(idToNames[sample.Instruction.OpscodeID], matchingOpscodes)
		}
	}

	// finding mapping opscode name / opscode id by removing already found ones
	// from possible list of others
	matchedOpscodes := make(map[string]int)
	for len(matchedOpscodes) != len(opscodes) {
		for id, names := range idToNames {
			if len(names) == 1 {
				// found
				matchedOpscodes[names[0]] = id
				delete(idToNames, id)
			} else {
				// removing already found
				var namesAfterRemovals []string
				for _, name := range names {
					if _, ok := matchedOpscodes[name]; !ok {
						namesAfterRemovals = append(namesAfterRemovals, name)
					}
				}
				idToNames[id] = namesAfterRemovals
			}
		}
	}

	// updating global opscode ids
	for i, opscode := range opscodes {
		opscodes[i].ID = matchedOpscodes[opscode.Name]
	}
	return matchMoreThan3Opscode
}

func checkSamples(samples []Sample) bool {
	// building opscode map
	opscodesMap := make(map[int]Opscode)
	for _, opscode := range opscodes {
		opscodesMap[opscode.ID] = opscode
	}
	// running instruction on registers
	for _, sample := range samples {
		output := opscodesMap[sample.Instruction.OpscodeID].Func(
			sample.Before,
			sample.Instruction.InputA,
			sample.Instruction.InputB,
			sample.Instruction.Output,
		)
		if !compareOutputs(output, sample.After) {
			return false
		}
	}
	return true
}

func computeInstructions(instructions []Instruction) [4]int {
	// building opscode map
	opscodesMap := make(map[int]Opscode)
	for _, opscode := range opscodes {
		opscodesMap[opscode.ID] = opscode
	}
	var registers [4]int
	for _, instruction := range instructions {
		registers = opscodesMap[instruction.OpscodeID].Func(registers, instruction.InputA, instruction.InputB, instruction.Output)
	}
	return registers
}

// comparing two arrays
func compareOutputs(o1 [4]int, o2 [4]int) bool {
	for i := range o1 {
		if o1[i] != o2[i] {
			return false
		}
	}
	return true
}

// unsafe string -> integer parsing
func atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}

// gives the intersection of two slices
func intersect(a []string, b []string) []string {
	var res []string
	for _, ka := range a {
		for _, kb := range b {
			if ka == kb {
				res = append(res, ka)
			}
		}
	}
	return res
}
