package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	_ "time"

	"uiosample/e1000"
	"uiosample/hugetlb"
	"uiosample/pci"
)

func main() {
	prog := path.Base(os.Args[0])
	if len(os.Args) < 2 {
		fmt.Printf("usage: %v <PCI ID>\n", prog)
		os.Exit(1)
	}
	pciid, err := strconv.ParseUint(os.Args[1], 0, 8)
	if err != nil {
		log.Fatal(err)
	}
	hugetlb.SetPages(128)
	hugetlb.Init()

	addr := &pci.Addr{ID: uint8(pciid)}

	c, err := pci.NewConfig(0)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = c.SetBusMaster()
	if err != nil {
		log.Fatal(err)
	}

	s, err := c.Dump()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Config:\n%v\n", s)

	dev, err := pci.NewDevice(addr, c)
	if err != nil {
		log.Fatal(err)
	}

	rxn := 2
	txn := 2
	d := e1000.NewDriver(dev, rxn, txn)
	d.Init()

	ch := make(chan []byte, 1)
	defer close(ch)
	go d.Serve(ch)
	/*
	go func() {
		for {
			pkt := Packet{
				EtherHdr{
					Dst:  []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
					Src:  d.Mac,
					Type: EtherTypeARP,
				},
				ARPHdr{
					HType: 1,
					PType: 0x800,
					HLen:  6,
					PLen:  4,
					Op:    ARPRequest,
					SMac:  d.Mac,
					SIP:   []byte{30,30,0,2},
					TMac:  []byte{0,0,0,0,0,0},
					TIP:   []byte{30,30,0,1},
				},
			}
			n := pkt.Len()
			b := make([]byte, n)
			err = pkt.Encode(b)
			if err != nil {
				continue
			}
			log.Printf("Tx: %x\n", b)
			d.Tx(b)
			time.Sleep(time.Second * 2)
		}
	}()
	*/
	for pkt := range ch {
		log.Printf("Recv: %x\n", pkt)
		eth, err := DecodeEtherHdr(pkt)
		if err != nil {
			continue
		}
		log.Printf("EtherHdr: %#+v\n", eth)
		switch eth.Type {
		case EtherTypeIPv4:
			procIPv4(d, &eth, pkt[eth.Len():])
		case EtherTypeARP:
			procARP(d, &eth, pkt[eth.Len():])
		}
	}
}

func procIPv4(d *e1000.Driver, eth *EtherHdr, payload []byte) error {
	ip, err := DecodeIPv4Hdr(payload)
	if err != nil {
		return err
	}
	log.Printf("IPv4Hdr: %#+v\n", ip)

	switch ip.Proto {
	case IPProtoICMP:
		procICMP(d, eth, &ip, payload[ip.Len():])
	}
	return nil
}

func procICMP(d *e1000.Driver, eth *EtherHdr, ip *IPv4Hdr, payload []byte) error {
	icmp, err := DecodeICMPHdr(payload)
	if err != nil {
		return err
	}
	log.Printf("ICMPHdr: %#+v\n", icmp)
	switch icmp.Type {
	case ICMPTypeEchoRequest:
	default:
		return nil
	}
	echo, err := DecodeICMPEchoHdr(payload[icmp.Len():])
	if err != nil {
		return err
	}
	log.Printf("ICMPEchoHdr: %#+v\n", echo)
	txicmp := Packet{
		ICMPHdr{
			Type:   ICMPTypeEchoReply,
			Code:   0,
			Chksum: 0,
		},
		ICMPEchoHdr{
			ID:  echo.ID,
			Seq: echo.Seq,
		},
		Data(payload[icmp.Len()+echo.Len():]),
	}
	txicmp2 := Packet{
		ICMPHdr{
			Type:   ICMPTypeEchoReply,
			Code:   0,
			Chksum: uint16(^txicmp.Sum()),
		},
		ICMPEchoHdr{
			ID:  echo.ID,
			Seq: echo.Seq,
		},
		Data(payload[icmp.Len()+echo.Len():]),
	}
	ipv4 := IPv4Hdr{
		Version: ip.Version,
		Hdrlen:  ip.Hdrlen,
		ToS:     ip.ToS,
		Length:  ip.Length,
		ID:      0,
		Flags:   0,
		Offset:  0,
		TTL:     64,
		Proto:   IPProtoICMP,
		Chksum:  0,
		Src:     ip.Dst,
		Dst:     ip.Src,
	}
	ipv4.Chksum = uint16(^ipv4.Sum())
	pkt := Packet{
		EtherHdr{
			Dst:  eth.Src,
			Src:  d.Mac,
			Type: eth.Type,
		},
		ipv4,
		txicmp2,
	}

	n := pkt.Len()
	b := make([]byte, n)
	err = pkt.Encode(b)
	if err != nil {
		return err
	}
	log.Printf("Tx: %x\n", b)
	d.Tx(b)
	return nil
}

func procARP(d *e1000.Driver, eth *EtherHdr, payload []byte) error {
	arp, err := DecodeARPHdr(payload)
	if err != nil {
		return err
	}
	log.Printf("ARPHdr: %#+v\n", arp)
	pkt := Packet{
		EtherHdr{
			Dst:  eth.Src,
			Src:  d.Mac,
			Type: eth.Type,
		},
		ARPHdr{
			HType: arp.HType,
			PType: arp.PType,
			HLen:  arp.HLen,
			PLen:  arp.PLen,
			Op:    ARPReply,
			SMac:  d.Mac,
			SIP:   arp.TIP,
			TMac:  arp.SMac,
			TIP:   arp.SIP,
		},
	}

	n := pkt.Len()
	b := make([]byte, n)
	err = pkt.Encode(b)
	if err != nil {
		return err
	}
	log.Printf("Tx: %x\n", b)
	d.Tx(b)
	return nil
}
