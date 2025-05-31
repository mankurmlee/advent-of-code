package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load() string {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	return s.Text()
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func partOne(s string) int {
	re := regexp.MustCompile(`-?\d+`)
	data := re.FindAllString(s, -1)
	sum := 0
	for _, s := range data {
		sum += atoi(s)
	}
	return sum
}

func partTwo(s string) int {
	var data any
	err := json.Unmarshal([]byte(s), &data)
	check(err)
	return parse(data)
}

func parse(data any) int {
	switch v := data.(type) {
	case map[string]interface{}:
		sum := 0
		for _, value := range v {
			switch s := value.(type) {
			case string:
				if s == "red" {
					return 0
				}
			}
			sum += parse(value)
		}
		return sum
	case []interface{}:
		sum := 0
		for _, item := range v {
			sum += parse(item)
		}
		return sum
	case float64:
		return int(v)
	case string:
		return 0
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
	return 0
}

func main() {
	s := load()
	fmt.Println("Part 1:", partOne(s))
	fmt.Println("Part 2:", partTwo(s))
}
