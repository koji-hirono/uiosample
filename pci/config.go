package pci

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	Command              = 0x04
	CommandMaster uint16 = 0x4

	CapList = 0x34
)

func LoadVendorID(addr *Addr) (uint16, error) {
	s, err := loadString(addr, "vendor")
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(id), nil
}

func LoadDeviceID(addr *Addr) (uint16, error) {
	s, err := loadString(addr, "device")
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(id), nil
}

func LoadSubsystemVendorID(addr *Addr) (uint16, error) {
	s, err := loadString(addr, "subsystem_vendor")
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(id), nil
}

func LoadSubsystemDeviceID(addr *Addr) (uint16, error) {
	s, err := loadString(addr, "subsystem_device")
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(id), nil
}

func LoadRevisionID(addr *Addr) (uint8, error) {
	s, err := loadString(addr, "revision")
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		return 0, err
	}
	return uint8(id), nil
}

func loadString(addr *Addr, name string) (string, error) {
	fname := fmt.Sprintf("/sys/bus/pci/devices/%s/%s", addr, name)
	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	if !s.Scan() {
		return "", io.ErrUnexpectedEOF
	}
	return s.Text(), nil
}

type Config struct {
	VendorID          uint16
	DeviceID          uint16
	SubsystemVendorID uint16
	SubsystemDeviceID uint16
	RevisionID        uint8
	addr              *Addr
	fd                int
}

func OpenConfig(addr *Addr) (*Config, error) {
	c := new(Config)
	c.addr = addr
	fname := fmt.Sprintf("/sys/bus/pci/devices/%s/config", addr)
	fd, err := syscall.Open(fname, syscall.O_RDWR, 0)
	if err != nil {
		return c, err
	}
	c.fd = fd
	vendor, err := LoadVendorID(addr)
	if err != nil {
		return c, err
	}
	c.VendorID = vendor
	device, err := LoadDeviceID(addr)
	if err != nil {
		return c, err
	}
	c.DeviceID = device
	rev, err := LoadRevisionID(addr)
	if err != nil {
		return c, err
	}
	c.RevisionID = rev
	subvendor, err := LoadSubsystemVendorID(addr)
	if err != nil {
		return c, err
	}
	c.SubsystemVendorID = subvendor
	subdevice, err := LoadSubsystemDeviceID(addr)
	if err != nil {
		return c, err
	}
	c.SubsystemDeviceID = subdevice
	return c, nil
}

func (c *Config) Close() error {
	return syscall.Close(c.fd)
}

func (c *Config) Read8(off int) (uint8, error) {
	b := [1]byte{}
	_, err := syscall.Pread(c.fd, b[:], int64(off))
	if err != nil {
		return 0, err
	}
	return uint8(b[0]), nil
}

func (c *Config) Write8(off int, val, mask uint8) error {
	d, err := c.Read8(off)
	d = (d & ^mask) | (d & mask)
	b := [1]byte{byte(d)}
	_, err = syscall.Pwrite(c.fd, b[:], int64(off))
	return err
}

func (c *Config) Read16(off int) (uint16, error) {
	b := make([]byte, 2)
	_, err := syscall.Pread(c.fd, b, int64(off))
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func (c *Config) Write16(off int, val, mask uint16) error {
	d, err := c.Read16(off)
	if err != nil {
		return err
	}
	d = (d & ^mask) | (d & mask)
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, d)
	_, err = syscall.Pwrite(c.fd, b, int64(off))
	return err
}

func (c *Config) Dump() (string, error) {
	buf := [64]byte{}
	n, err := syscall.Pread(c.fd, buf[:], 0)
	if err != nil {
		return "", err
	}
	if n != 64 {
		return "", err
	}

	return hex.Dump(buf[:]), nil
}

func (c *Config) SetBusMaster() error {
	buf := [2]byte{}
	n, err := syscall.Pread(c.fd, buf[:], int64(Command))
	if err != nil {
		return err
	}
	if n != 2 {
		return fmt.Errorf("n != 2")
	}

	reg := (*uint16)(unsafe.Pointer(&buf[0]))

	if *reg&CommandMaster == 0 {
		*reg |= CommandMaster
		_, err := syscall.Pwrite(c.fd, buf[:], int64(Command))
		if err != nil {
			return err
		}
	}

	return nil
}
