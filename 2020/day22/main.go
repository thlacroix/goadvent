package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

type Deck interface {
	Draw() int
	Add(int)
	L() int
	Copy(int) Deck
	Score() int
	fmt.Stringer
}

type DeckChan chan int

func (d DeckChan) Draw() int {
	return <-d
}

func (d DeckChan) Add(i int) {
	d <- i
}

func (d DeckChan) L() int {
	return len(d)
}

func (d DeckChan) Copy(length int) Deck {
	dc := make(DeckChan, cap(d))

	d <- -1
	var i int
	for {
		v := <-d

		if v == -1 {
			break
		}

		d <- v
		if length < 0 || i < length {
			dc <- v
		}
		i++
	}

	return dc
}

func (d DeckChan) String() string {
	var s strings.Builder

	d <- -1

	for {
		v := <-d

		if v == -1 {
			break
		}

		s.WriteString(fmt.Sprintf("%d+", v))
		d <- v
	}

	return s.String()
}

func (d DeckChan) Score() int {
	var score int

	close(d)

	L := len(d)
	var i int

	for v := range d {
		score += (L - i) * v
		i++
	}
	return score
}

func main() {
	var part1, part2 int
	player1, player2 := make([]int, 0, 25), make([]int, 0, 25)
	p1, p2 := make(DeckChan, 100), make(DeckChan, 100)
	err := helpers.ScanGroup("input.txt", func(ss []string) error {
		var p []int
		var c chan int
		if len(p1) == 0 {
			p = player1
			c = p1
		} else {
			p = player2
			c = p2
		}

		for _, s := range ss[1:] {
			v, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			p = append(p, v)
			c <- v
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	part1 = play(p1.Copy(-1), p2.Copy(-1))
	_, part2 = play2(p1.Copy(-1), p2.Copy(-1), true)
	fmt.Println(part1, part2)
}

func play(p1, p2 Deck) int {
	var winner Deck
playLoop:
	for {
		var v1, v2 int

		if p1.L() == 0 {
			winner = p2
			break playLoop
		}

		if p2.L() == 0 {
			winner = p1
			break playLoop
		}

		v1 = p1.Draw()
		v2 = p2.Draw()

		if v1 > v2 {
			p1.Add(v1)
			p1.Add(v2)
		}

		if v2 > v1 {
			p2.Add(v2)
			p2.Add(v1)
		}
	}

	return winner.Score()
}

type Player bool

const (
	P1 Player = true
	P2 Player = false
)

func play2(p1, p2 Deck, count bool) (Player, int) {
	var winner Player
	history := make(map[string]bool, p1.L()+p2.L())
playLoop:
	for {
		var v1, v2 int
		repr := fmt.Sprintf("%s|%s", p1, p2)
		if history[repr] {
			winner = P1
			break
		}

		history[repr] = true

		if p1.L() == 0 {
			winner = P2
			break playLoop
		}

		if p2.L() == 0 {
			winner = P1
			break playLoop
		}

		v1 = p1.Draw()
		v2 = p2.Draw()

		var roundWinner Player

		if v1 <= p1.L() && v2 <= p2.L() {
			p1c := p1.Copy(v1)
			p2c := p2.Copy(v2)
			roundWinner, _ = play2(p1c, p2c, false)
		} else {
			roundWinner = Player(v1 > v2)
		}

		if roundWinner == P1 {
			p1.Add(v1)
			p1.Add(v2)
		}

		if roundWinner == P2 {
			p2.Add(v2)
			p2.Add(v1)
		}
	}

	if !count {
		return winner, 0
	}

	winnerDeck := p1
	if winner == P2 {
		winnerDeck = p2
	}

	return winner, winnerDeck.Score()

}
