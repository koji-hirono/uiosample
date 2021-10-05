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
	log.Printf("off: %x\n", off)
	buf := [8]byte{}
	n, err := f.ReadAt(buf[:], int64(off))
	if err != nil {
		return 0, err
	}
	if n != 8 {
		return 0, fmt.Errorf("to few read")
	}

	log.Printf("buf: %x\n", buf)

	addr := (*uintptr)(unsafe.Pointer(&buf[0]))
	log.Printf("addr: %x\n", *addr)
	return (*addr & 0x007fffffffffffff) * pagesize + p % pagesize, nil
}

/*
func PageAlloc(n int) ([]byte, error) {
	size := n * 2 * 1024 * 1024
	*/
func Alloc(size int) ([]byte, error) {

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

func Free(b []byte) {
	syscall.Munmap(b)
}

var PoolSizeList = [...]int{64, 128, 256, 512, 1024, 2048, 4096, 8192}
var PoolTable map[int]*Pool

func Init() {
	/*
	PoolTable = make(map[int]*Pool)
	for _, unit := range PoolSizeList {
		b, err := PageAlloc(1)
		if err != nil {
			continue
		}
		PoolTable[unit] = NewPool(b, unit)
	}
	*/
}

func hogeAlloc(size int) ([]byte, error) {
	/*
	for _, unit := range PoolSizeList {
		if size <= unit {
			p, ok := PoolTable[unit]
			if !ok {
				return nil, ErrOutOfMemory
			}
			b, ok := p.Get()
			if !ok {
				return nil, ErrOutOfMemory
			}
			return b, nil
		}
	}
	*/
	/*
	n := (size / (2 * 1024 * 1024)) + 1
	return PageAlloc(n)
	*/
	return nil, nil
}

func hogeFree(b []byte) {
	/*
	size := cap(b)
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
	PageFree(b)
	*/
}
