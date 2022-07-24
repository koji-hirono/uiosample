package znet

import (
	"reflect"
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
	IPProtoTCP        = 16
	IPProtoUDP        = 17
)

func DecodeIPv4Hdr(b []byte) (*IPv4Hdr, int) {
	h := (*IPv4Hdr)(unsafe.Pointer(&b[0]))
	return h, 20
}

func (h IPv4Hdr) Bytes() []byte {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = h.Len()
	hdr.Len = h.Len()
	hdr.Data = uintptr(unsafe.Pointer(&h))
	return b
}

func (h IPv4Hdr) Len() int {
	return 20
}

type IPv4PseudoHdr struct {
	Src    IPv4Addr
	Dst    IPv4Addr
	Pad    Uint8
	Proto  Uint8
	Length Uint16
}

func DecodeIPv4PseudoHdr(b []byte) (*IPv4PseudoHdr, int) {
	h := (*IPv4PseudoHdr)(unsafe.Pointer(&b[0]))
	return h, 12
}

func (h IPv4PseudoHdr) Bytes() []byte {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = h.Len()
	hdr.Len = h.Len()
	hdr.Data = uintptr(unsafe.Pointer(&h))
	return b
}

func (h IPv4PseudoHdr) Len() int {
	return 12
}
