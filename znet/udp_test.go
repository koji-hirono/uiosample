package znet

import (
	"testing"
)

func TestUDPChecksum(t *testing.T) {
	want := uint16(0xe481)

	data := []byte{
		0x20, 0x05, 0x00, 0x1a, 0x00, 0x00, 0x01, 0x00,
		0x00, 0x3c, 0x00, 0x05, 0x00, 0x0a, 0x00, 0x00,
		0x0c, 0x00, 0x60, 0x00, 0x04, 0xe5, 0xca, 0xf0,
		0x0b, 0x00, 0x59, 0x00, 0x01, 0x00,
	}

	udp := &UDPHdr{}
	udp.SrcPort.Set(8805)
	udp.DstPort.Set(8805)
	udplen := len(data) + udp.Len()
	udp.Length.Set(uint16(udplen))
	udp.Chksum.Set(0)

	pip := &IPv4PseudoHdr{}
	pip.Src.Set([]byte{10, 0, 0, 12})
	pip.Dst.Set([]byte{10, 0, 0, 10})
	pip.Pad.Set(0)
	pip.Proto.Set(IPProtoUDP)
	pip.Length.Set(uint16(udplen))

	calc := NewCalc()
	calc.Append(pip.Bytes())
	calc.Append(udp.Bytes())
	calc.Append(data)
	chksum := calc.Checksum()
	if chksum != want {
		t.Errorf("want: %x; but got %x\n", want, chksum)
	}
}
