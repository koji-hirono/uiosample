package main

import (
	"encoding/binary"
)

type Encoder interface {
	Len() int
	Sum() uint32
	Encode([]byte) error
}

type Packet []Encoder

func (p Packet) Len() int {
	n := 0
	for _, e := range p {
		n += e.Len()
	}
	return n
}

func (p Packet) Sum() uint32 {
	sum := uint32(0)
	for _, e := range p {
		sum += e.Sum()
	}
	return sum
}

func (p Packet) Encode(b []byte) error {
	off := 0
	for _, e := range p {
		err := e.Encode(b[off:])
		if err != nil {
			return err
		}
		off += e.Len()
	}
	return nil
}

type Data []byte

func (d Data) Len() int {
	return len(d)
}

func (d Data) Sum() uint32 {
	var sum uint32
	n := len(d)
	for i := 0; i < n; i += 2 {
		x := uint32(d[i]) << 8
		if i + 1 < n {
			x |= uint32(d[i+1])
		}
		sum += x
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return sum
}

func (d Data) Encode(b []byte) error {
	copy(b, d)
	return nil
}

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

type ICMPHdr struct {
	Type   uint8
	Code   uint8
	Chksum uint16
}

const (
	ICMPTypeEchoReply   uint8 = 0
	ICMPTypeEchoRequest uint8 = 8
)

func DecodeICMPHdr(b []byte) (ICMPHdr, error) {
	var h ICMPHdr
	h.Type = b[0]
	h.Code = b[1]
	h.Chksum = binary.BigEndian.Uint16(b[2:4])
	return h, nil
}

func (h ICMPHdr) Len() int {
	return 4
}

func (h ICMPHdr) Sum() uint32 {
	n := h.Len()
	b := make([]byte, n)
	h.Encode(b)
	return Data(b).Sum()
}

func (h ICMPHdr) Encode(b []byte) error {
	b[0] = h.Type
	b[1] = h.Code
	binary.BigEndian.PutUint16(b[2:], h.Chksum)
	return nil
}

type ICMPEchoHdr struct {
	ID  uint16
	Seq uint16
}

func DecodeICMPEchoHdr(b []byte) (ICMPEchoHdr, error) {
	var h ICMPEchoHdr
	h.ID = binary.BigEndian.Uint16(b[0:2])
	h.Seq = binary.BigEndian.Uint16(b[2:4])
	return h, nil
}

func (h ICMPEchoHdr) Len() int {
	return 4
}

func (h ICMPEchoHdr) Sum() uint32 {
	n := h.Len()
	b := make([]byte, n)
	h.Encode(b)
	return Data(b).Sum()
}

func (h ICMPEchoHdr) Encode(b []byte) error {
	binary.BigEndian.PutUint16(b[0:2], h.ID)
	binary.BigEndian.PutUint16(b[2:4], h.Seq)
	return nil
}
