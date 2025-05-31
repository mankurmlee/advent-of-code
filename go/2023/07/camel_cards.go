package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"slices"
	"cmp"
)

type Hand string

type Play struct {
	hand    Hand
	bid     int
}

type Puzzle struct {
	plays []Play
}

func (ctx Puzzle) GetWinnings() (total int) {
	for i, p := range ctx.plays {
		total += (i + 1) * p.bid
	}
	return total
}

func (h Hand) TypeOrder(jokers bool) int {
	hm := make(map[rune]int)
	for _, r := range []rune(h) {
		hm[r]++
	}

	var counts []int
	var jokerCount int
	for k, v := range hm {
		if jokers && k == 'J' {
			jokerCount = v
			continue
		}
		counts = append(counts, v)
	}
	if len(counts) == 0 {
		counts = []int{0}
	}

	slices.Sort(counts)
	slices.Reverse(counts)

	if jokers {
		counts[0] += jokerCount
	}

	if counts[0] == 5 {
		return 6
	}

	if counts[0] == 4 {
		return 5
	}

	if counts[0] == 3 {
		if counts[1] == 2 {
			return 4
		}
		return 3
	}

	if counts[0] == 2 {
		if counts[1] == 2 {
			return 2
		}
		return 1
	}

	return 0
}

func cardOrder(r rune, jokers bool) int {
	switch r {
	case 'T':
		return 10
	case 'J':
		if jokers {
			return 1
		}
		return 11
	case 'Q':
		return 12
	case 'K':
		return 13
	case 'A':
		return 14
	}

	i, err := strconv.Atoi(string(r))
	check(err)
	return i
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	s := bufio.NewScanner(f)
	for s.Scan() {
		playData := strings.Fields(s.Text())
		bid, err := strconv.Atoi(playData[1])
		check(err)
		p.plays = append(p.plays, Play{Hand(playData[0]), bid})
	}
	check(s.Err())
	return p
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	// Sort the plays
	slices.SortFunc(p.plays, func(a, b Play) int {
		if n := cmp.Compare(a.hand.TypeOrder(false), b.hand.TypeOrder(false)); n != 0 {
			return n
		}
		handA := []rune(a.hand)
		handB := []rune(b.hand)
		for i, u := range handA {
			v := handB[i]
			n := cmp.Compare(cardOrder(u, false), cardOrder(v, false))
			if n != 0 {
				return n
			}
		}
		return 0
	})
	answer1 := p.GetWinnings()

	// Sort the plays with joker rules
	slices.SortFunc(p.plays, func(a, b Play) int {
		if n := cmp.Compare(a.hand.TypeOrder(true), b.hand.TypeOrder(true)); n != 0 {
			return n
		}
		handA := []rune(a.hand)
		handB := []rune(b.hand)
		for i, u := range handA {
			v := handB[i]
			n := cmp.Compare(cardOrder(u, true), cardOrder(v, true))
			if n != 0 {
				return n
			}
		}
		return 0
	})
	answer2 := p.GetWinnings()

	for _, p := range p.plays {
		fmt.Println(p.hand, "type ordering is", p.hand.TypeOrder(true))
	}

	fmt.Println("Part 1 answer is", answer1)
	fmt.Println("Part 2 answer is", answer2)
}
