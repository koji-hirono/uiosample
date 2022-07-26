package main

import (
	"uiosample/e1000"
	"uiosample/pci"
)

type Port struct {
	c      *pci.Config
	dev    *pci.Device
	driver *e1000.Driver
}

func OpenPort(unit int, addr *pci.Addr) (*Port, error) {
	c, err := pci.OpenConfig(unit)
	if err != nil {
		return nil, err
	}

	err = c.SetBusMaster()
	if err != nil {
		c.Close()
		return nil, err
	}

	dev, err := pci.OpenDevice(addr, c)
	if err != nil {
		c.Close()
		return nil, err
	}

	// rxn >= 8
	// txn >= 8
	rxn := 64
	txn := 64
	driver, err := e1000.NewDriver(dev, rxn, txn, nil)
	if err != nil {
		dev.Close()
		c.Close()
		return nil, err
	}
	driver.Init()

	return &Port{
		c:      c,
		dev:    dev,
		driver: driver,
	}, nil
}

func (p *Port) Close() {
	p.dev.Close()
	p.c.Close()
}

func (p *Port) Mac() []byte {
	return p.driver.Mac
}

func (p *Port) RxBurst(pkts [][]byte) int {
	return p.driver.RxBurst(pkts)
}

func (p *Port) TxBurst(pkts [][]byte) int {
	return p.driver.TxBurst(pkts)
}
