package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"slices"
)

type Node struct {
	pos  string
	last *Node
}

type Puzzle struct {
	links map[string][]string
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
	links := map[string][]string{}
	for s.Scan() {
		data := strings.Split(s.Text(), "-")
		links[data[0]] = append(links[data[0]], data[1])
		links[data[1]] = append(links[data[1]], data[0])
	}
	check(s.Err())
	p.links = links
	return p
}

func (n Node) Been() (been map[string]bool, hasDouble bool) {
	been = map[string]bool{}
	cur := n
	for {
		if cur.pos[0] >= 'a' && cur.pos[0] <= 'z' && been[cur.pos] {
			hasDouble = true
		}
		been[cur.pos] = true
		if cur.last == nil {
			break
		}
		cur = *cur.last
	}
	return been, hasDouble
}

func (p Puzzle) CountPaths(start, exit string, oneDouble bool) (count int) {
	q := []Node{Node{pos: start}}
	for len(q) > 0 {
		n := len(q)
		s := q[n-1]
		q = q[:n-1]
		if s.pos == exit {
			count++
			continue
		}
		been, hasDouble := s.Been()
		for _, adj := range p.links[s.pos] {
			if adj == start {
				continue
			}
			if adj[0] >= 'a' && adj[0] <= 'z' && been[adj] {
				if !oneDouble || hasDouble {
					continue
				}
			}
			q = append(q, Node{adj, &s})
		}
	}
	return count
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	c := p.CountPaths("start", "end", false)
	fmt.Println("Part 1:", c)
	d := p.CountPaths("start", "end", true)
	fmt.Println("Part 2:", d)
}
