package em

import (
	"uiosample/pci"
)

const (
	// Device Control - RW
	CTRL = 0x0000

	// Device Control Duplicate (Shadow) - RW
	CTRL_DUP = 0x0004

	// Device Status - RO
	STATUS = 0x0008

	// EEPROM/Flash Control - RW
	EECD = 0x0010

	// Extended Device Control - RW
	CTRL_EXT = 0x0018

	// MDI Control - RW
	MDIC = 0x0020

	// Flow Control Address Low - RW
	FCAL = 0x0028

	// Flow Control Address High - RW
	FCAH = 0x002c

	// Flow Control Type - RW
	FCT = 0x0030

	// VLAN Ether Type - RW
	VET = 0x0038

	// Interrupt Cause Read - R/clr
	ICR = 0x00c0

	// Interrupt Throttling Rate - RW
	ITR = 0x00c4

	// Interrupt Mask Set - RW
	IMS = 0x00d0

	// Interrupt Mask Clear - WO
	IMC = 0x00d8

	// Rx Control - RW
	RCTL = 0x0100

	// Flow Control Transmit Timer Value - RW
	FCTTV = 0x0170

	// Tx Configuration Word - RW
	TXCW = 0x0178

	// Rx Configuration Word - RO
	RXCW = 0x0180

	// Transmit Control
	TCTL = 0x0400

	// LED Control - RW
	LEDCTL = 0x0e00

	// TX corrupted data
	IOSFPC = 0x0f28

	// Packet Buffer Allocation - RW
	PBA = 0x1000

	// Flow Control Receive Threshold Low - RW
	FCRTL = 0x2160

	// Flow Control Receive Threshold High - RW
	FCRTH = 0x2168

	// Receive Descriptor Base Address
	// RDBAL0 = 0x2800
	// RDBAH0 = 0x2804

	// Receive Descriptor Length
	// RDLEN0 = 0x2808
	// RDH0 = 0x2810
	// RDT0 = 0x2818

	// Rx Delay Timer - RW
	RDTR = 0x2820

	// Receive Descriptor Control
	// RXDCTL0 = 0x2828

	// Transmit Descriptor Base Address
	// TDBAL0 = 0x3800
	// TDBAH0 = 0x3804

	// Transmit Descriptor Length
	// TDLEN0 = 0x3808
	// TDH0 = 0x3810
	// TDT0 = 0x3818

	// Transmit Interrupt Delay Value
	TIDV = 0x3820

	// TXDCTL0 = 0x3828

	// TARC0 = 0x3840

	// some statistics register

	// CRC Error Count - R/clr
	CRCERRS = 0x4000

	// Alignment Error Count - R/clr
	ALGNERRC = 0x4004

	// Symbol Error Count - R/clr
	SYMERRS = 0x4008

	// Receive Error Count - R/clr
	RXERRC = 0x400c

	// Missed Packet Count - R/clr
	MPC = 0x4010

	// Single Collision Count - R/clr
	SCC = 0x4014

	// Excessive Collision Count - R/clr
	ECOL = 0x4018

	// Multiple Collision Count - R/clr
	MCC = 0x401c

	// Late Collision Count - R/clr
	LATECOL = 0x4020

	// Collision Count - R/clr
	COLC = 0x4028

	// Defer Count - R/clr
	DC = 0x4030

	// Tx-No CRS - R/clr
	TNCRS = 0x4034

	// Sequence Error Count - R/clr
	SEC = 0x4038

	// Carrier Extension Error Count - R/clr
	CEXTERR = 0x403c

	// Receive Length Error Count - R/clr
	RLEC = 0x4040

	// XON Rx Count - R/clr
	XONRXC = 0x4048

	// XON Tx Count - R/clr
	XONTXC = 0x404C

	// XOFF Rx Count - R/clr
	XOFFRXC = 0x4050

	// XOFF Tx Count - R/clr
	XOFFTXC = 0x4054

	// Flow Control Rx Unsupported Count- R/clr
	FCRUC = 0x4058

	// Packets Rx (64 bytes) - R/clr
	PRC64 = 0x405c

	// Packets Rx (65-127 bytes) - R/clr
	PRC127 = 0x4060

	// Packets Rx (128-255 bytes) - R/clr
	PRC255 = 0x4064

	// Packets Rx (255-511 bytes) - R/clr
	PRC511 = 0x4068

	// Packets Rx (512-1023 bytes) - R/clr
	PRC1023 = 0x406c

	// Packets Rx (1024-1522 bytes) - R/clr
	PRC1522 = 0x4070

	// Good Packets Rx Count - R/clr
	GPRC = 0x4074

	// Broadcast Packets Rx Count - R/clr
	BPRC = 0x4078

	// Multicast Packets Rx Count - R/clr
	MPRC = 0x407c

	// Good Packets Tx Count - R/clr
	GPTC = 0x4080

	// Good Octets Rx Count Low - R/clr
	GORCL = 0x4088

	// Good Octets Rx Count High - R/clr
	GORCH = 0x408c

	// Good Octets Tx Count Low - R/clr
	GOTCL = 0x4090

	// Good Octets Tx Count High - R/clr
	GOTCH = 0x4094

	// Rx No Buffers Count - R/clr
	RNBC = 0x40a0

	// Rx Undersize Count - R/clr
	RUC = 0x40a4

	// Rx Fragment Count - R/clr
	RFC = 0x40a8

	// Rx Oversize Count - R/clr
	ROC = 0x40ac

	// Rx Jabber Count - R/clr
	RJC = 0x40b0

	// Management Packets Rx Count - R/clr
	MGTPRC = 0x40b4

	// Management Packets Dropped Count - R/clr
	MGTPDC = 0x40b8

	// Management Packets Tx Count - R/clr
	MGTPTC = 0x40bc

	// Total Octets Rx Low - R/clr
	TORL = 0x40c0

	// Total Octets Rx High - R/clr
	TORH = 0x40c4

	// Total Octets Tx Low - R/clr
	TOTL = 0x40c8

	// Total Octets Tx High - R/clr
	TOTH = 0x40cc

	// Total Packets Rx - R/clr
	TPR = 0x40d0

	// Total Packets Tx - R/clr
	TPT = 0x40d4

	// Packets Tx (64 bytes) - R/clr
	PTC64 = 0x40d8

	// Packets Tx (65-127 bytes) - R/clr
	PTC127 = 0x40dc

	// Packets Tx (128-255 bytes) - R/clr
	PTC255 = 0x40e0

	// Packets Tx (256-511 bytes) - R/clr
	PTC511 = 0x40e4

	// Packets Tx (512-1023 bytes) - R/clr
	PTC1023 = 0x40e8

	// Packets Tx (1024-1522 Bytes) - R/clr
	PTC1522 = 0x40ec

	// Multicast Packets Tx Count - R/clr
	MPTC = 0x040f0

	// Broadcast Packets Tx Count - R/clr
	BPTC = 0x040f4

	// TCP Segmentation Context Tx - R/clr
	TSCTC = 0x40f8

	// TCP Segmentation Context Tx Fail - R/clr
	TSCTFC = 0x40fc

	// PCS Configuration 0 - RW
	PCS_CFG0 = 0x4200

	// PCS Link Control - RW
	PCS_LCTL = 0x4208

	// PCS Link Status - RO
	PCS_LSTAT = 0x420c

	// AN advertisement - RW
	PCS_ANADV = 0x4218

	// Link Partner Ability - RW
	PCS_LPAB = 0x0421c

	// Rx Checksum Control - RW
	RXCSUM = 0x05000

	// Receive Filter Control
	RFCTL = 0x05008

	// Multicast Table Array - RW Array
	MTA = 0x5200

	// Receive Address (MAC address)
	RAL0 = 0x5400
	RAH0 = 0x5404

	// VLAN Filter Table Array - RW Array
	VFTA = 0x5600

	// Wakeup Control - RW
	WUC = 0x5800

	// Management Control - RW
	MANC = 0x5820

	// Management Control To Host - RW
	MANC2H = 0x5860

	// 3GIO Control Register
	GCR = 0x5b00

	// Function Active and Power State to MNG
	FACTPS = 0x5b30

	// SW Semaphore
	SWSM = 0x5b50

	// FW Semaphore
	FWSM = 0x5b54
)

