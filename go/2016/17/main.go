package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
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

func main() {
	flag.Parse()
	p := readFile(flag.Arg(0))[0]
	solve(p, false)
	solve(p, true)
}

func solve(salt string, longest bool) {
	var best int
	q := []PathState{{Vec{0, 0}, ""}}
	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		if s.pos.x == 3 && s.pos.y == 3 {
			if !longest {
				fmt.Println(s.path)
				return
			}
			if len(s.path) > best {
				best = len(s.path)
			}
			continue
		}
		locked := getLocked(salt + s.path)
		for i, d := range DIRS {
			if locked[i] {
				continue
			}
			s1 := PathState{
				s.pos.Add(d),
				s.path + DIRN[i],
			}
			if s1.pos.x < 0 || s1.pos.x > 3 || s1.pos.y < 0 || s1.pos.y > 3 {
				continue
			}
			q = append(q, s1)
		}
	}
	fmt.Println(best)
}

func getLocked(s string) [4]bool {
	hash := md5.Sum([]byte(s))
	h := hex.EncodeToString(hash[:])
	return [4]bool{
		h[0] < 'b',
		h[1] < 'b',
		h[2] < 'b',
		h[3] < 'b',
	}
}

type Vec struct{ x, y int }

func (v Vec) Add(o Vec) Vec {
	return Vec{v.x + o.x, v.y + o.y}
}

var DIRS = [4]Vec{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
var DIRN = [4]string{"U", "D", "L", "R"}

type PathState struct {
	pos  Vec
	path string
}
