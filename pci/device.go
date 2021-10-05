package pci

import (
	"bufio"
	"fmt"
	"os"
)

type Device struct {
	Addr   Addr
	Config *Config
	Infos  []ResourceInfo
	Ress   []Resource
}

func NewDevice(addr *Addr, c *Config) (*Device, error) {
	d := new(Device)
	d.Addr = *addr
	d.Config = c
	infos, err := scanResourceInfo(addr)
	if err != nil {
		return d, err
	}
	d.Infos = infos
	d.Ress = make([]Resource, len(infos))
	for i, info := range d.Infos {
		switch info.Type() {
		case ResourceTypeIO:
			r, err := NewIOResource(addr, i, &info)
			if err != nil {
				return d, err
			}
			d.Ress[i] = r
		case ResourceTypeMem:
			r, err := NewMemResource(addr, i, &info)
			if err != nil {
				return d, err
			}
			d.Ress[i] = r
		case ResourceTypeReg:
		case ResourceTypeIRQ:
		case ResourceTypeDMA:
		case ResourceTypeBus:
		default:
		}
	}
	return d, nil
}

func scanResourceInfo(addr *Addr) ([]ResourceInfo, error) {
	rs := []ResourceInfo{}
	fname := "/sys/bus/pci/devices/" + addr.String() + "/resource"
	f, err := os.Open(fname)
	if err != nil {
		return rs, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for i := 0; s.Scan(); i++ {
		r := &ResourceInfo{}
		n, err := fmt.Sscanf(s.Text(), "0x%x 0x%x 0x%x",
			&r.Phys, &r.End, &r.Flags)
		if err != nil {
			return rs, err
		}
		if n != 3 {
			return rs, fmt.Errorf("n != 3")
		}
		rs = append(rs, *r)
	}

	return rs, nil
}
