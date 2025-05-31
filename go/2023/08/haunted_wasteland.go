package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
)

type NodeKey struct {
	node      string
	direction rune
}

type Loc struct {
	node string
	pos  int
}

type State struct {
	Loc
	steps int
}

type Puzzle struct {
	route   string
	nodeMap map[NodeKey]string
	next    map[Loc]State
}

func (p *Puzzle) GhostWalk() (steps []int) {
	var states []State

	// Get starting positions
	for k, _ := range p.nodeMap {
		if strings.HasSuffix(k.node, "A") && k.direction == 'R' {
			states = append(states, State{Loc{node: k.node}, 0})
		}
	}

	for _, s := range states {
		steps = append(steps, p.GetNext(s.Loc).steps)
	}

	return steps
}

func (p Puzzle) GetNext(l Loc) (s State) {
	s, ok := p.next[l]
	if ok {
		return s
	}

	curr  := l.node
	start := l.pos
	route := []rune(p.route)
	n := len(route)
	var steps int
	for {
		for i := start; i < n; i++ {
			r := route[i]
			next := p.nodeMap[NodeKey{curr, r}]
			steps++
			if strings.HasSuffix(next, "Z") {
				j := (i + 1) % n
				new := State{ Loc{next, j}, steps }
				p.next[l] = new
				return new
			}
			curr = next
		}
		start = 0
	}
}

func (p Puzzle) CountSteps() (steps int) {
	node := "AAA"
	for {
		for _, r := range p.route {
			next := p.nodeMap[NodeKey{node, r}]
			steps++
			if next == "ZZZ" {
				return steps
			}
			node = next
		}
	}
	return steps
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return (a * b) / gcd(a, b)
}

func calculateLCM(numbers []int) int {
	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		result = lcm(result, numbers[i])
	}
	return result
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	p.route = s.Text()
	s.Scan()
	check(s.Err())
	r := strings.NewReplacer("(", " ", ")", " ", ",", " ")
	p.nodeMap = make(map[NodeKey]string)
	for s.Scan() {
		nodeData := strings.Fields(r.Replace(s.Text()))
		p.nodeMap[NodeKey{nodeData[0], 'L'}] = nodeData[2]
		p.nodeMap[NodeKey{nodeData[0], 'R'}] = nodeData[3]
	}
	check(s.Err())
	p.next = make(map[Loc]State)
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
	steps := p.CountSteps()
	fmt.Println("Answer to part 1 is", steps)

	ghoststeps := p.GhostWalk()
	fmt.Println("Answer to part 2 is", calculateLCM(ghoststeps))
}
