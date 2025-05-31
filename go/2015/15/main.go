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

type Ingredient struct {
	name  string
	props map[string]int
}

type Mixer struct {
	ingredients []Ingredient
}

func (m Mixer) FindBest(n int, bowl []Ingredient, limit int) (out int, kcal int) {
	idx := len(bowl)
	raw := m.ingredients[idx]
	if len(m.ingredients)-idx == 1 {
		s, c := Score(Mix(bowl, raw, n))
		if limit > 0 && c != limit {
			return 0, 0
		}
		return s, c
	}
	for i := range n {
		s, c := m.FindBest(n-i, Mix(bowl, raw, i), limit)
		if limit == 0 || c == limit {
			if s > out {
				out = s
				kcal = c
			}
		}
	}
	return out, kcal
}

func Mix(bowl []Ingredient, raw Ingredient, n int) []Ingredient {
	cooked := Ingredient{raw.name, map[string]int{}}
	for k, v := range raw.props {
		cooked.props[k] = v * n
	}
	out := make([]Ingredient, len(bowl))
	copy(out, bowl)
	return append(out, cooked)
}

func Score(bowl []Ingredient) (score int, kcal int) {
	score = 1
	for k := range bowl[0].props {
		sum := 0
		for _, v := range bowl {
			sum += v.props[k]
		}
		if sum <= 0 {
			return 0, 0
		}
		if k == "calories" {
			kcal = sum
		} else {
			score *= sum
		}
	}
	return score, kcal
}

func load() (out []Ingredient) {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	re := regexp.MustCompile(`[:,]`)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := re.ReplaceAllLiteralString(s.Text(), "")
		words := strings.Fields(line)
		n := len(words)
		props := map[string]int{}
		for i := 1; i < n; i += 2 {
			props[words[i]] = atoi(words[i+1])
		}
		out = append(out, Ingredient{words[0], props})
	}
	check(s.Err())
	return out
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
	m := Mixer{load()}
	score, _ := m.FindBest(100, []Ingredient{}, 0)
	fmt.Println("Part 1:", score)
	score, _ = m.FindBest(100, []Ingredient{}, 500)
	fmt.Println("Part 2:", score)
}
