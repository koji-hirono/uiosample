package znet

import (
	"unsafe"
)

type ICMPHdr struct {
	Type   Uint8
	Code   Uint8
	Chksum Uint16
}

const (
	ICMPTypeEchoReply   uint8 = 0
	ICMPTypeEchoRequest       = 8
)

func DecodeICMPHdr(b []byte) (*ICMPHdr, int) {
	h := (*ICMPHdr)(unsafe.Pointer(&b[0]))
	return h, 4
}

func (h ICMPHdr) Len() int {
	return 4
}

type ICMPEchoHdr struct {
	ID  Uint16
	Seq Uint16
}

func DecodeICMPEchoHdr(b []byte) (*ICMPEchoHdr, int) {
	h := (*ICMPEchoHdr)(unsafe.Pointer(&b[0]))
	return h, 4
}

func (h ICMPEchoHdr) Len() int {
	return 4
}
