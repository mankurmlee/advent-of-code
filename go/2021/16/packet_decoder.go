package main

import (
	"flag"
	"os"
	"bufio"
	"fmt"
	"encoding/hex"
	"strconv"
	"encoding/binary"
	"strings"
)

type ByteStream struct {
	data  []byte
	pos   int
}

type Packet struct {
	version  int
	typeid   int
	value    int
	children []*Packet
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func load(filename string) string {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan()
	check(s.Err())
	return s.Text()
}

func decode(hexStr string) []byte {
	bytes, err := hex.DecodeString(hexStr)
	check(err)
	return bytes
}

func (b *ByteStream) Read(bits int) (out int) {
	endpos := b.pos + bits - 1
	start  := b.pos >> 3
	end    := endpos >> 3 + 1
	pad    := 4 + start - end
	bytes := make([]byte, 4)
	copy(bytes[pad:], b.data[start:end])
	out = int(binary.BigEndian.Uint32(bytes))
	out >>= 7 - endpos & 7
	out &= (1 << bits) - 1
	b.pos += bits
	return out
}

func (b *ByteStream) ParsePacket() (out Packet) {
	out.version = b.Read(3)
	out.typeid  = b.Read(3)

	if out.typeid == 4 {
		out.value = b.ParseLiteral()
		return out
	}

	mode := b.Read(1)
	if mode == 0 {
		l := b.Read(15)
		e := b.pos + l
		for b.pos < e {
			sub := b.ParsePacket()
			out.children = append(out.children, &sub)
		}
	} else {
		n := b.Read(11)
		for i := 0; i < n; i++ {
			sub := b.ParsePacket()
			out.children = append(out.children, &sub)
		}
	}

	return out
}

func (b *ByteStream) ParseLiteral() (out int) {
	for {
		v := b.Read(5)
		out |= v & 15
		if v < 16 {
			break
		}
		out <<= 4
	}
	return out
}

func (p Packet) String() string {
	list := make([]string, 3)
	list[0] = strconv.Itoa(p.version)
	list[1] = strconv.Itoa(p.typeid)
	if p.typeid == 4 {
		list[2] = strconv.Itoa(p.value)
	} else {
		var children []string
		for _, c := range p.children {
			children = append(children, c.String())
		}
		list[2] = strings.Join(children, " ")
	}
	return "{"+ strings.Join(list, " ") + "}"
}

func (p Packet) VersionSum() (sum int) {
	sum = p.version
	for _, c := range p.children {
		sum += c.VersionSum()
	}
	return sum
}

func (p Packet) Eval() (out int) {
	switch p.typeid {
	case 0:
		for _, c := range p.children {
			out += c.Eval()
		}
	case 1:
		out = p.children[0].Eval()
		for _, c := range p.children[1:] {
			out *= c.Eval()
		}
	case 2:
		out = p.children[0].Eval()
		for _, c := range p.children[1:] {
			out = min(out, c.Eval())
		}
	case 3:
		out = p.children[0].Eval()
		for _, c := range p.children[1:] {
			out = max(out, c.Eval())
		}
	case 4:
		out = p.value
	case 5:
		a := p.children[0].Eval()
		b := p.children[1].Eval()
		if a > b { out = 1 }
	case 6:
		a := p.children[0].Eval()
		b := p.children[1].Eval()
		if a < b { out = 1 }
	case 7:
		a := p.children[0].Eval()
		b := p.children[1].Eval()
		if a == b { out = 1 }
	}
	return out
}

func main() {
	flag.Parse()
	p := load(flag.Arg(0))
	fmt.Println(p)

	b := ByteStream{ data: decode(p) }
	pack := b.ParsePacket()

	fmt.Println(pack)
	fmt.Println("Part 1:", pack.VersionSum())
	fmt.Println("Part 2:", pack.Eval())
}
