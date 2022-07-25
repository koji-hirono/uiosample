package main

import (
	"uiosample/hugetlb"
	"uiosample/znet"
)

func procICMP(port *Port, eth *znet.EtherHdr, ip *znet.IPv4Hdr, payload []byte) error {
	icmp, _ := znet.DecodeICMPHdr(payload)
	switch icmp.Type {
	case znet.ICMPTypeEchoRequest:
	default:
		return nil
	}
	echo, _ := znet.DecodeICMPEchoHdr(payload[icmp.Len():])

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

	markipv4 := n
	ipv4, m := znet.DecodeIPv4Hdr(b[n:])
	ipv4.VerHL = ip.VerHL
	ipv4.ToS = ip.ToS
	ipv4.Length = ip.Length
	ipv4.ID.Set(0)
	ipv4.FlgOff.Set(0)
	ipv4.TTL.Set(64)
	ipv4.Proto.Set(znet.IPProtoICMP)
	ipv4.Chksum.Set(0)
	ipv4.Src = ip.Dst
	ipv4.Dst = ip.Src
	n += m
	markipv4end := n

	markicmp := n
	txicmp, m := znet.DecodeICMPHdr(b[n:])
	txicmp.Type.Set(znet.ICMPTypeEchoReply)
	txicmp.Code.Set(0)
	txicmp.Chksum.Set(0)
	n += m

	txecho, m := znet.DecodeICMPEchoHdr(b[n:])
	txecho.ID = echo.ID
	txecho.Seq = echo.Seq
	n += m

	m = copy(b[n:], payload[icmp.Len()+echo.Len():])
	n += m

	txicmp.Chksum.Set(znet.CalcChecksum(b[markicmp:n]))

	ipv4.Chksum.Set(znet.CalcChecksum(b[markipv4:markipv4end]))

	for port.TxBurst([][]byte{b[:n]}) == 0 {
	}
	return nil
}
