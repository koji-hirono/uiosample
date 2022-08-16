package main

import (
	"net"

	"uiosample/znet"
)

type FlowDesc struct {
	Action   uint8
	Dir      uint8
	Proto    uint8
	Src      net.IPNet
	Dst      net.IPNet
	SrcPorts [][]uint16
	DstPorts [][]uint16
}

type SDFFilter struct {
	FD  *FlowDesc
	TTC *uint16
	SPI *uint32
	FL  *uint32
	BID *uint32
}

type FTEID struct {
	TEID     uint32
	GTPuAddr net.IP
}

type PDI struct {
	UEAddr net.IP
	FTEID  *FTEID
	SDF    *SDFFilter
}

func (r PDI) Match(key *PDRKey) bool {
	outer := key.Outer
	if outer != nil {
		f := r.FTEID
		if f == nil {
			return false
		}
		if f.TEID != outer.GTP.TEID.Get() {
			return false
		}
		outerip := net.IP(outer.IP.Dst[:])
		if !f.GTPuAddr.Equal(outerip) {
			return false
		}
		var ip net.IP
		// TODO:
		/*
			switch outer.Sess.PDUType() {
			case znet.GTPPDUTypeDL:
				ip = net.IPv4(key.IP.Dst[:]...)
			case znet.GTPPDUTypeUL:
				ip = net.IPv4(key.IP.Src[:]...)
			}
		*/
		ip = net.IP(key.IP.Src[:])
		if !r.UEAddr.Equal(ip) {
			return false
		}
	} else {
		if r.FTEID != nil {
			return false
		}
		ip := net.IP(key.IP.Dst[:])
		if !r.UEAddr.Equal(ip) {
			return false
		}
	}
	return true
}

type PDROuterKey struct {
	IP   *znet.IPv4Hdr
	UDP  *znet.UDPHdr
	GTP  *znet.GTPv1Hdr
	Sess *znet.GTPExtPDUSess
}

type PDRKey struct {
	Outer *PDROuterKey
	IP    *znet.IPv4Hdr
}

type PDR struct {
	SEID            uint64
	ID              uint16
	Precedence      uint32
	PDI             *PDI
	OuterHdrRemoval *uint8
	FARID           uint32
	QERID           uint32
	URRID           uint32
}

type PDRTable struct {
	s map[uint64]map[uint16]*PDR
}

func NewPDRTable() *PDRTable {
	t := new(PDRTable)
	t.s = make(map[uint64]map[uint16]*PDR)
	return t
}

func (t *PDRTable) Put(seid uint64, id uint16, pdr *PDR) {
	_, ok := t.s[seid]
	if !ok {
		t.s[seid] = make(map[uint16]*PDR)
	}
	t.s[seid][id] = pdr
}

func (t *PDRTable) Delete(seid uint64, id uint16) {
	_, ok := t.s[seid]
	if !ok {
		return
	}
	delete(t.s[seid], id)
	if len(t.s[seid]) > 0 {
		return
	}
	delete(t.s, seid)
}

func (t *PDRTable) Find(key *PDRKey) *PDR {
	for _, pdrs := range t.s {
		for _, pdr := range pdrs {
			if pdr.PDI == nil {
				continue
			}
			if pdr.PDI.Match(key) {
				return pdr
			}
		}
	}
	return nil
}
