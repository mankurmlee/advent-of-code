package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var DIR = map[Vec]rune{
	{0, -1}: '^',
	{1, 0}:  '>',
	{0, 1}:  'v',
	{-1, 0}: '<',
}
var NUMPAD = map[Vec]rune{
	{0, 0}: '7', {1, 0}: '8', {2, 0}: '9',
	{0, 1}: '4', {1, 1}: '5', {2, 1}: '6',
	{0, 2}: '1', {1, 2}: '2', {2, 2}: '3',
	{1, 3}: '0', {2, 3}: 'A',
}
var DIRPAD = map[Vec]rune{
	{1, 0}: '^', {2, 0}: 'A',
	{0, 1}: '<', {1, 1}: 'v', {2, 1}: '>',
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	check(err)
	return i
}

func val(a string) int {
	re := regexp.MustCompile(`[1-9]\d*`)
	return atoi(re.FindString(a))
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

func main() {
	flag.Parse()
	p := Puzzle{readFile(flag.Arg(0))}
	kp := Keypad{make(map[Route]int)}
	for i := range 25 {
		kp = NewKeypad(DIRPAD, kp)
		if i == 1 {
			p.Checksum(NewKeypad(NUMPAD, kp))
		}
	}
	p.Checksum(NewKeypad(NUMPAD, kp))
}

type Vec struct {
	x, y int
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

func (v Vec) Equals(o Vec) bool {
	return v.x == o.x && v.y == o.y
}

type State struct {
	Vec
	cost int
	last rune
}

func findRoute(
	grid map[Vec]rune,
	start, end Vec,
	pad Keypad,
) (best int) {
	pq := NewPriorityQueue[State]()
	pq.Enqueue(State{start, 0, 'A'}, 0)
	for len(pq) > 0 {
		s := pq.Dequeue()
		if s.Equals(end) {
			c := pad.RouteCost(Route{s.last, 'A'})
			if best == 0 || s.cost+c < best {
				best = s.cost + c
			}
			continue
		}
		for v, d := range DIR {
			v1 := s.Add(v)
			if _, ok := grid[v1]; !ok {
				continue
			}
			c := pad.RouteCost(Route{s.last, d})
			s1 := State{
				v1,
				s.cost + c,
				d,
			}
			if best > 0 && s1.cost > best {
				continue
			}
			pq.Enqueue(s1, s1.cost)
		}
	}
	return best
}

type Route struct {
	a, b rune
}

type Puzzle struct {
	codes []string
}

func (p Puzzle) Checksum(pad Keypad) {
	sum := 0
	for _, c := range p.codes {
		sum += val(c) * pad.Cost(c)
	}
	fmt.Println(sum)
}

type Keypad struct {
	dict map[Route]int
}

func (k Keypad) RouteCost(route Route) int {
	if b, ok := k.dict[route]; ok {
		return b
	}
	return 1
}

func NewKeypad(layout map[Vec]rune, control Keypad) Keypad {
	out := make(map[Route]int)
	buttons := make(map[rune]Vec)

	for k, v := range layout {
		buttons[v] = k
	}

	for a := range buttons {
		for b := range buttons {
			if a == b {
				continue
			}
			out[Route{a, b}] = findRoute(
				layout,
				buttons[a],
				buttons[b],
				control,
			)
		}
	}

	return Keypad{out}
}

func (k Keypad) Cost(code string) (cost int) {
	o := 'A'
	for _, ch := range code {
		cost += k.dict[Route{o, ch}]
		o = ch
	}
	return cost
}
