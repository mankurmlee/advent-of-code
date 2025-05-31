package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
)

type Grid struct {
	base    byte
	size    int
	data    [][]byte
}

func (g Grid) Get(x, y int) byte {
	if x < 0 || y < 0 || x >= g.size || y >= g.size {
		return g.base
	}
	return g.data[y][x]
}

func (g Grid) CountLit() (n int) {
	for j, _ := range g.data {
		for _, b := range g.data[j] {
			if b == '#' { n++ }
		}
	}
	return n
}

type Key [9]byte

type IEA map[Key]byte

type Puzzle struct {
	iea     IEA
	grid    Grid
}

func (f IEA) Enhance(old Grid) (new Grid) {
	var k Key
	var i, j, u, v int

	n := old.size + 2

	for i, _ := range k {
		k[i] = old.base
	}

	new.base = f[k]
	new.size = n

	new.data = make([][]byte, n)
	for j, _ = range new.data {
		new.data[j] = make([]byte, n)
		v = j - 1
		for i, _ = range new.data[j] {
			u = i - 1
			k[0] = old.Get(u - 1, v - 1)
			k[1] = old.Get(u    , v - 1)
			k[2] = old.Get(u + 1, v - 1)
			k[3] = old.Get(u - 1, v    )
			k[4] = old.Get(u    , v    )
			k[5] = old.Get(u + 1, v    )
			k[6] = old.Get(u - 1, v + 1)
			k[7] = old.Get(u    , v + 1)
			k[8] = old.Get(u + 1, v + 1)
			new.data[j][i] = f[k]
		}
	}

	return new
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) (p Puzzle) {
	var g Grid
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		bytes := []byte(s.Text())
		if p.iea == nil {
			p.iea = make(map[Key]byte)
			for i, b := range bytes {
				p.iea[itok(i)] = b
			}
			s.Scan()
			continue
		}
		g.data = append(g.data, bytes)
	}
	check(s.Err())
	g.base = '.'
	g.size = len(g.data)
	p.grid = g
	return p
}

func itok(n int) (out Key) {
	for i, _ := range out {
		if n & (1 << (8 - i)) > 0 {
			out[i] = '#'
		} else {
			out[i] = '.'
		}
	}
	return out
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	f := p.iea
	g := p.grid

	for _ = range 2 {
		g = f.Enhance(g)
	}
	fmt.Println("Part 1:", g.CountLit())

	for _ = range 48 {
		g = f.Enhance(g)
	}
	fmt.Println("Part 2:", g.CountLit())
}
