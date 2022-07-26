package pci

import (
	"fmt"
)

type Device struct {
	Addr   Addr
	Config *Config
	Infos  []ResourceInfo
	Ress   []Resource
}

func OpenDevice(addr *Addr, c *Config) (*Device, error) {
	d := new(Device)
	d.Addr = *addr
	d.Config = c
	infos, err := ScanResourceInfo(addr)
	if err != nil {
		return d, err
	}
	d.Infos = infos
	d.Ress = make([]Resource, len(infos))
	return d, nil
}

func (d *Device) Close() {
	for _, r := range d.Ress {
		if r != nil {
			r.Close()
		}
	}
}

func (d *Device) GetResource(i int) (Resource, error) {
	if i >= len(d.Ress) {
		return nil, fmt.Errorf("illegal resource id:%v", i)
	}
	r := d.Ress[i]
	if r != nil {
		return r, nil
	}
	info := &d.Infos[i]
	switch info.Type() {
	case ResourceTypeIO:
		r, err := OpenIOResource(&d.Addr, i, info)
		if err != nil {
			return nil, err
		}
		d.Ress[i] = r
		return r, nil
	case ResourceTypeMem:
		r, err := OpenMemResource(&d.Addr, i, info)
		if err != nil {
			return nil, err
		}
		d.Ress[i] = r
		return r, nil
	default:
		return nil, fmt.Errorf("not supported type: %v", info.Type())
	}
}
