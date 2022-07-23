package hugetlb

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

var (
	ErrOutOfMemory = errors.New("out of memory")
)

func SetPages(n int) error {
	fname := "/proc/sys/vm/nr_hugepages"
	f, err := os.OpenFile(fname, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%v\n", n)
	if err != nil {
		return err
	}
	return nil
}

func VirtToPhys(p uintptr) (uintptr, error) {
	fname := "/proc/self/pagemap"
	f, err := os.Open(fname)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	pagesize := uintptr(os.Getpagesize())
	off := p / pagesize * 8
	buf := [8]byte{}
	n, err := f.ReadAt(buf[:], int64(off))
	if err != nil {
		return 0, err
	}
	if n != 8 {
		return 0, fmt.Errorf("to few read")
	}

	addr := (*uintptr)(unsafe.Pointer(&buf[0]))
	return (*addr&0x007fffffffffffff)*pagesize + p%pagesize, nil
}

func PageAlloc(n int) ([]byte, error) {
	size := n * 2 * 1024 * 1024
	prot := syscall.PROT_READ | syscall.PROT_WRITE
	flags := syscall.MAP_PRIVATE | syscall.MAP_ANONYMOUS | syscall.MAP_HUGETLB
	buf, err := syscall.Mmap(-1, 0, size, prot, flags)
	if err != nil {
		return buf, err
	}
	err = syscall.Mlock(buf)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func PageFree(b []byte) {
	syscall.Munmap(b)
}

var PoolSizeList = [...]int{512, 1024, 2048, 4096, 8192}
var PoolTable map[int]*Pool

func Init() {
	PoolTable = make(map[int]*Pool)
	for _, unit := range PoolSizeList {
		b, err := PageAlloc(1)
		if err != nil {
			continue
		}
		virt := uintptr(unsafe.Pointer(&b[0]))
		phys, err := VirtToPhys(virt)
		if err != nil {
			PageFree(b)
			continue
		}
		PoolTable[unit] = NewPool(b, phys, unit)
	}
}

func Stat() {
	for _, unit := range PoolSizeList {
		p, ok := PoolTable[unit]
		if !ok {
			continue
		}

		log.Printf("== Pool[%v]\n", unit)
		log.Printf("len : %v\n", len(p.b))
		log.Printf("used: %v\n", p.used)
		log.Printf("free: %v\n", p.free)
	}
}

func Alloc(size int) ([]byte, uintptr, error) {
	for _, unit := range PoolSizeList {
		if size <= unit {
			p, ok := PoolTable[unit]
			if !ok {
				return nil, 0, ErrOutOfMemory
			}
			b, ok := p.Get()
			if !ok {
				return nil, 0, ErrOutOfMemory
			}
			phys, err := p.PhysAddr(b)
			if err != nil {
				p.Put(b)
				return nil, 0, err
			}
			return b, phys, nil
		}
	}
	log.Printf("Alloc: unknown pool size: %v\n", size)
	n := (size / (2 * 1024 * 1024)) + 1
	b, err := PageAlloc(n)
	if err != nil {
		return nil, 0, err
	}
	virt := uintptr(unsafe.Pointer(&b[0]))
	phys, err := VirtToPhys(virt)
	if err != nil {
		PageFree(b)
		return nil, 0, err
	}
	return b, phys, nil
}

func Free(b []byte) {
	size := cap(b)
	/*
		for _, unit := range PoolSizeList {
			if size <= unit {
				p, ok := PoolTable[unit]
				if !ok {
					return
				}
				p.Put(b)
				return
			}
		}
	*/
	p, ok := PoolTable[size]
	if !ok {
		log.Printf("Free: unknown pool size: %v\n", size)
		PageFree(b)
		return
	}
	p.Put(b)
}

func PhysAddr(b []byte) (uintptr, error) {
	size := cap(b)
	for _, unit := range PoolSizeList {
		if size <= unit {
			p, ok := PoolTable[unit]
			if !ok {
				continue
			}
			return p.PhysAddr(b)
		}
	}
	virt := uintptr(unsafe.Pointer(&b[0]))
	return VirtToPhys(virt)
}
