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

	c, err := pci.OpenConfig(addr)
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

	d, err := AttachDriver(dev, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	err = d.Configure(1, 1, nil)
	if err != nil {
		t.Fatal(err)
	}

	rxd := 8
	txd := 8
	err = d.RxQueueSetup(0, rxd, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = d.TxQueueSetup(0, txd, nil)
	if err != nil {
		t.Fatal(err)
	}

	d.Start()
	defer d.Stop()

	rxq := d.RxQueue(0)

	pkts := make([][]byte, 8, 8)
	for {
		n := rxq.Do(pkts)
		for i := 0; i < n; i++ {
			log.Printf("pkt: %x\n", pkts[i])
		}
	}
}
