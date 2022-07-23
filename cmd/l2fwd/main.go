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
	bRx1 = bench.New("Rx1 Packet")
	bRx2 = bench.New("Rx2 Packet")
)

type Device struct {
	c      *pci.Config
	dev    *pci.Device
	driver *e1000.Driver
}

func OpenDevice(unit int, addr *pci.Addr, rx, tx chan []byte) (*Device, error) {
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
	rxn := 8
	txn := 8
	driver := e1000.NewDriver(dev, rxn, txn, nil)
	driver.Init()

	go driver.Serve(rx)
	go driver.ServeTx(tx)

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

	ch12 := make(chan []byte, 10)
	defer close(ch12)
	ch21 := make(chan []byte, 10)
	defer close(ch21)

	pciid1, err := strconv.ParseUint(os.Args[1], 0, 8)
	if err != nil {
		log.Fatal(err)
	}
	addr1 := &pci.Addr{ID: uint8(pciid1)}

	pciid2, err := strconv.ParseUint(os.Args[2], 0, 8)
	if err != nil {
		log.Fatal(err)
	}
	addr2 := &pci.Addr{ID: uint8(pciid2)}

	dev1, err := OpenDevice(0, addr1, ch21, ch12)
	if err != nil {
		log.Fatal(err)
	}
	defer dev1.Close()

	dev2, err := OpenDevice(1, addr2, ch12, ch21)
	if err != nil {
		log.Fatal(err)
	}
	defer dev2.Close()

	var stat e1000.Stat
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-sig:
			bRx1.Print()
			bRx2.Print()
			dev1.driver.UpdateStat(&stat)
			PrintStat(&stat)
			dev2.driver.UpdateStat(&stat)
			PrintStat(&stat)
			hugetlb.Stat()
			os.Exit(0)
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
