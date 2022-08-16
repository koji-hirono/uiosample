package main

import (
	"net"
	"testing"

	"uiosample/znet"
)

func TestPDRTable_Find(t *testing.T) {
	t.Run("N3 match", func(t *testing.T) {
		tbl := NewPDRTable()
		tbl.Put(1, 1, &PDR{
			PDI: &PDI{
				UEAddr: net.IPv4(60, 60, 0, 2),
				FTEID: &FTEID{
					TEID:     78,
					GTPuAddr: net.IPv4(30, 30, 0, 1),
				},
			},
		})
		key := &PDRKey{
			Outer: &PDROuterKey{
				IP: &znet.IPv4Hdr{
					Src: znet.IPv4Addr([4]byte{30, 30, 0, 2}),
					Dst: znet.IPv4Addr([4]byte{30, 30, 0, 1}),
				},
				UDP: &znet.UDPHdr{},
				GTP: &znet.GTPv1Hdr{
					TEID: znet.Uint32([4]byte{0, 0, 0, 78}),
				},
				Sess: &znet.GTPExtPDUSess{
					TypeSpare: znet.Uint8(znet.GTPPDUTypeUL << 4),
				},
			},
			IP: &znet.IPv4Hdr{
				Src: znet.IPv4Addr([4]byte{60, 60, 0, 2}),
				Dst: znet.IPv4Addr([4]byte{70, 70, 0, 2}),
			},
		}
		pdr := tbl.Find(key)
		if pdr == nil {
			t.Fatal("not found")
		}
	})
	t.Run("N6 match", func(t *testing.T) {
		tbl := NewPDRTable()
		tbl.Put(1, 1, &PDR{
			PDI: &PDI{
				UEAddr: net.IPv4(60, 60, 0, 2),
			},
		})
		key := &PDRKey{
			IP: &znet.IPv4Hdr{
				Src: znet.IPv4Addr([4]byte{70, 70, 0, 2}),
				Dst: znet.IPv4Addr([4]byte{60, 60, 0, 2}),
			},
		}
		pdr := tbl.Find(key)
		if pdr == nil {
			t.Fatal("not found")
		}
	})
}
