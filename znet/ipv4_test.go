package znet

import (
	"testing"
)

func TestIPv4Checksum(t *testing.T) {
	b := make([]byte, 64)
	want := uint16(0x3e6b)
	ip, n := DecodeIPv4Hdr(b)
	ip.VerHL = 4 << 4 | 5;
	ip.ToS = 0
	ip.Length.Set(84)
	ip.ID.Set(0)
	ip.FlgOff.Set(0)
	ip.TTL = 64
	ip.Proto = IPProtoICMP
	ip.Chksum.Set(0)
	ip.Src = [4]byte{30, 30, 0, 2}
	ip.Dst = [4]byte{30, 30, 0, 1}
	chksum := CalcChecksum(b[:n])
	if chksum != want {
		t.Errorf("want: %x; but got %x\n", want, chksum)
	}
}

func TestIPv4Checksum2(t *testing.T) {
	want := uint16(0x3e6b)
	ip := IPv4Hdr{
		VerHL: 4 << 4 | 5,
		ToS:     0,
		Length:  Uint16([2]byte{0, 84}),
		ID:      Uint16([2]byte{0, 0}),
		FlgOff:  Uint16([2]byte{0, 0}),
		TTL:     64,
		Proto:   IPProtoICMP,
		Chksum:  Uint16([2]byte{0, 0}),
		Src:     [4]byte{30, 30, 0, 2},
		Dst:     [4]byte{30, 30, 0, 1},
	}
	b := ip.Bytes()
	chksum := CalcChecksum(b)
	if chksum != want {
		t.Errorf("want: %x; but got %x\n", want, chksum)
	}
}
