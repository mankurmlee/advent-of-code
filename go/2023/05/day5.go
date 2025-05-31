package main


import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)


type Rule struct {
	start  int
	end    int
	offset int
}


type Almanac [][]Rule

func (a Almanac) GetLocation(seed int) int {
	s := seed
	for _, v := range a {
		for _, r := range v {
			if s >= r.start && s <= r.end {
				s += r.offset
				break
			}
		}
	}
	return s
}


// Parses the input file
func parseInput(filename string) ([]int, Almanac) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	s := bufio.NewScanner(f)

	var seeds   []int
	var alm     Almanac
	var row     []Rule

	for s.Scan() {
		l := s.Text()

		// Empty line
		if l == "" {
			continue
		}

		// Seeds
		if strings.HasPrefix(l, "seeds:") {
			for _, part := range strings.Split(l[7:], " ") {
				seed, _ := strconv.Atoi(part)
				seeds = append(seeds, seed)
			}
			continue
		}

		// Maps
		if strings.HasSuffix(l, "map:") {
			alm = append(alm, row)
			row = []Rule{}
			continue
		}

		// Ranges
		rule := parseRule(l)
		row = append(row, rule)
	}

	// Make sure to add the last row
	alm = append(alm, row)

	check(s.Err())

	return seeds, alm[1:]
}


// Parses and returns a rule
func parseRule(text string) Rule {
	parts := strings.Split(text, " ")

	dest, err := strconv.Atoi(parts[0])
	check(err)

	src, err := strconv.Atoi(parts[1])
	check(err)

	length, err := strconv.Atoi(parts[2])
	check(err)

	rule := Rule{
		start:  src,
		end:    src+length-1,
		offset: dest-src,
	}

	return rule
}


// Lazy error checking
func check(err error) {
	if err != nil {
		panic(err)
	}
}


func main() {
	//flag.Parse()
	//filename := flag.Args()[0]
	filename := "sample.txt"
	//filename := "myinput.txt"

	seeds, alm := parseInput(filename)

	var min, loc int
	for i, seed := range seeds {
		loc = alm.GetLocation(seed)
		if i == 0 || loc < min {
			min = loc
		}
	}

	fmt.Println("Lowest Location:", min)
}
