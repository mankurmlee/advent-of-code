package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

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

func parseInts(s string) (nums []int) {
	r := regexp.MustCompile(`\d+`)
	for _, v := range r.FindAllString(s, -1) {
		nums = append(nums, atoi(v))
	}
	return nums
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

func load(s string) (p Puzzle) {
	p.node = make(map[Vec]Info)
	for _, l := range readFile(s)[2:] {
		ints := parseInts(l)
		p.node[Vec{ints[0], ints[1]}] = Info{ints[2], ints[3], ints[4]}
	}
	return p
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	p.PartOne()
	p.PartTwo()
}

type Vec struct {
	x, y int
}

func (v Vec) Equals(o Vec) bool {
	return v.x == o.x && v.y == o.y
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

var DIRS = [4]Vec{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

type Info struct {
	size, used, avail int
}

type Puzzle struct {
	node map[Vec]Info
}

type GoalKey struct {
	goal, gap Vec
}

type State struct {
	goal, gap Vec
	cost      int
}

func (s State) GapGuess() int {
	return s.cost + s.gap.x + s.gap.y
}

func (s State) GoalGuess() int {
	return s.cost + s.goal.x + s.goal.y
}

func (p Puzzle) PartOne() {
	var count int
	for a, adata := range p.node {
		for b, bdata := range p.node {
			if a == b || adata.used == 0 || adata.used > bdata.avail {
				continue
			}
			count++
		}
	}
	fmt.Println(count)
}

func (p Puzzle) PartTwo() {
	seen := make(map[GoalKey]int)
	start := p.GoalPos()
	minsize := p.node[start].used
	dest := Vec{}
	s0 := State{start, p.GapPos(), 0}
	pq := NewPriorityQueue[State]()
	pq.Enqueue(s0, s0.GoalGuess())
	for pq.Len() > 0 {
		s := pq.Dequeue()
		if s.goal.Equals(dest) {
			fmt.Println(s.cost)
			return
		}
		for _, d := range DIRS {
			v1 := s.goal.Add(d)
			if p.node[v1].size < minsize {
				continue
			}
			s1 := p.MoveGoal(s, v1)
			if s1.cost == 0 {
				continue
			}
			gk := GoalKey{s1.goal, s1.gap}
			best, ok := seen[gk]
			if ok && best <= s1.cost {
				continue
			}
			seen[gk] = s1.cost
			pq.Enqueue(s1, s1.GoalGuess())
		}
	}
	fmt.Println("No solution found!")
}

func (p Puzzle) MoveGoal(s0 State, dest Vec) State {
	seen := make(map[Vec]int)
	pq := NewPriorityQueue[State]()
	pq.Enqueue(s0, s0.GapGuess())
	for pq.Len() > 0 {
		s := pq.Dequeue()
		if s.gap.Equals(dest) {
			return State{
				s.gap,
				s.goal,
				s.cost + 1,
			}
		}
		sdata := p.node[s.gap]
		for _, d := range DIRS {
			v1 := s.gap.Add(d)
			if v1.Equals(s0.goal) {
				continue
			}
			vdata := p.node[v1]
			if vdata.size == 0 || sdata.size < vdata.used {
				continue
			}
			s1 := State{s0.goal, v1, s.cost + 1}
			best, ok := seen[v1]
			if ok && best <= s1.cost {
				continue
			}
			seen[v1] = s1.cost
			pq.Enqueue(s1, s1.GapGuess())
		}
	}
	return State{}
}

func (p Puzzle) GapPos() Vec {
	for k, v := range p.node {
		if v.used == 0 {
			return k
		}
	}
	panic("No gap found!")
}

func (p Puzzle) GoalPos() (topright Vec) {
	for pos := range p.node {
		if pos.y != 0 {
			continue
		}
		if pos.x > topright.x {
			topright.x = pos.x
		}
	}
	return topright
}