// CTRL
const (
	CTRL_FD                 uint32 = 0x00000001 // Full duplex.0=half; 1=full
	CTRL_PRIOR              uint32 = 0x00000004 // Priority on PCI. 0=rx,1=fair
	CTRL_GIO_MASTER_DISABLE uint32 = 0x00000004 //Blocks new Master reqs
	CTRL_LRST               uint32 = 0x00000008 // Link reset. 0=normal,1=reset
	CTRL_ASDE               uint32 = 0x00000020 // Auto-speed detect enable
	CTRL_SLU                uint32 = 0x00000040 // Set link up (Force Link)
	CTRL_ILOS               uint32 = 0x00000080 // Invert Loss-Of Signal
	CTRL_SPD_SEL            uint32 = 0x00000300 // Speed Select Mask
	CTRL_SPD_10             uint32 = 0x00000000 // Force 10Mb
	CTRL_SPD_100            uint32 = 0x00000100 // Force 100Mb
	CTRL_SPD_1000           uint32 = 0x00000200 // Force 1Gb
	CTRL_FRCSPD             uint32 = 0x00000800 // Force Speed
	CTRL_FRCDPX             uint32 = 0x00001000 // Force Duplex
	CTRL_LANPHYPC_OVERRIDE  uint32 = 0x00010000 // SW control of LANPHYPC
	CTRL_LANPHYPC_VALUE     uint32 = 0x00020000 // SW value of LANPHYPC
	CTRL_MEHE               uint32 = 0x00080000 // Memory Error Handling Enable
	CTRL_SWDPIN0            uint32 = 0x00040000 // SWDPIN 0 value
	CTRL_SWDPIN1            uint32 = 0x00080000 // SWDPIN 1 value
	CTRL_SWDPIN2            uint32 = 0x00100000 // SWDPIN 2 value
	CTRL_ADVD3WUC           uint32 = 0x00100000 // D3 WUC
	CTRL_EN_PHY_PWR_MGMT    uint32 = 0x00200000 // PHY PM enable
	CTRL_SWDPIN3            uint32 = 0x00200000 // SWDPIN 3 value
	CTRL_SWDPIO0            uint32 = 0x00400000 // SWDPIN 0 Input or output
	CTRL_SWDPIO2            uint32 = 0x01000000 // SWDPIN 2 input or output
	CTRL_SWDPIO3            uint32 = 0x02000000 // SWDPIN 3 input or output
	CTRL_DEV_RST            uint32 = 0x20000000 // Device reset
	CTRL_RST                uint32 = 0x04000000 // Global reset
	CTRL_RFCE               uint32 = 0x08000000 // Receive Flow Control enable
	CTRL_TFCE               uint32 = 0x10000000 // Transmit flow control enable
	CTRL_VME                uint32 = 0x40000000 // IEEE VLAN mode enable
	CTRL_PHY_RST            uint32 = 0x80000000 // PHY Reset
	CTRL_I2C_ENA            uint32 = 0x02000000 // I2C enable

	CTRL_MDIO_DIR = CTRL_SWDPIO2
	CTRL_MDIO     = CTRL_SWDPIN2
	CTRL_MDC_DIR  = CTRL_SWDPIO3
	CTRL_MDC      = CTRL_SWDPIN3
)

