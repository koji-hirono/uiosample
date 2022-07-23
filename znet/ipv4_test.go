package znet

import (
	"testing"
)

func TestIPv4Checksum(t *testing.T) {
	b := make([]byte, 64)
	want := uint16(0x3e6b)
	ip, n := DecodeIPv4Hdr(b)
	ip.VerHL.Set(4<<4 | 5)
	ip.ToS.Set(0)
	ip.Length.Set(84)
	ip.ID.Set(0)
	ip.FlgOff.Set(0)
	ip.TTL.Set(64)
	ip.Proto.Set(IPProtoICMP)
	ip.Chksum.Set(0)
	ip.Src.Set([]byte{30, 30, 0, 2})
	ip.Dst.Set([]byte{30, 30, 0, 1})
	chksum := CalcChecksum(b[:n])
	if chksum != want {
		t.Errorf("want: %x; but got %x\n", want, chksum)
	}
}
