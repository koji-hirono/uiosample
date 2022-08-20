package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"uiosample/ethdev"
	"uiosample/hugetlb"
	"uiosample/pci"
)

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

	port1, err := OpenPort(addr1)
	if err != nil {
		log.Fatal(err)
	}
	defer port1.Close()

	port2, err := OpenPort(addr2)
	if err != nil {
		log.Fatal(err)
	}
	defer port2.Close()

	s := NewServer(port1, port2)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	s.Serve(sig)

	PrintCounters(port1.driver.CounterGroup())
	PrintCounters(port2.driver.CounterGroup())

	hugetlb.Stat()
}

func PrintCounter(name string, c ethdev.Counter) {
	if c == nil {
		return
	}
	fmt.Printf("%s: %v\n", name, c.Value())
}

func PrintCounters(g *ethdev.CounterGroup) {
	PrintCounter("RxPackets", g.RxPackets)
	PrintCounter("TxPackets", g.TxPackets)
	PrintCounter("RxOctets ", g.RxOctets)
	PrintCounter("TxOctets ", g.TxOctets)
	PrintCounter("RxMissed ", g.RxMissed)
	PrintCounter("RxErrors ", g.RxErrors)
	PrintCounter("TxErrors ", g.TxErrors)
	for name, c := range g.Ext {
		PrintCounter(name, c)
	}
}
