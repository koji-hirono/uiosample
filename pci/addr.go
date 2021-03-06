package pci

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
)

type Addr struct {
	Domain uint32
	Bus    uint8
	ID     uint8
	Func   uint8
}

func (a *Addr) String() string {
	return fmt.Sprintf("%04x:%02x:%02x.%01x", a.Domain, a.Bus, a.ID, a.Func)
}

func (a *Addr) ScanResources() ([]Resource, error) {
	rs := []Resource{}
	fname := "/sys/bus/pci/devices/" + a.String() + "/resource"
	f, err := os.Open(fname)
	if err != nil {
		return rs, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	i := 0
	for s.Scan() {
		r := &Resource{Addr: a, Index: i}
		n, err := fmt.Sscanf(s.Text(), "0x%x 0x%x 0x%x",
			&r.Phys, &r.End, &r.Flags)
		if err != nil {
			return rs, err
		}
		if n != 3 {
			return rs, fmt.Errorf("n != 3")
		}
		if r.Phys != 0 {
			rs = append(rs, *r)
		}
		i++
	}

	return rs, nil
}

type Resource struct {
	Addr  *Addr
	Index int
	Phys  uint64
	End   uint64
	Flags uint64
}

type ResourceType uint64

const (
	ResourceTypeIO  ResourceType = 0x0100
	ResourceTypeMem              = 0x0200
	ResourceTypeReg              = 0x0300
	ResourceTypeIRQ              = 0x0400
	ResourceTypeDMA              = 0x0800
	ResourceTypeBus              = 0x1000
)

func (t ResourceType) String() string {
	switch t {
	case ResourceTypeIO:
		return "IO"
	case ResourceTypeMem:
		return "Mem"
	case ResourceTypeReg:
		return "Reg"
	case ResourceTypeIRQ:
		return "IRQ"
	case ResourceTypeDMA:
		return "DMA"
	case ResourceTypeBus:
		return "Bus"
	default:
		return ""
	}
}

func (r *Resource) Type() ResourceType {
	return ResourceType(r.Flags & 0x1f00)
}

func (r *Resource) Map() ([]byte, error) {
	fname := fmt.Sprintf("/sys/bus/pci/devices/%s/resource%v", r.Addr, r.Index)
	fd, err := syscall.Open(fname, syscall.O_RDWR, 0)
	if err != nil {
		return []byte{}, err
	}
	defer syscall.Close(fd)

	size := int(r.End - r.Phys + 1)
	prot := syscall.PROT_READ | syscall.PROT_WRITE
	return syscall.Mmap(fd, 0, size, prot, syscall.MAP_SHARED)
}
