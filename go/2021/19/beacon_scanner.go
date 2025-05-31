package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Vertex [3]int32

func (or Vertex) GetRotation() Matrix {
	var i int32
	rot := identity
	for i = 0; i < or[0]; i++ {
		rot = rot.MulM(xrot)
	}
	for i = 0; i < or[1]; i++ {
		rot = rot.MulM(yrot)
	}
	for i = 0; i < or[2]; i++ {
		rot = rot.MulM(zrot)
	}
	return rot
}

func (a Vertex) Add(b Vertex) (c Vertex) {
	c[0] = a[0] + b[0]
	c[1] = a[1] + b[1]
	c[2] = a[2] + b[2]
	return c
}

func (a Vertex) Sub(b Vertex) (c Vertex) {
	c[0] = a[0] - b[0]
	c[1] = a[1] - b[1]
	c[2] = a[2] - b[2]
	return c
}

func (a Vertex) Dist(b Vertex) int32 {
	return abs(a[0] - b[0]) + abs(a[1] - b[1]) + abs(a[2] - b[2])
}

func (v Vertex) Move(in Scanner) Scanner {
	out := Scanner {
		in.id, v, make([]Vertex, len(in.objs)),
	}
	for i, o := range in.objs {
		out.objs[i] = Vertex{
			o[0] + v[0], o[1] + v[1], o[2] + v[2],
		}
	}
	return out
}

type Matrix [9]int32

func (m Matrix) MulV(u Vertex) (v Vertex) {
	v[0] = m[0] * u[0] + m[1] * u[1] + m[2] * u[2]
	v[1] = m[3] * u[0] + m[4] * u[1] + m[5] * u[2]
	v[2] = m[6] * u[0] + m[7] * u[1] + m[8] * u[2]
	return v
}

func (a Matrix) MulM(b Matrix) (out Matrix) {
	out[0] = a[0] * b[0] + a[3] * b[1] + a[6] * b[2]
	out[1] = a[1] * b[0] + a[4] * b[1] + a[7] * b[2]
	out[2] = a[2] * b[0] + a[5] * b[1] + a[8] * b[2]
	out[3] = a[0] * b[3] + a[3] * b[4] + a[6] * b[5]
	out[4] = a[1] * b[3] + a[4] * b[4] + a[7] * b[5]
	out[5] = a[2] * b[3] + a[5] * b[4] + a[8] * b[5]
	out[6] = a[0] * b[6] + a[3] * b[7] + a[6] * b[8]
	out[7] = a[1] * b[6] + a[4] * b[7] + a[7] * b[8]
	out[8] = a[2] * b[6] + a[5] * b[7] + a[8] * b[8]
	return out
}

type Scanner struct {
	id      int32
	offset  Vertex
	objs    []Vertex
}

func (s Scanner) Clip(offset Vertex) []Vertex {
	out := make([]Vertex, 0, len(s.objs))
	lo := offset.Add(Vertex{-1000, -1000, -1000})
	hi := offset.Add(Vertex{ 1000,  1000,  1000})
	for _, v := range s.objs {
		if  v[0] > lo[0] && v[0] < hi[0] &&
			v[1] > lo[1] && v[1] < hi[1] &&
			v[2] > lo[2] && v[2] < hi[2] {
			out = append(out, v)
		}
	}
	return out
}

func (a Scanner) Overlaps(b Scanner) bool {
	clipb := b.Clip(a.offset)
	if len(clipb) < 12 {
		return false
	}
	clipa := a.Clip(b.offset)
	if len(clipa) != len(clipb) {
		return false
	}
	lut := make(map[Vertex]bool)
	for _, v := range clipa {
		lut[v] = true
	}
	for _, v := range clipb {
		if !lut[v] {
			return false
		}
	}
	return true
}

func (a Scanner) FindDisplacement(b Scanner) (out Scanner, ok bool) {
	for _, va := range a.objs {
		for _, vb := range b.objs {
			move := va.Sub(vb)
			out := move.Move(b)
			if a.Overlaps(out) {
				return out, true
			}
		}
	}
	return out, false
}

