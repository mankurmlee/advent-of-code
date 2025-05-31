package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"slices"
)

type Card struct {
	wins []int
	nums []int
}

func (c Card) GetMatches() (matches int) {
	for _, n := range c.nums {
		if slices.Contains(c.wins, n) {
			matches++
		}
	}
	return matches
}

func (c Card) GetWorth() (worth int) {
	for _, n := range c.nums {
		if slices.Contains(c.wins, n) {
			if worth > 0 {
				worth *= 2
			} else {
				worth = 1
			}
		}
	}
	return worth
}

type Puzzle struct {
	cards []Card
}

func (p Puzzle) CountScratches() (scratches int) {
	copies := make([]int, len(p.cards))
	for i, c := range p.cards {
		fmt.Printf("Card %d has %d instances\n", i + 1, copies[i] + 1)
		scratches += copies[i] + 1
		matches := c.GetMatches()
		for j := i + 1; j < i + 1 + matches; j++ {
			copies[j] += copies[i] + 1
		}
	}

	return scratches
}

func (p Puzzle) GetWorth() (total int) {
	for i, c := range p.cards {
		worth := c.GetWorth()
		fmt.Println("Card", i+1, "is worth", worth)
		total += worth
	}
	return total
}

func parseCard(txt string) (c Card) {
	cardData := strings.Split(txt, ":")
	numData := strings.Split(cardData[1], "|")

	for _, w := range strings.Fields(numData[0]) {
		c.wins = append(c.wins, Atoi(w))
	}
	for _, n := range strings.Fields(numData[1]) {
		c.nums = append(c.nums, Atoi(n))
	}
	return c
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	s := bufio.NewScanner(f)
	for s.Scan() {
		card := parseCard(s.Text())
		p.cards = append(p.cards, card)
	}
	check(s.Err())
	return p
}

func Atoi(A string) int {
	i, err := strconv.Atoi(A)
	check(err)
	return i
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	//~ worth := p.GetWorth()
	//~ fmt.Println("Answer to part 1:", worth)
	scratches := p.CountScratches()
	fmt.Println("Answer to part 2:", scratches)
}
