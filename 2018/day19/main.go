package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const registerCount = 6

const target = 10551428

var rInstruction = regexp.MustCompile(`(\w+) (\d+) (\d+) (\d+)`)

var opscodes = map[string]Opscode{
	"addr": Opscode{
		Name: "addr",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
			result := registers
			result[output] = registers[inputa] + registers[inputb]
			return result
		},
	},
	"addi": Opscode{
		Name: "addi",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
			result := registers
			result[output] = registers[inputa] + inputb
			return result
		},
	},
	"mulr": Opscode{
		Name: "mulr",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
			result := registers
			result[output] = registers[inputa] * registers[inputb]
			return result
		},
	},
	"muli": Opscode{
		Name: "muli",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
			result := registers
			result[output] = registers[inputa] * inputb
			return result
		},
	},
	"banr": Opscode{
		Name: "banr",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
			result := registers
			result[output] = registers[inputa] & registers[inputb]
			return result
		},
	},
	"bani": Opscode{
		Name: "bani",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
			result := registers
			result[output] = registers[inputa] & inputb
			return result
		},
	},
	"borr": Opscode{
		Name: "borr",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
			result := registers
			result[output] = registers[inputa] | registers[inputb]
			return result
		},
	},
	"bori": Opscode{
		Name: "bori",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
			result := registers
			result[output] = registers[inputa] | inputb
			return result
		},
	},
	"setr": Opscode{
		Name: "setr",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
			result := registers
			result[output] = registers[inputa]
			return result
		},
	},
	"seti": Opscode{
		Name: "seti",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
			result := registers
			result[output] = inputa
			return result
		},
	},
	"gtir": Opscode{
		Name: "gtir",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
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
	"gtri": Opscode{
		Name: "gtri",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
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
	"gtrr": Opscode{
		Name: "gtrr",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
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
	"eqir": Opscode{
		Name: "eqir",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
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
	"eqri": Opscode{
		Name: "eqri",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
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
	"eqrr": Opscode{
		Name: "eqrr",
		Func: func(registers [registerCount]int, inputa int, inputb int, output int) [registerCount]int {
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
	if ip, instructions, err := getInstructions(fileName); err != nil {
		log.Fatal(err)
	} else {
		res := processInstructions(ip, instructions, 0)
		fmt.Println("Part1 result is", res)

		// For Part2, computing the solution with raw power would be too long, so
		// an analytic solution is nedeed. When running the simulation for a while,
		// we can notice a long-running loop between instructions 3 and 11.
		// By looking at the instructions, we can easily understand what they do.
		// Basically, we have a target, a multiplier, an increment, and a counter.
		// The instructions increases the counter when the multiplier times the
		// increment equals the target, and then try with the next multiplier.
		// This means that the counter is the sum of all divisors of the target
		fmt.Println(countDividers(target))
	}
}

type Instruction struct {
	Opscode Opscode
	InputA  int
	InputB  int
	Output  int
}

type OpscodeFunc func([registerCount]int, int, int, int) [registerCount]int

type Opscode struct {
	ID   int
	Name string
	Func OpscodeFunc
}

func getInstructions(fileName string) (int, []Instruction, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()
	ip := atoi(strings.Split(firstLine, "#ip ")[1])
	var instructions []Instruction
	for scanner.Scan() {
		line := scanner.Text()
		extracts := rInstruction.FindAllStringSubmatch(line, -1)
		if len(extracts) != 1 || len(extracts[0]) != 5 {
			return 0, nil, errors.New("Can't parse instruction line " + line)
		}
		instruction := Instruction{
			Opscode: opscodes[extracts[0][1]],
			InputA:  atoi(extracts[0][2]),
			InputB:  atoi(extracts[0][3]),
			Output:  atoi(extracts[0][4]),
		}
		instructions = append(instructions, instruction)
	}
	return ip, instructions, nil
}

func processInstructions(ip int, instructions []Instruction, firstRegisterValue int) int {
	var registers [registerCount]int
	registers[0] = firstRegisterValue
	for registers[ip] >= 0 && registers[ip] < len(instructions) {
		instruction := instructions[registers[ip]]
		registers = instruction.Opscode.Func(registers, instruction.InputA, instruction.InputB, instruction.Output)
		// increasing ip
		registers[ip]++
	}
	return registers[0]
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s %d %d %d", i.Opscode.Name, i.InputA, i.InputB, i.Output)
}

// unsafe string -> integer parsing
func atoi(s string) int {
	d, _ := strconv.Atoi(s)
	return d
}

// basic implementation to sum the dividers of the target, in O(√N)
func countDividers(v int) int {
	var count int
	for i := 1; i <= int(math.Sqrt(float64(v))); i++ {
		if v%i == 0 {
			count += i
			count += v / i
		}
	}
	return count
}
