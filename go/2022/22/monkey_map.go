package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Vertex struct {
	x, y int
}

type Range struct {
	lo, hi int
}

type Puzzle struct {
	grid        [][]byte
	path        []string
	width       int
	height      int
	rowrange    []Range
	colrange    []Range
	offmap      func(Scout) Scout
}

type Scout struct {
	pos    Vertex
	facing int
}

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

func abs(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	// Load the map
	for s.Scan() {
		line := s.Text()
		if line == "" {
			break
		}
		runes := []byte(line)
		if len(runes) > p.width {
			p.width = len(runes)
		}
		p.grid = append(p.grid, runes)
	}
	check(s.Err())
	p.height = len(p.grid)

	// Load the path
	r := strings.NewReplacer("L", " L ", "R", " R ")
	s.Scan()
	check(s.Err())
	p.path = strings.Fields(r.Replace(s.Text()))

	// Get row ranges
	for _, row := range p.grid {
		lo := -1
		hi := -1
		for x, tile := range row {
			if tile == ' ' {
				continue
			}
			if lo == -1 {
				lo = x
			}
			hi = x
		}
		p.rowrange = append(p.rowrange, Range{lo, hi})
	}

	// Get col ranges
	for x := 0; x < p.width; x++ {
		lo := -1
		hi := -1
		for y := 0; y < p.height; y++ {
			if x >= len(p.grid[y]) {
				continue
			}
			if p.grid[y][x] == ' ' {
				continue
			}
			if lo == -1 {
				lo = y
			}
			hi = y
		}
		p.colrange = append(p.colrange, Range{lo, hi})
	}
	p.offmap = p.Straight

	return p
}

func (p Puzzle) MapWalk() (s Scout) {
	s.pos = Vertex{ p.rowrange[0].lo, 0 }
	for i, v := range p.path {
		if i % 2 == 1 {
			// change direction
			if v == "R" {
				s.facing = (s.facing + 1) % 4
			} else {
				s.facing = (s.facing + 3) % 4
			}
			continue
		}
		// move forward in the direction of facing
		steps := atoi(v)
		for j := 0; j < steps; j++ {
			peek := s
			switch s.facing {
			case 0:
				peek.pos.x++
				if peek.pos.x > p.rowrange[peek.pos.y].hi {
					peek = p.offmap(s)
				}
			case 1:
				peek.pos.y++
				if peek.pos.y > p.colrange[peek.pos.x].hi {
					peek = p.offmap(s)
				}
			case 2:
				peek.pos.x--
				if peek.pos.x < p.rowrange[peek.pos.y].lo {
					peek = p.offmap(s)
				}
			case 3:
				peek.pos.y--
				if peek.pos.y < p.colrange[peek.pos.x].lo {
					peek = p.offmap(s)
				}
			}
			if p.grid[peek.pos.y][peek.pos.x] == '#' {
				fmt.Println("Wall at", peek.pos)
				break
			}
			s = peek
			fmt.Println("Moved to", s.pos)
		}
	}
	return s
}

func (s Scout) GetPassword() int {
	return 1000 * (s.pos.y + 1) + 4 * (s.pos.x + 1) + s.facing
}

func (p Puzzle) Straight(s Scout) Scout {
	peek := s
	switch s.facing {
	case 0:
		peek.pos.x = p.rowrange[peek.pos.y].lo
	case 1:
		peek.pos.y = p.colrange[peek.pos.x].lo
	case 2:
		peek.pos.x = p.rowrange[peek.pos.y].hi
	case 3:
		peek.pos.y = p.colrange[peek.pos.x].hi
	}
	return peek
}

