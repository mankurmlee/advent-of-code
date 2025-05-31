package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type Dir struct {
	parent      *Dir
	children    map[string]*Dir
	filesize    int
}

type SizeTable map[*Dir]int

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func load(filename string) (root *Dir) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)

	root = &Dir{}
	cwd := root

	for s.Scan() {
		data := strings.Fields(s.Text())
		switch data[0] {
		case "$":
			if data[1] == "cd" {
				switch data[2] {
				case "/":
					cwd = root
				case "..":
					cwd = cwd.parent
				default:
					cwd, _ = cwd.children[data[2]]
				}
			}
		case "dir":
			if cwd.children == nil {
				cwd.children = make(map[string]*Dir)
			}
			cwd.children[data[1]] = &Dir{ parent: cwd }
		default: // file listing
			cwd.filesize += Atoi(data[0])
		}
	}

	check(s.Err())
	return root
}

func (t *SizeTable) Create(dir *Dir) {
	if *t == nil {
		*t = make(map[*Dir]int)
	}
	t.GetSize(dir)
}

func (t *SizeTable) GetSize(dir *Dir) (size int) {
	size = dir.filesize
	for _, d := range dir.children {
		size += t.GetSize(d)
	}
	(*t)[dir] = size
	return size
}

func main() {
	flag.Parse()
	root := load(flag.Arg(0))

	var sizes SizeTable
	sizes.Create(root)

	var sum int
	for _, v := range sizes {
		if v <= 100000 {
			sum += v
		}
	}
	fmt.Println("Part 1:", sum)

	freeup := sizes[root] - 40000000
	var min int
	for _, s := range sizes {
		if s >= freeup && (min == 0 || min > s) {
			min = s
		}
	}
	fmt.Println("Part 2:", min)
}
