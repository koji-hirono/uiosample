package main

import (
	"encoding/binary"
)

type ARPHdr struct {
	HType uint16
	PType uint16
	HLen  uint8
	PLen  uint8
	Op    uint16
	SMac  []byte
	SIP   []byte
	TMac  []byte
	TIP   []byte
}

// ARPHdr.Op
const (
	ARPRequest = 1
	ARPReply   = 2
)

func DecodeARPHdr(b []byte) (ARPHdr, error) {
	var h ARPHdr
	h.HType = binary.BigEndian.Uint16(b[0:2])
	h.PType = binary.BigEndian.Uint16(b[2:4])
	h.HLen = b[4]
	h.PLen = b[5]
	h.Op = binary.BigEndian.Uint16(b[6:8])
	h.SMac = b[8:14]
	h.SIP = b[14:18]
	h.TMac = b[18:24]
	h.TIP = b[24:28]
	return h, nil
}

func (h ARPHdr) Len() int {
	return 28
}

func (h ARPHdr) Sum() uint32 {
	n := h.Len()
	b := make([]byte, n)
	h.Encode(b)
	return Data(b).Sum()
}

func (h ARPHdr) Encode(b []byte) error {
	binary.BigEndian.PutUint16(b[0:], h.HType)
	binary.BigEndian.PutUint16(b[2:], h.PType)
	b[4] = h.HLen
	b[5] = h.PLen
	binary.BigEndian.PutUint16(b[6:], h.Op)
	copy(b[8:], h.SMac)
	copy(b[14:], h.SIP)
	copy(b[18:], h.TMac)
	copy(b[24:], h.TIP)
	return nil
}