// STATUS
const (
	STATUS_FD                uint32 = 0x00000001 // Duplex 0=half 1=full
	STATUS_LU                uint32 = 0x00000002 // Link up.0=no,1=link
	STATUS_FUNC_MASK         uint32 = 0x0000000C // PCI Function Mask
	STATUS_FUNC_SHIFT        uint32 = 2
	STATUS_FUNC_1            uint32 = 0x00000004 // Function 1
	STATUS_TXOFF             uint32 = 0x00000010 // transmission paused
	STATUS_SPEED_MASK        uint32 = 0x000000C0
	STATUS_SPEED_10          uint32 = 0x00000000 // Speed 10Mb/s
	STATUS_SPEED_100         uint32 = 0x00000040 // Speed 100Mb/s
	STATUS_SPEED_1000        uint32 = 0x00000080 // Speed 1000Mb/s
	STATUS_LAN_INIT_DONE     uint32 = 0x00000200 // Lan Init Compltn by NVM
	STATUS_PHYRA             uint32 = 0x00000400 // PHY Reset Asserted
	STATUS_GIO_MASTER_ENABLE uint32 = 0x00080000 // Master request status
	STATUS_PCI66             uint32 = 0x00000800 // In 66Mhz slot
	STATUS_BUS64             uint32 = 0x00001000 // In 64 bit slot
	STATUS_2P5_SKU           uint32 = 0x00001000 // Val of 2.5GBE SKU strap
	STATUS_2P5_SKU_OVER      uint32 = 0x00002000 // Val of 2.5GBE SKU Over
	STATUS_PCIX_MODE         uint32 = 0x00002000 // PCI-X mode
	STATUS_PCIX_SPEED        uint32 = 0x0000C000 // PCI-X bus speed

	// Constants used to interpret the masked PCI-X bus speed.
	STATUS_PCIX_SPEED_66  uint32 = 0x00000000 // PCI-X bus spd 50-66MHz
	STATUS_PCIX_SPEED_100 uint32 = 0x00004000 // PCI-X bus spd 66-100MHz
	STATUS_PCIX_SPEED_133 uint32 = 0x00008000 // PCI-X bus spd 100-133MHz
	STATUS_PCIM_STATE     uint32 = 0x40000000 // PCIm function state
)

