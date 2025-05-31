package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type Wordsearch struct {
	width  int
	height int
	data   []rune
}

var directions = [][2]int{
	{-1, -1}, {-1, 0}, {-1, 1}, {0, -1},
	{0, 1}, {1, -1}, {1, 0}, {1, 1},
}

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

func load(filename string) (ws Wordsearch) {
	d := readFile(filename)
	for _, s := range d {
		ws.data = append(ws.data, []rune(s)...)
	}
	ws.width = len(d[0])
	ws.height = len(d)
	return ws
}

func main() {
	flag.Parse()
	ws := load(flag.Arg(0))
	fmt.Println(ws.PartOne())
	fmt.Println(ws.PartTwo())
}

func (ws Wordsearch) PartTwo() (tot int) {
	w := ws.width
	for y := range ws.height - 2 {
		for x := range ws.width - 2 {
			s0 := string([]rune{ws.data[x+y*w], ws.data[(x+1)+(y+1)*w], ws.data[(x+2)+(y+2)*w]})
			s1 := string([]rune{ws.data[(x+2)+y*w], ws.data[(x+1)+(y+1)*w], ws.data[x+(y+2)*w]})
			if (s0 == "MAS" || s0 == "SAM") && (s1 == "MAS" || s1 == "SAM") {
				tot++
			}
		}
	}
	return tot
}

func (ws Wordsearch) PartOne() (tot int) {
	for y := range ws.height {
		for x := range ws.width {
			for _, d := range directions {
				if ws.Search(x, y, d[0], d[1]) {
					tot++
				}
			}
		}
	}
	return tot
}

func (ws Wordsearch) Search(x, y, xs, ys int) bool {
	for _, r := range "XMAS" {
		if x < 0 || x >= ws.width || y < 0 || y >= ws.height || ws.data[x+y*ws.width] != r {
			return false
		}
		x += xs
		y += ys
	}
	return true
}
