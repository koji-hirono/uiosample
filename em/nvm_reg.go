package em

// NVM Word Offsets
const (
	NVM_COMPAT             uint16 = 0x0003
	NVM_ID_LED_SETTINGS           = 0x0004
	NVM_VERSION                   = 0x0005
	NVM_SERDES_AMPLITUDE          = 0x0006
	NVM_PHY_CLASS_WORD            = 0x0007
	I210_NVM_FW_MODULE_PTR        = 0x0010
	I350_NVM_FW_MODULE_PTR        = 0x0051
	NVM_FUTURE_INIT_WORD1         = 0x0019
	NVM_ETRACK_WORD               = 0x0042
	NVM_ETRACK_HIWORD             = 0x0043
	NVM_COMB_VER_OFF              = 0x0083
	NVM_COMB_VER_PTR              = 0x003d
)

const (
	NVM_MAC_ADDR    uint16 = 0x0000
	NVM_SUB_DEV_ID         = 0x000B
	NVM_SUB_VEN_ID         = 0x000C
	NVM_DEV_ID             = 0x000D
	NVM_VEN_ID             = 0x000E
	NVM_INIT_CTRL_2        = 0x000F
	NVM_INIT_CTRL_4        = 0x0013
	NVM_LED_1_CFG          = 0x001C
	NVM_LED_0_2_CFG        = 0x001F
)

const (
	NVM_INIT_CONTROL2_REG      uint16 = 0x000F
	NVM_INIT_CONTROL3_PORT_B          = 0x0014
	NVM_INIT_3GIO_3                   = 0x001A
	NVM_SWDEF_PINS_CTRL_PORT_0        = 0x0020
	NVM_INIT_CONTROL3_PORT_A          = 0x0024
	NVM_CFG                           = 0x0012
	NVM_ALT_MAC_ADDR_PTR              = 0x0037
	NVM_CHECKSUM_REG                  = 0x003F
	NVM_COMPATIBILITY_REG_3           = 0x0003
	NVM_COMPATIBILITY_BIT_MASK        = 0x8000
)

// For checksumming, the sum of all words in the NVM should equal 0xBABA.
const NVM_SUM = 0xBABA

// PBA (printed board assembly) number words
const (
	NVM_PBA_OFFSET_0          = 8
	NVM_PBA_OFFSET_1          = 9
	NVM_PBA_PTR_GUARD         = 0xFAFA
	NVM_RESERVED_WORD         = 0xFFFF
	NVM_PHY_CLASS_A           = 0x8000
	NVM_SERDES_AMPLITUDE_MASK = 0x000F
	NVM_SIZE_MASK             = 0x1C00
	NVM_SIZE_SHIFT            = 10
	NVM_WORD_SIZE_BASE_SHIFT  = 6
	NVM_SWDPIO_EXT_SHIFT      = 4
)

// NVM Commands - Microwire
const (
	NVM_READ_OPCODE_MICROWIRE  uint16 = 0x6  // NVM read opcode
	NVM_WRITE_OPCODE_MICROWIRE uint16 = 0x5  // NVM write opcode
	NVM_ERASE_OPCODE_MICROWIRE uint16 = 0x7  // NVM erase opcode
	NVM_EWEN_OPCODE_MICROWIRE  uint16 = 0x13 // NVM erase/write enable
	NVM_EWDS_OPCODE_MICROWIRE  uint16 = 0x10 // NVM erase/write disable
)

// NVM Commands - SPI
const (
	NVM_MAX_RETRY_SPI = 5000 // Max wait of 5ms, for RDY signal

	NVM_READ_OPCODE_SPI  uint16 = 0x03 // NVM read opcode
	NVM_WRITE_OPCODE_SPI uint16 = 0x02 // NVM write opcode
	NVM_A8_OPCODE_SPI    uint16 = 0x08 // opcode bit-3 = address bit-8
	NVM_WREN_OPCODE_SPI  uint16 = 0x06 // NVM set Write Enable latch
	NVM_RDSR_OPCODE_SPI  uint16 = 0x05 // NVM read Status register
)

// SPI NVM Status Register
const NVM_STATUS_RDY_SPI = 0x01

// Mask bits for fields in Word 0x0f of the NVM
const (
	NVM_WORD0F_PAUSE_MASK      = 0x3000
	NVM_WORD0F_PAUSE           = 0x1000
	NVM_WORD0F_ASM_DIR         = 0x2000
	NVM_WORD0F_SWPDIO_EXT_MASK = 0x00F0
)

func NVM_82580_LAN_FUNC_OFFSET(n uint16) uint16 {
	if n != 0 {
		return 0x40 + (0x40 * n)
	} else {
		return 0
	}
}

// Mask bits for fields in Word 0x24 of the NVM
const (
	NVM_WORD24_COM_MDIO uint16 = 0x0008 // MDIO interface shared
	NVM_WORD24_EXT_MDIO uint16 = 0x0004 // MDIO accesses routed extrnl

	// Offset of Link Mode bits for 82575/82576
	NVM_WORD24_LNK_MODE_OFFSET = 8
	// Offset of Link Mode bits for 82580 up
	NVM_WORD24_82580_LNK_MODE_OFFSET = 4
)
