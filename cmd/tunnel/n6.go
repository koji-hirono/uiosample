package main

import (
	"uiosample/hugetlb"
	"uiosample/znet"
)

func (s *Server) procN6(pkt []byte) {
	eth, _ := znet.DecodeEtherHdr(pkt)
	switch eth.Type.Get() {
	case znet.EtherTypeARP:
		s.procN6ARP(eth, pkt[eth.Len():])
	case znet.EtherTypeIPv4:
		s.procN6IPv4(eth, pkt[eth.Len():])
	}
	hugetlb.Free(pkt)
}

func (s *Server) procN6ARP(eth *znet.EtherHdr, payload []byte) error {
	return procARP(s.port2, eth, payload)
}

func (s *Server) procN6IPv4(eth *znet.EtherHdr, payload []byte) error {
	ip, n := znet.DecodeIPv4Hdr(payload)

	key := &PDRKey{IP: ip}
	pdr := s.pdrtbl.Find(key)
	if pdr == nil {
		return nil
	}
	qer := s.qertbl.Get(pdr.SEID, pdr.QERID)
	far := s.fartbl.Get(pdr.SEID, pdr.FARID)
	if far == nil {
		return nil
	}
	if far.Action&ApplyActionDROP != 0 {
		return nil
	} else if far.Action&ApplyActionFORW != 0 {
		param := far.Param
		if param == nil {
			return nil
		}
		creation := param.Creation
		if creation == nil {
			return nil
		}
		return s.procN6FAR(far, qer, ip, payload[n:])
	} else if far.Action&ApplyActionBUFF != 0 {
		// TODO:
		return nil
	} else {
		return nil
	}
}

func (s *Server) procN6FAR(far *FAR, qer *QER, ip *znet.IPv4Hdr, payload []byte) error {
	creation := far.Param.Creation
	dstmac, ok := LookupMAC(creation.PeerAddr)
	if !ok {
		// TODO: wait for ARP reply
		sendARPRequest(s.port1, creation.PeerAddr, s.GTPAddr)
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
	hdr.Src.Set(s.port1.Mac())
	hdr.Type.Set(znet.EtherTypeIPv4)
	n += m

	// encap gtp header
	// outer ip
	offouterip := n
	outerip, m := znet.DecodeIPv4Hdr(b[n:])
	outerip.VerHL.Set(4<<4 | 5)
	outerip.ToS.Set(0)
	outerip.Length.Set(0)
	outerip.ID.Set(0)
	outerip.FlgOff.Set(0)
	outerip.TTL.Set(64)
	outerip.Proto.Set(znet.IPProtoUDP)
	outerip.Chksum.Set(0)
	outerip.Src.Set(s.GTPAddr)
	outerip.Dst.Set(creation.PeerAddr)
	n += m

	// udp
	offudp := n
	udp, m := znet.DecodeUDPHdr(b[n:])
	udp.SrcPort.Set(s.GTPPort)
	udp.DstPort.Set(creation.Port)
	udp.Length.Set(0)
	udp.Chksum.Set(0)
	n += m

	// gtp
	offgtp := n
	gtp, m := znet.DecodeGTPv1Hdr(b[n:])
	gtp.Flags.Set(1<<5 | 1<<4 | 1<<2)
	gtp.Type.Set(znet.GTPTypeTPDU)
	gtp.Length.Set(0)
	gtp.TEID.Set(creation.TEID)
	gtp.Seq.Set(0)
	gtp.NPDU.Set(0)
	if qer != nil {
		gtp.Ext.Set(znet.GTPExtTypePDUSess)
		n += m

		// ext length
		b[n] = 1
		n++

		// gtp ext
		gtpext, m := znet.DecodeGTPExtPDUSess(b[n:])
		gtpext.TypeSpare.Set(znet.GTPPDUTypeDL << 4)
		gtpext.FlagsQFI.Set(qer.QFI)
		n += m

		b[n] = znet.GTPExtTypeNone
		n++
	} else {
		gtp.Ext.Set(znet.GTPExtTypeNone)
		n += m
	}

	m = copy(b[n:], ip.Bytes())
	n += m
	m = copy(b[n:], payload)
	n += m

	// update gtp length
	gtp.Length.Set(uint16(n - offgtp - 8))

	// update udp length
	udp.Length.Set(uint16(n - offudp))

	// update outerip length
	outerip.Length.Set(uint16(n - offouterip))

	// update outerip checksum
	chksum := znet.CalcChecksum(b[offouterip : offouterip+outerip.Len()])
	outerip.Chksum.Set(chksum)

	for s.port1.TxBurst([][]byte{b[:n]}) == 0 {
	}

	return nil
}
