package main

import (
	"container/list"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/thlacroix/goadvent/helpers"
)

// Deck represents a player deck with all methods necessary for both parts
// It has been implemented first with a chan (not effecient), then with a container/list (better)
type Deck interface {
	Draw() int
	Add(int)
	L() int
	Copy(int) Deck
	Score() int
	fmt.Stringer
}

func main() {
	var part1, part2 int
	var player1, player2 []int
	err := helpers.ScanGroup("input.txt", func(ss []string) error {
		p := make([]int, 0, 25)

		for _, s := range ss[1:] {
			v, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			p = append(p, v)
		}

		if len(player1) == 0 {
			player1 = p
		} else {
			player2 = p
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	part1 = play(NewDeckL(player1), NewDeckL(player2))
	_, part2 = play2(NewDeckL(player1), NewDeckL(player2), true)
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

// =================== DeckL =====================

// DeckL implements Deck with a container/list
type DeckL struct {
	list *list.List
}

func NewDeckL(cards []int) DeckL {
	d := DeckL{list: list.New()}

	for _, c := range cards {
		d.list.PushBack(c)
	}

	return d
}

func (d DeckL) Draw() int {
	c := d.list.Front()
	d.list.Remove(c)
	return c.Value.(int)
}

func (d DeckL) Add(i int) {
	d.list.PushBack(i)
}

func (d DeckL) L() int {
	return d.list.Len()
}

func (d DeckL) Copy(length int) Deck {
	dc := DeckL{list: list.New()}

	e := dc.list.Front()

	for i := 0; i < length; i++ {
		if e == nil {
			break
		}
		dc.list.PushBack(e.Value)

		e = e.Next()
	}

	return dc
}

func (d DeckL) String() string {
	var s strings.Builder

	e := d.list.Front()

	for e != nil {
		s.WriteString(fmt.Sprintf("%v", e.Value))
		e = e.Next()
	}

	return s.String()
}

func (d DeckL) Score() int {
	var score int

	L := d.L()
	e := d.list.Front()

	var i int
	for e != nil {
		score += (L - i) * e.Value.(int)
		i++
		e = e.Next()
	}
	return score
}

// =================== DeckChan =====================

// DeckChan implements Deck with a channel
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
		if i < length {
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

func NewDeckChan(cards []int) DeckChan {
	d := make(DeckChan, len(cards)*3)

	for _, c := range cards {
		d <- c
	}

	return d
}
