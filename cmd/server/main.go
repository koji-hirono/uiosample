package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	_ "time"

	"uiosample/bench"
	"uiosample/e1000"
	"uiosample/hugetlb"
	"uiosample/pci"
)

var (
	bRx         = bench.New("Rx Packet")
	bDecodeIPv4 = bench.New("Decode IPv4")
	bDecodeICMP = bench.New("Decode ICMP")
	bDecodeARP  = bench.New("Decode ARP")
	bEncodeICMP = bench.New("Encode ICMP")
	bEncodeARP  = bench.New("Encode ARP")
	bTxICMP     = bench.New("Tx ICMP")
	bTxARP      = bench.New("Tx ARP")
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

	// rxn >= 8
	// txn >= 8
	rxn := 32
	txn := 32
	d := e1000.NewDriver(dev, rxn, txn, nil)
	d.Init()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	Serve(d, sig)

	bRx.Print()
	bDecodeIPv4.Print()
	bDecodeICMP.Print()
	bDecodeARP.Print()
	bEncodeICMP.Print()
	bEncodeARP.Print()
	bTxICMP.Print()
	bTxARP.Print()

	var stat e1000.Stat
	d.UpdateStat(&stat)

	PrintStat(&stat)

	hugetlb.Stat()
}

func Serve(d *e1000.Driver, sig chan os.Signal) {
	pkts := make([][]byte, 32, 32)
	for {
		select {
		case <-sig:
			return
		default:
		}
		bRx.Start()
		n := d.RxBurst(pkts)
		bRx.End()
		for i := 0; i < n; i++ {
			pkt := pkts[i]
			//log.Printf("Recv: %x\n", pkt)
			eth, err := DecodeEtherHdr(pkt)
			if err != nil {
				hugetlb.Free(pkt)
				continue
			}
			//log.Printf("EtherHdr: %#+v\n", eth)
			switch eth.Type {
			case EtherTypeIPv4:
				err := procIPv4(d, &eth, pkt[eth.Len():])
				if err != nil {
					hugetlb.Free(pkt)
					log.Fatal(err)
				}
			case EtherTypeARP:
				err := procARP(d, &eth, pkt[eth.Len():])
				if err != nil {
					hugetlb.Free(pkt)
					log.Fatal(err)
				}
			}
			hugetlb.Free(pkt)
		}
	}
}

func PrintStat(stat *e1000.Stat) {
	fmt.Printf("MPC : %v\n", stat.MPC)
	fmt.Printf("GPRC: %v\n", stat.GPRC)
	fmt.Printf("GPTC: %v\n", stat.GPTC)
	fmt.Printf("GORC: %v\n", stat.GORC)
	fmt.Printf("GOTC: %v\n", stat.GOTC)
}

func sendARPRequest(d *e1000.Driver) error {
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
			SIP:   []byte{30, 30, 0, 2},
			TMac:  []byte{0, 0, 0, 0, 0, 0},
			TIP:   []byte{30, 30, 0, 1},
		},
	}
	n := pkt.Len()
	b, _, err := hugetlb.Alloc(n)
	if err != nil {
		return err
	}
	err = pkt.Encode(b)
	if err != nil {
		return err
	}
	log.Printf("Tx: %x\n", b)
	for d.TxBurst([][]byte{b}) == 0 {
	}
	return nil
}

func procIPv4(d *e1000.Driver, eth *EtherHdr, payload []byte) error {
	bDecodeIPv4.Start()
	ip, err := DecodeIPv4Hdr(payload)
	if err != nil {
		return err
	}
	//log.Printf("IPv4Hdr: %#+v\n", ip)
	bDecodeIPv4.End()

	switch ip.Proto {
	case IPProtoICMP:
		return procICMP(d, eth, &ip, payload[ip.Len():])
	}
	return nil
}

func procICMP(d *e1000.Driver, eth *EtherHdr, ip *IPv4Hdr, payload []byte) error {
	bDecodeICMP.Start()
	icmp, err := DecodeICMPHdr(payload)
	if err != nil {
		return err
	}
	//log.Printf("ICMPHdr: %#+v\n", icmp)
	switch icmp.Type {
	case ICMPTypeEchoRequest:
	default:
		return nil
	}
	echo, err := DecodeICMPEchoHdr(payload[icmp.Len():])
	if err != nil {
		return err
	}
	//log.Printf("ICMPEchoHdr: %#+v\n", echo)
	bDecodeICMP.End()
	bEncodeICMP.Start()
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
	b, _, err := hugetlb.Alloc(n)
	if err != nil {
		return err
	}
	err = pkt.Encode(b)
	if err != nil {
		return err
	}
	bEncodeICMP.End()
	//log.Printf("Tx: %x\n", b)
	bTxICMP.Start()
	for d.TxBurst([][]byte{b}) == 0 {
	}
	bTxICMP.End()
	return nil
}

func procARP(d *e1000.Driver, eth *EtherHdr, payload []byte) error {
	bDecodeARP.Start()
	arp, err := DecodeARPHdr(payload)
	if err != nil {
		return err
	}
	//log.Printf("ARPHdr: %#+v\n", arp)
	bDecodeARP.End()
	bEncodeARP.Start()
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
	b, _, err := hugetlb.Alloc(n)
	if err != nil {
		return err
	}
	err = pkt.Encode(b)
	if err != nil {
		return err
	}
	bEncodeARP.End()
	//log.Printf("Tx: %x\n", b)
	bTxARP.Start()
	for d.TxBurst([][]byte{b}) == 0 {
	}
	bTxARP.End()
	return nil
}
