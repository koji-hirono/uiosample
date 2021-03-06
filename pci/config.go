package pci

import (
	"encoding/hex"
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	Command = 0x04
	CommandMaster uint16 = 0x4

	CapList = 0x34
)

func ConfigOpen() (int, error) {
	fname := "/sys/class/uio/uio0/device/config"
	return unix.Open(fname, unix.O_RDWR, 0)
}

func ConfigDump(fd int) (string, error) {
	buf := [64]byte{}
	n, err := unix.Pread(fd, buf[:], 0)
	if err != nil {
		return "", err
	}
	if n != 64 {
		return "", err
	}

	return hex.Dump(buf[:]), nil
}

func SetBusMaster(fd int) error {
	buf := [2]byte{}
	n, err := unix.Pread(fd, buf[:], int64(Command))
	if err != nil {
		return err
	}
	if n != 2 {
		return fmt.Errorf("n != 2")
	}

	reg := (*uint16)(unsafe.Pointer(&buf[0]))

	if *reg & CommandMaster == 0 {
		*reg |= CommandMaster
		_, err := unix.Pwrite(fd, buf[:], int64(Command))
		if err != nil {
			return err
		}
	}

	return nil
}
