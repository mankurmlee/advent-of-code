package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Reindeer struct {
	name  string
	speed int
	fly   int
	rest  int
}

func (r Reindeer) DistanceAfter(t int) (dist int) {
	periodDist := r.speed * r.fly
	period := r.fly + r.rest
	quot := t / period
	rem := t % period
	dist = quot * periodDist
	if rem >= r.fly {
		dist += periodDist
	} else {
		dist += r.speed * rem
	}
	return dist
}

func load() (out []Reindeer, after int) {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	re := regexp.MustCompile(`\d+`)
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	after = atoi(s.Text())
	for s.Scan() {
		line := s.Text()
		name := strings.Fields(line)[0]
		data := re.FindAllString(line, -1)
		speed := atoi(data[0])
		fly := atoi(data[1])
		rest := atoi(data[2])
		out = append(out, Reindeer{name, speed, fly, rest})
	}
	check(s.Err())
	return out, after
}

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

func main() {
	r, t := load()
	dist, _ := statusAfter(r, t)
	fmt.Println("Part 1:", dist)
	fmt.Println("Part 2:", partTwo(r, t))
}

func partTwo(all []Reindeer, t int) (out int) {
	scores := map[string]int{}
	for i := 1; i < t+1; i++ {
		_, l := statusAfter(all, i)
		for _, name := range l {
			s, ok := scores[name]
			if !ok {
				s = 0
			}
			scores[name] = s + 1
		}
	}
	for _, v := range scores {
		if v > out {
			out = v
		}
	}
	return out
}

func statusAfter(all []Reindeer, t int) (out int, leaders []string) {
	for _, r := range all {
		d := r.DistanceAfter(t)
		if d > out {
			out = d
			leaders = []string{r.name}
		} else if d == out {
			leaders = append(leaders, r.name)
		}
	}
	return out, leaders
}
