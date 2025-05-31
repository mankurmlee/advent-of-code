package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"encoding/json"
	"cmp"
	"slices"
)

type ElemType int
const (
	UNKNOWN ElemType = iota
	INT
	OBJECT
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) (p [][]interface{}) {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if line != "" {
			p = append(p, unserialise(line))
		}
	}
	check(s.Err())
	return p
}

func unserialise(input string) (l []interface{}) {
	err := json.Unmarshal([]byte(input), &l)
	check(err)
	for i, o := range l {
		l[i] = parseObject(o)
	}
	return l
}

func parseObject(o interface{}) interface{} {
	switch value := o.(type) {
	case float64:
		return int(value)
	case []interface{}:
		list := make([]interface{}, len(value))
		for i, o := range value {
			list[i] = parseObject(o)
		}
		return list
	default:
		return value
	}
}

func typeof(o interface{}) ElemType {
	switch o.(type) {
	case int:
		return INT
	case interface{}:
		return OBJECT
	default:
		return UNKNOWN
	}
}

func compare(a, b []interface{}) int {
	n1, n2 := len(a), len(b)
	for i := 0; i < n1 && i < n2; i++ {
		u, v := a[i], b[i]
		uType, vType := typeof(u), typeof(v)
		if uType == INT && vType == INT {
			if r := cmp.Compare(u.(int), v.(int)); r != 0 {
				return r
			}
			continue
		}
		if uType == INT {
			u = []interface{}{ u }
		} else if vType == INT {
			v = []interface{}{ v }
		}
		if r := compare(u.([]interface{}), v.([]interface{})); r != 0 {
			return r
		}
	}
	return cmp.Compare(n1, n2)
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))

	var sumIndices int
	n := len(p)
	for i := 0; i < n; i += 2 {
		if compare(p[i], p[i+1]) < 0 {
			sumIndices += i / 2 + 1
		}
	}
	fmt.Println("Part 1:", sumIndices)

	div1 := []interface{}{[]interface{}{ 2 }}
	div2 := []interface{}{[]interface{}{ 6 }}
	p = append(p, div1, div2)
	slices.SortFunc(p, compare)
	i1 := slices.IndexFunc(p, func(x []interface{}) bool {
		return compare(x, div1) == 0
	})
	i2 := slices.IndexFunc(p, func(x []interface{}) bool {
		return compare(x, div2) == 0
	})
	fmt.Println("Part 2:", (i1 + 1) * (i2 + 1))
}
