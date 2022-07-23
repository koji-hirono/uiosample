package ethernet

import (
	"testing"
)

func TestICMPChecksum(t *testing.T) {
	want := uint16(0x1ae2)
	icmp := Packet{
		ICMPHdr{Type: 0x8, Code: 0x0, Chksum: 0},
		ICMPEchoHdr{ID: 0x257b, Seq: 0x5},
		Data([]byte{
			0x65, 0x95, 0x59, 0x61, 0x00, 0x00, 0x00, 0x00,
			0x31, 0xd4, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
			0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
			0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
			0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
		}),
	}
	sum := icmp.Sum()
	chksum := uint16(^sum)
	if chksum != want {
		t.Errorf("want: %x; but got %x\n", want, chksum)
	}
}