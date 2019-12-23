package main

import (
	"fmt"
	"log"
	"time"

	"github.com/thlacroix/goadvent/2019/intcode"
	"github.com/thlacroix/goadvent/helpers"
)

func main() {
	ints, err := helpers.GetInts("day23input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(runMachines(ints, 50))
}

func runMachines(ints []int, N int) int {
	machines := make([]*intcode.Machine, N)
	res := make(chan int)
	nat := make(chan int, 100)

	// creating the machines and starts the consumers
	for i := 0; i < 50; i++ {
		m := intcode.NewBufferedMachine(ints, 1000, 1000)
		m = m.WithDefaultInput(-1)
		machines[i] = m
		go func(c chan int) {
			for {
				id := <-c
				x := <-c
				y := <-c
				var otherC chan int
				if id == 255 {
					otherC = nat
				} else {
					otherC = machines[id].Input
				}
				otherC <- x
				otherC <- y
			}
		}(m.Output)
	}

	// starting the machines
	for i, m := range machines {
		go m.Run()
		m.AddInput(i)
	}

	// starting the NAT
	go func() {
		var x, y int
		var idleCount int
		var lastSent int
		var sentOnce bool
		for {
			select {
			case x = <-nat:
				y = <-nat
				idleCount = 0
				if !sentOnce {
					// printing part 1 solution
					fmt.Println(y)
					sentOnce = true
				}
			default:
				// idle detection, by looking at channel lengths for
				// X consecutive periods of T microseconds
				idle := true
				for _, m := range machines {
					if len(m.Input) != 0 {
						idle = false
						break
					}
					if len(m.Output) != 0 {
						idle = false
						break
					}
				}
				if idle {
					idleCount++
					time.Sleep(time.Second / (100 * 1000))
				} else {
					idleCount = 0
				}
				if idleCount == 10 {
					idleCount = 0
					if sentOnce && y == lastSent {
						res <- y
						return
					}
					lastSent = y
					c := machines[0].Input
					c <- x
					c <- y
				}
			}
		}
	}()

	// waiting to get final input
	return <-res
}
