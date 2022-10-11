package e1000

import (
	"uiosample/ethdev"
)

type Counter struct {
	reg  Reg
	addr int
	val  uint64
}

func NewCounter(reg Reg, addr int) *Counter {
	return &Counter{reg: reg, addr: addr}
}

func (c *Counter) Value() uint64 {
	c.val += uint64(c.reg.Read(c.addr))
	return c.val
}

func (c *Counter) Clear() {
	c.reg.Read(c.addr)
	c.val = 0
}

type DCounter struct {
	reg  Reg
	addr [2]int
	val  uint64
}

func NewDCounter(reg Reg, high, low int) *DCounter {
	return &DCounter{reg: reg, addr: [...]int{high, low}}
}

func (c *DCounter) Value() uint64 {
	val := uint64(c.reg.Read(c.addr[0]))
	val <<= 32
	val += uint64(c.reg.Read(c.addr[1]))
	c.val += val
	return c.val
}

func (c *DCounter) Clear() {
	c.reg.Read(c.addr[0])
	c.reg.Read(c.addr[1])
	c.val = 0
}

type DummyCounter struct{}

func (DummyCounter) Value() uint64 {
	return 0
}

func (DummyCounter) Clear() {
}

func NewCounterGroup(reg Reg) *ethdev.CounterGroup {
	g := new(ethdev.CounterGroup)
	g.Ext = make(map[string]ethdev.Counter)
	g.Ext["MPC"] = NewCounter(reg, MPC)
	g.Ext["GPRC"] = NewCounter(reg, GPRC)
	g.Ext["GPTC"] = NewCounter(reg, GPTC)
	g.Ext["GORC"] = NewDCounter(reg, GORCH, GORCL)
	g.Ext["GOTC"] = NewDCounter(reg, GOTCH, GOTCL)

	g.RxPackets = g.Ext["GPRC"]
	g.TxPackets = g.Ext["GPTC"]
	g.RxOctets = g.Ext["GORC"]
	g.TxOctets = g.Ext["GOTC"]
	g.RxMissed = g.Ext["MPC"]
	g.RxErrors = DummyCounter{}
	g.TxErrors = DummyCounter{}

	return g
}
