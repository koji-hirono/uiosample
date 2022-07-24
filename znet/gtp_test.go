package znet

import (
	"bytes"
	"testing"
)

func TestDecodeGTPv1_withPDUTypeUL(t *testing.T) {
	b := []byte{
		0x34, 0xff, 0x00, 0x5c, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x85, 0x01, 0x10, 0x09, 0x00,
	}

	hdr, n := DecodeGTPv1Hdr(b)
	if n != 12 {
		t.Errorf("want %v; but got %v\n", 12, n)
	}

	if hdr.Version() != 1 {
		t.Errorf("want %v; but got %v\n", 1, hdr.Version())
	}

	if hdr.PT() != 1 {
		t.Errorf("want %v; but got %v\n", 1, hdr.PT())
	}

	if !hdr.HasExt() {
		t.Errorf("want %v; but got %v\n", true, hdr.HasExt())
	}

	if hdr.HasSeq() {
		t.Errorf("want %v; but got %v\n", false, hdr.HasSeq())
	}

	if hdr.HasNPDU() {
		t.Errorf("want %v; but got %v\n", false, hdr.HasNPDU())
	}

	if hdr.Type.Get() != GTPTypeTPDU {
		t.Errorf("want %v; but got %v\n", GTPTypeTPDU, hdr.Type.Get())
	}

	if hdr.Length.Get() != uint16(92) {
		t.Errorf("want %v; but got %v\n", uint16(92), hdr.Length.Get())
	}

	if hdr.TEID.Get() != uint32(1) {
		t.Errorf("want %v; but got %v\n", uint32(1), hdr.TEID.Get())
	}

	if hdr.Ext.Get() != GTPExtTypePDUSess {
		t.Errorf("want %v; but got %v\n", GTPExtTypePDUSess, hdr.Ext.Get())
	}

	extlen := b[n]
	if extlen != 1 {
		t.Errorf("want %v; but got %v\n", 1, extlen)
	}
	n++

	ext, m := DecodeGTPExtPDUSess(b[n:])
	if m != 2 {
		t.Errorf("want %v; but got %v\n", 2, m)
	}
	n += m

	if ext.PDUType() != GTPPDUTypeUL {
		t.Errorf("want %v; but got %v\n", GTPPDUTypeUL, ext.PDUType())
	}

	if ext.QFI() != 9 {
		t.Errorf("want %v; but got %v\n", 9, ext.QFI())
	}

	nextext := b[n]
	if nextext != GTPExtTypeNone {
		t.Errorf("want %v; but got %v\n", GTPExtTypeNone, nextext)
	}
	n++

	if n != len(b) {
		t.Errorf("want %v; but got %v\n", len(b), n)
	}
}

func TestDecodeGTPv1_withPDUTypeDL(t *testing.T) {
	b := []byte{
		0x34, 0xff, 0x00, 0x5c, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x85, 0x01, 0x00, 0x09, 0x00,
	}

	hdr, n := DecodeGTPv1Hdr(b)
	if n != 12 {
		t.Errorf("want %v; but got %v\n", 12, n)
	}

	if hdr.Version() != 1 {
		t.Errorf("want %v; but got %v\n", 1, hdr.Version())
	}

	if hdr.PT() != 1 {
		t.Errorf("want %v; but got %v\n", 1, hdr.PT())
	}

	if !hdr.HasExt() {
		t.Errorf("want %v; but got %v\n", true, hdr.HasExt())
	}

	if hdr.HasSeq() {
		t.Errorf("want %v; but got %v\n", false, hdr.HasSeq())
	}

	if hdr.HasNPDU() {
		t.Errorf("want %v; but got %v\n", false, hdr.HasNPDU())
	}

	if hdr.Type.Get() != GTPTypeTPDU {
		t.Errorf("want %v; but got %v\n", GTPTypeTPDU, hdr.Type.Get())
	}

	if hdr.Length.Get() != uint16(92) {
		t.Errorf("want %v; but got %v\n", uint16(92), hdr.Length.Get())
	}

	if hdr.TEID.Get() != uint32(1) {
		t.Errorf("want %v; but got %v\n", uint32(1), hdr.TEID.Get())
	}

	if hdr.Ext.Get() != GTPExtTypePDUSess {
		t.Errorf("want %v; but got %v\n", GTPExtTypePDUSess, hdr.Ext.Get())
	}

	extlen := b[n]
	if extlen != 1 {
		t.Errorf("want %v; but got %v\n", 1, extlen)
	}
	n++

	ext, m := DecodeGTPExtPDUSess(b[n:])
	if m != 2 {
		t.Errorf("want %v; but got %v\n", 2, m)
	}
	n += m

	if ext.PDUType() != GTPPDUTypeDL {
		t.Errorf("want %v; but got %v\n", GTPPDUTypeDL, ext.PDUType())
	}

	if ext.PPI() != 0 {
		t.Errorf("want %v; but got %v\n", 0, ext.PPI())
	}

	if ext.RQI() != 0 {
		t.Errorf("want %v; but got %v\n", 0, ext.RQI())
	}

	if ext.QFI() != 9 {
		t.Errorf("want %v; but got %v\n", 9, ext.QFI())
	}

	nextext := b[n]
	if nextext != GTPExtTypeNone {
		t.Errorf("want %v; but got %v\n", GTPExtTypeNone, nextext)
	}
	n++

	if n != len(b) {
		t.Errorf("want %v; but got %v\n", len(b), n)
	}
}

func TestEncodeGTPv1_withPDUTypeDL(t *testing.T) {
	want := []byte{
		0x34, 0xff, 0x00, 0x5c, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x85, 0x01, 0x00, 0x09, 0x00,
	}

	b := make([]byte, 16)

	hdr, n := DecodeGTPv1Hdr(b)
	hdr.Flags.Set(1<<5 | 1<<4 | 1<<2)
	hdr.Type.Set(GTPTypeTPDU)
	hdr.Length.Set(92)
	hdr.TEID.Set(1)
	hdr.Seq.Set(0)
	hdr.NPDU.Set(0)
	hdr.Ext.Set(GTPExtTypePDUSess)

	// length
	b[n] = 1
	n++

	ext, m := DecodeGTPExtPDUSess(b[n:])
	ext.TypeSpare.Set(GTPPDUTypeDL << 4)
	ext.FlagsQFI.Set(9)
	n += m

	b[n] = GTPExtTypeNone
	n++

	if !bytes.Equal(b[:n], want) {
		t.Errorf("want %x; but got %x\n", b[:n], want)
	}
}
