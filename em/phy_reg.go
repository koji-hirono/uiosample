package em

// PHY Registers defined by IEEE
const (
	PHY_CONTROL      = 0x00 // Control Register
	PHY_STATUS       = 0x01 // Status Register
	PHY_ID1          = 0x02 // Phy Id Reg (word 1)
	PHY_ID2          = 0x03 // Phy Id Reg (word 2)
	PHY_AUTONEG_ADV  = 0x04 // Autoneg Advertisement
	PHY_LP_ABILITY   = 0x05 // Link Partner Ability (Base Page)
	PHY_AUTONEG_EXP  = 0x06 // Autoneg Expansion Reg
	PHY_NEXT_PAGE_TX = 0x07 // Next Page Tx
	PHY_LP_NEXT_PAGE = 0x08 // Link Partner Next Page
	PHY_1000T_CTRL   = 0x09 // 1000Base-T Control Reg
	PHY_1000T_STATUS = 0x0A // 1000Base-T Status Reg
	PHY_EXT_STATUS   = 0x0F // Extended Status Reg
)

const PHY_REVISION_MASK uint32 = 0xFFFFFFF0

const MAX_PHY_REG_ADDRESS = 0x1F // 5 bit address bus (0-0x1F)

// PHY Control Register
const (
	MII_CR_SPEED_SELECT_MSB = 0x0040 // bits 6,13: 10=1000, 01=100, 00=10
	MII_CR_COLL_TEST_ENABLE = 0x0080 // Collision test enable
	MII_CR_FULL_DUPLEX      = 0x0100 // FDX =1, half duplex =0
	MII_CR_RESTART_AUTO_NEG = 0x0200 // Restart auto negotiation
	MII_CR_ISOLATE          = 0x0400 // Isolate PHY from MII
	MII_CR_POWER_DOWN       = 0x0800 // Power down
	MII_CR_AUTO_NEG_EN      = 0x1000 // Auto Neg Enable
	MII_CR_SPEED_SELECT_LSB = 0x2000 // bits 6,13: 10=1000, 01=100, 00=10
	MII_CR_LOOPBACK         = 0x4000 // 0 = normal, 1 = loopback
	MII_CR_RESET            = 0x8000 // 0 = normal, 1 = PHY reset
	MII_CR_SPEED_1000       = 0x0040
	MII_CR_SPEED_100        = 0x2000
	MII_CR_SPEED_10         = 0x0000
)

// PHY Status Register
const (
	MII_SR_EXTENDED_CAPS     = 0x0001 // Extended register capabilities
	MII_SR_JABBER_DETECT     = 0x0002 // Jabber Detected
	MII_SR_LINK_STATUS       = 0x0004 // Link Status 1 = link
	MII_SR_AUTONEG_CAPS      = 0x0008 // Auto Neg Capable
	MII_SR_REMOTE_FAULT      = 0x0010 // Remote Fault Detect
	MII_SR_AUTONEG_COMPLETE  = 0x0020 // Auto Neg Complete
	MII_SR_PREAMBLE_SUPPRESS = 0x0040 // Preamble may be suppressed
	MII_SR_EXTENDED_STATUS   = 0x0100 // Ext. status info in Reg 0x0F
	MII_SR_100T2_HD_CAPS     = 0x0200 // 100T2 Half Duplex Capable
	MII_SR_100T2_FD_CAPS     = 0x0400 // 100T2 Full Duplex Capable
	MII_SR_10T_HD_CAPS       = 0x0800 // 10T   Half Duplex Capable
	MII_SR_10T_FD_CAPS       = 0x1000 // 10T   Full Duplex Capable
	MII_SR_100X_HD_CAPS      = 0x2000 // 100X  Half Duplex Capable
	MII_SR_100X_FD_CAPS      = 0x4000 // 100X  Full Duplex Capable
	MII_SR_100T4_CAPS        = 0x8000 // 100T4 Capable
)

// 1000BASE-T Status Register
const (
	SR_1000T_IDLE_ERROR_CNT   = 0x00FF // Num idle err since last rd
	SR_1000T_ASYM_PAUSE_DIR   = 0x0100 // LP asym pause direction bit
	SR_1000T_LP_HD_CAPS       = 0x0400 // LP is 1000T HD capable
	SR_1000T_LP_FD_CAPS       = 0x0800 // LP is 1000T FD capable
	SR_1000T_REMOTE_RX_STATUS = 0x1000 // Remote receiver OK
	SR_1000T_LOCAL_RX_STATUS  = 0x2000 // Local receiver OK
	SR_1000T_MS_CONFIG_RES    = 0x4000 // 1=Local Tx Master, 0=Slave
	SR_1000T_MS_CONFIG_FAULT  = 0x8000 // Master/Slave config fault

	SR_1000T_PHY_EXCESSIVE_IDLE_ERR_COUNT = 5
)

const (
	I82578_EPSCR_DOWNSHIFT_ENABLE       = 0x0020
	I82578_EPSCR_DOWNSHIFT_COUNTER_MASK = 0x001C
)

// BME1000 PHY Specific Control Register
const (
	BME1000_PSCR_ENABLE_DOWNSHIFT = 0x0800 // 1 = enable downshift
)
