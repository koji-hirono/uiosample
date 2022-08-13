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

	port1, err := OpenPort(0, addr1)
	if err != nil {
		log.Fatal(err)
	}
	defer port1.Close()

	port2, err := OpenPort(1, addr2)
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
