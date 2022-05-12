package hugetlb

import (
	"reflect"
	"unsafe"
)

type Entry struct {
	next *Entry
}

type Pool struct {
	b    []byte
	free *Entry
	unit int
	used int
}

func NewPool(b []byte, unit int) *Pool {
	return &Pool{b: b, unit: unit}
}

func (p *Pool) Get() ([]byte, bool) {
	if e := p.free; e != nil {
		p.free = e.next
		var b []byte
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
		hdr.Data = uintptr(unsafe.Pointer(e))
		hdr.Cap = p.unit
		hdr.Len = p.unit
		return b, true
	}
	if p.used+p.unit <= len(p.b) {
		b := p.b[p.used : p.used+p.unit]
		p.used += p.unit
		return b, true
	}
	return nil, false
}

func (p *Pool) Put(b []byte) {
	// TODO: double free
	e := (*Entry)(unsafe.Pointer(&b[0]))
	e.next = p.free
	p.free = e
}
