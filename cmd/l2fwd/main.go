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
	driver *e1000.Driver
}

func OpenDevice(unit int, addr *pci.Addr) (*Device, error) {
	c, err := pci.NewConfig(unit)
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

	dev, err := pci.NewDevice(addr, c)
	if err != nil {
		c.Close()
		return nil, err
	}

	// rxn >= 8
	// txn >= 8
	rxn := 64
	txn := 64
	driver := e1000.NewDriver(dev, rxn, txn, nil)
	driver.Init()

	return &Device{
		c:      c,
		dev:    dev,
		driver: driver,
	}, nil
}

func (d *Device) Close() {
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

	dev2, err := OpenDevice(1, addr2)
	if err != nil {
		log.Fatal(err)
	}
	defer dev2.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	Serve(dev1.driver, dev2.driver, sig)

	bRx1.Print()
	bRx2.Print()
	bTx1.Print()
	bTx2.Print()

	var stat1 e1000.Stat
	dev1.driver.UpdateStat(&stat1)
	PrintStat(&stat1)

	var stat2 e1000.Stat
	dev2.driver.UpdateStat(&stat2)

	PrintStat(&stat2)

	hugetlb.Stat()
}

func Serve(d1 *e1000.Driver, d2 *e1000.Driver, sig chan os.Signal) {
	pkts := make([][]byte, 32, 32)
	for {
		select {
		case <-sig:
			return
		default:
		}
		bRx1.Start()
		n := d1.RxBurst(pkts)
		bRx1.End()
		for off := 0; off < n; {
			bTx2.Start()
			m := d2.TxBurst(pkts[off:n])
			bTx2.End()
			off += m
			if off < n {
				log.Printf("dev2 send [%v/%v]\n", off, n)
			}
		}
		bRx2.Start()
		n = d2.RxBurst(pkts)
		bRx2.End()
		for off := 0; off < n; {
			bTx1.Start()
			m := d1.TxBurst(pkts[off:n])
			bTx1.End()
			off += m
			if off < n {
				log.Printf("dev1 send [%v/%v]\n", off, n)
			}
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
