package main

import (
	"flag"
	"os"
	"bufio"
	"strings"
	"fmt"
)

type State int
const (
	OFF State = iota
	ON
)

type Pulse int
const (
	LOW Pulse = iota
	HIGH
)

type Msg struct {
	src  int
	pulse Pulse
	dest int
}

type Module interface {
	Propagate(rcv Msg) (outbox []Msg)
	Print()
	GetName() string
}

type ModuleData struct {
	index *ModuleList
	id    int
	name  string
	dests []int
}

type ModuleDict map[string]int
type ModuleList []Module

type FlipFlop struct {
	ModuleData
	state State
}

type Conjunction struct {
	ModuleData
	lastPulse map[int]Pulse
}

type Broadcast struct {
	ModuleData
}

type Button struct {
	ModuleData
}

func (m ModuleData) GetName() string {
	return m.name
}

func (m *FlipFlop) Propagate(rcv Msg) (outbox []Msg) {
	if rcv.pulse == HIGH {
		return outbox
	}

	var pulse Pulse
	if m.state == OFF {
		m.state = ON
		pulse = HIGH
	} else {
		m.state = OFF
		pulse = LOW
	}

	for _, dest := range m.dests {
		outbox = append(outbox, Msg{m.id, pulse, dest})
	}
	return outbox
}

func (m *Conjunction) Propagate(rcv Msg) (outbox []Msg) {
	m.lastPulse[rcv.src] = rcv.pulse

	pulse := LOW
	for _, v := range m.lastPulse {
		if v == LOW {
			pulse = HIGH
			break
		}
	}

	for _, dest := range m.dests {
		outbox = append(outbox, Msg{m.id, pulse, dest})
	}
	return outbox
}

func (m *Broadcast) Propagate(rcv Msg) (outbox []Msg) {
	for _, dest := range m.dests {
		outbox = append(outbox, Msg{m.id, rcv.pulse, dest})
	}
	return outbox
}

func (m *Button) Propagate(rcv Msg) (outbox []Msg) {
	outbox = append(outbox, Msg{m.id, rcv.pulse, m.dests[0]})
	return outbox
}

func (m *FlipFlop) Print() {
	fmt.Println("%", m.id, "=", m.state)
}

func (m *Conjunction) Print() {
	fmt.Println("&", m.id, "=", m.lastPulse)
}

func (m *Broadcast) Print() {}

func (m *Button) Print() {}

func createModules(mods [][]string) (lut ModuleDict, modules ModuleList) {
	var broadcastid int

	lut = make(ModuleDict)
	for id, m := range mods {
		if m[0] == "b" {
			broadcastid = id
		}
		name := m[1]
		lut[name] = id
	}

	srcs := make(map[int][]int)

	for id, m := range mods {
		dests := m[2:]
		for _, name := range dests {
			d, ok := lut[name]
			if !ok {
				continue
			}
			srcs[d] = append(srcs[d], id)
		}
	}

	n := len(mods)
	modules = make(ModuleList, n + 1)
	for id, m := range mods {
		typ   := m[0]
		name  := m[1]

		dests := m[2:]
		var dids []int
		for _, name := range dests {
			id, ok := lut[name]
			if !ok {
				id = -1
			}
			dids = append(dids, id)
		}


		md := ModuleData{
			index : &modules,
			id    : id,
			name  : name,
			dests : dids,
		}

		switch typ {
		case "b":
			modules[id] = &Broadcast{ ModuleData: md }
		case "%":
			modules[id] = &FlipFlop{ ModuleData: md }
		case "&":
			lp := make(map[int]Pulse)
			for _, src := range srcs[id] {
				lp[src] = LOW
			}
			modules[id] = &Conjunction{ ModuleData: md, lastPulse: lp }
		}
	}

	id := n
	lut["button"] = id
	modules[id] = &Button{ ModuleData: ModuleData{
		index : &modules,
		id    : id,
		name  : "button",
		dests : []int{broadcastid},
	} }

	return lut, modules
}

func parseModule(text string) []string {
	io := strings.Split(text, " -> ")
	typ := string(io[0][0])
	id := io[0][1:]
	dests := strings.Split(io[1], ", ")
	if typ == "b" {
		id = io[0]
	}
	out := []string{ typ, id }
	out = append(out, dests...)
	return out
}

func load(filename string) (lut ModuleDict, modules ModuleList) {
	var mods [][]string
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		mods = append(mods, parseModule(s.Text()))
	}
	check(s.Err())
	return createModules(mods)
}

func check(err error) { if err != nil { panic(err) } }

func main () {
	flag.Parse()
	lut, modules := load(flag.Arg(0))

	button    := lut["button"]
	broadcast := lut["broadcaster"]
	start := Msg{button, LOW, broadcast}

	run := 0
	var done bool
	for {
		done = modules.PushButton(start, run+1)
		run++

		if done || run >= 100000 {
			break
		}
	}

	fmt.Println("Run:", run)
	for _, m := range modules {
		m.Print()
	}

	fmt.Println("Answer to part 2 is", run)
}

func (modules ModuleList) Print(m Msg) {
	src := modules[m.src].GetName()
	dest := "output"
	if m.dest >= 0 {
		 dest = modules[m.dest].GetName()
	}
	if m.pulse == LOW {
		fmt.Println(src, "-low->", dest)
	} else {
		fmt.Println(src, "-high->", dest)
	}
}

func (modules ModuleList) PushButton(start Msg, run int) bool {
	var queue []Msg
	var head Msg

	queue = append(queue, start)

	for len(queue) > 0 {
		head, queue = queue[0], queue[1:]
		if head.dest == 33 && head.pulse == LOW {
			fmt.Println("module 4", run)
		}
		if head.dest < 0 {
			continue
		}
		newMsgs := modules[head.dest].Propagate(head)
		queue = append(queue, newMsgs...)
	}
	return false
}
