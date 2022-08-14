package em

// IGP01E1000 Specific Registers
const (
	IGP01E1000_PHY_PORT_CONFIG = 0x10 // Port Config
	IGP01E1000_PHY_PORT_STATUS = 0x11 // Status
	IGP01E1000_PHY_PORT_CTRL   = 0x12 // Control
	IGP01E1000_PHY_LINK_HEALTH = 0x13 // PHY Link Health
	IGP01E1000_GMII_FIFO       = 0x14 // GMII FIFO
	IGP02E1000_PHY_POWER_MGMT  = 0x19 // Power Management
	IGP01E1000_PHY_PAGE_SELECT = 0x1F // Page Select

	BM_PHY_PAGE_SELECT = 22 // Page Select for BM

	IGP_PAGE_SHIFT = 5
	PHY_REG_MASK   = 0x1f
)

const (
	IGP01E1000_PHY_PCS_INIT_REG  = 0x00B4
	IGP01E1000_PHY_POLARITY_MASK = 0x0078
)

const (
	IGP01E1000_PSCR_AUTO_MDIX      = 0x1000
	IGP01E1000_PSCR_FORCE_MDI_MDIX = 0x2000 // 0=MDI, 1=MDIX
)

const (
	IGP01E1000_PSCFR_SMART_SPEED = 0x0080
)

// Enable flexible speed on link-up
const (
	IGP01E1000_GMII_FLEX_SPD = 0x0010
	IGP01E1000_GMII_SPD      = 0x0020 // Enable SPD
)

const (
	IGP02E1000_PM_SPD     = 0x0001 // Smart Power Down
	IGP02E1000_PM_D0_LPLU = 0x0002 // For D0a states
	IGP02E1000_PM_D3_LPLU = 0x0004 // For all other states
)

const (
	IGP01E1000_PLHR_SS_DOWNGRADE = 0x8000
)

const (
	IGP01E1000_PSSR_POLARITY_REVERSED = 0x0002
	IGP01E1000_PSSR_MDIX              = 0x0800
	IGP01E1000_PSSR_SPEED_MASK        = 0xc000
	IGP01E1000_PSSR_SPEED_1000MBPS    = 0xc000
)

const (
	IGP02E1000_PHY_CHANNEL_NUM = 4
	IGP02E1000_PHY_AGC_A       = 0x11b1
	IGP02E1000_PHY_AGC_B       = 0x12b1
	IGP02E1000_PHY_AGC_C       = 0x14b1
	IGP02E1000_PHY_AGC_D       = 0x18b1
)

const (
	IGP02E1000_AGC_LENGTH_SHIFT = 9 // Course=15:13, Fine=12:9
	IGP02E1000_AGC_LENGTH_MASK  = 0x7f
	IGP02E1000_AGC_RANGE        = 15
)

func ReadPHYRegIGP(hw *HW, offset uint32) (uint16, error) {
	phy := &hw.PHY
	err := phy.Op.Acquire()
	if err != nil {
		return 0, err
	}
	defer phy.Op.Release()

	if offset > MAX_PHY_MULTI_PAGE_REG {
		err := WritePHYRegMDIC(hw, IGP01E1000_PHY_PAGE_SELECT, uint16(offset))
		if err != nil {
			return 0, err
		}
	}
	return ReadPHYRegMDIC(hw, MAX_PHY_REG_ADDRESS&offset)
}

func WritePHYRegIGP(hw *HW, offset uint32, data uint16) error {
	phy := &hw.PHY
	err := phy.Op.Acquire()
	if err != nil {
		return err
	}
	defer phy.Op.Release()

	if offset > MAX_PHY_MULTI_PAGE_REG {
		err := WritePHYRegMDIC(hw, IGP01E1000_PHY_PAGE_SELECT, uint16(offset))
		if err != nil {
			return err
		}
	}
	return WritePHYRegMDIC(hw, MAX_PHY_REG_ADDRESS&offset, data)
}

func CheckPolarityIGP(hw *HW) error {
	return nil
}

func PHYForceSpeedDuplexIGP(hw *HW) error {
	return nil
}

func GetCableLengthIGP2(hw *HW) error {
	return nil
}

func GetPHYInfoIGP(hw *HW) error {
	return nil
}
