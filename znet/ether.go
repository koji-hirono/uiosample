package znet

import (
	"unsafe"
)

type EtherHdr struct {
	Dst  MacAddr
	Src  MacAddr
	Type Uint16
}

const (
	EtherTypeIPv4 uint16 = 0x0800
	EtherTypeARP         = 0x0806
)

func DecodeEtherHdr(b []byte) (*EtherHdr, int) {
	h := (*EtherHdr)(unsafe.Pointer(&b[0]))
	return h, 14
}

func (h EtherHdr) Len() int {
	return 14
}
