package pci

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	Command              = 0x04
	CommandMaster uint16 = 0x4

	CapList = 0x34
)

type Config struct {
	id int
	fd int
}

func NewConfig(id int) (*Config, error) {
	c := new(Config)
	c.id = id
	fname := fmt.Sprintf("/sys/class/uio/uio%v/device/config", id)
	fd, err := syscall.Open(fname, syscall.O_RDWR, 0)
	if err != nil {
		return c, err
	}
	c.fd = fd
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
