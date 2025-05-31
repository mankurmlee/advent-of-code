package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type Block struct {
	pos  int
	size int
}

type File struct {
	Block
	id int
}

type Disk struct {
	files []File
	space []Block
}

func (d Disk) Defrag() {
	for j := len(d.files) - 1; j >= 0; j-- {
		size := d.files[j].size
		for i, s := range d.space {
			if s.pos >= d.files[j].pos {
				break
			}
			if s.size >= size {
				d.files[j].pos = s.pos
				d.space[i].pos += size
				d.space[i].size -= size
				break
			}
		}
	}
}

func (d Disk) Checksum() {
	sum := 0
	for _, f := range d.files {
		for i := range f.size {
			sum += f.id * (f.pos + i)
		}
	}
	fmt.Println(sum)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	check(err)
	return i
}

func readFile(filename string) (lines []string) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	check(s.Err())
	return lines
}

func load(filename string) (puzzle []int) {
	d := readFile(filename)
	for _, l := range d {
		for _, r := range l {
			puzzle = append(puzzle, atoi(string(r)))
		}
	}
	return puzzle
}

func main() {
	flag.Parse()
	puzzle := load(flag.Arg(0))
	partOne(puzzle)
	partTwo(puzzle)
}

func partTwo(puzzle []int) {
	d := createDisk(puzzle)
	d.Defrag()
	d.Checksum()
}

func createDisk(puzzle []int) (d Disk) {
	j := 0
	for i, v := range puzzle {
		if v == 0 {
			continue
		}
		l := Block{j, v}
		j += v
		if i%2 == 0 {
			d.files = append(d.files, File{l, i >> 1})
		} else {
			d.space = append(d.space, l)
		}
	}
	return d
}

func partOne(puzzle []int) {
	dm := diskmap(puzzle)
	defrag(dm)
	checksum(dm)
}

func checksum(dm map[int]int) {
	sum := 0
	for k, v := range dm {
		sum += k * v
	}
	fmt.Println(sum)
}

func defrag(dm map[int]int) {
	keys := mapkeys(dm)
	slices.Sort(keys)
	end := len(keys) - 1
	j := keys[end]
	for i := 0; i < j; i++ {
		if _, ok := dm[i]; ok {
			continue
		}
		dm[i] = dm[j]
		delete(dm, j)
		end--
		j = keys[end]
	}
}

func mapkeys(m map[int]int) (keys []int) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func diskmap(puzzle []int) map[int]int {
	dm := make(map[int]int)
	j := 0
	for i, v := range puzzle {
		if i%2 == 0 {
			id := i >> 1
			for range v {
				dm[j] = id
				j++
			}
		} else {
			j += v
		}
	}
	return dm
}