// EECD
const (
	EECD_SK        uint32 = 0x00000001 // NVM Clock
	EECD_CS        uint32 = 0x00000002 // NVM Chip Select
	EECD_DI        uint32 = 0x00000004 // NVM Data In
	EECD_DO        uint32 = 0x00000008 // NVM Data Out
	EECD_REQ       uint32 = 0x00000040 // NVM Access Request
	EECD_GNT       uint32 = 0x00000080 // NVM Access Grant
	EECD_PRES      uint32 = 0x00000100 // NVM Present
	EECD_SIZE      uint32 = 0x00000200 // NVM Size (0=64 word 1=256 word)
	EECD_BLOCKED   uint32 = 0x00008000 // Bit banging access blocked flag
	EECD_ABORT     uint32 = 0x00010000 // NVM operation aborted flag
	EECD_TIMEOUT   uint32 = 0x00020000 // NVM read operation timeout flag
	EECD_ERROR_CLR uint32 = 0x00040000 // NVM error status clear bit

	// NVM Addressing bits based on type 0=small, 1=large
	EECD_ADDR_BITS uint32 = 0x00000400
	EECD_TYPE      uint32 = 0x00002000 // NVM Type (1-SPI, 0-Microwire)
)

// CTRL_EXT
const (
	CTRL_EXT_LPCD      uint32 = 0x00000004 // LCD Power Cycle Done
	CTRL_EXT_SDP4_DATA uint32 = 0x00000010 // SW Definable Pin 4 data
	CTRL_EXT_SDP6_DATA uint32 = 0x00000040 // SW Definable Pin 6 data
	CTRL_EXT_SDP3_DATA uint32 = 0x00000080 // SW Definable Pin 3 data

	// SDP 4/5 (bits 8,9) are reserved in >= 82575

	CTRL_EXT_SDP4_DIR       uint32 = 0x00000100 // Direction of SDP4 0=in 1=out
	CTRL_EXT_SDP6_DIR       uint32 = 0x00000400 // Direction of SDP6 0=in 1=out
	CTRL_EXT_SDP3_DIR       uint32 = 0x00000800 // Direction of SDP3 0=in 1=out
	CTRL_EXT_FORCE_SMBUS    uint32 = 0x00000800 // Force SMBus mode
	CTRL_EXT_EE_RST         uint32 = 0x00002000 // Reinitialize from EEPROM
	CTRL_EXT_PFRSTD         uint32 = 0x00004000 // Physical Func Reset Done Indication
	CTRL_EXT_SDLPE          uint32 = 0x00040000 // SerDes Low Power Enable
	CTRL_EXT_SPD_BYPS       uint32 = 0x00008000 // Speed Select Bypass
	CTRL_EXT_RO_DIS         uint32 = 0x00020000 // Relaxed Ordering disable
	CTRL_EXT_DMA_DYN_CLK_EN uint32 = 0x00080000 // DMA Dynamic Clk Gating
	CTRL_EXT_LINK_MODE_MASK uint32 = 0x00c00000

	// Offset of the link mode field in Ctrl Ext register
	CTRL_EXT_LINK_MODE_OFFSET = 22

	CTRL_EXT_LINK_MODE_1000BASE_KX uint32 = 0x00400000
	CTRL_EXT_LINK_MODE_GMII        uint32 = 0x00000000
	CTRL_EXT_LINK_MODE_PCIE_SERDES uint32 = 0x00C00000
	CTRL_EXT_LINK_MODE_SGMII       uint32 = 0x00800000
	CTRL_EXT_EIAME                 uint32 = 0x01000000
	CTRL_EXT_IRCA                  uint32 = 0x00000001
	CTRL_EXT_DRV_LOAD              uint32 = 0x10000000 // Drv loaded bit for FW
	CTRL_EXT_IAME                  uint32 = 0x08000000 // Int ACK Auto-mask
	CTRL_EXT_PBA_CLR               uint32 = 0x80000000 // PBA Clear
	CTRL_EXT_LSECCK                uint32 = 0x00001000
	CTRL_EXT_PHYPDEN               uint32 = 0x00100000
)

