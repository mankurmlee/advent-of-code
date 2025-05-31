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

func fileinput() (out []string) {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		out = append(out, s.Text())
	}
	check(s.Err())
	return out
}

type GameState struct {
	ManaSpent     int
	BossHitPoints int
	BossAttack    int
	HitPoints     int
	Mana          int
	Shield        int
	Poison        int
	Recharge      int
}

func (s GameState) BuffTurn() GameState {
	out := s
	if out.Shield > 0 {
		out.Shield--
	}
	if out.Poison > 0 {
		out.BossHitPoints -= 3
		out.Poison--
	}
	if out.Recharge > 0 {
		out.Mana += 101
		out.Recharge--
	}
	return out
}

func (s GameState) BossTurn() GameState {
	out := s
	out = out.BuffTurn()
	if out.BossHitPoints <= 0 {
		return out
	}
	atk := out.BossAttack
	if out.Shield > 0 {
		atk -= 7
		if atk < 1 {
			atk = 1
		}
	}
	out.HitPoints -= atk
	return out
}

func (s GameState) CastRecharge() GameState {
	out := s
	out.Recharge = 5
	out.Mana -= 229
	out.ManaSpent += 229
	return out
}

func (s GameState) CastPoison() GameState {
	out := s
	out.Poison = 6
	out.Mana -= 173
	out.ManaSpent += 173
	return out
}

func (s GameState) CastShield() GameState {
	out := s
	out.Shield = 6
	out.Mana -= 113
	out.ManaSpent += 113
	return out
}

func (s GameState) CastDrain() GameState {
	out := s
	out.BossHitPoints -= 2
	out.HitPoints += 2
	out.Mana -= 73
	out.ManaSpent += 73
	return out
}

func (s GameState) CastMissile() GameState {
	out := s
	out.BossHitPoints -= 4
	out.Mana -= 53
	out.ManaSpent += 53
	return out
}

type Game struct {
	BossHitPoints int
	BossAttack    int
}

func NewGame() Game {
	lines := fileinput()
	hpData := strings.Fields(lines[0])
	atkData := strings.Fields(lines[1])
	ehp := atoi(hpData[len(hpData)-1])
	eatk := atoi(atkData[len(atkData)-1])
	return Game{ehp, eatk}
}

func (g Game) Play() int {
	seen := map[GameState]bool{}
	pq := NewPriorityQueue[GameState]()
	s := g.InitialState()
	pq.Enqueue(s, s.ManaSpent)
	for pq.Len() > 0 {
		s = pq.Dequeue()
		if s.BossHitPoints <= 0 {
			break
		}
		for _, n := range g.NextStates(s) {
			if seen[n] {
				continue
			}
			seen[n] = true
			pq.Enqueue(n, n.ManaSpent)
		}
	}
	return s.ManaSpent
}

func (g Game) InitialState() GameState {
	var s GameState
	s.BossHitPoints = g.BossHitPoints
	s.BossAttack = g.BossAttack
	s.HitPoints = 50
	s.Mana = 500
	return s
}

func (g Game) NextStates(s GameState) (out []GameState) {
	s = s.BuffTurn()
	if s.BossHitPoints <= 0 {
		return append(out, s)
	}
	if s.Mana >= 229 && s.Recharge == 0 {
		n := s.CastRecharge()
		n = n.BossTurn()
		if n.HitPoints > 0 {
			out = append(out, n)
		}
	}
	if s.Mana >= 173 && s.Poison == 0 {
		n := s.CastPoison()
		n = n.BossTurn()
		if n.HitPoints > 0 {
			out = append(out, n)
		}
	}
	if s.Mana >= 113 && s.Shield == 0 {
		n := s.CastShield()
		n = n.BossTurn()
		if n.HitPoints > 0 {
			out = append(out, n)
		}
	}
	if s.Mana >= 73 {
		n := s.CastDrain()
		if n.BossHitPoints <= 0 {
			out = append(out, n)
		} else {
			n = n.BossTurn()
			if n.HitPoints > 0 {
				out = append(out, n)
			}
		}
	}
	if s.Mana >= 53 {
		n := s.CastMissile()
		if n.BossHitPoints <= 0 {
			out = append(out, n)
		} else {
			n = n.BossTurn()
			if n.HitPoints > 0 {
				out = append(out, n)
			}
		}
	}
	return out
}

func main() {
	g := NewGame()
	fmt.Println("Part 1:", g.Play())
	g.BossAttack++
	fmt.Println("Part 2:", g.Play())
}
