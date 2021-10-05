package pci

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
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

type ResourceInfo struct {
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

func (r *ResourceInfo) Type() ResourceType {
	return ResourceType(r.Flags & 0x1f00)
}

type Resource interface {
	Read32(int) uint32
	Write32(int, uint32, uint32)
	Close() error
}

type MemResource struct {
	b []byte
}

func NewMemResource(addr *Addr, index int, info *ResourceInfo) (*MemResource, error) {
	r := new(MemResource)
	fname := fmt.Sprintf("/sys/bus/pci/devices/%s/resource%v", addr, index)
	fd, err := syscall.Open(fname, syscall.O_RDWR, 0)
	if err != nil {
		return r, err
	}
	defer syscall.Close(fd)

	size := int(info.End - info.Phys + 1)
	prot := syscall.PROT_READ | syscall.PROT_WRITE
	b, err := syscall.Mmap(fd, 0, size, prot, syscall.MAP_SHARED)
	if err != nil {
		return r, err
	}
	r.b = b

	return r, nil
}

func (r *MemResource) Close() error {
	return syscall.Munmap(r.b)
}

func (r *MemResource) Read32(off int) uint32 {
	return *(*uint32)(unsafe.Pointer(&r.b[off]))
}

func (r *MemResource) Write32(off int, val uint32, mask uint32) {
	d := (*uint32)(unsafe.Pointer(&r.b[off]))
	*d = (*d & ^mask) | (val & mask)
}


type IOResource struct {
	fd int
}

func NewIOResource(addr *Addr, index int, info *ResourceInfo) (*IOResource, error) {
	r := new(IOResource)
	fname := fmt.Sprintf("/sys/bus/pci/devices/%s/resource%v", addr, index)
	fd, err := syscall.Open(fname, syscall.O_RDWR, 0)
	if err != nil {
		return r, err
	}
	r.fd = fd
	return r, nil
}

func (r *IOResource) Close() error {
	return syscall.Close(r.fd)
}

func (r *IOResource) Read32(off int) uint32 {
	b := make([]byte, 4)
	_, err := syscall.Pread(r.fd, b, int64(off))
	if err != nil {
		return 0
	}
	return binary.LittleEndian.Uint32(b)
}

func (r *IOResource) Write32(off int, val uint32, mask uint32) {
	d := r.Read32(off)
	d = (d & ^mask) | (val & mask)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, d)
	syscall.Pwrite(r.fd, b, int64(off))
}