// MDIC
const (
	MDIC_REG_MASK  uint32 = 0x001F0000
	MDIC_REG_SHIFT        = 16
	MDIC_PHY_MASK  uint32 = 0x03E00000
	MDIC_PHY_SHIFT        = 21
	MDIC_OP_WRITE  uint32 = 0x04000000
	MDIC_OP_READ   uint32 = 0x08000000
	MDIC_READY     uint32 = 0x10000000
	MDIC_ERROR     uint32 = 0x40000000
	MDIC_DEST      uint32 = 0x80000000

	VFTA_BLOCK_SIZE = 8
)

// SerDes Control
const (
	GEN_CTL_READY uint32 = 0x80000000

	GEN_CTL_ADDRESS_SHIFT = 8
	GEN_POLL_TIMEOUT      = 640
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
	RCTL_RST        uint32 = 0x00000001 // Software reset
	RCTL_EN         uint32 = 0x00000002 // enable
	RCTL_SBP        uint32 = 0x00000004 // store bad packet
	RCTL_UPE        uint32 = 0x00000008 // unicast promisc enable
	RCTL_MPE        uint32 = 0x00000010 // multicast promisc enable
	RCTL_LPE        uint32 = 0x00000020 // long packet enable
	RCTL_LBM_NO     uint32 = 0x00000000 // no loopback mode
	RCTL_LBM_MAC    uint32 = 0x00000040 // MAC loopback mode
	RCTL_LBM_TCVR   uint32 = 0x000000C0 // tcvr loopback mode
	RCTL_DTYP_PS    uint32 = 0x00000400 // Packet Split descriptor
	RCTL_RDMTS_HALF uint32 = 0x00000000 // Rx desc min thresh size
	RCTL_RDMTS_HEX  uint32 = 0x00010000

	RCTL_DTYP_MASK uint32 = 0x00000c00

	RCTL_RDMTS1_HEX = RCTL_RDMTS_HEX

	RCTL_MO_SHIFT = 12 // multicast offset shift

	RCTL_MO_3 uint32 = 0x00003000 // multicast offset 15:4
	RCTL_BAM  uint32 = 0x00008000 // broadcast enable

	// these buffer sizes are valid if E1000_RCTL_BSEX is 0
	RCTL_SZ_2048 uint32 = 0x00000000 // Rx buffer size 2048
	RCTL_SZ_1024 uint32 = 0x00010000 // Rx buffer size 1024
	RCTL_SZ_512  uint32 = 0x00020000 // Rx buffer size 512
	RCTL_SZ_256  uint32 = 0x00030000 // Rx buffer size 256

	// these buffer sizes are valid if E1000_RCTL_BSEX is 1
	RCTL_SZ_16384 uint32 = 0x00010000 // Rx buffer size 16384
	RCTL_SZ_8192  uint32 = 0x00020000 // Rx buffer size 8192
	RCTL_SZ_4096  uint32 = 0x00030000 // Rx buffer size 4096
	RCTL_VFE      uint32 = 0x00040000 // vlan filter enable
	RCTL_CFIEN    uint32 = 0x00080000 // canonical form enable
	RCTL_CFI      uint32 = 0x00100000 // canonical form indicator
	RCTL_DPF      uint32 = 0x00400000 // discard pause frames
	RCTL_PMCF     uint32 = 0x00800000 // pass MAC control frames
	RCTL_BSEX     uint32 = 0x02000000 // Buffer size extension
	RCTL_SECRC    uint32 = 0x04000000 // Strip Ethernet CRC
)

