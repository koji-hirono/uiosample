package main

import (
	"testing"
)

func TestIPv4Checksum(t *testing.T) {
	want := uint16(0x3e6b)
	ip := IPv4Hdr{
		Version: 4,
		Hdrlen:  5,
		ToS:     0,
		Length:  84,
		ID:      0,
		Flags:   0,
		Offset:  0,
		TTL:     64,
		Proto:   IPProtoICMP,
		Chksum:  0,
		Src:     []byte{30, 30, 0, 2},
		Dst:     []byte{30, 30, 0, 1},
	}
	sum := ip.Sum()
	chksum := uint16(^sum)
	if chksum != want {
		t.Errorf("want: %x; but got %x\n", want, chksum)
	}
}
