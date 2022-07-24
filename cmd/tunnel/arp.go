package main

import (
	"uiosample/e1000"
	"uiosample/hugetlb"
	"uiosample/znet"
)

func procARP(d *e1000.Driver, eth *znet.EtherHdr, payload []byte) error {
	arp, _ := znet.DecodeARPHdr(payload)

	b, _, err := hugetlb.Alloc(2048)
	if err != nil {
		return err
	}
	n := 0
	hdr, m := znet.DecodeEtherHdr(b)
	hdr.Dst = eth.Src
	hdr.Src.Set(d.Mac)
	hdr.Type = eth.Type
	n += m

	txarp, m := znet.DecodeARPHdr(b[n:])
	txarp.HType = arp.HType
	txarp.PType = arp.PType
	txarp.HLen = arp.HLen
	txarp.PLen = arp.PLen
	txarp.Op.Set(znet.ARPReply)
	txarp.SMac.Set(d.Mac)
	txarp.SIP = arp.TIP
	txarp.TMac = arp.SMac
	txarp.TIP = arp.SIP
	n += m

	for d.TxBurst([][]byte{b[:n]}) == 0 {
	}
	return nil
}
