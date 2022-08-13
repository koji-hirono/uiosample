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
	"uiosample/em"
	"uiosample/ethdev"
	"uiosample/hugetlb"
	"uiosample/pci"
)

var (
	bRx1 = bench.New("Rx1 Packet")
	bRx2 = bench.New("Rx2 Packet")
	bTx1 = bench.New("Tx1 Packet")
	bTx2 = bench.New("Tx2 Packet")
)

type Device struct {
	c      *pci.Config
	dev    *pci.Device
	driver ethdev.Port
}

func OpenDevice(unit int, addr *pci.Addr) (*Device, error) {
	c, err := pci.OpenConfig(unit)
	if err != nil {
		return nil, err
	}

	err = c.SetBusMaster()
	if err != nil {
		c.Close()
		return nil, err
	}

	s, err := c.Dump()
	if err != nil {
		c.Close()
		return nil, err
	}
	log.Printf("Config:\n%v\n", s)

	dev, err := pci.OpenDevice(addr, c)
	if err != nil {
		c.Close()
		return nil, err
	}

	driver, err := em.AttachDriver(dev, nil)
	if err != nil {
		dev.Close()
		c.Close()
		return nil, err
	}

	config := &ethdev.Config{
		VNIC: true,
	}
	err = driver.Configure(1, 1, config)
	if err != nil {
		driver.Detach()
		dev.Close()
		c.Close()
		return nil, err
	}

	// rxn >= 8
	// txn >= 8
	rxn := 64
	txn := 64

	rxconfig := &ethdev.RxConfig{
		Threshold: ethdev.RingThreshold{
			Prefetch: 0x20,
			Host: 4,
			Writeback: 4,
		},
	}
	err = driver.RxQueueSetup(0, rxn, rxconfig)
	if err != nil {
		driver.Close()
		dev.Close()
		c.Close()
		return nil, err
	}

	txconfig := &ethdev.TxConfig{
		Threshold: ethdev.RingThreshold{
			Prefetch: 0x1f,
			Host: 1,
			Writeback: 1,
		},
	}
	err = driver.TxQueueSetup(0, txn, txconfig)
	if err != nil {
		driver.Close()
		dev.Close()
		c.Close()
		return nil, err
	}

	driver.Start()

	mac, _ := driver.GetMACAddr()
	log.Printf("MAC Address: %x\n", mac)

	driver.SetPromisc(true, true)

	return &Device{
		c:      c,
		dev:    dev,
		driver: driver,
	}, nil
}

func (d *Device) Close() {
	d.driver.Close()
	d.dev.Close()
	d.c.Close()
}

func main() {
	prog := path.Base(os.Args[0])
	if len(os.Args) < 3 {
		fmt.Printf("usage: %v <PCI ID> <PCI ID2>\n", prog)
		os.Exit(1)
	}
	hugetlb.SetPages(128)
	hugetlb.Init()

	addr1, err := pci.ParseAddr(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	addr2, err := pci.ParseAddr(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	dev1, err := OpenDevice(0, addr1)
	if err != nil {
		log.Fatal(err)
	}
	defer dev1.Close()
	log.Println("device 0 open.")

	dev2, err := OpenDevice(1, addr2)
	if err != nil {
		log.Fatal(err)
	}
	defer dev2.Close()
	log.Println("device 1 open.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	Serve(dev1.driver, dev2.driver, sig)

	bRx1.Print()
	bRx2.Print()
	bTx1.Print()
	bTx2.Print()

	PrintCounters(dev1.driver.CounterGroup())
	PrintCounters(dev2.driver.CounterGroup())

	hugetlb.Stat()
}

func Serve(d1 ethdev.Port, d2 ethdev.Port, sig chan os.Signal) {
	pkts := make([][]byte, 32, 32)
	for {
		select {
		case <-sig:
			return
		default:
		}
		bRx1.Start()
		n := d1.RxQueue(0).Do(pkts)
		bRx1.End()
		for off := 0; off < n; {
			bTx2.Start()
			m := d2.TxQueue(0).Do(pkts[off:n])
			bTx2.End()
			off += m
			if off < n {
				log.Printf("dev2 send [%v/%v]\n", off, n)
			}
		}
		bRx2.Start()
		n = d2.RxQueue(0).Do(pkts)
		bRx2.End()
		for off := 0; off < n; {
			bTx1.Start()
			m := d1.TxQueue(0).Do(pkts[off:n])
			bTx1.End()
			off += m
			if off < n {
				log.Printf("dev1 send [%v/%v]\n", off, n)
			}
		}
	}
}

func PrintCounters(g *ethdev.CounterGroup) {
	fmt.Printf("RxPackets: %v\n", g.RxPackets.Value())
	fmt.Printf("TxPackets: %v\n", g.TxPackets.Value())
	fmt.Printf("RxOctets : %v\n", g.RxOctets.Value())
	fmt.Printf("TxOctets : %v\n", g.TxOctets.Value())
	fmt.Printf("RxMissed : %v\n", g.RxMissed.Value())
	fmt.Printf("RxErrors : %v\n", g.RxErrors.Value())
	fmt.Printf("TxErrors : %v\n", g.TxErrors.Value())

	for name, c := range g.Ext {
		fmt.Printf("%s: %v\n", name, c.Value())
	}
}
