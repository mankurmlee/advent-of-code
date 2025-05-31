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

type Node struct {
	tunnels []string
}

type Valve struct {
	rate    int
}

type CostKey struct {
	a, b string
}

type Puzzle struct {
	nodes        map[string]Node
	valves       map[string]Valve
	cost         map[CostKey]int
	combis       map[string]int
}

type NodeFinder struct {
	pos  string
	cost int
}

type Car struct {
	pos     string
	fuel    int
	score   int
	been    map[string]int
}

type Combi struct {
	valves map[string]bool
	score  int
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
	r := strings.NewReplacer(";", " ", "=", " ", ",", " ")

	p.nodes   = make(map[string]Node)
	p.valves  = make(map[string]Valve)
	p.combis  = make(map[string]int)

	for s.Scan() {
		data := strings.Fields(r.Replace(s.Text()))
		name := data[1]
		rate, err := strconv.Atoi(data[5])
		check(err)
		p.nodes[name] = Node{ data[10:] }
		if rate > 0 {
			p.valves[name] = Valve{ rate }
		}
	}
	check(s.Err())

	p.CalcTravelCosts()

	return p
}

func (p *Puzzle) CalcTravelCosts() {
	p.cost = make(map[CostKey]int)

	v := make([]string, 0, len(p.valves))
	for name := range p.valves {
		v = append(v, name)
	}
	if !slices.Contains(v, "AA") {
		v = append(v, "AA")
	}
	slices.Sort(v)
	n := len(v)
	for i := 0; i < n - 1; i++ {
		for j := i + 1; j < n; j++ {
			a, b := v[i], v[j]
			p.cost[CostKey{a, b}] = p.BFS(a, b)
		}
	}
}

func (p Puzzle) BFS(a, b string) int {
	q := []*NodeFinder{ &NodeFinder{pos: a} }
	been := make(map[string]int)
	been[a] = 0
	for len(q) > 0 {
		f := q[0]
		q[0] = nil
		q = q[1:]
		if f.pos == b {
			return f.cost
		}
		for _, pos := range p.nodes[f.pos].tunnels {
			adj := NodeFinder{ pos, f.cost + 1 }
			if best, ok := been[pos]; ok && best <= adj.cost {
				continue
			}
			been[pos] = adj.cost
			q = append(q, &adj)
		}
	}
	return 99
}

func (p Puzzle) GetCost(a, b string) (cost int) {
	if a < b {
		return p.cost[CostKey{a, b}]
	}
	return p.cost[CostKey{b, a}]
}

func (p Puzzle) GetAdj(f *Car) (adjs []*Car) {
	for k, v := range p.valves {
		if _, ok := f.been[k]; ok {
			continue
		}

		fuel := f.fuel - p.GetCost(f.pos, k) - 1 // 1 to open valve
		if fuel <= 0 {
			continue
		}

		score := f.score + fuel * v.rate

		been := make(map[string]int)
		for key, step := range f.been {
			been[key] = step
		}
		been[k] = len(been)

		adjs = append(adjs, &Car{ k, fuel, score, been })
	}
	return adjs
}

func (c Car) GetKey() string {
	var opened []string
	for k := range c.been {
		opened = append(opened, k)
	}
	slices.Sort(opened)
	return strings.Join(opened, ",")
}

func (p Puzzle) GetMost(initialFuel int) (most int) {
	q := []*Car{ &Car{
		pos: "AA",
		fuel: initialFuel,
		been: make(map[string]int),
	} }
	for len(q) > 0 {
		n := len(q)
		f := q[n-1]
		q[n-1] = nil
		q = q[:n-1]

		key := f.GetKey()
		if best, ok := p.combis[key]; !ok || f.score > best {
			p.combis[key] = f.score
		}

		if f.score > most {
			if initialFuel == 30 {
				fmt.Println(f)
			}
			most = f.score
		}
		q = append(q, p.GetAdj(f)...)
	}
	return most
}

func (p Puzzle) GetCombisByScoreDescending() (combis []Combi) {
	p.combis = make(map[string]int)
	p.GetMost(26)

	for k, score := range p.combis {
		valves := make(map[string]bool)
		for _, v := range strings.Split(k, ",") {
			valves[v] = true
		}
		combis = append(combis, Combi{ valves, score })
	}

	slices.SortFunc(combis, func(a, b Combi) int {
		return -cmp.Compare(a.score, b.score)
	})
	return combis
}

func (p Puzzle) FindBestPair() (pairscore int) {

	combis := p.GetCombisByScoreDescending()
	hiscore := combis[0].score

	n := len(combis)
	for j := 1; j < n; j++ {
		for i := 0; i < j; i++ {
			u, v := combis[i], combis[j]
			if v.score + hiscore <= pairscore {
				return pairscore
			}
			if u.Disjoint(v) {
				if u.score + v.score > pairscore {
					fmt.Println(u, v)
					pairscore = u.score + v.score
				}
			}
		}
	}

	return pairscore
}

func (p Combi) Disjoint(q Combi) bool {
	for k := range q.valves {
		if p.valves[k] {
			return false
		}
	}
	return true
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	most := p.GetMost(30)
	fmt.Println("Part 1:", most)
	fmt.Println("Part 2:", p.FindBestPair())
}
