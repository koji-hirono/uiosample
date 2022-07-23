package main

import (
	"encoding/binary"
)

type EtherHdr struct {
	Dst  []byte
	Src  []byte
	Type uint16
}

const (
	EtherTypeIPv4 = 0x0800
	EtherTypeARP  = 0x0806
)

func DecodeEtherHdr(b []byte) (EtherHdr, error) {
	var h EtherHdr
	h.Dst = b[0:6]
	h.Src = b[6:12]
	h.Type = binary.BigEndian.Uint16(b[12:14])
	return h, nil
}

func (h EtherHdr) Len() int {
	return 14
}

func (h EtherHdr) Sum() uint32 {
	n := h.Len()
	b := make([]byte, n)
	h.Encode(b)
	return Data(b).Sum()
}

func (h EtherHdr) Encode(b []byte) error {
	copy(b[0:], h.Dst)
	copy(b[6:], h.Src)
	binary.BigEndian.PutUint16(b[12:], h.Type)
	return nil
}