// TXCW
const (
	TXCW_FD         uint32 = 0x00000020 // TXCW full duplex
	TXCW_PAUSE      uint32 = 0x00000080 // TXCW sym pause request
	TXCW_ASM_DIR    uint32 = 0x00000100 // TXCW astm pause direction
	TXCW_PAUSE_MASK uint32 = 0x00000180 // TXCW pause request mask
	TXCW_ANE        uint32 = 0x80000000 // Auto-neg enable
)

// RXCW
const (
	RXCW_CW    uint32 = 0x0000ffff // RxConfigWord mask
	RXCW_IV    uint32 = 0x08000000 // Receive config invalid
	RXCW_C     uint32 = 0x20000000 // Receive config
	RXCW_SYNCH uint32 = 0x40000000 // Receive config synch
)

// TCTL
const (
	TCTL_EN   uint32 = 0x00000002 // enable Tx
	TCTL_PSP  uint32 = 0x00000008 // pad short packets
	TCTL_CT   uint32 = 0x00000ff0 // collision threshold
	TCTL_COLD uint32 = 0x003ff000 // collision distance
	TCTL_RTLC uint32 = 0x01000000 // Re-transmit on late collision
	TCTL_MULR uint32 = 0x10000000 // Multiple request support
)

// Collision related configuration parameters
const (
	CT_SHIFT            = 4
	COLLISION_THRESHOLD = 15
	COLLISION_DISTANCE  = 63
	COLD_SHIFT          = 12
)

// LEDCTL
const (
	LEDCTL_LED0_MODE_MASK uint32 = 0xf

	LEDCTL_LED0_MODE_SHIFT = 0

	LEDCTL_LED0_IVRT  uint32 = 0x00000040
	LEDCTL_LED0_BLINK uint32 = 0x00000080

	LEDCTL_MODE_LINK_UP uint32 = 0x2
	LEDCTL_MODE_LED_ON  uint32 = 0xe
	LEDCTL_MODE_LED_OFF uint32 = 0xf
)

// Flow Control
const (
	FCRTH_RTH  uint32 = 0x0000FFF8 // Mask Bits[15:3] for RTH
	FCRTL_RTL  uint32 = 0x0000FFF8 // Mask Bits[15:3] for RTL
	FCRTL_XONE uint32 = 0x80000000 // Enable XON frame transmission
)

func RDBAL(n int) int {
	if n < 4 {
		return 0x2800 + (n * 0x100)
	} else {
		return 0xc000 + (n * 0x40)
	}
}

func RDBAH(n int) int {
	if n < 4 {
		return 0x2804 + (n * 0x100)
	} else {
		return 0xc004 + (n * 0x40)
	}
}

