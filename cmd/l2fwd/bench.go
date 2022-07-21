package main

import (
	"fmt"
)

type Bench struct {
	ident string
	n     uint64
	total uint64
	min   uint64
	max   uint64
	start uint64
}

func NewBench(ident string) *Bench {
	b := new(Bench)
	b.ident = ident
	b.Reset()
	return b
}

func (b *Bench) Reset() {
	b.total = 0
	b.min = ^uint64(0)
	b.max = 0
	b.start = 0
}

func (b *Bench) Start() {
	b.start = Rdtsc()
}

func (b *Bench) End() {
	end := Rdtsc()
	d := end - b.start
	b.total += d
	b.n++
	if d < b.min {
		b.min = d
	}
	if d > b.max {
		b.max = d
	}
}

func (b *Bench) Ave() uint64 {
	return b.total / b.n
}

func (b *Bench) Print() {
	fmt.Printf("=== %v\n", b.ident)
	fmt.Printf("N    : %v\n", b.n)
	fmt.Printf("Total: %v\n", b.total)
	fmt.Printf("Min  : %v\n", b.min)
	fmt.Printf("Max  : %v\n", b.max)
	fmt.Printf("Ave  : %v\n", b.Ave())
}
