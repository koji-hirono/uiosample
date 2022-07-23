package hugetlb

import (
	"reflect"
	"sync"
	"unsafe"
)

type Entry struct {
	next *Entry
}

type Pool struct {
	b    []byte
	phys uintptr
	free *Entry
	unit int
	used int
	m    sync.Mutex
}

func NewPool(b []byte, phys uintptr, unit int) *Pool {
	return &Pool{b: b, phys: phys, unit: unit}
}

func (p *Pool) Get() ([]byte, bool) {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = p.unit
	hdr.Len = p.unit
	p.m.Lock()
	defer p.m.Unlock()
	if e := p.free; e != nil {
		p.free = e.next
		hdr.Data = uintptr(unsafe.Pointer(e))
		return b, true
	}
	if p.used+p.unit <= len(p.b) {
		hdr.Data = uintptr(unsafe.Pointer(&p.b[p.used]))
		p.used += p.unit
		return b, true
	}
	return nil, false
}

func (p *Pool) Put(b []byte) {
	// TODO: double free
	e := (*Entry)(unsafe.Pointer(&b[0]))
	p.m.Lock()
	defer p.m.Unlock()
	e.next = p.free
	p.free = e
}

func (p *Pool) PhysAddr(b []byte) (uintptr, error) {
	start := uintptr(unsafe.Pointer(&p.b[0]))
	end := uintptr(unsafe.Pointer(&b[0]))
	return p.phys + (end - start), nil
}
