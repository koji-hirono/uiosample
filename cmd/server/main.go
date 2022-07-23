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
	"uiosample/znet"
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
			eth, _ := znet.DecodeEtherHdr(pkt)
			//log.Printf("EtherHdr: %#+v\n", eth)
			switch eth.Type.Get() {
			case znet.EtherTypeIPv4:
				err := procIPv4(d, eth, pkt[eth.Len():])
				if err != nil {
					hugetlb.Free(pkt)
					log.Fatal(err)
				}
			case znet.EtherTypeARP:
				err := procARP(d, eth, pkt[eth.Len():])
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
	b, _, err := hugetlb.Alloc(2048)
	if err != nil {
		return err
	}
	n := 0
	hdr, m := znet.DecodeEtherHdr(b)
	copy(hdr.Dst[:], []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	copy(hdr.Src[:], d.Mac)
	hdr.Type.Set(znet.EtherTypeARP)
	n += m

	arp, m := znet.DecodeARPHdr(b[m:])
	arp.HType.Set(1)
	arp.PType.Set(0x800)
	arp.HLen = 6
	arp.PLen = 4
	arp.Op.Set(znet.ARPRequest)
	copy(arp.SMac[:], d.Mac)
	copy(arp.SIP[:], []byte{30, 30, 0, 2})
	copy(arp.TMac[:], []byte{0, 0, 0, 0, 0, 0})
	copy(arp.TIP[:], []byte{30, 30, 0, 1})
	n += m

	log.Printf("Tx: %x\n", b[:n])
	for d.TxBurst([][]byte{b[:n]}) == 0 {
	}
	return nil
}

func procIPv4(d *e1000.Driver, eth *znet.EtherHdr, payload []byte) error {
	bDecodeIPv4.Start()
	ip, _ := znet.DecodeIPv4Hdr(payload)
	//log.Printf("IPv4Hdr: %#+v\n", ip)
	bDecodeIPv4.End()

	switch ip.Proto {
	case znet.IPProtoICMP:
		return procICMP(d, eth, ip, payload[ip.Len():])
	}
	return nil
}

func procICMP(d *e1000.Driver, eth *znet.EtherHdr, ip *znet.IPv4Hdr, payload []byte) error {
	bDecodeICMP.Start()
	icmp, _ := znet.DecodeICMPHdr(payload)
	//log.Printf("ICMPHdr: %#+v\n", icmp)
	switch icmp.Type {
	case znet.ICMPTypeEchoRequest:
	default:
		return nil
	}
	echo, _ := znet.DecodeICMPEchoHdr(payload[icmp.Len():])
	//log.Printf("ICMPEchoHdr: %#+v\n", echo)
	bDecodeICMP.End()

	bEncodeICMP.Start()
	b, _, err := hugetlb.Alloc(2048)
	if err != nil {
		return err
	}
	n := 0
	hdr, m := znet.DecodeEtherHdr(b)
	copy(hdr.Dst[:], eth.Src[:])
	copy(hdr.Src[:], d.Mac)
	copy(hdr.Type[:], eth.Type[:])
	n += m

	markipv4 := n
	ipv4, m := znet.DecodeIPv4Hdr(b[n:])
	ipv4.VerHL = ip.VerHL
	ipv4.ToS = ip.ToS
	copy(ipv4.Length[:], ip.Length[:])
	ipv4.ID.Set(0)
	ipv4.FlgOff.Set(0)
	ipv4.TTL = 64
	ipv4.Proto = znet.IPProtoICMP
	ipv4.Chksum.Set(0)
	copy(ipv4.Src[:], ip.Dst[:])
	copy(ipv4.Dst[:], ip.Src[:])
	n += m
	markipv4end := n

	markicmp := n
	txicmp, m := znet.DecodeICMPHdr(b[n:])
	txicmp.Type = znet.ICMPTypeEchoReply
	txicmp.Code = 0
	txicmp.Chksum.Set(0)
	n += m

	txecho, m := znet.DecodeICMPEchoHdr(b[n:])
	copy(txecho.ID[:], echo.ID[:])
	copy(txecho.Seq[:], echo.Seq[:])
	n += m

	m = copy(b[n:], payload[icmp.Len()+echo.Len():])
	n += m

	txicmp.Chksum.Set(znet.CalcChecksum(b[markicmp:n]))

	ipv4.Chksum.Set(znet.CalcChecksum(b[markipv4:markipv4end]))

	bEncodeICMP.End()
	//log.Printf("Tx: %x\n", b[:n])
	bTxICMP.Start()
	for d.TxBurst([][]byte{b[:n]}) == 0 {
	}
	bTxICMP.End()
	return nil
}

func procARP(d *e1000.Driver, eth *znet.EtherHdr, payload []byte) error {
	bDecodeARP.Start()
	arp, _ := znet.DecodeARPHdr(payload)
	//log.Printf("ARPHdr: %#+v\n", arp)
	bDecodeARP.End()

	bEncodeARP.Start()
	b, _, err := hugetlb.Alloc(2048)
	if err != nil {
		return err
	}
	n := 0
	hdr, m := znet.DecodeEtherHdr(b)
	copy(hdr.Dst[:], eth.Src[:])
	copy(hdr.Src[:], d.Mac[:])
	copy(hdr.Type[:], eth.Type[:])
	n += m

	txarp, m := znet.DecodeARPHdr(b[n:])
	copy(txarp.HType[:], arp.HType[:])
	copy(txarp.PType[:], arp.PType[:])
	txarp.HLen = arp.HLen
	txarp.PLen = arp.PLen
	txarp.Op.Set(znet.ARPReply)
	copy(txarp.SMac[:], d.Mac[:])
	copy(txarp.SIP[:], arp.TIP[:])
	copy(txarp.TMac[:], arp.SMac[:])
	copy(txarp.TIP[:], arp.SIP[:])
	n += m

	bEncodeARP.End()
	//log.Printf("Tx: %x\n", b[:n])
	bTxARP.Start()
	for d.TxBurst([][]byte{b[:n]}) == 0 {
	}
	bTxARP.End()
	return nil
}