func RDLEN(n int) int {
	if n < 4 {
		return 0x2808 + (n * 0x100)
	} else {
		return 0xc008 + (n * 0x40)
	}
}

func SRRCTL(n int) int {
	if n < 4 {
		return 0x280c + (n * 0x100)
	} else {
		return 0xc00c + (n * 0x40)
	}
}

func RDH(n int) int {
	if n < 4 {
		return 0x2810 + (n * 0x100)
	} else {
		return 0xc010 + (n * 0x40)
	}
}

func RDCTL(n int) int {
	if n < 4 {
		return 0x2814 + (n * 0x100)
	} else {
		return 0xc014 + (n * 0x40)
	}
}

func RDT(n int) int {
	if n < 4 {
		return 0x2818 + (n * 0x100)
	} else {
		return 0xc018 + (n * 0x40)
	}
}

func RXDCTL(n int) int {
	if n < 4 {
		return 0x2828 + (n * 0x100)
	} else {
		return 0xc028 + (n * 0x40)
	}
}

// RXDCTL
const (
	RXDCTL_GRAN uint32 = 0x01000000 // RXDCTL Granularity
)

func RQDPC(n int) int {
	if n < 4 {
		return 0x2830 + (n * 0x100)
	} else {
		return 0xc030 + (n * 0x40)
	}
}

func TDBAL(n int) int {
	if n < 4 {
		return 0x3800 + (n * 0x100)
	} else {
		return 0xe000 + (n * 0x40)
	}
}

func TDBAH(n int) int {
	if n < 4 {
		return 0x3804 + (n * 0x100)
	} else {
		return 0xe004 + (n * 0x40)
	}
}

func TDLEN(n int) int {
	if n < 4 {
		return 0x3808 + (n * 0x100)
	} else {
		return 0xe008 + (n * 0x40)
	}
}

func TDH(n int) int {
	if n < 4 {
		return 0x3810 + (n * 0x100)
	} else {
		return 0xe010 + (n * 0x40)
	}
}

func TXCTL(n int) int {
	if n < 4 {
		return 0x3814 + (n * 0x100)
	} else {
		return 0xe014 + (n * 0x40)
	}
}

func TDT(n int) int {
	if n < 4 {
		return 0x3818 + (n * 0x100)
	} else {
		return 0xe018 + (n * 0x40)
	}
}

func TXDCTL(n int) int {
	if n < 4 {
		return 0x03828 + (n * 0x100)
	} else {
		return 0x0e028 + (n * 0x40)
	}
}

// TXDCTL
const (
	TXDCTL_PTHRESH      uint32 = uint32(0x0000003f)
	TXDCTL_HTHRESH      uint32 = uint32(0x00003f00)
	TXDCTL_WTHRESH      uint32 = uint32(0x003f0000)
	TXDCTL_GRAN         uint32 = uint32(0x01000000)
	TXDCTL_QUEUE_ENABLE uint32 = uint32(0x02000000)

	TXDCTL_FULL_TX_DESC_WB      uint32 = uint32(0x01010000)
	TXDCTL_MAX_TX_DESC_PREFETCH uint32 = uint32(0x0100001f)
	TXDCTL_COUNT_DESC           uint32 = uint32(0x00400000)
)

func TDWBAL(n int) int {
	if n < 4 {
		return 0x03838 + (n * 0x100)
	} else {
		return 0x0e038 + (n * 0x40)
	}
}

func TDWBAH(n int) int {
	if n < 4 {
		return 0x0383c + (n * 0x100)
	} else {
		return 0x0e03c + (n * 0x40)
	}
}

func TARC(n int) int {
	return 0x03840 + (n * 0x100)
}

// TARC
const (
	TARC0_ENABLE          uint32 = 0x00000400 // Enable Tx Queue 0
	TARC0_CB_MULTIQ_3_REQ uint32 = 0x30000000
	TARC0_CB_MULTIQ_2_REQ uint32 = 0x20000000
)

func RAL(n int) int {
	if n <= 15 {
		return 0x05400 + (n * 8)
	} else {
		return 0x054e0 + (n-16)*8
	}
}

