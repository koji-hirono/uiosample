package znet

import (
	"unsafe"
)

type IPv4Hdr struct {
	VerHL  Uint8
	ToS    Uint8
	Length Uint16
	ID     Uint16
	FlgOff Uint16
	TTL    Uint8
	Proto  Uint8
	Chksum Uint16
	Src    IPv4Addr
	Dst    IPv4Addr
}

const (
	IPProtoICMP uint8 = 1
)

func DecodeIPv4Hdr(b []byte) (*IPv4Hdr, int) {
	h := (*IPv4Hdr)(unsafe.Pointer(&b[0]))
	return h, 20
}

func (h IPv4Hdr) Len() int {
	return 20
}