func (p Puzzle) SampleCube(s Scout) Scout {
	peek := s
	size := abs(p.width - p.height)
	x := s.pos.x / size
	y := s.pos.y / size
	if        x == 2 && y == 0 {
		peek.pos.x -= 2 * size
		switch s.facing {
		case 0:
			peek.pos.x = 3 * size
			peek.pos.y = 3 * size - peek.pos.y - 1
			peek.facing = 2
		case 2:
			peek.pos.x = size + peek.pos.y
			peek.pos.y = size
			peek.facing = 1
		case 3:
			peek.pos.x = size - peek.pos.x - 1
			peek.pos.y = size
			peek.facing = 1
		}
	} else if x == 0 && y == 1 {
		peek.pos.y -= size
		switch s.facing {
		case 1:
			peek.pos.x = 3 * size - peek.pos.x - 1
			peek.pos.y = 3 * size - 1
			peek.facing = 3
		case 2:
			peek.pos.x = 4 * size - peek.pos.y - 1
			peek.pos.y = 3 * size - 1
			peek.facing = 3
		case 3:
			peek.pos.x = 3 * size - peek.pos.x - 1
			peek.pos.y = 0
			peek.facing = 1
		}
	} else if x == 1 && y == 1 {
		peek.pos.x -= size
		peek.pos.y -= size
		switch s.facing {
		case 1:
			peek.pos.y = 3 * size - peek.pos.x - 1
			peek.pos.x = 2 * size
			peek.facing = 0
		case 3:
			peek.pos.y = peek.pos.x
			peek.pos.x = 2 * size
			peek.facing = 0
		}

	} else if x == 2 && y == 1 {
		peek.pos.x -= 2 * size
		peek.pos.y -= size
		if s.facing == 0 {
			peek.pos.x = 4 * size - peek.pos.y - 1
			peek.pos.y = 2 * size
			peek.facing = 1
		}
	} else if x == 2 && y == 2 {
		peek.pos.x -= size * 2
		peek.pos.y -= size * 2
		switch s.facing {
		case 1:
			peek.pos.x = size - peek.pos.x - 1
			peek.pos.y = size * 2 - 1
			peek.facing = 3
		case 2:
			peek.pos.x = 2* size - peek.pos.x - 1
			peek.pos.y = size * 2 - 1
			peek.facing = 3
		}
	} else if x == 3 && y == 2 {
		peek.pos.x -= size * 3
		peek.pos.y -= size * 2
		switch s.facing {
		case 0:
			peek.pos.x = size * 3 - 1
			peek.pos.y = size - peek.pos.y - 1
			peek.facing = 2
		case 1:
			peek.pos.y = 2 * size - peek.pos.x - 1
			peek.pos.x = 0
			peek.facing = 0
		case 3:
			peek.pos.y = 2 * size - peek.pos.x - 1
			peek.pos.x = size * 3 - 1
			peek.facing = 2
		}
	}
	return peek
}

func (p Puzzle) PuzzleCube(s Scout) Scout {
	ox := s.pos.x
	oy := s.pos.y
	of := s.facing
	size := abs(p.width - p.height)
	cx := ox / size
	cy := oy / size
	x := ox % size
	y := oy % size

	if        cx == 1 && cy == 0 {
		switch of {
		case 2:
			ox = 0
			oy = size * 2 + (size - y - 1)
			of = 0
		case 3:
			ox = 0
			oy = size * 3 + x
			of = 0
		}
	} else if cx == 2 && cy == 0 {
		switch of {
		case 0:
			ox = size * 2 - 1
			oy = size * 2 + (size - y - 1)
			of = 2
		case 1:
			ox = size * 2 - 1
			oy = size + x
			of = 2
		case 3:
			ox = x
			oy = size * 4 - 1
			of = 3
		}
	} else if cx == 1 && cy == 1 {
		switch of {
		case 0:
			ox = size * 2 + y
			oy = size - 1
			of = 3
		case 2:
			ox = y
			oy = size * 2
			of = 1
		}
	} else if cx == 0 && cy == 2 {
		switch of {
		case 2:
			ox = size
			oy = size - y - 1
			of = 0
		case 3:
			ox = size
			oy = size + x
			of = 0
		}
	} else if cx == 1 && cy == 2 {
		switch of {
		case 0:
			ox = size * 3 - 1
			oy = size - y - 1
			of = 2
		case 1:
			ox = size - 1
			oy = size * 3 + x
			of = 2
		}
	} else if cx == 0 && cy == 3 {
		switch of {
		case 0:
			ox = size + y
			oy = size * 3 - 1
			of = 3
		case 1:
			ox = size * 2 + x
			oy = 0
			of = 1
		case 2:
			ox = size + y
			oy = 0
			of = 1
		}
	}
	return Scout{ Vertex{ ox, oy }, of }
}

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	p := load(filename)
	scout := p.MapWalk()
	fmt.Println("Scout position:", scout)
	fmt.Println("Part 1:", scout.GetPassword())

	if filename != "puzzle.txt" {
		p.offmap = p.SampleCube
	} else {
		p.offmap = p.PuzzleCube
	}
	scout = p.MapWalk()
	fmt.Println("Scout position:", scout)
	fmt.Println("Part 2:", scout.GetPassword())
}
