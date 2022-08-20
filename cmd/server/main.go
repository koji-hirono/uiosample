package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
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
	addr, err := pci.ParseAddr(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	hugetlb.SetPages(128)
	hugetlb.Init()

	c, err := pci.OpenConfig(addr)
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

	dev, err := pci.OpenDevice(addr, c)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	// rxn >= 8
	// txn >= 8
	rxn := 32
	txn := 32
	d, err := e1000.NewDriver(dev, rxn, txn, nil)
	if err != nil {
		log.Fatal(err)
	}
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
	hdr.Dst.Set([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	hdr.Src.Set(d.Mac)
	hdr.Type.Set(znet.EtherTypeARP)
	n += m

	arp, m := znet.DecodeARPHdr(b[m:])
	arp.HType.Set(1)
	arp.PType.Set(0x800)
	arp.HLen.Set(6)
	arp.PLen.Set(4)
	arp.Op.Set(znet.ARPRequest)
	arp.SMac.Set(d.Mac)
	arp.SIP.Set([]byte{30, 30, 0, 2})
	arp.TMac.Set([]byte{0, 0, 0, 0, 0, 0})
	arp.TIP.Set([]byte{30, 30, 0, 1})
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

	switch ip.Proto.Get() {
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
	hdr.Dst = eth.Src
	hdr.Src.Set(d.Mac)
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

	bEncodeARP.End()
	//log.Printf("Tx: %x\n", b[:n])
	bTxARP.Start()
	for d.TxBurst([][]byte{b[:n]}) == 0 {
	}
	bTxARP.End()
	return nil
}
