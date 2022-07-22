package e1000

import (
	"log"
	"testing"

	"uiosample/hugetlb"
	"uiosample/pci"
)

func TestDriver(t *testing.T) {
	hugetlb.SetPages(128)
	hugetlb.Init()

	addr := &pci.Addr{ID: 17}

	c, err := pci.NewConfig(0)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	err = c.SetBusMaster()
	if err != nil {
		t.Fatal(err)
	}

	s, err := c.Dump()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Config:\n%v\n", s)

	dev, err := pci.NewDevice(addr, c)
	if err != nil {
		t.Fatal(err)
	}

	rxn := 8
	txn := 8
	d := NewDriver(dev, rxn, txn, nil)
	d.Init()
	ch := make(chan []byte, 10)
	defer close(ch)
	go d.Serve(ch)
	for pkt := range ch {
		log.Printf("pkt: %x\n", pkt)
	}
}
