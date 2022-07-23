package ethernet

import (
	"encoding/binary"
)

type ICMPHdr struct {
	Type   uint8
	Code   uint8
	Chksum uint16
}

const (
	ICMPTypeEchoReply   uint8 = 0
	ICMPTypeEchoRequest uint8 = 8
)

func DecodeICMPHdr(b []byte) (ICMPHdr, error) {
	var h ICMPHdr
	h.Type = b[0]
	h.Code = b[1]
	h.Chksum = binary.BigEndian.Uint16(b[2:4])
	return h, nil
}

func (h ICMPHdr) Len() int {
	return 4
}

func (h ICMPHdr) Sum() uint32 {
	n := h.Len()
	b := make([]byte, n)
	h.Encode(b)
	return Data(b).Sum()
}

func (h ICMPHdr) Encode(b []byte) error {
	b[0] = h.Type
	b[1] = h.Code
	binary.BigEndian.PutUint16(b[2:], h.Chksum)
	return nil
}

type ICMPEchoHdr struct {
	ID  uint16
	Seq uint16
}

func DecodeICMPEchoHdr(b []byte) (ICMPEchoHdr, error) {
	var h ICMPEchoHdr
	h.ID = binary.BigEndian.Uint16(b[0:2])
	h.Seq = binary.BigEndian.Uint16(b[2:4])
	return h, nil
}

func (h ICMPEchoHdr) Len() int {
	return 4
}

func (h ICMPEchoHdr) Sum() uint32 {
	n := h.Len()
	b := make([]byte, n)
	h.Encode(b)
	return Data(b).Sum()
}

func (h ICMPEchoHdr) Encode(b []byte) error {
	binary.BigEndian.PutUint16(b[0:2], h.ID)
	binary.BigEndian.PutUint16(b[2:4], h.Seq)
	return nil
}
