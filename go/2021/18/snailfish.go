package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"encoding/json"
	"strconv"
)

type Snailfish struct {
	isLeaf bool
	value  int
	pair   [2]*Snailfish
}

type Puzzle []Snailfish

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func unserialise(input string) Snailfish {
	var j []interface{}
	err := json.Unmarshal([]byte(input), &j)
	check(err)
	return parseSnailfish(j)
}

func parseSnailfish(o interface{}) Snailfish {
	switch value := o.(type) {
	case float64:
		return Snailfish{ isLeaf: true, value: int(value) }
	case []interface{}:
		var pair [2]*Snailfish
		for i, o := range value {
			sf := parseSnailfish(o)
			pair[i] = &sf
		}
		return Snailfish{ isLeaf: false, pair: pair }
	default:
		fmt.Println("Unexpected object:", o)
		return Snailfish{}
	}
}

func (s Snailfish) String() string {
	if s.isLeaf {
		return strconv.Itoa(s.value)
	}
	p := s.pair[0].String()
	q := s.pair[1].String()
	return fmt.Sprintf("[%s %s]", p, q)
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		p = append(p, unserialise(s.Text()))
	}
	check(s.Err())
	return p
}

func (s Snailfish) Magnitude() int {
	if s.isLeaf {
		return s.value
	}
	p := s.pair[0].Magnitude()
	q := s.pair[1].Magnitude()
	return 3 * p + 2 * q
}

func (a Snailfish) Add(b Snailfish) (c Snailfish) {
	c = Snailfish{ isLeaf: false, pair: [2]*Snailfish{&a, &b} }
	for {
		if _, ok := c.Explode(0); ok {
			continue
		}
		if c.Split() {
			continue
		}
		break
	}
	return c
}

func (s *Snailfish) Explode(depth int) (v [2]int, ok bool) {
	if s.isLeaf {
		return [2]int{}, false
	}

	v, ok = s.pair[0].Explode(depth + 1)
	if ok {
		if v[1] > 0 {
			s.pair[1].LeftAdd(v[1])
			v[1] = 0
		}
		return v, true
	}

	v, ok = s.pair[1].Explode(depth + 1)
	if ok {
		if v[0] > 0 {
			s.pair[0].RightAdd(v[0])
			v[0] = 0
		}
		return v, true
	}

	if depth >= 4 {
		s.isLeaf = true
		s.value = 0
		return [2]int{s.pair[0].value, s.pair[1].value}, true
	}

	return [2]int{}, false
}

func (s *Snailfish) LeftAdd(n int) {
	if s.isLeaf {
		s.value += n
	} else {
		s.pair[0].LeftAdd(n)
	}
}

func (s *Snailfish) RightAdd(n int) {
	if s.isLeaf {
		s.value += n
	} else {
		s.pair[1].RightAdd(n)
	}
}

func (s *Snailfish) Split() bool {
	if s.isLeaf {
		if s.value < 10 {
			return false
		}
		s.isLeaf = false
		i := s.value / 2
		j := s.value - i
		s.pair[0] = &Snailfish{ isLeaf: true, value: i }
		s.pair[1] = &Snailfish{ isLeaf: true, value: j }
		return true
	}
	if s.pair[0].Split() {
		return true
	}
	if s.pair[1].Split() {
		return true
	}
	return false
}

func (p Puzzle) FinalSum() (out Snailfish) {
	out = p[0]
	for _, num := range p[1:] {
		out = out.Add(num)
	}
	return out
}

func (s Snailfish) Clone() (o Snailfish) {
	o = s
	if o.isLeaf {
		return o
	}
	p := s.pair[0].Clone()
	q := s.pair[1].Clone()
	o.pair[0] = &p
	o.pair[1] = &q
	return o
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	for _, n := range p {
		fmt.Println(n)
	}
	sum := p.FinalSum()
	fmt.Println("Final sum:", sum)
	fmt.Println("Part 1:", sum.Magnitude())

	p = load(flag.Arg(0))
	var (
		best int
		best_i, best_j Snailfish
	)
	for _, j := range p {
		for _, i := range p {
			if i == j { continue }
			u, v := i.Clone(), j.Clone()
			m := u.Add(v).Magnitude()
			if m > best {
				best = m
				best_i = i
				best_j = j
			}
		}
	}
	fmt.Println(best_i, "+")
	fmt.Println(best_j, "=")
	fmt.Println("Part 2:", best)
}