func (base Scanner) FindRotation(in Scanner) (out Scanner, ok bool) {
	for i, rot := range rotations {
		arr := in.Reorient(rot)
		out, ok := base.FindDisplacement(arr)
		if ok {
			fmt.Println("FindRotation():", base.id, "x", out.id, ", o", out.offset, "r", i+1)
			return out, true
		}
	}
	return out, false
}

func (in Scanner) Reorient(rot Matrix) (out Scanner) {
	out.id = in.id
	for _, v := range in.objs {
		out.objs = append(out.objs, rot.MulV(v))
	}
	return out
}

type Puzzle []Scanner

func getRotations(data []Vertex) (out []Matrix) {
	for _, o := range data {
		out = append(out, o.GetRotation())
	}
	return out
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	var sc Scanner
	for s.Scan() {
		line := s.Text()
		if line == "" { continue }
		data := strings.Fields(line)
		if data[0] == "---" {
			if len(sc.objs) > 0 {
				p = append(p, sc)
			}
			sc = Scanner{ id: atoi(data[2]) }
			continue
		}
		bdata := strings.Split(data[0], ",")
		sc.objs = append(sc.objs, Vertex{
			atoi(bdata[0]), atoi(bdata[1]), atoi(bdata[2]),
		})
	}
	check(s.Err())
	p = append(p, sc)
	return p
}

func abs(n int32) int32 {
	if n >= 0 { return n }
	return -n
}

func atoi(a string) int32 {
	i, err := strconv.Atoi(a)
	check(err)
	return int32(i)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	p := load(flag.Arg(0))
	q := partOne(p)
	partTwo(q)
}

func partTwo(q Puzzle) {
	n := len(q)
	var m int32
	for i := 0; i < n - 1; i++ {
		a := q[i].offset
		for j := i + 1; j < n; j++ {
			b := q[j].offset
			d := a.Dist(b)
			if d > m { m = d }
		}
	}
	fmt.Println("Part 2:", m)
}

func partOne(p Puzzle) Puzzle {
	var todo, done Puzzle
	untested := p[:1]
	unmerged := p[1:]
	for len(untested) > 0 && len(unmerged) > 0 {
		n := len(untested)
		test := untested[n-1]
		untested = untested[:n-1]
		done = append(done, test)
		todo = make(Puzzle, 0, len(unmerged))
		for _, s := range unmerged {
			m, ok := test.FindRotation(s)
			if ok {
				untested = append(untested, m)
			} else {
				todo = append(todo, s)
			}
		}
		unmerged = todo
	}
	done = append(done, untested...)

	lut := make(map[Vertex]bool)
	for _, s := range done {
		for _, v := range s.objs {
			lut[v] = true
		}
	}
	fmt.Println("Part 1:", len(lut))

	return done
}

var (
	rotations = getRotations([]Vertex{
		Vertex{ 0, 0, 0 }, Vertex{ 0, 0, 1 }, Vertex{ 0, 0, 2 }, Vertex{ 0, 0, 3 },
		Vertex{ 0, 1, 0 }, Vertex{ 0, 1, 1 }, Vertex{ 0, 1, 2 }, Vertex{ 0, 1, 3 },
		Vertex{ 0, 2, 0 }, Vertex{ 0, 2, 1 }, Vertex{ 0, 2, 2 }, Vertex{ 0, 2, 3 },
		Vertex{ 0, 3, 0 }, Vertex{ 0, 3, 1 }, Vertex{ 0, 3, 2 }, Vertex{ 0, 3, 3 },
		Vertex{ 1, 0, 0 }, Vertex{ 1, 0, 1 }, Vertex{ 1, 0, 2 }, Vertex{ 1, 0, 3 },
		Vertex{ 3, 0, 0 }, Vertex{ 3, 0, 1 }, Vertex{ 3, 0, 2 }, Vertex{ 3, 0, 3 },
	})
	xrot = Matrix{ 1, 0, 0, 0, 0, 1, 0,-1, 0 }
	yrot = Matrix{ 0, 0, 1, 0, 1, 0, -1, 0, 0 }
	zrot = Matrix{ 0, 1, 0, -1, 0, 0, 0, 0, 1 }
	identity = Matrix{ 1, 0, 0, 0, 1, 0, 0, 0, 1 }
)
