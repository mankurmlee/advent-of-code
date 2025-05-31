package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
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

func fileinput() (out []int) {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		out = append(out, atoi(s.Text()))
	}
	check(s.Err())
	return out
}

func total(weights []int) (out int) {
	for _, w := range weights {
		out += w
	}
	return out
}

type CombiKey struct {
	target  int
	weights string
}

func NewCombiKey(target int, weights []int) CombiKey {
	s := make([]string, 0, len(weights))
	for _, w := range weights {
		s = append(s, fmt.Sprint(w))
	}
	return CombiKey{target, strings.Join(s, ",")}
}

type Balancer struct {
	Weights []int
	cache   map[CombiKey][][]int
}

func NewBalancer() Balancer {
	return Balancer{fileinput(), map[CombiKey][][]int{}}
}

// weights should be in ascending order
func (s Balancer) FindCombis(target int, weights []int) (out [][]int) {
	key := NewCombiKey(target, weights)
	if cached, ok := s.cache[key]; ok {
		return cached
	}
	n := len(weights)
	for i, w := range weights {
		rem := target - w
		if rem < 0 {
			break
		} else if rem == 0 {
			out = append(out, []int{w})
			break
		}
		j := i + 1
		if j == n {
			break
		}
		for _, c := range s.FindCombis(rem, weights[j:]) {
			out = append(out, append(c, w))
		}
	}
	s.cache[key] = out
	return out
}

func (s Balancer) SplitGroup(target int, weights []int) (none []int) {
	combis := s.FindCombis(target, weights)
	if len(combis) == 0 {
		return none
	}
	if total(weights) == target+target {
		return combis[0]
	}
	slices.SortFunc(combis, cmpGroup)
	for _, grpOne := range combis {
		rem := removeGroup(weights, grpOne)
		if len(s.SplitGroup(target, rem)) != 0 {
			return grpOne
		}
	}
	return none
}

func (s Balancer) Balance(parts int) int {
	target := total(s.Weights) / parts
	grpOne := s.SplitGroup(target, s.Weights)
	return quantumEntanglement(grpOne)
}

func main() {
	b := NewBalancer()
	fmt.Println("Part 1:", b.Balance(3))
	fmt.Println("Part 2:", b.Balance(4))
}

func quantumEntanglement(weights []int) int {
	prod := 1
	for _, w := range weights {
		prod *= w
	}
	return prod
}

func cmpGroup(a, b []int) int {
	i := len(a)
	j := len(b)
	if i == j {
		i = quantumEntanglement(a)
		j = quantumEntanglement(b)
	}
	if i < j {
		return -1
	}
	if i > j {
		return 1
	}
	return 0
}

func set(list []int) map[int]bool {
	out := map[int]bool{}
	for _, e := range list {
		out[e] = true
	}
	return out
}

func removeGroup(a, b []int) []int {
	s := set(b)
	out := make([]int, 0, len(a))
	for _, e := range a {
		if !s[e] {
			out = append(out, e)
		}
	}
	return out
}
