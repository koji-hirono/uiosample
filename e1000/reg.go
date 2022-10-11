package e1000

import (
	"uiosample/pci"
)

const (
	CTRL = 0x0000

	STATUS = 0x0008

	// Flow Control Address
	FCAL = 0x0028
	FCAH = 0x002c

	// Flow Control Type
	FCT = 0x0030

	// Interrupt Cause Read Register
	ICR = 0x00c0

	// Interrupt Mask Set/Read Register
	IMS = 0x00d0

	// Interrupt Mask Clear Register
	IMC = 0x00d8

	// Receive control
	RCTL = 0x0100

	// Flow Control Transmit Timer Value
	FCTTV = 0x0170

	// Transmit Control
	TCTL = 0x0400

	// Receive Descriptor Base Address
	RDBAL = 0x2800
	RDBAH = 0x2804

	// Receive Descriptor Length
	RDLEN = 0x2808
	RDH   = 0x2810
	RDT   = 0x2818

	// Receive Descriptor Control
	RXDCTL = 0x2828

	// Transmit Descriptor Base Address
	TDBAL = 0x3800
	TDBAH = 0x3804

	// Transmit Descriptor Length
	TDLEN = 0x3808
	TDH   = 0x3810
	TDT   = 0x3818

	// Transmit Interrupt Delay Value
	TIDV = 0x3820

	// some statistics register

	RXERRC = 0x400c

	// Missed Packets Count
	MPC = 0x4010

	// Good Packets Received Counts
	GPRC = 0x4074

	// Good Packets Transmitted Count
	GPTC = 0x4080

	// Good Octets Received Count
	GORCL = 0x4088
	GORCH = 0x408c

	// Good Octets Transmitted Count
	GOTCL = 0x4090
	GOTCH = 0x4094

	// Receive Address (MAC address)
	RAL0 = 0x5400
	RAH0 = 0x5404

	// 3GIO Control Register
	GCR = 0x5b00
)

// CTRL
const (
	CTRL_FD      uint32 = uint32(1) << 0
	CTRL_LRST    uint32 = uint32(1) << 3 // reserved
	CTRL_ASDE    uint32 = uint32(1) << 5
	CTRL_SLU     uint32 = uint32(1) << 6
	CTRL_ILOS    uint32 = uint32(1) << 7 // reserved
	CTRL_RST     uint32 = uint32(1) << 26
	CTRL_VME     uint32 = uint32(1) << 30
	CTRL_PHY_RST uint32 = uint32(1) << 31
)

// IMS
const (
	IMS_TXDW  uint32 = uint32(1) << 0
	IMS_TXQE  uint32 = uint32(1) << 1
	IMS_LSC   uint32 = uint32(1) << 2
	IMS_RXSEQ uint32 = uint32(1) << 3
	IMS_RXDMT uint32 = uint32(1) << 4
	IMS_RXO   uint32 = uint32(1) << 6
	IMS_RXT   uint32 = uint32(1) << 7
	IMS_RXQ0  uint32 = uint32(1) << 20
	IMS_RXQ1  uint32 = uint32(1) << 21
	IMS_TXQ0  uint32 = uint32(1) << 22
	IMS_TXQ1  uint32 = uint32(1) << 23
	IMS_OTHER uint32 = uint32(1) << 24
)

// RCTL
const (
	RCTL_EN     uint32 = uint32(1) << 1
	RCTL_UPE    uint32 = uint32(1) << 3
	RCTL_MPE    uint32 = uint32(1) << 4
	RCTL_LPE    uint32 = uint32(1) << 5
	RCTL_LBM    uint32 = uint32(1)<<6 | uint32(1)<<7
	RCTL_BAM    uint32 = uint32(1) << 15
	RCTL_BSIZE1 uint32 = uint32(1) << 16
	RCTL_BSIZE2 uint32 = uint32(1) << 17
	RCTL_BSEX   uint32 = uint32(1) << 25
	RCTL_SECRC  uint32 = uint32(1) << 26
)

// TCTL
const (
	TCTL_EN  uint32 = uint32(1) << 1
	TCTL_PSP uint32 = uint32(1) << 3
)

type Reg struct {
	res pci.Resource
}

func (r Reg) Read(reg int) uint32 {
	return r.res.Read32(reg)
}

func (r Reg) Write(reg int, val uint32) {
	r.res.Write32(reg, val)
}

func (r Reg) MaskWrite(reg int, val, mask uint32) {
	r.res.MaskWrite32(reg, val, mask)
}
