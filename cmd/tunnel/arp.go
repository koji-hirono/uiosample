package main

import (
	"sync"

	"uiosample/hugetlb"
	"uiosample/znet"
)

type ARPEntry struct {
	MAC  [6]byte
	Port *Port
}

var ARPTable sync.Map

func LookupMAC(ip []byte) ([]byte, bool) {
	var key [4]byte
	copy(key[:], ip)
	i, ok := ARPTable.Load(key)
	if !ok {
		return nil, false
	}
	e, ok := i.(*ARPEntry)
	if !ok {
		return nil, false
	}
	return e.MAC[:], true
}

func AddARPEntry(ip []byte, mac [6]byte, port *Port) {
	var key [4]byte
	copy(key[:], ip)
	e := &ARPEntry{MAC: mac, Port: port}
	ARPTable.Store(key, e)
}

func procARP(port *Port, eth *znet.EtherHdr, payload []byte) error {
	arp, _ := znet.DecodeARPHdr(payload)
	switch arp.Op.Get() {
	case znet.ARPRequest:
	case znet.ARPReply:
		AddARPEntry(arp.SIP[:], arp.SMac, port)
		return nil
	default:
		return nil
	}

	AddARPEntry(arp.SIP[:], arp.SMac, port)

	if arp.Op.Get() == znet.ARPReply {
		return nil
	}

	b, _, err := hugetlb.Alloc(2048)
	if err != nil {
		return err
	}
	n := 0
	hdr, m := znet.DecodeEtherHdr(b)
	hdr.Dst = eth.Src
	hdr.Src.Set(port.Mac())
	hdr.Type = eth.Type
	n += m

	txarp, m := znet.DecodeARPHdr(b[n:])
	txarp.HType = arp.HType
	txarp.PType = arp.PType
	txarp.HLen = arp.HLen
	txarp.PLen = arp.PLen
	txarp.Op.Set(znet.ARPReply)
	txarp.SMac.Set(port.Mac())
	txarp.SIP = arp.TIP
	txarp.TMac = arp.SMac
	txarp.TIP = arp.SIP
	n += m

	for port.TxBurst([][]byte{b[:n]}) == 0 {
	}
	return nil
}

func sendARPRequest(port *Port, dst, src []byte) error {
	b, _, err := hugetlb.Alloc(2048)
	if err != nil {
		return err
	}
	n := 0
	hdr, m := znet.DecodeEtherHdr(b)
	hdr.Dst.Set([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	hdr.Src.Set(port.Mac())
	hdr.Type.Set(znet.EtherTypeARP)
	n += m

	arp, m := znet.DecodeARPHdr(b[m:])
	arp.HType.Set(1)
	arp.PType.Set(0x800)
	arp.HLen.Set(6)
	arp.PLen.Set(4)
	arp.Op.Set(znet.ARPRequest)
	arp.SMac.Set(port.Mac())
	arp.SIP.Set(src)
	arp.TMac.Set([]byte{0, 0, 0, 0, 0, 0})
	arp.TIP.Set(dst)
	n += m

	for port.TxBurst([][]byte{b[:n]}) == 0 {
	}
	return nil
}
