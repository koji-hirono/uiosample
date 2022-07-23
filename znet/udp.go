package znet

import (
	"reflect"
	"unsafe"
)

type UDPHdr struct {
	SrcPort Uint16
	DstPort Uint16
	Length  Uint16
	Chksum  Uint16
}

func DecodeUDPHdr(b []byte) (*UDPHdr, int) {
	h := (*UDPHdr)(unsafe.Pointer(&b[0]))
	return h, 8
}

func (h UDPHdr) Bytes() []byte {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = h.Len()
	hdr.Len = h.Len()
	hdr.Data = uintptr(unsafe.Pointer(&h))
	return b
}

func (h UDPHdr) Len() int {
	return 8
}
