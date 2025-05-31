package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Dice struct {
	num   int
	max   int
	rolls int
}

func (d *Dice) Roll() int {
	num := d.num + 1
	d.num = num % d.max
	d.rolls++
	return num
}

type Player struct {
	space int
	score int
}

func (p *Player) TakeTurn(d *Dice) {
	p.space += d.Roll() + d.Roll() + d.Roll()
	p.space %= 10
	if p.space == 0 {
		p.score += 10
	} else {
		p.score += p.space
	}
}

type GameState struct {
	players [2]Player
	turn    int
}

func (g GameState) PlayGame(d *Dice) int {
	p := g.players
	for {
		p[0].TakeTurn(d)
		if p[0].score >= 1000 {
			return p[1].score * d.rolls
		}
		p[1].TakeTurn(d)
		if p[1].score >= 1000 {
			return p[0].score * d.rolls
		}
	}
}

func (g GameState) Next(roll int) GameState {
	out := g
	p := out.players[g.turn]
	p.space = (p.space + roll) % 10
	if p.space == 0 {
		p.score += 10
	} else {
		p.score += p.space
	}
	out.players[g.turn] = p
	out.turn = (g.turn + 1) % 2
	return out
}

type Memo struct {
	cache map[GameState][2]int
	hits  int
}

func (m *Memo) CountWins(s GameState) [2]int {
	if cache, ok := m.cache[s]; ok {
		m.hits++
		return cache
	}
	var total [2]int
	for steps, count := range dirac {
		next := s.Next(steps)
		if next.players[s.turn].score >= 21 {
			total[s.turn] += count
			continue
		}
		wins := m.CountWins(next)
		total[0] += wins[0] * count
		total[1] += wins[1] * count
	}
	m.cache[s] = total
	return total
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) (g GameState) {
	var i int
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		data := strings.Fields(s.Text())
		n := len(data)
		g.players[i].space, err = strconv.Atoi(data[n-1])
		i++
		check(err)
	}
	check(s.Err())
	return g
}

func main() {
	flag.Parse()

	// Make LUT for 3-roll 3-dice total frequency
	dirac = make(map[int]int)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				dirac[i+j+k]++
			}
		}
	}

	g := load(flag.Arg(0))
	fmt.Println("Part 1:", g.PlayGame(&Dice{max: 100}))

	var m Memo
	m.cache = make(map[GameState][2]int)
	res := m.CountWins(g)
	fmt.Println("Cached results:", len(m.cache))
	fmt.Println("Cache hits:", m.hits)
	fmt.Println("Part 2:", max(res[0], res[1]))
}

var dirac map[int]int
