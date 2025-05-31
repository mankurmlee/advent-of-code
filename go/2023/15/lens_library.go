package main

import (
	"flag"
	"os"
	"bufio"
	"strings"
	"fmt"
	"slices"
	"strconv"
)

type Lens struct {
	label       string
	focalLength int
}

type Box []Lens

type Conduit [256]Box

func (b *Box) AddLens(label string, focalLength int) {
	i := slices.IndexFunc(*b, func(l Lens) bool { return l.label == label })
	if i >= 0 {
		(*b)[i].focalLength = focalLength
	} else {
		*b = append(*b, Lens{label, focalLength})
	}
}

func (b *Box) RemoveLens(label string) {
	i := slices.IndexFunc(*b, func(l Lens) bool { return l.label == label })
	if i >= 0 {
		*b = append((*b)[:i], (*b)[i+1:]...)
	}
}

func (c *Conduit) HASHMAP(step string) {
	label, fl, found := strings.Cut(step, "=")
	if found {
		focalLength, _ := strconv.Atoi(fl)
		(*c)[Hash(label)].AddLens(label, focalLength)
	} else {
		label, _ = strings.CutSuffix(step, "-")
		(*c)[Hash(label)].RemoveLens(label)
	}
}

func (c *Conduit) GetFocusingPower() (total int) {
	for box, b := range *c {
		for lens, l := range b {
			total += (box + 1) * (lens + 1) * l.focalLength
		}
	}
	return total
}

func Hash(seq string) (hash int) {
	for _, ch := range seq {
		hash = (hash + int(ch)) * 17 % 256
	}
	return hash
}

func LoadSequence(filename string) []string {
	f, _ := os.Open(filename)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	return strings.Split(s.Text(), ",")
}

func main() {
	flag.Parse()

	var hashsum int
	c := &Conduit{}
	for _, step := range LoadSequence(flag.Arg(0)) {
		hashsum += Hash(step)
		c.HASHMAP(step)
	}

	fmt.Println("Part 1 answer:", hashsum)
	fmt.Println("Part 2 answer:", c.GetFocusingPower())
}
