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

func fileinput() (out []string) {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		out = append(out, s.Text())
	}
	check(s.Err())
	return out
}

type Sue struct {
	props map[string]int
}

func (left Sue) Contains(right Sue) bool {
	props := left.props
	for k, v := range right.props {
		value, exists := props[k]
		if !exists || v != value {
			return false
		}
	}
	return true
}

func (left Sue) CustomMatch(right Sue) bool {
	props := left.props
	for k, v := range right.props {
		value, exists := props[k]
		if !exists {
			return false
		}
		switch k {
		case "cats", "trees":
			if v <= value {
				return false
			}
		case "pomeranians", "goldfish":
			if v >= value {
				return false
			}
		default:
			if v != value {
				return false
			}
		}
	}
	return true
}

func parse(inp []string) (out []Sue) {
	re := regexp.MustCompile(`[:,]`)
	for _, l := range inp {
		l = strings.TrimSpace(l)
		l = re.ReplaceAllLiteralString(l, "")
		words := strings.Fields(l)
		n := len(words)
		props := map[string]int{}
		for i := 2; i < n; i += 2 {
			props[words[i]] = atoi(words[i+1])
		}
		out = append(out, Sue{props})
	}
	return out
}

func partOne(inp []Sue) int {
	sue0 := inp[0]
	for i, s := range inp {
		if i != 0 && sue0.Contains(s) {
			return i
		}
	}
	return 0
}

func partTwo(inp []Sue) int {
	sue0 := inp[0]
	for i, s := range inp {
		if i != 0 && sue0.CustomMatch(s) {
			return i
		}
	}
	return 0
}

func main() {
	sues := parse(fileinput())
	fmt.Println("Part 1:", partOne(sues))
	fmt.Println("Part 2:", partTwo(sues))
}
