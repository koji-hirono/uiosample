package em

import (
	"uiosample/ethdev"
)

type Counter struct {
	hw  *HW
	reg int
	val uint64
}

func NewCounter(hw *HW, reg int) *Counter {
	return &Counter{hw: hw, reg: reg}
}

func (c *Counter) Value() uint64 {
	c.val += uint64(c.hw.RegRead(c.reg))
	return c.val
}

func (c *Counter) Clear() {
	c.hw.RegRead(c.reg)
	c.val = 0
}

type DCounter struct {
	hw  *HW
	reg [2]int
	val uint64
}

func NewDCounter(hw *HW, high, low int) *DCounter {
	return &DCounter{hw: hw, reg: [...]int{high, low}}
}

func (c *DCounter) Value() uint64 {
	val := uint64(c.hw.RegRead(c.reg[0]))
	val <<= 32
	val += uint64(c.hw.RegRead(c.reg[1]))
	c.val += val
	return c.val
}

func (c *DCounter) Clear() {
	c.hw.RegRead(c.reg[0])
	c.hw.RegRead(c.reg[1])
	c.val = 0
}

func NewCounterGroup(hw *HW) *ethdev.CounterGroup {
	g := new(ethdev.CounterGroup)
	var rxerrors []ethdev.Counter
	g.Ext = make(map[string]ethdev.Counter)
	if hw.PHY.MediaType == MediaTypeCopper {
		g.Ext["SYMERRS"] = NewCounter(hw, SYMERRS)
		g.Ext["SEC"] = NewCounter(hw, SEC)
	}
	g.Ext["CRCERRS"] = NewCounter(hw, CRCERRS)
	rxerrors = append(rxerrors, g.Ext["CRCERRS"])
	g.Ext["MPC"] = NewCounter(hw, MPC)
	g.Ext["SCC"] = NewCounter(hw, SCC)
	g.Ext["ECOL"] = NewCounter(hw, ECOL)

	g.Ext["MCC"] = NewCounter(hw, MCC)
	g.Ext["LATECOL"] = NewCounter(hw, LATECOL)
	g.Ext["COLC"] = NewCounter(hw, COLC)
	g.Ext["DC"] = NewCounter(hw, DC)
	g.Ext["RLEC"] = NewCounter(hw, RLEC)
	rxerrors = append(rxerrors, g.Ext["RLEC"])
	g.Ext["XONRXC"] = NewCounter(hw, XONRXC)
	g.Ext["XONTXC"] = NewCounter(hw, XONTXC)

	// For watchdog management we need to know if we have been
	// paused during the last interval, so capture that here.
	pauseFrame := NewCounter(hw, XOFFRXC)
	g.Ext["XOFFRXC"] = pauseFrame
	g.Ext["XOFFTXC"] = NewCounter(hw, XOFFTXC)
	g.Ext["FCRUC"] = NewCounter(hw, FCRUC)
	g.Ext["PRC64"] = NewCounter(hw, PRC64)
	g.Ext["PRC127"] = NewCounter(hw, PRC127)
	g.Ext["PRC255"] = NewCounter(hw, PRC255)
	g.Ext["PRC511"] = NewCounter(hw, PRC511)
	g.Ext["PRC1023"] = NewCounter(hw, PRC1023)
	g.Ext["PRC1522"] = NewCounter(hw, PRC1522)
	g.Ext["GPRC"] = NewCounter(hw, GPRC)
	g.Ext["BPRC"] = NewCounter(hw, BPRC)
	g.Ext["MPRC"] = NewCounter(hw, MPRC)
	g.Ext["GPTC"] = NewCounter(hw, GPTC)

	// For the 64-bit byte counters the low dword must be read first.
	// Both registers clear on the read of the high dword.
	g.Ext["GORC"] = NewDCounter(hw, GORCH, GORCL)
	g.Ext["GOTC"] = NewDCounter(hw, GOTCH, GOTCL)

	g.Ext["RNBC"] = NewCounter(hw, RNBC)
	g.Ext["RUC"] = NewCounter(hw, RUC)
	g.Ext["RFC"] = NewCounter(hw, RFC)
	g.Ext["ROC"] = NewCounter(hw, ROC)
	g.Ext["RJC"] = NewCounter(hw, RJC)

	g.Ext["TORH"] = NewCounter(hw, TORH)
	g.Ext["TOTH"] = NewCounter(hw, TOTH)

	g.Ext["TPR"] = NewCounter(hw, TPR)
	g.Ext["TPT"] = NewCounter(hw, TPT)
	g.Ext["PTC64"] = NewCounter(hw, PTC64)
	g.Ext["PTC127"] = NewCounter(hw, PTC127)
	g.Ext["PTC255"] = NewCounter(hw, PTC255)
	g.Ext["PTC511"] = NewCounter(hw, PTC511)
	g.Ext["PTC1023"] = NewCounter(hw, PTC1023)
	g.Ext["PTC1522"] = NewCounter(hw, PTC1522)
	g.Ext["MPTC"] = NewCounter(hw, MPTC)
	g.Ext["BPTC"] = NewCounter(hw, BPTC)

	if hw.MAC.Type >= MACType82571 {
		g.Ext["IAC"] = NewCounter(hw, IAC)
		g.Ext["ICRXPTC"] = NewCounter(hw, ICRXPTC)
		g.Ext["ICRXATC"] = NewCounter(hw, ICRXATC)
		g.Ext["ICTXPTC"] = NewCounter(hw, ICTXPTC)
		g.Ext["ICTXATC"] = NewCounter(hw, ICTXATC)
		g.Ext["ICTXQEC"] = NewCounter(hw, ICTXQEC)
		g.Ext["ICTXQMTC"] = NewCounter(hw, ICTXQMTC)
		g.Ext["ICRXDMTC"] = NewCounter(hw, ICRXDMTC)
		g.Ext["ICRXOC"] = NewCounter(hw, ICRXOC)
	}

	if hw.MAC.Type >= MACType82543 {
		g.Ext["ALGNERRC"] = NewCounter(hw, ALGNERRC)
		rxerrors = append(rxerrors, g.Ext["ALGNERRC"])
		g.Ext["RXERRC"] = NewCounter(hw, RXERRC)
		rxerrors = append(rxerrors, g.Ext["RXERRC"])
		g.Ext["TNCRS"] = NewCounter(hw, TNCRS)
		g.Ext["CEXTERR"] = NewCounter(hw, CEXTERR)
		rxerrors = append(rxerrors, g.Ext["CEXTERR"])
		g.Ext["TSCTC"] = NewCounter(hw, TSCTC)
		g.Ext["TSCTFC"] = NewCounter(hw, TSCTFC)
	}

	g.RxPackets = g.Ext["GPRC"]
	g.TxPackets = g.Ext["GPTC"]
	g.RxOctets = g.Ext["GORC"]
	g.TxOctets = g.Ext["GOTC"]
	g.RxMissed = g.Ext["MPC"]
	g.RxErrors = ethdev.NewSumCounter(rxerrors...)
	g.TxErrors = ethdev.NewSumCounter(g.Ext["ECOL"], g.Ext["LATECOL"])

	return g
}
