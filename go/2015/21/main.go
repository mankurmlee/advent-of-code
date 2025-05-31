package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func fileinput() (out [][]string) {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	var chunk []string
	for s.Scan() {
		t := strings.TrimSpace(s.Text())
		if t != "" {
			chunk = append(chunk, t)
		} else {
			out = append(out, chunk)
			chunk = []string{}
		}
	}
	check(s.Err())
	out = append(out, chunk)
	return out
}

type Equipment struct {
	Slot   string
	Name   string
	Cost   int
	Damage int
	Armor  int
}

type Shop struct {
	wares map[string][]Equipment
}

func newShop(data [][]string) Shop {
	s := Shop{map[string][]Equipment{}}
	s.Parse(data[0])
	s.Parse(data[1])
	s.Parse(data[2])
	return s
}

func (s *Shop) Parse(lines []string) {
	var eq []Equipment
	r := strings.NewReplacer(":", " ")
	slot := strings.Fields(r.Replace(lines[0]))[0]
	if slot != "Weapons" {
		eq = append(eq, Equipment{slot, "Unequipped", 0, 0, 0})
	}
	for _, l := range lines[1:] {
		fields := strings.Fields(l)
		n := len(fields)
		name := strings.Join(fields[:n-3], " ")
		cost := atoi(fields[n-3])
		damage := atoi(fields[n-2])
		armor := atoi(fields[n-1])
		eq = append(eq, Equipment{slot, name, cost, damage, armor})
	}
	s.wares[slot] = eq
}

type Actor struct {
	HitPoints int
	Damage    int
	Armor     int
}

func (a Actor) Attack(b *Actor) {
	b.HitPoints -= a.Damage - b.Armor
	if b.HitPoints < 0 {
		b.HitPoints = 0
	}
}

func (a Actor) Dead() bool {
	return a.HitPoints <= 0
}

type Loadout []Equipment

func (l Loadout) Cost() (c int) {
	for _, e := range l {
		c += e.Cost
	}
	return c
}

type Simulator struct {
	shop   Shop
	player Actor
	boss   Actor
}

func newSim(data [][]string) Simulator {
	var s Simulator
	s.shop = newShop(data)
	pdata := data[3]
	s.player.HitPoints = s.Parse(pdata[0])
	s.boss.HitPoints = s.Parse(pdata[1])
	s.boss.Damage = s.Parse(pdata[2])
	s.boss.Armor = s.Parse(pdata[3])
	return s
}

func (s Simulator) Parse(line string) int {
	f := strings.Fields(line)
	return atoi(f[1])
}

func (s Simulator) OrderLoadouts(reverse bool) PriorityQueue[Loadout] {
	pq := NewPriorityQueue[Loadout]()
	rings := s.shop.wares["Rings"]
	n := len(rings)
	for _, w := range s.shop.wares["Weapons"] {
		for _, a := range s.shop.wares["Armor"] {
			for i := range n - 1 {
				for j := i; j < n; j++ {
					if i != 0 && i == j {
						continue
					}
					l := Loadout{w, a, rings[i], rings[j]}
					c := l.Cost()
					if reverse {
						c = -c
					}
					pq.Enqueue(l, c)
				}
			}
		}
	}
	return pq
}

func (s Simulator) Simulate(reverse bool) int {
	pq := s.OrderLoadouts(reverse)
	for pq.Len() > 0 {
		l := pq.Dequeue()
		if s.CanWin(l) != reverse {
			fmt.Println(l)
			return l.Cost()
		}
	}
	return 0
}

func (s Simulator) CanWin(l Loadout) bool {
	p := s.player
	b := s.boss
	for _, e := range l {
		p.Damage += e.Damage
		p.Armor += e.Armor
	}
	for {
		p.Attack(&b)
		if b.Dead() {
			return true
		}
		b.Attack(&p)
		if p.Dead() {
			return false
		}
	}
}

func main() {
	s := newSim(fileinput())
	fmt.Println("Part 1:", s.Simulate(false))
	fmt.Println("Part 2:", s.Simulate(true))
}
