package znet

import (
	"reflect"
	"unsafe"
)

const (
	GTPTypeTPDU uint8 = 255
)

const (
	GTPExtTypeNone    uint8 = 0
	GTPExtTypePDUSess       = 0x85
)

type GTPv1Hdr struct {
	Flags  Uint8
	Type   Uint8
	Length Uint16
	TEID   Uint32
	Seq    Uint16
	NPDU   Uint8
	Ext    Uint8
}

func DecodeGTPv1Hdr(b []byte) (*GTPv1Hdr, int) {
	h := (*GTPv1Hdr)(unsafe.Pointer(&b[0]))
	return h, 12
}

func (h GTPv1Hdr) Len() int {
	return 12
}

func (h GTPv1Hdr) Bytes() []byte {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = h.Len()
	hdr.Len = h.Len()
	hdr.Data = uintptr(unsafe.Pointer(&h))
	return b
}

func (h GTPv1Hdr) Version() uint8 {
	return h.Flags.Get() >> 5
}

func (h GTPv1Hdr) PT() uint8 {
	return (h.Flags.Get() >> 4) & 1
}

func (h GTPv1Hdr) HasExt() bool {
	return h.Flags.Get()&0x4 != 0
}

func (h GTPv1Hdr) HasSeq() bool {
	return h.Flags.Get()&0x2 != 0
}

func (h GTPv1Hdr) HasNPDU() bool {
	return h.Flags.Get()&0x1 != 0
}

const (
	GTPPDUTypeDL uint8 = 0
	GTPPDUTypeUL       = 1
)

type GTPExtPDUSess struct {
	TypeSpare Uint8
	FlagsQFI  Uint8
}

func DecodeGTPExtPDUSess(b []byte) (*GTPExtPDUSess, int) {
	e := (*GTPExtPDUSess)(unsafe.Pointer(&b[0]))
	return e, 2
}

func (e GTPExtPDUSess) Len() int {
	return 2
}

func (e GTPExtPDUSess) PDUType() uint8 {
	return (e.TypeSpare.Get() >> 4) & 0xf
}

func (e GTPExtPDUSess) PPI() uint8 {
	return (e.FlagsQFI.Get() >> 7) & 1
}

func (e GTPExtPDUSess) RQI() uint8 {
	return (e.FlagsQFI.Get() >> 6) & 1
}

func (e GTPExtPDUSess) QFI() uint8 {
	return e.FlagsQFI.Get() & 0x3f
}