func RAH(n int) int {
	if n <= 15 {
		return 0x05404 + (n * 8)
	} else {
		return 0x054e4 + (n-16)*8
	}
}

// PCS_CFG
const (
	PCS_CFG_PCS_EN = 8
)

// PCS_LCTL
const (
	PCS_LCTL_FLV_LINK_UP = 1
	PCS_LCTL_FSV_10      = 0
	PCS_LCTL_FSV_100     = 2
	PCS_LCTL_FSV_1000    = 4
	PCS_LCTL_FDV_FULL    = 8
	PCS_LCTL_FSD         = 0x10
	PCS_LCTL_FORCE_LINK  = 0x20
	PCS_LCTL_FORCE_FCTRL = 0x80
	PCS_LCTL_AN_ENABLE   = 0x10000
	PCS_LCTL_AN_RESTART  = 0x20000
	PCS_LCTL_AN_TIMEOUT  = 0x40000
)

// PCS_LSTAT
const (
	PCS_LSTS_LINK_OK     = 1
	PCS_LSTS_SPEED_100   = 2
	PCS_LSTS_SPEED_1000  = 4
	PCS_LSTS_DUPLEX_FULL = 8
	PCS_LSTS_SYNK_OK     = 0x10
	PCS_LSTS_AN_COMPLETE = 0x10000
)

// RXCSUM
const (
	RXCSUM_IPOFL  uint32 = 0x00000100 // IPv4 checksum offload
	RXCSUM_TUOFL  uint32 = 0x00000200 // TCP / UDP checksum offload
	RXCSUM_CRCOFL uint32 = 0x00000800 // CRC32 offload enable
	RXCSUM_IPPCSE uint32 = 0x00001000 // IP payload checksum enable
	RXCSUM_PCSD   uint32 = 0x00002000 // packet checksum disabled
)

// RFCTL
const (
	RFCTL_NFSW_DIS         uint32 = 0x00000040
	RFCTL_NFSR_DIS         uint32 = 0x00000080
	RFCTL_ACK_DIS          uint32 = 0x00001000
	RFCTL_EXTEN            uint32 = 0x00008000
	RFCTL_IPV6_EX_DIS      uint32 = 0x00010000
	RFCTL_NEW_IPV6_EXT_DIS uint32 = 0x00020000
	RFCTL_LEF              uint32 = 0x00040000
)

// MANC
const (
	MANC_SMBUS_EN           uint32 = 0x00000001 // SMBus Enabled - RO
	MANC_ASF_EN             uint32 = 0x00000002 // ASF Enabled - RO
	MANC_ARP_EN             uint32 = 0x00002000 // Enable ARP Request Filtering
	MANC_RCV_TCO_EN         uint32 = 0x00020000 // Receive TCO Packets Enabled
	MANC_BLK_PHY_RST_ON_IDE uint32 = 0x00040000 // Block phy resets
	MANC_EN_MAC_ADDR_FILTER uint32 = 0x00100000 // Enable MAC address filtering
	MANC_EN_MNG2HOST        uint32 = 0x00200000 // Enable MNG packets to host memory
)

// MANC2H
const (
	MANC2H_PORT_623 uint32 = 0x00000020 // Port 0x26f
	MANC2H_PORT_664 uint32 = 0x00000040 // Port 0x298
)

// FACTPS
const (
	FACTPS_MNGCG uint32 = 0x20000000
)

// SWSM
const (
	SWSM_SMBI     uint32 = 0x00000001 // Driver Semaphore bit
	SWSM_SWESMBI  uint32 = 0x00000002 // FW Semaphore bit
	SWSM_DRV_LOAD uint32 = 0x00000008 // Driver Loaded Bit
)

// FWSM
const (
	FWSM_MODE_MASK uint32 = 0xe

	FWSM_MODE_SHIFT = 1

	FWSM_FW_VALID        uint32 = 0x00008000
	FWSM_HI_EN_ONLY_MODE uint32 = 0x4
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

func (r Reg) WriteFlush() {
	r.Read(STATUS)
}
