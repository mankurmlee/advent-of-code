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
	x, y, z int
}

type Bounds struct {
	u, v Vertex
}

type Volume map[Vertex]bool

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

func load(filename string) (d Volume) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	d = make(map[Vertex]bool)
	for s.Scan() {
		data := strings.Split(s.Text(), ",")
		d[Vertex{ atoi(data[0]), atoi(data[1]), atoi(data[2]) }] = true
	}
	check(s.Err())
	return d
}

func (v Vertex) Add(u Vertex) Vertex {
	return Vertex{u.x + v.x, u.y + v.y, u.z + v.z}
}

func (vol Volume) GetSurfaceArea() (area int) {
	for v := range vol {
		for _, u := range v.GetAdj() {
			if !vol[u] {
				area++
			}
		}
	}
	return area
}

func (vol Volume) GetBounds() (bounds Bounds) {
	for v := range vol {
		bounds = Bounds{v, v}
		break
	}
	for v := range vol {
		bounds.u.x = min(bounds.u.x, v.x)
		bounds.v.x = max(bounds.v.x, v.x)
		bounds.u.y = min(bounds.u.y, v.y)
		bounds.v.y = max(bounds.v.y, v.y)
		bounds.u.z = min(bounds.u.z, v.z)
		bounds.v.z = max(bounds.v.z, v.z)
	}
	return bounds
}

func (vol Volume) GetBoundingVolume(b Bounds) (outer Volume) {
	var zv Vertex
	outer = make(map[Vertex]bool)
	outer[zv] = true
	q := []Vertex{ zv }
	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		for _, v := range s.GetAdj() {
			if outer[v] || vol[v] || v.x < b.u.x || v.y < b.u.y || v.z < b.u.z || v.x > b.v.x || v.y > b.v.y || v.z > b.v.z {
				continue
			}
			outer[v] = true
			q = append(q, v)
		}
	}
	return outer
}

func (v Vertex) GetAdj() []Vertex {
	return []Vertex{
		Vertex{v.x - 1, v.y, v.z},
		Vertex{v.x, v.y - 1, v.z},
		Vertex{v.x, v.y, v.z - 1},
		Vertex{v.x + 1, v.y, v.z},
		Vertex{v.x, v.y + 1, v.z},
		Vertex{v.x, v.y, v.z + 1},
	}
}

func main() {
	flag.Parse()
	vol := load(flag.Arg(0))
	area := vol.GetSurfaceArea()
	fmt.Println("Part 1:", area)

	bounds := vol.GetBounds()
	bounds.u = bounds.u.Add(Vertex{-1, -1, -1})
	bounds.v = bounds.v.Add(Vertex{ 1,  1,  1})
	vBounds := vol.GetBoundingVolume(bounds)

	var pockets Volume = make(map[Vertex]bool)
	for z := bounds.u.z + 1; z < bounds.v.z; z++ {
		for y := bounds.u.y + 1; y < bounds.v.y; y++ {
			for x := bounds.u.x + 1; x < bounds.v.x; x++ {
				v := Vertex{x, y, z}
				if !vBounds[v] && !vol[v] {
					pockets[v] = true
				}
			}
		}
	}

	innerArea := pockets.GetSurfaceArea()
	fmt.Println("Part 2:", area - innerArea)
}
