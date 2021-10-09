package pci

import (
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
	c, err := NewConfig(0)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	s, err := c.Dump()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("buf:\n%s", s)
}

func TestCap(t *testing.T) {
	c, err := NewConfig(0)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	err = c.SetBusMaster()
	if err != nil {
		t.Fatal(err)
	}

	pos, err := c.Read8(CapList)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("pos: %x\n", pos)

	vndr, err := c.Read8(int(pos))
	if err != nil {
		t.Fatal(err)
	}

	next, err := c.Read8(int(pos + 1))
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("vndr: %x\n", vndr)
	log.Printf("next: %x\n", next)
}

func TestMapResource(t *testing.T) {
	addr := &Addr{
		Domain: 0,
		Bus:    0,
		ID:     10,
		Func:   0,
	}
	device, err := NewDevice(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i, res := range device.Ress {
		log.Printf("Res[%v]:\n", i)
		if res != nil {
			log.Printf("0000: 0x%08x\n", res.Read32(0))
		}
	}
}
