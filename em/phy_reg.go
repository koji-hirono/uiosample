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

// Autoneg Advertisement Register
const (
	NWAY_AR_SELECTOR_FIELD = 0x0001 // indicates IEEE 802.3 CSMA/CD
	NWAY_AR_10T_HD_CAPS    = 0x0020 // 10T   Half Duplex Capable
	NWAY_AR_10T_FD_CAPS    = 0x0040 // 10T   Full Duplex Capable
	NWAY_AR_100TX_HD_CAPS  = 0x0080 // 100TX Half Duplex Capable
	NWAY_AR_100TX_FD_CAPS  = 0x0100 // 100TX Full Duplex Capable
	NWAY_AR_100T4_CAPS     = 0x0200 // 100T4 Capable
	NWAY_AR_PAUSE          = 0x0400 // Pause operation desired
	NWAY_AR_ASM_DIR        = 0x0800 // Asymmetric Pause Direction bit
	NWAY_AR_REMOTE_FAULT   = 0x2000 // Remote Fault detected
	NWAY_AR_NEXT_PAGE      = 0x8000 // Next Page ability supported
)

// Link Partner Ability Register (Base Page)
const (
	NWAY_LPAR_SELECTOR_FIELD = 0x0000 // LP protocol selector field
	NWAY_LPAR_10T_HD_CAPS    = 0x0020 // LP 10T Half Dplx Capable
	NWAY_LPAR_10T_FD_CAPS    = 0x0040 // LP 10T Full Dplx Capable
	NWAY_LPAR_100TX_HD_CAPS  = 0x0080 // LP 100TX Half Dplx Capable
	NWAY_LPAR_100TX_FD_CAPS  = 0x0100 // LP 100TX Full Dplx Capable
	NWAY_LPAR_100T4_CAPS     = 0x0200 // LP is 100T4 Capable
	NWAY_LPAR_PAUSE          = 0x0400 // LP Pause operation desired
	NWAY_LPAR_ASM_DIR        = 0x0800 // LP Asym Pause Direction bit
	NWAY_LPAR_REMOTE_FAULT   = 0x2000 // LP detected Remote Fault
	NWAY_LPAR_ACKNOWLEDGE    = 0x4000 // LP rx'd link code word
	NWAY_LPAR_NEXT_PAGE      = 0x8000 // Next Page ability supported
)

// Autoneg Expansion Register
const (
	NWAY_ER_LP_NWAY_CAPS      = 0x0001 // LP has Auto Neg Capability
	NWAY_ER_PAGE_RXD          = 0x0002 // LP 10T Half Dplx Capable
	NWAY_ER_NEXT_PAGE_CAPS    = 0x0004 // LP 10T Full Dplx Capable
	NWAY_ER_LP_NEXT_PAGE_CAPS = 0x0008 // LP 100TX Half Dplx Capable
	NWAY_ER_PAR_DETECT_FAULT  = 0x0010 // LP 100TX Full Dplx Capable
)

// 1000BASE-T Control Register
const (
	CR_1000T_ASYM_PAUSE = 0x0080 // Advertise asymmetric pause bit
	CR_1000T_HD_CAPS    = 0x0100 // Advertise 1000T HD capability
	CR_1000T_FD_CAPS    = 0x0200 // Advertise 1000T FD capability
	// 1=Repeater/switch device port 0=DTE device
	CR_1000T_REPEATER_DTE = 0x0400
	// 1=Configure PHY as Master 0=Configure PHY as Slave
	CR_1000T_MS_VALUE = 0x0800
	// 1=Master/Slave manual config value 0=Automatic Master/Slave config
	CR_1000T_MS_ENABLE        = 0x1000
	CR_1000T_TEST_MODE_NORMAL = 0x0000 // Normal Operation
	CR_1000T_TEST_MODE_1      = 0x2000 // Transmit Waveform test
	CR_1000T_TEST_MODE_2      = 0x4000 // Master Transmit Jitter test
	CR_1000T_TEST_MODE_3      = 0x6000 // Slave Transmit Jitter test
	CR_1000T_TEST_MODE_4      = 0x8000 // Transmitter Distortion test
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
	HV_INTC_FC_PAGE_START       = 768
	I82578_ADDR_REG             = 29
	I82577_ADDR_REG             = 16
	I82577_CFG_REG              = 22
	I82577_CFG_ASSERT_CRS_ON_TX = 1 << 15
	I82577_CFG_ENABLE_DOWNSHIFT = 3 << 10 // auto downshift
	I82577_CTRL_REG             = 23
)

// 82577 specific PHY registers
const (
	I82577_PHY_CTRL_2      = 18
	I82577_PHY_LBK_CTRL    = 19
	I82577_PHY_STATUS_2    = 26
	I82577_PHY_DIAG_STATUS = 31
)

// I82577 PHY Status 2
const (
	I82577_PHY_STATUS2_REV_POLARITY   = 0x0400
	I82577_PHY_STATUS2_MDIX           = 0x0800
	I82577_PHY_STATUS2_SPEED_MASK     = 0x0300
	I82577_PHY_STATUS2_SPEED_1000MBPS = 0x0200
)

// I82577 PHY Control 2
const (
	I82577_PHY_CTRL2_MANUAL_MDIX   = 0x0200
	I82577_PHY_CTRL2_AUTO_MDI_MDIX = 0x0400
	I82577_PHY_CTRL2_MDIX_CFG_MASK = 0x0600
)

// I82577 PHY Diagnostics Status
const (
	I82577_DSTATUS_CABLE_LENGTH       = 0x03FC
	I82577_DSTATUS_CABLE_LENGTH_SHIFT = 2
)

// 82580 PHY Power Management
const (
	I82580_PHY_POWER_MGMT = 0xE14
	I82580_PM_SPD         = 0x0001 // Smart Power Down
	I82580_PM_D0_LPLU     = 0x0002 // For D0a states
	I82580_PM_D3_LPLU     = 0x0004 // For all other states
	I82580_PM_GO_LINKD    = 0x0020 // Go Link Disconnect
)

const (
	I347AT4_PSCR_DOWNSHIFT_ENABLE = 0x0800
	I347AT4_PSCR_DOWNSHIFT_MASK   = 0x7000
	I347AT4_PSCR_DOWNSHIFT_1X     = 0x0000
	I347AT4_PSCR_DOWNSHIFT_2X     = 0x1000
	I347AT4_PSCR_DOWNSHIFT_3X     = 0x2000
	I347AT4_PSCR_DOWNSHIFT_4X     = 0x3000
	I347AT4_PSCR_DOWNSHIFT_5X     = 0x4000
	I347AT4_PSCR_DOWNSHIFT_6X     = 0x5000
	I347AT4_PSCR_DOWNSHIFT_7X     = 0x6000
	I347AT4_PSCR_DOWNSHIFT_8X     = 0x7000
)

const (
	I82578_EPSCR_DOWNSHIFT_ENABLE       = 0x0020
	I82578_EPSCR_DOWNSHIFT_COUNTER_MASK = 0x001C
)

// BME1000 PHY Specific Control Register
const (
	BME1000_PSCR_ENABLE_DOWNSHIFT = 0x0800 // 1 = enable downshift
)
