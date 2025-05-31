package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
)

type Pair struct {
	chip int
	rtg  int
}

type State struct {
	elevator int
	pairs    []Pair
	cost     int
}

func (s0 State) Search() {
	q := []State{s0}
	seen := make(map[int]struct{})
	seen[s0.Key()] = struct{}{}
	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		if s.Finished() {
			fmt.Println(s.cost)
			return
		}
		for _, s1 := range s.Next() {
			h := s1.Key()
			if _, ok := seen[h]; ok {
				continue
			}
			seen[h] = struct{}{}
			q = append(q, s1)
		}
	}
	fmt.Println("No solution found!")
}

func (s State) Next() (out []State) {
	e := s.elevator
	var movable []int
	for i, p := range s.pairs {
		if p.chip == e || p.rtg == e {
			movable = append(movable, i)
		}
	}
	if e < 3 {
		out = append(out, s.MoveTo(movable, e+1)...)
	}
	if e > 0 {
		out = append(out, s.MoveTo(movable, e-1)...)
	}
	return out
}

func (s State) MoveTo(movable []int, e1 int) (out []State) {
	e := s.elevator
	last := len(movable) - 1
	for k, i := range movable {
		p := s.pairs[i]
		if p.chip == e && p.rtg == e {
			s1 := s.CloneNext(e1)
			s1.pairs[i].chip = e1
			s1.pairs[i].rtg = e1
			if s1.Safe() {
				out = append(out, s1)
			}
		}
		if p.chip == e {
			s1 := s.CloneNext(e1)
			s1.pairs[i].chip = e1
			if s1.Safe() {
				out = append(out, s1)
			}
		}
		if p.rtg == e {
			s1 := s.CloneNext(e1)
			s1.pairs[i].rtg = e1
			if s1.Safe() {
				out = append(out, s1)
			}
		}
		if k == last {
			break
		}
		for _, j := range movable[k+1:] {
			if p.chip == e && s.pairs[j].chip == e {
				s1 := s.CloneNext(e1)
				s1.pairs[i].chip = e1
				s1.pairs[j].chip = e1
				if s1.Safe() {
					out = append(out, s1)
				}
			}
			if p.rtg == e && s.pairs[j].rtg == e {
				s1 := s.CloneNext(e1)
				s1.pairs[i].rtg = e1
				s1.pairs[j].rtg = e1
				if s1.Safe() {
					out = append(out, s1)
				}
			}
		}
	}
	return out
}

func (s State) CloneNext(e1 int) State {
	p1 := make([]Pair, len(s.pairs))
	copy(p1, s.pairs)
	return State{e1, p1, s.cost + 1}
}

func (s State) Safe() bool {
	vulnerable := make(map[int]struct{})
	dangerous := make(map[int]struct{})
	for _, p := range s.pairs {
		if p.chip != p.rtg {
			vulnerable[p.chip] = struct{}{}
		}
		dangerous[p.rtg] = struct{}{}
	}
	for f := range vulnerable {
		if _, ok := dangerous[f]; ok {
			return false
		}
	}
	return true
}

func (s State) Key() int {
	out := s.elevator
	for _, p := range s.Sorted() {
		out = (out << 4) | (p.chip << 2) | p.rtg
	}
	return out
}

func (s State) Sorted() []Pair {
	out := make([]Pair, len(s.pairs))
	copy(out, s.pairs)
	slices.SortFunc(out, func(a, b Pair) int {
		c := cmp.Compare(a.chip, b.chip)
		if c == 0 {
			return cmp.Compare(a.rtg, b.rtg)
		}
		return c
	})
	return out
}

func (s State) Finished() bool {
	for _, p := range s.pairs {
		if p.chip != 3 || p.rtg != 3 {
			return false
		}
	}
	return true
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readFile(filename string) (lines []string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	check(s.Err())
	return lines
}

func parseWords(s string) (out []string) {
	re := regexp.MustCompile(`\w+`)
	return re.FindAllString(s, -1)
}

func load(filename string) (out State) {
	elems := make(map[string]Pair)
	for i, l := range readFile(filename) {
		words := parseWords(l)
		for j, w := range words {
			a := w == "microchip"
			if !a && w != "generator" {
				continue
			}
			var k string
			if a {
				k = words[j-2]
			} else {
				k = words[j-1]
			}
			p := elems[k]
			if a {
				p.chip = i
			} else {
				p.rtg = i
			}
			elems[k] = p
		}
	}
	for _, p := range elems {
		out.pairs = append(out.pairs, p)
	}
	return out
}

func main() {
	flag.Parse()
	s := load(flag.Arg(0))
	s.Search()
	s.pairs = append(s.pairs, Pair{0, 0}, Pair{0, 0})
	s.Search()
}
