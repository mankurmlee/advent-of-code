package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
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

func findWords(s string) (out []string) {
	re := regexp.MustCompile(`\w+`)
	return re.FindAllString(s, -1)
}

func load(filename string) map[string][]string {
	out := make(map[string][]string)
	for _, l := range readFile(filename) {
		tokens := findWords(l)
		a, b := tokens[0], tokens[1]
		out[a] = append(out[a], b)
		out[b] = append(out[b], a)
	}
	return out
}

func main() {
	flag.Parse()
	conns := load(flag.Arg(0))
	partOne(conns)
	partTwo(conns)
}

func partOne(conns map[string][]string) {
	res := make(map[string]struct{})
	for k := range conns {
		if k[0] != 't' {
			continue
		}
		for _, t := range triplets(conns, k) {
			res[t] = struct{}{}
		}
	}
	fmt.Println(len(res))
}

func triplets(conns map[string][]string, k string) (out []string) {
	ends := conns[k]
	n := len(ends)
	for i, v1 := range ends[:n] {
		e := conns[v1]
		for _, v2 := range ends[i+1:] {
			if slices.Contains(e, v2) {
				k2 := []string{k, v1, v2}
				slices.Sort(k2)
				out = append(out, strings.Join(k2, ","))
			}
		}
	}
	return out
}

func partTwo(conns map[string][]string) {
	var comps []string
	cache := make(map[string][][]string)
	for k := range conns {
		comps = append(comps, k)
		cache[k] = [][]string{{k}}
	}
	slices.Sort(comps)
	memo := Memo{
		conns,
		cache,
	}
	lan := memo.BiggestLAN(comps)
	if len(lan) != 1 {
		fmt.Println("Unexpected result. Bad input?")
		panic("Unexpected")
	}
	fmt.Println(strings.Join(lan[0], ","))
}

type Memo struct {
	conns map[string][]string
	memo  map[string][][]string
}

func (m Memo) BiggestLAN(comps []string) (out [][]string) {
	key := strings.Join(comps, ",")
	if out, ok := m.memo[key]; ok {
		return out
	}
	out = [][]string{{}}
	for i, v := range comps {
		o := len(out[0])
		shared := intersect(comps[i:], m.conns[v])
		if len(shared) < o-1 {
			continue
		}
		loops := m.BiggestLAN(shared)
		if len(loops) == 0 {
			continue
		}
		n := len(loops[0]) + 1
		if o > n {
			continue
		}
		if o < n {
			out = [][]string{}
		}
		for _, l := range loops {
			l2 := make([]string, n)
			l2[0] = v
			copy(l2[1:], l)
			out = append(out, l2)
		}
	}
	m.memo[key] = out
	return out
}

func intersect(a, b []string) (out []string) {
	for _, v := range a {
		if slices.Contains(b, v) {
			out = append(out, v)
		}
	}
	return out
}
