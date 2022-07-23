package znet

import (
	"unsafe"
)

type ARPHdr struct {
	HType Uint16
	PType Uint16
	HLen  Uint8
	PLen  Uint8
	Op    Uint16
	SMac  MacAddr
	SIP   IPv4Addr
	TMac  MacAddr
	TIP   IPv4Addr
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
