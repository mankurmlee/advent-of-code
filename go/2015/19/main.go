package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
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

type DResult struct {
	steps int
	after string
}

type Sub struct {
	Find string
	Repl string
}

type Puzzle struct {
	base  string
	subs  []Sub
	seen  map[string]struct{}
	cache map[string]DResult
}

func newPuzzle() (out Puzzle) {
	lines := fileinput()
	cache := map[string]DResult{}
	for _, l := range lines {
		if l == "" {
			break
		}
		words := strings.Fields(l)
		k, v := words[0], words[2]
		out.subs = append(out.subs, Sub{k, v})
		cache[k] = DResult{0, k}
	}
	out.base = lines[len(out.subs)+1]
	out.seen = out.GenSubs(out.base)
	out.cache = cache
	return out
}

func (p Puzzle) GenSubs(base string) map[string]struct{} {
	seen := map[string]struct{}{}
	l := len(base)
	for _, s := range p.subs {
		n := len(s.Find)
		for i := 0; i+n <= l; i++ {
			if base[i:i+n] == s.Find {
				c := base[:i] + s.Repl + base[i+n:]
				seen[c] = struct{}{}
			}
		}
	}
	return seen
}

// Recursive substitution until it hits one of the base cases
func (p Puzzle) Deconstruct(base string) (out DResult) {
	if cached, exists := p.cache[base]; exists {
		return cached
	}

	out.steps = 1 << 32
	l := len(base)
	for _, s := range p.subs {
		if !strings.Contains(base, s.Repl) {
			continue
		}
		n := len(s.Repl)
		for i := 0; i+n <= l; i++ {
			if base[i:i+n] != s.Repl {
				continue
			}
			c := base[:i] + s.Find + base[i+n:]
			r := p.Deconstruct(c)
			r.steps++
			if r.steps < out.steps {
				out = r
			}
		}
	}
	p.cache[base] = out
	return out
}

func (p Puzzle) partTwo() (steps int) {
	var h string
	var s int
	// "Ar" only appears at the end of a pattern so split the input base on that
	toks := strings.Split(p.base, "Ar")
	for _, tok := range toks[:len(toks)-1] {
		// Reconstruct molecule based on:
		// h: Unmatched leftovers from last iteration
		// tok: The new chunk
		// Ar: Re-insert Ar as it was removed from the split operation
		m := h + tok + "Ar"
		h, s = p.Minimise(m)
		steps += s
	}
	return steps
}

func (p Puzzle) Minimise(inp string) (string, int) {
	bad := 1 << 32
	l := len(inp)
	for i := range l {
		b := inp[i]
		// Ignore molecules that start with lower case as none of the patterns will match that
		if b >= 'a' && b <= 'z' {
			continue
		}
		res := p.Deconstruct(inp[i:])
		if res.steps < bad {
			// s is the remainder after pattern substitution
			s := inp[:i] + res.after
			return s, res.steps
		}
	}
	return inp, bad
}

func main() {
	p := newPuzzle()
	fmt.Println("Part 1:", len(p.seen))
	fmt.Println("Part 2:", p.partTwo())
}
