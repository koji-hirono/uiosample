package main

import (
	"encoding/binary"
)

type IPv4Hdr struct {
	Version uint8
	Hdrlen  uint8
	ToS     uint8
	Length  uint16
	ID      uint16
	Flags   uint8
	Offset  uint16
	TTL     uint8
	Proto   uint8
	Chksum  uint16
	Src     []byte
	Dst     []byte
}

const (
	IPProtoICMP = 1
)

func DecodeIPv4Hdr(b []byte) (IPv4Hdr, error) {
	var h IPv4Hdr
	h.Version = (b[0] >> 4) & 0xf
	h.Hdrlen = b[0] & 0xf
	h.ToS = b[1]
	h.Length = binary.BigEndian.Uint16(b[2:4])
	h.ID = binary.BigEndian.Uint16(b[4:6])
	h.Flags = (b[6] >> 5) & 0x7
	h.Offset = uint16(b[6] & 0x1f)
	h.Offset <<= 8
	h.Offset |= uint16(b[7])
	h.TTL = b[8]
	h.Proto = b[9]
	h.Chksum = binary.BigEndian.Uint16(b[10:12])
	h.Src = b[12:16]
	h.Dst = b[16:20]
	return h, nil
}

func (h IPv4Hdr) Len() int {
	return 20
}

func (h IPv4Hdr) Sum() uint32 {
	n := h.Len()
	b := make([]byte, n)
	h.Encode(b)
	return Data(b).Sum()
}

func (h IPv4Hdr) Encode(b []byte) error {
	b[0] = h.Version<<4 | h.Hdrlen
	b[1] = h.ToS
	binary.BigEndian.PutUint16(b[2:4], h.Length)
	binary.BigEndian.PutUint16(b[4:6], h.ID)
	b[6] = h.Flags<<5 | byte(h.Offset>>8)
	b[7] = byte(h.Offset)
	b[8] = h.TTL
	b[9] = h.Proto
	binary.BigEndian.PutUint16(b[10:12], h.Chksum)
	copy(b[12:16], h.Src)
	copy(b[16:20], h.Dst)
	return nil
}
