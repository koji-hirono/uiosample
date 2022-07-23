package znet

import (
	"reflect"
	"unsafe"
)

type IPv4Hdr struct {
	VerHL   uint8
	ToS     uint8
	Length  Uint16
	ID      Uint16
	FlgOff  Uint16
	TTL     uint8
	Proto   uint8
	Chksum  Uint16
	Src     [4]byte
	Dst     [4]byte
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

func (h IPv4Hdr) Bytes() []byte {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = h.Len()
	hdr.Len = h.Len()
	hdr.Data = uintptr(unsafe.Pointer(&h))
	return b
}

func (h IPv4Hdr) Version() uint8 {
	return (h.VerHL >> 4) & 0xf
}

func (h IPv4Hdr) Hdrlen() uint8 {
	return h.VerHL & 0xf
}

func (h IPv4Hdr) Flags() uint8 {
	return h.FlgOff[0] >> 5
}

func (h IPv4Hdr) Offset() uint16 {
	x := uint16(h.FlgOff[0] & 0x1f) << 8
	x |= uint16(h.FlgOff[1])
	return x
}
