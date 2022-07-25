package main

import (
	"bytes"

	"uiosample/hugetlb"
	"uiosample/znet"
)

func (s *Server) procN3(pkt []byte) {
	eth, _ := znet.DecodeEtherHdr(pkt)
	switch eth.Type.Get() {
	case znet.EtherTypeARP:
		s.procN3ARP(eth, pkt[eth.Len():])
	case znet.EtherTypeIPv4:
		s.procN3IPv4(eth, pkt[eth.Len():])
	}
	hugetlb.Free(pkt)
}

func (s *Server) procN3IPv4(eth *znet.EtherHdr, payload []byte) error {
	ip, _ := znet.DecodeIPv4Hdr(payload)
	switch ip.Proto.Get() {
	case znet.IPProtoICMP:
		return s.procN3ICMP(eth, ip, payload[ip.Len():])
	case znet.IPProtoUDP:
		return s.procN3UDP(eth, ip, payload[ip.Len():])
	}
	return nil
}

func (s *Server) procN3ARP(eth *znet.EtherHdr, payload []byte) error {
	return procARP(s.port1, eth, payload)
}

func (s *Server) procN3ICMP(eth *znet.EtherHdr, ip *znet.IPv4Hdr, payload []byte) error {
	return procICMP(s.port1, eth, ip, payload)
}

func (s *Server) procN3UDP(eth *znet.EtherHdr, ip *znet.IPv4Hdr, payload []byte) error {
	n := 0
	udp, m := znet.DecodeUDPHdr(payload)
	if !MatchTunnel(ip, udp) {
		return nil
	}
	n += m

	gtp, m := znet.DecodeGTPv1Hdr(payload[n:])
	n += m

	// length
	n++

	gtpext, m := znet.DecodeGTPExtPDUSess(payload[n:])
	n += m

	// next ext
	n++

	innerip, m := znet.DecodeIPv4Hdr(payload[n:])
	n += m

	if !MatchN3PDR(gtp, gtpext, innerip) {
		return nil
	}

	return s.procN3FAR(gtp, innerip, payload[n:])
}

func (s *Server) procN3FAR(gtp *znet.GTPv1Hdr, ip *znet.IPv4Hdr, payload []byte) error {
	dstmac, ok := LookupMAC(ip.Dst[:])
	if !ok {
		// TODO: wait for ARP reply
		sendARPRequest(s.port2, ip.Dst[:], ip.Src[:])
		return nil
	}

	b, _, err := hugetlb.Alloc(2048)
	if err != nil {
		return err
	}

	n := 0

	// new ether header
	hdr, m := znet.DecodeEtherHdr(b)
	hdr.Dst.Set(dstmac)
	hdr.Src.Set(s.port2.driver.Mac)
	hdr.Type.Set(znet.EtherTypeIPv4)
	n += m

	// decap gtp header
	m = copy(b[n:], ip.Bytes())
	n += m
	m = copy(b[n:], payload)
	n += m

	for s.port2.TxBurst([][]byte{b[:n]}) == 0 {
	}

	return nil
}

func MatchTunnel(ip *znet.IPv4Hdr, udp *znet.UDPHdr) bool {
	if udp.SrcPort.Get() != 2152 {
		return false
	}
	if udp.DstPort.Get() != 2152 {
		return false
	}
	if !bytes.Equal(ip.Src[:], []byte{30, 30, 0, 2}) {
		return false
	}
	if !bytes.Equal(ip.Dst[:], []byte{30, 30, 0, 1}) {
		return false
	}
	return true
}

func MatchN3PDR(gtp *znet.GTPv1Hdr, gtpext *znet.GTPExtPDUSess, ip *znet.IPv4Hdr) bool {
	if gtp.TEID.Get() != 78 {
		return false
	}
	if !bytes.Equal(ip.Src[:], []byte{60, 60, 0, 2}) {
		return false
	}
	return true
}
