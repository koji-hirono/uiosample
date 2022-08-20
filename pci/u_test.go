package pci

import (
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
	addr := &Addr{
		Domain: 0,
		Bus:    0,
		ID:     10,
		Func:   0,
	}
	c, err := OpenConfig(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	s, err := c.Dump()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("buf:\n%s", s)

	t.Logf("Vendor ID: %x\n", c.VendorID)
	t.Logf("Device ID: %x\n", c.DeviceID)
	t.Logf("Revision ID: %x\n", c.RevisionID)
	t.Logf("Subsystem Vendor ID: %x\n", c.SubsystemVendorID)
	t.Logf("Subsystem Device ID: %x\n", c.SubsystemDeviceID)
}

func TestCap(t *testing.T) {
	addr := &Addr{
		Domain: 0,
		Bus:    0,
		ID:     10,
		Func:   0,
	}
	c, err := OpenConfig(addr)
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
	t.Logf("pos: %x\n", pos)

	vndr, err := c.Read8(int(pos))
	if err != nil {
		t.Fatal(err)
	}

	next, err := c.Read8(int(pos + 1))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("vndr: %x\n", vndr)
	t.Logf("next: %x\n", next)
}

func TestMapResource(t *testing.T) {
	addr := &Addr{
		Domain: 0,
		Bus:    0,
		ID:     10,
		Func:   0,
	}
	device, err := OpenDevice(addr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer device.Close()
	for i, res := range device.Ress {
		t.Logf("Res[%v]:\n", i)
		if res != nil {
			t.Logf("0000: 0x%08x\n", res.Read32(0))
		}
	}
}
