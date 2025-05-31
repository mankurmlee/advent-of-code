package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
)

var directions = map[rune]Vec{
	'^': {0, -1}, '>': {1, 0},
	'v': {0, 1}, '<': {-1, 0},
}

type Vec struct {
	x, y int
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

type Room struct {
	w, h int
	data []rune
}

func (r Room) Move(g []Vec, d Vec) {
	slices.Reverse(g)
	for _, v := range g {
		t := r.Peek(v)
		v1 := v.Add(d)
		r.Poke(v1, t)
		r.Poke(v, '.')
	}
}

func (r Room) CanMove(robot Vec, d Vec) (out []Vec, ok bool) {
	been := make(map[Vec]bool)
	been[robot] = true
	q := []Vec{robot}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		out = append(out, v)
		v1 := v.Add(d)
		t := r.Peek(v1)
		if t == '#' {
			return []Vec{}, false
		}
		if t != '[' && t != ']' && t != 'O' || been[v1] {
			continue
		}
		q = append(q, v1)
		been[v1] = true
		if d.y == 0 || t != '[' && t != ']' {
			continue
		}
		a := '<'
		if t == '[' {
			a = '>'
		}
		v2 := v1.Add(directions[a])
		if been[v2] {
			continue
		}
		q = append(q, v2)
		been[v2] = true
	}
	return out, true
}

func (r Room) Widen() Room {
	var data []rune
	for _, r := range r.data {
		switch r {
		case '#', '.':
			data = append(data, r, r)
		case 'O':
			data = append(data, '[', ']')
		case '@':
			data = append(data, '@', '.')
		default:
			panic("Unexpected room tile")
		}
	}
	return Room{r.w << 1, r.h, data}
}

func (r Room) Checksum(c rune) (tot int) {
	w := r.w
	for i, v := range r.data {
		if v == c {
			tot += 100*(i/w) + (i % w)
		}
	}
	return tot
}

func (r Room) Peek(pos Vec) rune {
	if pos.x < 0 || pos.y < 0 || pos.x >= r.w || pos.y >= r.h {
		return '#'
	}
	return r.data[pos.y*r.w+pos.x]
}

func (r Room) Poke(pos Vec, t rune) {
	if pos.x < 0 || pos.y < 0 || pos.x >= r.w || pos.y >= r.h {
		return
	}
	r.data[pos.y*r.w+pos.x] = t
}

func (r Room) Draw() {
	i := 0
	for range r.h {
		fmt.Println(string(r.data[i : i+r.w]))
		i += r.w
	}
}

func (r Room) FindRobot() Vec {
	for i, v := range r.data {
		if v == '@' {
			return Vec{i % r.w, i / r.w}
		}
	}
	return Vec{}
}

func (r Room) Clone() Room {
	items := make([]rune, len(r.data))
	copy(items, r.data)
	return Room{r.w, r.h, items}
}

type Puzzle struct {
	Room
	moves []rune
}

func (p Puzzle) PartTwo(c rune) {
	var r Room
	if c == 'O' {
		r = p.Room.Clone()
	} else {
		r = p.Room.Widen()
	}
	bot := r.FindRobot()
	for _, m := range p.moves {
		d := directions[m]
		g, ok := r.CanMove(bot, d)
		if !ok {
			continue
		}
		r.Move(g, d)
		bot = bot.Add(d)
	}
	r.Draw()
	fmt.Println(r.Checksum(c))
}

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
		if l == "" {
			if len(chunk) > 0 {
				chunks = append(chunks, chunk)
				chunk = []string{}
			}
			continue
		}
		chunk = append(chunk, l)
	}
	check(s.Err())
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}

func load(filename string) (p Puzzle) {
	d := readChunkedFile(filename)
	p.w = len(d[0][0])
	p.h = len(d[0])
	for _, l := range d[0] {
		p.data = append(p.data, []rune(l)...)
	}
	for _, l := range d[1] {
		p.moves = append(p.moves, []rune(l)...)
	}
	return p
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	p.PartTwo('O')
	p.PartTwo('[')
}
