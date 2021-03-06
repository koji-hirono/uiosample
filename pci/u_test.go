package pci

import (
	"encoding/hex"
	"log"
	"syscall"
	"testing"
)

func TestDevUIO(t *testing.T) {
	fname := "/dev/uio0"
	fd, err := syscall.Open(fname, syscall.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fd)
}

func TestConfig(t *testing.T) {
	fd, err := ConfigOpen()
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fd)

	s, err := ConfigDump(fd)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("buf:\n%s", s)
}

func TestCap(t *testing.T) {
	cfd, err := ConfigOpen()
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(cfd)

	err = SetBusMaster(cfd)
	if err != nil {
		t.Fatal(err)
	}

	buf := [2]byte{}
	n, err := syscall.Pread(cfd, buf[:1], int64(CapList))
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("n != 1")
	}
	pos := buf[0]
	log.Printf("pos: %x\n", pos)

	n, err = syscall.Pread(cfd, buf[:], int64(pos))
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatal("n != 2")
	}

	vndr := buf[0]
	next := buf[1]
	log.Printf("vndr: %x\n", vndr)
	log.Printf("next: %x\n", next)
}

func TestMapResource(t *testing.T) {
	addr := &Addr{
		Domain: 0,
		Bus:    0,
		ID:     9,
		Func:   0,
	}
	rs, err := addr.ScanResources()
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range rs {
		typ := r.Type()
		log.Printf("Type: %s\n", typ)
		switch typ {
		case ResourceTypeIO:
		case ResourceTypeMem:
			buf, err := r.Map()
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("buf:\n%s", hex.Dump(buf))
			syscall.Munmap(buf)
		case ResourceTypeReg:
		case ResourceTypeIRQ:
		case ResourceTypeDMA:
		case ResourceTypeBus:
		}
	}
}
