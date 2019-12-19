package intcode

// Mode represent the mode of a parameter defined in the operatioh=n
type Mode byte

const (
	Position Mode = iota
	Immediate
	Relative
)

// Machine is the IntCode program and stores the state
// The Input and Output chans are used to communicate with the machine
type Machine struct {
	Ints   []int
	Index  int
	Base   int
	Input  chan int
	Output chan int
	Done   chan bool
}

// NewMachine returns a new machine with a copy of the input ints
func NewMachine(ints []int) *Machine {
	return NewBufferedMachine(ints, 0, 0)
}

// NewBufferedMachine returns a new machine with a copy of the input ints,
// with the possibility to bufferise the input and output channels
func NewBufferedMachine(ints []int, inBuffer, outBuffer int) *Machine {
	intsCopy := make([]int, len(ints))
	copy(intsCopy, ints)
	return &Machine{Ints: intsCopy, Input: make(chan int, inBuffer), Output: make(chan int, outBuffer), Done: make(chan bool)}
}

// AddInput sends the input on the Input channel and return true
// If the machine is finished, returns false instead of sending
func (m *Machine) AddInput(i int) bool {
	for {
		select {
		case m.Input <- i:
			return true
		case <-m.Done:
			return false
		}
	}
}

// GetOutput returns the first available ouput from the Output channel
func (m *Machine) GetOutput() int {
	return <-m.Output
}

// GetOutputOrAddInputOrEnd either adds and input (and return the first bool as true),
// get an output, or signal the end of the program (last bool as true)
func (m *Machine) GetOutputOrAddInputOrEnd(i int) (int, bool, bool) {
	for {
		select {
		case o := <-m.Output:
			return o, false, false
		case m.Input <- i:
			return 0, true, false
		case <-m.Done:
			return 0, false, true
		}
	}
}

// GetOuputOrEnd returns the first available ouput from the Output channel,
// and true if the machine exited
func (m *Machine) GetOutputOrEnd() (int, bool) {
	for {
		select {
		case o := <-m.Output:
			return o, false
		case <-m.Done:
			return 0, true
		}
	}
}

// Run runs the machine, and sends a signal on the Done chan when done
func (m *Machine) Run() {
	for m.Index < len(m.Ints) {
		operation := m.Ints[m.Index] % 100
		switch operation {
		case 1:
			modes, parameters := m.getModesParameters(3)
			a, b := m.getValue(parameters[0], modes[0]), m.getValue(parameters[1], modes[1])
			m.writeInt(parameters[2], a+b, modes[2])
			m.Index += 4
		case 2:
			modes, parameters := m.getModesParameters(3)
			a, b := m.getValue(parameters[0], modes[0]), m.getValue(parameters[1], modes[1])
			m.writeInt(parameters[2], a*b, modes[2])
			m.Index += 4
		case 3:
			modes, parameters := m.getModesParameters(1)
			input := <-m.Input
			m.writeInt(parameters[0], input, modes[0])
			m.Index += 2
		case 4:
			modes, parameters := m.getModesParameters(1)
			a := m.getValue(parameters[0], modes[0])
			m.Output <- a
			m.Index += 2
		case 5:
			modes, parameters := m.getModesParameters(2)
			a, b := m.getValue(parameters[0], modes[0]), m.getValue(parameters[1], modes[1])
			if a != 0 {
				m.Index = b
			} else {
				m.Index += 3
			}
		case 6:
			modes, parameters := m.getModesParameters(2)
			a, b := m.getValue(parameters[0], modes[0]), m.getValue(parameters[1], modes[1])
			if a == 0 {
				m.Index = b
			} else {
				m.Index += 3
			}
		case 7:
			modes, parameters := m.getModesParameters(3)
			a, b := m.getValue(parameters[0], modes[0]), m.getValue(parameters[1], modes[1])
			if a < b {
				m.writeInt(parameters[2], 1, modes[2])
			} else {
				m.writeInt(parameters[2], 0, modes[2])
			}
			m.Index += 4
		case 8:
			modes, parameters := m.getModesParameters(3)
			a, b := m.getValue(parameters[0], modes[0]), m.getValue(parameters[1], modes[1])
			if a == b {
				m.writeInt(parameters[2], 1, modes[2])
			} else {
				m.writeInt(parameters[2], 0, modes[2])
			}
			m.Index += 4
		case 9:
			modes, parameters := m.getModesParameters(1)
			a := m.getValue(parameters[0], modes[0])
			m.Base += a
			m.Index += 2
		case 99:
			m.Done <- true
			return
		}
	}
	m.Done <- true
	return
}

// Using a helper to write to the list, depending on the mode, and if the
// list is long enough
func (m *Machine) writeInt(index, value int, mode Mode) {
	if mode == Relative {
		index += m.Base
	}

	if index < len(m.Ints) {
		m.Ints[index] = value
		return
	}

	intsCopy := make([]int, index+1)
	copy(intsCopy, m.Ints)
	intsCopy[index] = value
	m.Ints = intsCopy

	return
}

// Takes a param and its mode, and return the values to use
func (m *Machine) getValue(a int, mode Mode) int {
	switch mode {
	case Position:
		if a >= len(m.Ints) {
			return 0
		}
		return m.Ints[a]
	case Immediate:
		return a
	case Relative:
		if m.Base+a >= len(m.Ints) {
			return 0
		}
		return m.Ints[m.Base+a]
	}
	return -1
}

// Takes a number of parameters to process and return the list of modes and parameters.
// A mode to false means by position, true means immediate
func (m *Machine) getModesParameters(count int) ([]Mode, []int) {
	ope := m.Ints[m.Index]
	modes := make([]Mode, count)
	parameters := make([]int, count)
	div := 100
	for i := 0; i < count; i++ {
		mode := ope / div % 10
		modes[i] = Mode(mode)
		div = div * 10
		parameters[i] = m.Ints[m.Index+i+1]
	}
	return modes, parameters
}
