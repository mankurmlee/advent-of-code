package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
)

type Block struct {
	rocks   []bool
	w, h    int
}

type Blocks struct {
	blocks  []Block
	i       int
	len     int
}

type Stream struct {
	patt    []rune
	i       int
	len     int
}

type MovingBlock struct {
	x, y    int
	block   *Block
}

type Grid struct {
	itBlock     *Blocks
	itStream    *Stream
	atRest      []*MovingBlock
	height      int
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) (p Stream) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	p.patt = []rune(s.Text())
	p.len  = len(p.patt)
	return p
}

func createBlocks() (list Blocks) {
	list.blocks = append(list.blocks, createBlock(`
####`))
	list.blocks = append(list.blocks, createBlock(`
 #
###
 #`))
	list.blocks = append(list.blocks, createBlock(`
  #
  #
###`))
	list.blocks = append(list.blocks, createBlock(`
#
#
#
#`))
	list.blocks = append(list.blocks, createBlock(`
##
##`))
	list.len = len(list.blocks)
	return list
}

func createBlock(data string) (b Block) {
	s := strings.Split(data, "\n")[1:]
	for _, v := range s {
		if len(v) > b.w {
			b.w = len(v)
		}
	}
	b.h = len(s)

	for y := 0; y < b.h; y++ {
		rowlen := len(s[y])
		for x := 0; x < b.w; x++ {
			b.rocks = append(b.rocks, x < rowlen && s[y][x] == '#')
		}
	}

	return b
}

func (it *Stream) Next() (x rune) {
	x = it.patt[it.i]
	it.i = (it.i + 1) % it.len
	return x
}

func (it *Blocks) Next() (x *Block) {
	x = &it.blocks[it.i]
	it.i = (it.i + 1) % it.len
	return x
}

func (b *MovingBlock) SlideLeft(g *Grid) {
	b.x--
	if b.x < 0 {
		b.x++
		return
	}
	for _, o := range g.atRest {
		if b.CollidesWith(o) {
			b.x++
			return
		}
	}
}

func (b *MovingBlock) SlideRight(g *Grid) {
	b.x++
	if b.x + b.block.w - 1 >= 7 {
		b.x--
		return
	}
	for _, o := range g.atRest {
		if b.CollidesWith(o) {
			b.x--
			return
		}
	}
}

func (b *MovingBlock) Drop(g *Grid) bool {
	b.y--
	for _, o := range g.atRest {
		if b.CollidesWith(o) {
			b.y++
			g.atRest = append(g.atRest, b)
			return false
		}
	}
	if b.y - b.block.h + 1 < 1 {
		b.y++
		g.atRest = append(g.atRest, b)
		return false
	}
	return true
}

func (b *MovingBlock) CollidesWith(o *MovingBlock) bool {
	if  b.y - b.block.h + 1 > o.y ||
		o.y - o.block.h + 1 > b.y ||
		b.x + b.block.w - 1 < o.x ||
		o.x + o.block.w - 1 < b.x {
		return false
	}

	surf := make(map[int]bool)

	w, h := b.block.w, b.block.h
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if b.block.rocks[y * w + x] {
				surf[b.x + x + (b.y - y) * 7] = true
			}
		}
	}

	w, h = o.block.w, o.block.h
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if o.block.rocks[y * w + x] {
				if surf[o.x + x + (o.y - y) * 7] {
					return true
				}
			}
		}
	}

	return false
}

func (g *Grid) DropBlock() {
	nextBlock := g.itBlock.Next()
	b := MovingBlock{ 2, g.height + nextBlock.h + 3, nextBlock }
	for {
		if g.itStream.Next() == '<' {
			b.SlideLeft(g)
		} else {
			b.SlideRight(g)
		}
		if !b.Drop(g) {
			g.height = max(g.height, b.y)
			if len(g.atRest) > 100 {
				g.atRest = g.atRest[50:]  // prune old blocks
			}
			break
		}
	}
}

func (g *Grid) Reset() {
	g.itStream.i = 0
	g.itBlock.i = 0
	g.atRest = []*MovingBlock{}
	g.height = 0
}

func (g *Grid) DropHeight(n int) int {
	g.Reset()
	for i := 0; i < n; i++ {
		g.DropBlock()
	}
	return g.height
}

func (g *Grid) FlipHeight(n int) (int, int) {
	g.Reset()
	var drops, flips, before int
	for flips < n {
		before = g.itStream.i
		g.DropBlock()
		drops++
		if before > g.itStream.i {
			flips++
		}
	}
	return drops, g.height
}

func (g *Grid) ScaledHeight(target int) int {
	nblocks := g.itBlock.len
	d1, h1 := g.FlipHeight(nblocks)
	fmt.Println("Height after", d1, "drops is", h1)
	d2, h2 := g.FlipHeight(nblocks*2)
	fmt.Println("Height after", d2, "drops is", h2)

	dDrops  := d2 - d1
	dHeight := h2 - h1
	fmt.Println("Every", dDrops, "drops, increases height by", dHeight)

	dropOff   := d1 - dDrops
	heightOff := h1 - dHeight
	fmt.Println("Initial drop offset of", dropOff, "offsets height by", heightOff)

	scale      := (target - dropOff) / dDrops // integer division
	dropNear   := scale * dDrops + dropOff
	heightNear := scale * dHeight + heightOff
	fmt.Println("Height after", dropNear, "drops is", heightNear)

	remDrop    := target - dropNear
	remHeight  := g.DropHeight(d1 + remDrop) - h1
	fmt.Println("Remaining", remDrop, "drops, increases height by", remHeight)

	return heightNear + remHeight
}

func main() {
	flag.Parse()
	itStream := load(flag.Arg(0))
	itBlock  := createBlocks()

	g := &Grid{ itStream : &itStream, itBlock: &itBlock }
	fmt.Println("Part 1: Height after 2022 drops is", g.DropHeight(2022))
	fmt.Println("Part 2: Height after a trillion drops is", g.ScaledHeight(1000000000000))
}
