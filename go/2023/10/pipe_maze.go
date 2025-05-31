package main

import (
	"fmt"
	"os"
	"bufio"
	"slices"
	"flag"
)

type Tile struct {
	symbol rune
	N, E, S, W bool
}

type Maze [][]Tile

type Vec struct {
	x, y int
}

func main() {
	flag.Parse()
	filename := flag.Args()[0]

	start, maze := parseInput(filename)
	maze.Print()

	// Part 1
	loop, farthest := getLoopAndFarthest(start, maze)
	fmt.Println("Farthest distance from start:", farthest)

	// Part 2

	// Create visual representation, not strictly needed but fun
	layer := maze.NewEnclosedLayer(loop)
	layer.Print()

	size := layer.CountTileType('I')
	fmt.Println(size, "tiles are enclosed")
}

func parseInput(filename string) (Vec, Maze) {
	var start Vec
	var maze Maze
	var lineNumber int

	startTile := Tile {
		symbol: 'S',
		N: true,
		E: true,
		S: true,
		W: true,
	}

	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		row := parseLine(line)
		maze = append(maze, row)

		if pos := slices.Index(row, startTile); pos >= 0 {
			start = Vec {
				x: pos,
				y: lineNumber,
			}
		}
		lineNumber++
	}
	check(s.Err())

	return start, maze
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func parseLine(text string) []Tile {
	var row []Tile
	for _, c := range text {
		tile := Tile {
			symbol: c,
		}
		switch c {
		case 'L':
			tile.N = true
			tile.E = true
		case '|':
			tile.N = true
			tile.S = true
		case 'J':
			tile.N = true
			tile.W = true
		case 'F':
			tile.E = true
			tile.S = true
		case '-':
			tile.E = true
			tile.W = true
		case '7':
			tile.S = true
			tile.W = true
		case 'S':
			tile.N = true
			tile.E = true
			tile.S = true
			tile.W = true
		}
		row = append(row, tile)
	}
	return row
}

func getLoopAndFarthest(s Vec, m Maze) ([]Vec, int) {
	var prev, current, next []Vec
	var distance int

	current = append(current, s)

	for {
		next = []Vec{}

		for _, o := range current {
			for _, d := range m.Explore(o) {
				if slices.Contains(prev, d) || slices.Contains(next, d) {
					continue
				}
				next = append(next, d)
			}
		}

		prev = append(prev, current...)

		if len(next) == 0 {
			return prev, distance
		}

		current = next
		distance++
	}
}

func (m Maze) Print() {
	for _, r := range m {
		runes := []rune{}
		for _, t := range r {
			runes = append(runes, t.symbol)
		}
		fmt.Println(string(runes))
	}
}

func (m Maze) Explore(v Vec) []Vec {
	var visitable []Vec

	width  := len(m[0])
	height := len(m)

	tile := m[v.y][v.x]
	if tile.N {
		n := Vec {
			x : v.x,
			y : v.y - 1,
		}
		if n.y >= 0 && m[n.y][n.x].S {
			visitable = append(visitable, n)
		}
	}
	if tile.E {
		n := Vec {
			x : v.x + 1,
			y : v.y,
		}
		if n.x < width && m[n.y][n.x].W {
			visitable = append(visitable, n)
		}
	}
	if tile.S {
		n := Vec {
			x : v.x,
			y : v.y + 1,
		}
		if n.y < height && m[n.y][n.x].N {
			visitable = append(visitable, n)
		}
	}
	if tile.W {
		n := Vec {
			x : v.x - 1,
			y : v.y,
		}
		if n.x >= 0 && m[n.y][n.x].E {
			visitable = append(visitable, n)
		}
	}

	return visitable
}

func (m Maze) NewEnclosedLayer(loop []Vec) Maze {
	layer := m

	var up, down bool
	for y, row := range m {
		for x, tile := range row {
			if slices.Contains(loop, Vec{ x: x, y: y }) {
				switch tile.symbol {
				case '|':
					up   = !up
					down = !down
				case 'L', 'J':
					up = !up
				case '7', 'F':
					down = !down
				case 'S':
					if y > 0 && m[y-1][x].S {
						up = !up
					}
					if y < len(m)-1 && m[y+1][x].N {
						down = !down
					}
				}
				continue
			}

			sym := ' '
			if up && down {
				sym = 'I'
			} else if !up && !down {
				sym = 'O'
			}

			tile.symbol = sym
			layer[y][x] = tile
		}
	}

	return layer
}

func (m Maze) CountTileType(symbol rune) int {
	var count int
	for _, r := range m {
		for _, t := range r {
			if t.symbol == symbol {
				count++
			}
		}
	}
	return count
}
