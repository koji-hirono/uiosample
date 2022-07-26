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

	c, err := pci.OpenConfig(0)
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

	dev, err := pci.OpenDevice(addr, c)
	if err != nil {
		t.Fatal(err)
	}
	defer dev.Close()

	rxn := 8
	txn := 8
	d, err := NewDriver(dev, rxn, txn, nil)
	if err != nil {
		t.Fatal(err)
	}
	d.Init()

	pkts := make([][]byte, 8, 8)
	for {
		n := d.RxBurst(pkts)
		for i := 0; i < n; i++ {
			log.Printf("pkt: %x\n", pkts[i])
		}
	}
}
