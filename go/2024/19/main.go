package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readChunkedFile(filename string) (chunks [][]string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	var chunk []string
	for s.Scan() {
		l := s.Text()
		if l != "" {
			chunk = append(chunk, l)
			continue
		}
		if len(chunk) == 0 {
			continue
		}
		chunks = append(chunks, chunk)
		chunk = []string{}
	}
	check(s.Err())
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}

func load(filename string) Puzzle {
	re := regexp.MustCompile(`\w+`)
	d := readChunkedFile(filename)
	cache := make(map[string]int)
	cache[""] = 1
	return Puzzle{re.FindAllString(d[0][0], -1), d[1], cache}
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	count := 0
	sum := 0
	for _, design := range p.designs {
		w := p.CountWays(design)
		if w > 0 {
			count++
		}
		sum += w
	}
	fmt.Println(count)
	fmt.Println(sum)
}

type Puzzle struct {
	patts   []string
	designs []string
	cache   map[string]int
}

func (p Puzzle) CountWays(s string) (ways int) {
	if w, ok := p.cache[s]; ok {
		return w
	}
	for _, pat := range p.patts {
		i := len(pat)
		if i > len(s) {
			continue
		}
		if s[:i] != pat {
			continue
		}
		ways += p.CountWays(s[i:])
	}
	p.cache[s] = ways
	return ways
}
