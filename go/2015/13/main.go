package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Relation struct {
	happiness map[string]int
	head      string
	memo      map[string]int
}

func (r Relation) Diners() (out []string) {
	diners := map[string]struct{}{}
	for k := range r.happiness {
		d := strings.Split(k, ",")
		diners[d[0]] = struct{}{}
		diners[d[1]] = struct{}{}
	}
	for k := range diners {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func (r *Relation) PartOne() int {
	diners := r.Diners()
	r.head = diners[0]
	r.memo = map[string]int{}
	return r.Happiest(r.head, diners[1:])
}

func (r *Relation) PartTwo() int {
	diners := r.Diners()
	r.head = "Myself"
	r.memo = map[string]int{}
	return r.Happiest(r.head, diners)
}

func (r Relation) GetLink(a, b string) int {
	if a == "Myself" || b == "Myself" {
		return 0
	}
	return r.happiness[a+","+b] + r.happiness[b+","+a]
}

func (r *Relation) Happiest(head string, others []string) (out int) {
	key := strings.Join(append([]string{head}, others...), ",")
	if cached, exists := r.memo[key]; exists {
		return cached
	}
	n := len(others)
	if n == 0 {
		out = r.GetLink(head, r.head)
	} else if n == 1 {
		out = r.GetLink(head, others[0]) + r.GetLink(others[0], r.head)
	} else {
		for i, h := range others {
			t := remove(others, i)
			score := r.GetLink(head, h) + r.Happiest(h, t)
			if score > out {
				out = score
			}
		}
	}
	r.memo[key] = out
	return out
}

func remove(list []string, r int) []string {
	out := make([]string, len(list)-1)
	j := 0
	for i, v := range list {
		if i != r {
			out[j] = v
			j++
		}
	}
	return out
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load() (out Relation) {
	happiness := map[string]int{}
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		words := strings.Fields(s.Text())
		n := len(words)
		a := words[0]
		b := strings.ReplaceAll(words[n-1], ".", "")
		h := atoi(words[3])
		if words[2] == "lose" {
			h = -h
		}
		happiness[a+","+b] = h
	}
	check(s.Err())
	out.happiness = happiness
	return out
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func main() {
	r := load()
	fmt.Println("Part 1:", r.PartOne())
	fmt.Println("Part 2:", r.PartTwo())
}
