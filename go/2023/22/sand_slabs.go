package main

import (
	"fmt"
	"flag"
	"os"
	"bufio"
	"strings"
	"strconv"
	"slices"
)

type Vec struct {
	x, y, z int
}

type Brick struct {
	id int
	u, v Vec
}

// Calculate new position of brick if it were to fall
func (old Brick) Fall() (new Brick) {
	new = old
	if new.u.z <= 0 || new.v.z <= 0 {
		// Brick is on the ground
		return new
	}
	new.u.z--; 	new.v.z--
	return new
}

func (p Brick) Collides(q Brick) bool {
	if  p.v.x < q.u.x || q.v.x < p.u.x ||
		p.v.y < q.u.y || q.v.y < p.u.y ||
		p.v.z < q.u.z || q.v.z < p.u.z {
		return false
	}
	return true
}

type Puzzle struct {
	bricks      []Brick
	topIndex    map[int][]int
	botIndex    map[int][]int
}

func (p Puzzle) IsChildOf(child int, parents []int) bool {

	// Create list of possible paramours (illegitimate parents)
	var others []int
	plvl := p.bricks[child].u.z - 1
	for _, pid := range p.topIndex[plvl] {
		if !slices.Contains(parents, pid) {
			others = append(others, pid)
		}
	}

	// If any of these 'other' bricks are a parent of the child then
	// this child is illegitimate
	fallen := p.bricks[child].Fall()
	for _, pid := range others {
		if fallen.Collides(p.bricks[pid]) {
			return false
		}
	}

	// If we have eliminated all paramours then the child must be
	// supported only by legitimate parents
	return true
}

// Gets legitimate children of parents
// A legitimate child is where the child is ONLY supported by the specified parent(s)
func (p Puzzle) GetChildren(parents []int) (children []int) {
	clvl := p.bricks[parents[0]].v.z + 1
	for _, cid := range p.botIndex[clvl] {
		if p.IsChildOf(cid, parents) {
			children = append(children, cid)
		}
	}
	return children
}

func (p Puzzle) GetDescendants(id int) (desc []int) {
	level := p.bricks[id].v.z
	queue := make(map[int][]int)
	queue[level] = append(queue[level], id)
	for len(queue) > 0 {
		level = getLowestLevel(&queue)
		for _, cid := range p.GetChildren(queue[level]) {
			clvl := p.bricks[cid].v.z
			if !slices.Contains(queue[clvl], cid) {
				queue[clvl] = append(queue[clvl], cid)
				desc = append(desc, cid)
			}
		}
		delete(queue, level)
	}
	return desc
}

func (p Puzzle) SumDescendants() (sum int) {
	for _, parent := range p.bricks {
		sum += len(p.GetDescendants(parent.id))
	}
	return sum
}

func (p *Puzzle) BuildIndex() {
	p.topIndex = make(map[int][]int)
	p.botIndex = make(map[int][]int)
	for _, b := range p.bricks {
		p.topIndex[b.v.z] = append(p.topIndex[b.v.z], b.id)
		p.botIndex[b.u.z] = append(p.botIndex[b.u.z], b.id)
	}
}

func (p *Puzzle) UpdateIndex(before, after Brick) {
	if len(p.topIndex[before.v.z]) == 1 {
		delete(p.topIndex, before.v.z)
	} else {
		i := slices.Index(p.topIndex[before.v.z], before.id)
		p.topIndex[before.v.z] = append(p.topIndex[before.v.z][:i], p.topIndex[before.v.z][i+1:]...)
	}
	if len(p.botIndex[before.u.z]) == 1 {
		delete(p.botIndex, before.u.z)
	} else {
		i := slices.Index(p.botIndex[before.u.z], before.id)
		p.botIndex[before.u.z] = append(p.botIndex[before.u.z][:i], p.botIndex[before.u.z][i+1:]...)
	}
	p.topIndex[after.v.z] = append(p.topIndex[after.v.z], after.id)
	p.botIndex[after.u.z] = append(p.botIndex[after.u.z], after.id)
}

func (p Puzzle) CountRemovable() (count int) {
	for _, b := range p.bricks {
		cids := p.GetChildren([]int{b.id})
		if len(cids) == 0 {
			count++
		}
	}
	return count
}

// Simulate all bricks falling and settling
func (p Puzzle) Settle() {
	var settled bool
	for {
		settled = true
		for i, b := range p.bricks {
			s := p.SettleBrick(b)
			if s != b {
				settled = false
				p.bricks[i] = s
				p.UpdateIndex(b, s)
			}
		}
		if settled {
			break
		}
	}
}

// Simulate a single brick falling as far as it can
func (p Puzzle) SettleBrick(b Brick) (s Brick) {
	var old Brick
	new := b
	for {
		old = new
		new = old.Fall()
		if old == new || p.Collides(new) {
			break
		}
	}
	return old
}

// Check if brick collides with anything other brick
func (p Puzzle) Collides(b Brick) bool {
	lvl := b.u.z
	for _, pid := range p.topIndex[lvl] {
		if b.id == pid {
			continue
		}
		if b.Collides(p.bricks[pid]) {
			return true
		}
	}
	return false
}

func getLowestLevel(dict *map[int][]int) int {
	lowest := -1
	for k := range *dict {
		if lowest < 0 || k < lowest {
			lowest = k
		}
	}
	return lowest
}

func load(filename string) (p Puzzle) {
	f, err := os.Open(filename)
	check(err)
	s := bufio.NewScanner(f)
	for id := 0; s.Scan(); id++ {
		brickData := strings.Split(s.Text(), "~")
		u := parseVec(strings.Split(brickData[0], ","))
		v := parseVec(strings.Split(brickData[1], ","))
		p.bricks = append(p.bricks, Brick{id, u, v})
	}
	check(s.Err())
	return p
}

func parseVec(vecData []string) (v Vec) {
	v.x = Atoi(vecData[0])
	v.y = Atoi(vecData[1])
	v.z = Atoi(vecData[2])
	return v
}

func Atoi(intData string) (n int) {
	n, err := strconv.Atoi(intData)
	check(err)
	return n
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	puzzle := load(flag.Arg(0))
	puzzle.BuildIndex()

	fmt.Println("Settling bricks...")
	puzzle.Settle()
	fmt.Println("Bricks settled")

	n := puzzle.CountRemovable()
	fmt.Println("Answer to part 1:", n)

	sum := puzzle.SumDescendants()
	fmt.Println("Answer to part 2:", sum)
}
