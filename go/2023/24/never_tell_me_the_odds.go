package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
	"gonum.org/v1/gonum/mat"
	"math"
)

type Vec struct {
	x, y, z int
}

type Hailstone struct {
	pos, vel Vec
}

func (h Hailstone) GetCoe() []int {
	return []int{
		h.vel.y,
		-h.vel.x,
		-h.pos.y,
		h.pos.x,
		h.pos.y * h.vel.x - h.pos.x * h.vel.y,
	}
}

type Puzzle struct {
	objs []Hailstone
}

func (p Puzzle) GetRow(i, j int) []int {
	a := p.objs[i].GetCoe()
	b := p.objs[j].GetCoe()
	return []int{
		a[0] - b[0],
		a[1] - b[1],
		a[2] - b[2],
		a[3] - b[3],
		a[4] - b[4],
	}
}

func (p Puzzle) GetAnswer() {
	a := p.GetRow(0, 1)
	b := p.GetRow(2, 3)
	c := p.GetRow(4, 5)
	d := p.GetRow(6, 7)
	data := append(a[:4], b[:4]...)
	data = append(data, c[:4]...)
	data = append(data, d[:4]...)
	floats := make([]float64, len(data))
	for i, v := range data {
		floats[i] = float64(v)
	}

	matA := mat.NewDense(4, 4, floats)
	matB := mat.NewDense(4, 1, []float64{
		float64(-a[4]),
		float64(-b[4]),
		float64(-c[4]),
		float64(-d[4]),
	})

	var aInv mat.Dense
	err := aInv.Inverse(matA)
	check(err)
	var matX mat.Dense
	matX.Mul(&aInv, matB)

	col := mat.Col(nil, 0, &matX)
	px := int(math.Round(col[0]))
	py := int(math.Round(col[1]))
	vx := int(math.Round(col[2]))
	vy := int(math.Round(col[3]))

	px0 := p.objs[0].pos.x
	vx0 := p.objs[0].vel.x
	pz0 := p.objs[0].pos.z
	vz0 := p.objs[0].vel.z
	px1 := p.objs[1].pos.x
	vx1 := p.objs[1].vel.x
	pz1 := p.objs[1].pos.z
	vz1 := p.objs[1].vel.z

	t0 := (px0 - px) / (vx - vx0)
	t1 := (px1 - px) / (vx - vx1)

	vz := (pz1 - pz0 + vz1 * t1 - vz0 * t0) / (t1 - t0)
	pz := pz1 + vz1 * t1 - vz * t1

	fmt.Println(px, ",", py, ",", pz, "@", vx, ",", vy, ",", vz)
	fmt.Println("Answer to part 2:", px + py + pz)
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	s := bufio.NewScanner(f)
	r := strings.NewReplacer(",", " ", "@", " ")
	for s.Scan() {
		objData := strings.Fields(r.Replace(s.Text()))
		o := Hailstone{
			Vec{ Atoi(objData[0]), Atoi(objData[1]), Atoi(objData[2]) },
			Vec{ Atoi(objData[3]), Atoi(objData[4]), Atoi(objData[5]) },
		}
		p.objs = append(p.objs, o)
	}
	check(s.Err())
	return p
}

func Atoi(txt string) int {
	i, err := strconv.Atoi(txt)
	check(err)
	return i
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	puzzle := load(flag.Arg(0))

	puzzle.GetAnswer()
}
