package znet

import (
	"unsafe"
)

type ARPHdr struct {
	HType Uint16
	PType Uint16
	HLen  uint8
	PLen  uint8
	Op    Uint16
	SMac  [6]byte
	SIP   [4]byte
	TMac  [6]byte
	TIP   [4]byte
}

// ARPHdr.Op
const (
	ARPRequest uint16 = 1
	ARPReply          = 2
)

func DecodeARPHdr(b []byte) (*ARPHdr, int) {
	h := (*ARPHdr)(unsafe.Pointer(&b[0]))
	return h, 28
}

func (h ARPHdr) Len() int {
	return 28
}
