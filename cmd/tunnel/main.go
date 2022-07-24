package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"

	"uiosample/e1000"
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

	var stat1 e1000.Stat
	port1.driver.UpdateStat(&stat1)
	PrintStat(&stat1)

	var stat2 e1000.Stat
	port2.driver.UpdateStat(&stat2)

	PrintStat(&stat2)

	hugetlb.Stat()
}

func PrintStat(stat *e1000.Stat) {
	fmt.Printf("MPC : %v\n", stat.MPC)
	fmt.Printf("GPRC: %v\n", stat.GPRC)
	fmt.Printf("GPTC: %v\n", stat.GPTC)
	fmt.Printf("GORC: %v\n", stat.GORC)
	fmt.Printf("GOTC: %v\n", stat.GOTC)
}
