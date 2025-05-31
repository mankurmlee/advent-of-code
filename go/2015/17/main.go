package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
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

func parse(inp []string) (target int, containers []int) {
	for _, c := range inp[1:] {
		containers = append(containers, atoi(c))
	}
	return atoi(inp[0]), containers
}

func hash(target int, stack []int) string {
	return fmt.Sprint(target) + fmt.Sprint(stack)
}

type Solver struct {
	Jugs []int
	Memo map[string][][]int
}

func (s Solver) CountCombis(target int, opts []int) (out [][]int) {
	key := hash(target, opts)
	if cached, exists := s.Memo[key]; exists {
		return cached
	}
	n := len(opts)
	for opt_i, jug_i := range opts {
		v := s.Jugs[jug_i]
		if v == target {
			out = append(out, []int{jug_i})
		} else if v < target && opt_i+1 < n {
			for _, combis := range s.CountCombis(target-v, opts[opt_i+1:]) {
				out = append(out, append([]int{jug_i}, combis...))
			}
		}
	}
	s.Memo[key] = out
	return out
}

func (s Solver) GetCombis(target int) (out [][]int) {
	n := len(s.Jugs)
	jugs := make([]int, n)
	for i := range n {
		jugs[i] = i
	}
	return s.CountCombis(target, jugs)
}

func shortest(combis [][]int) int {
	out := 1 << 32
	for _, c := range combis {
		n := len(c)
		if n < out {
			out = n
		}
	}
	return out
}

func countShortest(combis [][]int) (out int) {
	tar := shortest(combis)
	for _, c := range combis {
		if len(c) == tar {
			out++
		}
	}
	return out
}

func main() {
	target, jugs := parse(fileinput())
	s := Solver{jugs, make(map[string][][]int)}
	combis := s.GetCombis(target)
	fmt.Println("Part 1:", len(combis))
	fmt.Println("Part 2:", countShortest(combis))
}
