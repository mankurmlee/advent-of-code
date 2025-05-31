package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strconv"
	"slices"
	"strings"
)

type Stack []byte

type Move struct {
	num  int
	src  int
	dest int
}

type Puzzle struct {
	stacks []Stack
	moves  []Move
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	p.stacks = loadStacks(s)

	s.Scan(); check(s.Err()) // read and discard empty line

	p.moves = loadMoves(s)

	return p
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func loadStacks(s *bufio.Scanner) (stacks []Stack) {
	for s.Scan() {
		line := s.Text()
		if line[1] == '1' {
			break
		}
		n := len(line)

		nStacks := (n + 3) / 4
		for len(stacks) < nStacks {
			stacks = append(stacks, []byte{})
		}

		for i := 0; i < n; i += 4 {
			crate := line[i+1]
			if crate == ' ' {
				continue
			}
			j := i / 4
			stacks[j] = append(stacks[j], crate)
		}
	}
	check(s.Err())
	for i, _ := range stacks {
		slices.Reverse(stacks[i])
	}
	return stacks
}

func loadMoves(s *bufio.Scanner) (moves []Move) {
	for s.Scan() {
		data := strings.Fields(s.Text())
		moves = append(moves, Move{
			Atoi(data[1]),
			Atoi(data[3]),
			Atoi(data[5]),
		})
	}
	check(s.Err())
	return moves
}

func (p Puzzle) GetTopCrates() string {
	var top strings.Builder
	for _, s := range p.stacks {
		n := len(s)
		top.WriteByte(s[n-1])
	}
	return top.String()
}

func (p Puzzle) Clone() (q Puzzle) {
	q.stacks = make([]Stack, len(p.stacks))
	for i, s := range p.stacks {
		q.stacks[i] = make([]byte, len(s))
		copy(q.stacks[i], p.stacks[i])
	}
	q.moves = make([]Move, len(p.moves))
	copy(q.moves, p.moves)
	return q
}

func (p Puzzle) Rearrange() {
	for _, m := range p.moves {
		for i := 0; i < m.num; i++ {
			c := p.stacks[m.src - 1].Pop()
			p.stacks[m.dest - 1].Push(c)
		}
	}
}

func (s *Stack) Pop() (item byte) {
	n := len(*s) - 1
	item = (*s)[n]
	*s = (*s)[:n]
	return item
}

func (s *Stack) Push(item byte) {
	*s = append(*s, item)
}

func (p Puzzle) NewRearrange() {
	for _, m := range p.moves {
		var mini Stack
		for i := 0; i < m.num; i++ {
			c := p.stacks[m.src - 1].Pop()
			mini.Push(c)
		}
		for i := 0; i < m.num; i++ {
			c := mini.Pop()
			p.stacks[m.dest - 1].Push(c)
		}
	}
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	q := p.Clone()
	q.Rearrange()
	fmt.Println("Part 1:", q.GetTopCrates())

	r := p.Clone()
	r.NewRearrange()
	fmt.Println("Part 2:", r.GetTopCrates())
}
