package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
)

var dirVec = map[byte]Vertex{
	'^': Vertex{  0, -1 },
	'>': Vertex{  1,  0 },
	'v': Vertex{  0,  1 },
	'<': Vertex{ -1,  0 },
}

type Vertex struct { x, y int }

func (v1 Vertex) Add(v2 Vertex) Vertex {
	return Vertex{ v1.x + v2.x, v1.y + v2.y }
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) []byte {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	return []byte(s.Text())
}

func santaDelivery(dirs []byte) int {
	houses := make(map[Vertex]bool)
	pos := Vertex{0, 0}
	houses[pos] = true
	for _, d := range dirs {
		pos = pos.Add(dirVec[d])
		houses[pos] = true
	}
	return len(houses)
}

func roboDelivery(dirs []byte) int {
	houses := make(map[Vertex]bool)
	santa := Vertex{0, 0}
	robo  := Vertex{0, 0}
	houses[santa] = true
	for i, d := range dirs {
		if i % 2 == 0 {
			santa = santa.Add(dirVec[d])
			houses[santa] = true
		} else {
			robo = robo.Add(dirVec[d])
			houses[robo] = true
		}
	}
	return len(houses)
}

func main() {
	flag.Parse()
	dirs := load(flag.Arg(0))

	fmt.Println("Part 1:", santaDelivery(dirs))
	fmt.Println("Part 2:", roboDelivery(dirs))
}
