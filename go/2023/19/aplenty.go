package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"slices"
)

type Rule struct {
	category    rune
	operator    rune
	value       int
	dest        string
}

type Workflow struct {
	rules []Rule
	dflt  string
}

type Part map[rune]int

type Puzzle struct {
	flows   map[string]Workflow
	parts   []Part
}

type Path struct {
	node   string
	choice int
	flow   *Workflow
	last   *Path
}

type Range struct {
	min, max int
}

type Ratings map[rune]Range

func (r *Ratings) Apply(rule Rule, negate bool) {
	cat := rule.category
	gt  := rule.operator == '>'

	if negate {
		gt = !gt
	}

	val := rule.value
	if !negate {
		if gt {
			val++
		} else {
			val--
		}
	}

	rng, _ := (*r)[cat]
	if gt {
		if val > rng.min {
			rng.min = val
		}
	} else {
		if val < rng.max {
			rng.max = val
		}
	}
	(*r)[cat] = rng
}

func (ctx Puzzle) CountCombis(path Path) (count int) {
	rat := make(Ratings)
	rat['x'] = Range{1, 4000}
	rat['m'] = Range{1, 4000}
	rat['a'] = Range{1, 4000}
	rat['s'] = Range{1, 4000}

	for _, p := range path.Expand() {
		if p.node == "in" {
			continue
		}
		for j := 0; j < p.choice; j++ {
			rat.Apply(p.flow.rules[j], true)
		}
		if p.choice < len(p.flow.rules) {
			rat.Apply(p.flow.rules[p.choice], false)
		}
	}

	rng, _ := rat['x']
	count   = rng.max - rng.min + 1
	rng, _  = rat['m']
	count  *= rng.max - rng.min + 1
	rng, _  = rat['a']
	count  *= rng.max - rng.min + 1
	rng, _  = rat['s']
	count  *= rng.max - rng.min + 1

	return count
}

func (p Path) Expand() (journey []Path) {
	curr := p
	for {
		journey = append(journey, curr)
		if curr.last == nil {
			slices.Reverse(journey)
			return journey
		}
		curr = *curr.last
	}
	return journey
}

func (ctx Puzzle) CountAllCombis(paths []Path) (count int) {
	for _, p := range paths {
		count += ctx.CountCombis(p)
	}
	return count
}

// Path Search
func (p Puzzle) FindAcceptPaths() (apaths []Path) {
	var queue []Path
	queue = append(queue, Path{node: "in"})
	for len(queue) > 0 {
		curr := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if curr.node == "A" {
			apaths = append(apaths, curr)
			continue
		}
		paths := p.Explore(curr)
		queue = append(queue, paths...)
	}
	return apaths
}

func (p Puzzle) Explore(path Path) (paths []Path) {
	flow := p.flows[path.node]
	f := &flow
	for i, r := range f.rules {
		if r.dest == "R" {
			continue
		}
		if path.Seen(r.dest) {
			continue
		}
		paths = append(paths, Path{ r.dest, i, f, &path })
	}
	if !path.Seen(f.dflt) {
		i := len(f.rules)
		paths = append(paths, Path{ f.dflt, i, f, &path })
	}
	return paths
}

func (p Path) Seen(node string) bool {
	for {
		if p.node == node {
			return true
		}
		if p.last == nil {
			return false
		}
		p = *p.last
	}
	return false
}

// Part 1
func (ctx Puzzle) IsAccepted(p Part) bool {
	state := "in"
	for state != "A" && state != "R" {
		f := ctx.flows[state]
		dest := f.dflt
		for _, r := range f.rules {
			if  r.operator == '<' && p[r.category] < r.value ||
				r.operator == '>' && p[r.category] > r.value {
				dest = r.dest
				break
			}
		}
		state = dest
	}
	return state == "A"
}

func (ctx Puzzle) GetAcceptedParts() (accepted []Part) {
	for _, p := range ctx.parts {
		if ctx.IsAccepted(p) {
			accepted = append(accepted, p)
		}
	}
	return accepted
}

func (ctx Puzzle) SumAcceptedParts() (sum int) {
	for _, p := range ctx.GetAcceptedParts() {
		for _, v := range p {
			sum += v
		}
	}
	return sum
}

func parseFlow(txt string) (name string, w Workflow) {
	repl := strings.NewReplacer("{", " ", "}", " ")
	flowData := strings.Fields(repl.Replace(txt))
	name = flowData[0]
	rules := strings.Split(flowData[1], ",")
	nrules := len(rules)
	for i, rule := range rules {
		if i == nrules - 1 {
			w.dflt = rule
			continue
		}
		var r Rule
		ruleData := strings.Split(rule, ":")
		r.category = rune(ruleData[0][0])
		r.operator = rune(ruleData[0][1])
		r.value    = Atoi(ruleData[0][2:])
		r.dest     = ruleData[1]
		w.rules = append(w.rules, r)
	}

	// Remove last rule if dest is same as default
	for i := len(w.rules) - 1; i >= 0; i-- {
		if w.rules[i].dest != w.dflt {
			break
		}
		w.rules = w.rules[:i]
	}

	return name, w
}

func parsePart(txt string) (p Part) {
	p = make(map[rune]int)
	r := strings.NewReplacer("{", "", "}", "", ",", " ")
	for _, cat := range strings.Fields(r.Replace(txt)) {
		catData := strings.Split(cat, "=")
		r := []rune(catData[0])[0]
		v := Atoi(catData[1])
		p[r] = v
	}
	return p
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	// Parse workflows
	p.flows = make(map[string]Workflow)
	for s.Scan() {
		t := s.Text()
		if t == "" {
			break
		}
		name, flow := parseFlow(t)
		p.flows[name] = flow
	}
	check(s.Err())

	// Parse parts
	for s.Scan() {
		part := parsePart(s.Text())
		p.parts = append(p.parts, part)
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

	fmt.Println("Part 1 answer is", p.SumAcceptedParts())

	paths := p.FindAcceptPaths()
	count := p.CountAllCombis(paths)
	fmt.Println("Part 2 answer is", count)
}
