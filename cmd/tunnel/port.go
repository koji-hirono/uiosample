package main

import (
	"uiosample/em"
	"uiosample/ethdev"
	"uiosample/pci"
)

type Port struct {
	c      *pci.Config
	dev    *pci.Device
	driver ethdev.Port
}

func OpenPort(addr *pci.Addr) (*Port, error) {
	c, err := pci.OpenConfig(addr)
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
		driver.Close()
		dev.Close()
		c.Close()
		return nil, err
	}

	// rxn >= 8
	// txn >= 8
	rxn := 64
	txn := 64

	rxconfig := &ethdev.RxConfig{}
	err = driver.RxQueueSetup(0, rxn, rxconfig)
	if err != nil {
		driver.Close()
		dev.Close()
		c.Close()
		return nil, err
	}

	txconfig := &ethdev.TxConfig{}
	err = driver.TxQueueSetup(0, txn, txconfig)
	if err != nil {
		driver.Close()
		dev.Close()
		c.Close()
		return nil, err
	}

	err = driver.Start()
	if err != nil {
		driver.Close()
		dev.Close()
		c.Close()
		return nil, err
	}

	driver.SetPromisc(true, true)

	return &Port{
		c:      c,
		dev:    dev,
		driver: driver,
	}, nil
}

func (p *Port) Close() {
	p.driver.Close()
	p.dev.Close()
	p.c.Close()
}

func (p *Port) Mac() []byte {
	mac, _ := p.driver.GetMACAddr()
	return mac[:]
}

func (p *Port) RxBurst(pkts [][]byte) int {
	return p.driver.RxQueue(0).Do(pkts)
}

func (p *Port) TxBurst(pkts [][]byte) int {
	return p.driver.TxQueue(0).Do(pkts)
}
