package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Rule [2]int
type Manual []int

type ManIndex struct {
	pages map[int]int
	mid   int
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

func findNums(s string) (nums []int) {
	r := regexp.MustCompile(`\d+`)
	for _, v := range r.FindAllString(s, -1) {
		nums = append(nums, atoi(v))
	}
	return nums
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

func load(filename string) (rules []Rule, mans []Manual) {
	linebreaks := 0
	for _, s := range readFile(filename) {
		nums := findNums(s)
		if len(nums) == 0 {
			linebreaks++
			if linebreaks > 1 {
				break
			}
			continue
		}
		if linebreaks == 0 {
			rules = append(rules, Rule{nums[0], nums[1]})
		} else if linebreaks == 1 {
			mans = append(mans, nums)
		}
	}
	return rules, mans
}

func main() {
	flag.Parse()
	rules, mandata := load(flag.Arg(0))
	mans := index(mandata)

	fmt.Println(partOne(rules, mans))
	fmt.Println(partTwo(rules, mans))
}

func partTwo(rules []Rule, mans []ManIndex) (tot int) {
	for _, m := range mans {
		tot += fixed(rules, m)
	}
	return tot
}

func fixed(rules []Rule, m ManIndex) int {
	changed := false
	pages := clone(m.pages)
	for sortPages(rules, pages) {
		changed = true
	}
	if !changed {
		return 0
	}
	mp := len(pages) / 2
	for k, v := range pages {
		if v == mp {
			return k
		}
	}
	return 0
}

func sortPages(rules []Rule, pages map[int]int) (changed bool) {
	for _, r := range rules {
		p1, p2 := r[0], r[1]
		i, e1 := pages[p1]
		j, e2 := pages[p2]
		if e1 && e2 && i > j {
			changed = true
			pages[p1], pages[p2] = j, i
		}
	}
	return changed
}

func clone(orig map[int]int) map[int]int {
	res := make(map[int]int)
	for k, v := range orig {
		res[k] = v
	}
	return res
}

func partOne(rules []Rule, mans []ManIndex) (tot int) {
	for _, m := range mans {
		if ordered(rules, m) {
			tot += m.mid
		}
	}
	return tot
}

func ordered(rules []Rule, m ManIndex) bool {
	p := m.pages
	for _, r := range rules {
		i, e1 := p[r[0]]
		j, e2 := p[r[1]]
		if e1 && e2 && i > j {
			return false
		}
	}
	return true
}

func index(data []Manual) (mans []ManIndex) {
	for _, m := range data {
		var man ManIndex
		midpoint := len(m) / 2
		man.pages = make(map[int]int)
		for i, p := range m {
			man.pages[p] = i
			if i == midpoint {
				man.mid = p
			}
		}
		mans = append(mans, man)
	}
	return mans
}
