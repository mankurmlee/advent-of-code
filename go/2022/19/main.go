package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"strconv"
	"runtime"
)

const (
	ORE      = 0
	CLAY     = 1
	OBSIDIAN = 2
	GEODE    = 3
)

type Items struct {
	ore         int
	clay        int
	obsidian    int
	geode       int
}

type Blueprint struct {
	id          int
	costs       []Items
	maxspend    Items
}

type State struct {
	wallet  Items
	robots  Items
	fuel    int
	order   int
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func load(filename string) (p []Blueprint) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	re := regexp.MustCompile(`[^0-9]`)
	for s.Scan() {
		data := strings.Fields(re.ReplaceAllString(s.Text(), " "))
		b := Blueprint {
			id   : atoi(data[0]),
			costs: []Items{
				Items{ ore: atoi(data[1]) },
				Items{ ore: atoi(data[2]) },
				Items{ ore: atoi(data[3]), clay: atoi(data[4]) },
				Items{ ore: atoi(data[5]), obsidian: atoi(data[6]) },
			},
			maxspend: Items {
				ore : max(atoi(data[1]), atoi(data[2]), atoi(data[3]), atoi(data[5])),
				clay: atoi(data[4]),
				obsidian: atoi(data[6]),
			},
		}
		p = append(p, b)
	}
	check(s.Err())
	return p
}

func (bp Blueprint) MostGeodes(t int) (most int) {
	s1 := State{
		robots: Items{ ore: 1 },
		fuel  : t,
		order : ORE,
	}
	s2 := State{
		robots: Items{ ore: 1 },
		fuel  : t,
		order : CLAY,
	}
	q := []State{ s1, s2 }

	for len(q) > 0 {
		n := len(q)
		s := q[n-1]
		q = q[:n-1]
		if s.wallet.geode > most {
			//~ fmt.Println("BP", bp.id, ":", s)
			most = s.wallet.geode
		}
		if s.wallet.geode + s.fuel * s.robots.geode + ((s.fuel - 1) + (s.fuel-1)*(s.fuel-1)) / 2 > most {
			q = append(q, bp.NextStates(s)...)
		}
	}

	return most
}

func (bp Blueprint) NextStates(s State) (next []State) {
	if s.fuel <= 0 {
		return next
	}
	robot := s.order
	affords := true
	costs := bp.costs[robot]
	for s.wallet.ore < costs.ore || s.wallet.clay < costs.clay || s.wallet.obsidian < costs.obsidian {
		if s.fuel <= 1 {
			affords = false
			break
		}
		s.wallet.ore      += s.robots.ore
		s.wallet.clay     += s.robots.clay
		s.wallet.obsidian += s.robots.obsidian
		s.wallet.geode    += s.robots.geode
		s.fuel--
	}

	// spend resources
	if affords {
		s.wallet.ore      -= costs.ore
		s.wallet.clay     -= costs.clay
		s.wallet.obsidian -= costs.obsidian
	}

	// collect resources
	s.wallet.ore      += s.robots.ore
	s.wallet.clay     += s.robots.clay
	s.wallet.obsidian += s.robots.obsidian
	s.wallet.geode    += s.robots.geode

	// add robot
	if affords {
		switch robot {
		case ORE:
			s.robots.ore++
		case CLAY:
			s.robots.clay++
		case OBSIDIAN:
			s.robots.obsidian++
		case GEODE:
			s.robots.geode++
		}
	}

	// decrease remaining fuel
	s.fuel--

	if s.fuel > 0 {
		if s.robots.ore < bp.maxspend.ore {
			next = append(next, State{ s.wallet, s.robots, s.fuel, ORE })
		}
		if s.robots.clay < bp.maxspend.clay {
			next = append(next, State{ s.wallet, s.robots, s.fuel, CLAY })
		}
		if s.robots.clay > 0 && s.robots.obsidian < bp.maxspend.obsidian {
			next = append(next, State{ s.wallet, s.robots, s.fuel, OBSIDIAN })
		}
	}
	if s.robots.obsidian > 0 {
		next = append(next, State{ s.wallet, s.robots, s.fuel, GEODE })
	}

	return next
}

func worker(bps []Blueprint, fuelleft int) chan []int {
	ch := make(chan []int)
	go func () {
		defer close(ch)
		var results []int
		for _, bp := range bps {
			results = append(results, bp.MostGeodes(fuelleft))
		}
		ch <- results
	}()
	return ch
}

func process(bps []Blueprint, fuelleft int) (most []int) {
	ntasks := len(bps)
	ncpus := runtime.NumCPU()
	if ntasks < ncpus {
		ncpus = ntasks
	}
	batch := ntasks / ncpus

	var results []chan []int
	for i := 0; i < ncpus; i++ {
		end := (i + 1) * batch
		if i == ncpus - 1 {
			end = ntasks
		}
		results = append(results, worker(bps[i*batch:end], fuelleft))
	}
	for i := 0; i < ncpus; i++ {
		res := <- results[i]
		most = append(most, res...)
	}
	return most
}

func main() {
	flag.Parse()
	bps := load(flag.Arg(0))

	var sum int
	for i, v := range process(bps, 24) {
		sum += bps[i].id * v
	}
	fmt.Println("Part 1:", sum)

	prod := 1
	for _, v := range process(bps[:min(3,len(bps))], 32) {
		prod *= v
	}
	fmt.Println("Part 2:", prod)
}
