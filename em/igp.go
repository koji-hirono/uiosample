package em

import (
	"time"
)

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

func PHYInitScriptIGP3(hw *HW) error {
	phy := &hw.PHY
	// PHY init IGP 3
	// Enable rise/fall, 10-mode work in class-A
	phy.Op.WriteReg(0x2F5B, 0x9018)
	// Remove all caps from Replica path filter
	phy.Op.WriteReg(0x2F52, 0x0000)
	// Bias trimming for ADC, AFE and Driver (Default)
	phy.Op.WriteReg(0x2FB1, 0x8B24)
	// Increase Hybrid poly bias
	phy.Op.WriteReg(0x2FB2, 0xF8F0)
	// Add 4% to Tx amplitude in Gig mode
	phy.Op.WriteReg(0x2010, 0x10B0)
	// Disable trimming (TTT)
	phy.Op.WriteReg(0x2011, 0x0000)
	// Poly DC correction to 94.6% + 2% for all channels
	phy.Op.WriteReg(0x20DD, 0x249A)
	// ABS DC correction to 95.9%
	phy.Op.WriteReg(0x20DE, 0x00D3)
	// BG temp curve trim
	phy.Op.WriteReg(0x28B4, 0x04CE)
	// Increasing ADC OPAMP stage 1 currents to max
	phy.Op.WriteReg(0x2F70, 0x29E4)
	// Force 1000 ( required for enabling PHY regs configuration)
	phy.Op.WriteReg(0x0000, 0x0140)
	// Set upd_freq to 6
	phy.Op.WriteReg(0x1F30, 0x1606)
	// Disable NPDFE
	phy.Op.WriteReg(0x1F31, 0xB814)
	// Disable adaptive fixed FFE (Default)
	phy.Op.WriteReg(0x1F35, 0x002A)
	// Enable FFE hysteresis
	phy.Op.WriteReg(0x1F3E, 0x0067)
	// Fixed FFE for short cable lengths
	phy.Op.WriteReg(0x1F54, 0x0065)
	// Fixed FFE for medium cable lengths
	phy.Op.WriteReg(0x1F55, 0x002A)
	// Fixed FFE for long cable lengths
	phy.Op.WriteReg(0x1F56, 0x002A)
	// Enable Adaptive Clip Threshold
	phy.Op.WriteReg(0x1F72, 0x3FB0)
	// AHT reset limit to 1
	phy.Op.WriteReg(0x1F76, 0xC0FF)
	// Set AHT master delay to 127 msec
	phy.Op.WriteReg(0x1F77, 0x1DEC)
	// Set scan bits for AHT
	phy.Op.WriteReg(0x1F78, 0xF9EF)
	// Set AHT Preset bits
	phy.Op.WriteReg(0x1F79, 0x0210)
	// Change integ_factor of channel A to 3
	phy.Op.WriteReg(0x1895, 0x0003)
	// Change prop_factor of channels BCD to 8
	phy.Op.WriteReg(0x1796, 0x0008)
	// Change cg_icount + enable integbp for channels BCD
	phy.Op.WriteReg(0x1798, 0xD008)
	// Change cg_icount + enable integbp + change prop_factor_master
	// to 8 for channel A
	phy.Op.WriteReg(0x1898, 0xD918)
	// Disable AHT in Slave mode on channel A
	phy.Op.WriteReg(0x187A, 0x0800)
	// Enable LPLU and disable AN to 1000 in non-D0a states,
	// Enable SPD+B2B
	phy.Op.WriteReg(0x0019, 0x008D)
	// Enable restart AN on an1000_dis change
	phy.Op.WriteReg(0x001B, 0x2080)
	// Enable wh_fifo read clock in 10/100 modes
	phy.Op.WriteReg(0x0014, 0x0045)
	// Restart AN, Speed selection is 1000
	phy.Op.WriteReg(0x0000, 0x1340)
	return nil
}

func CopperLinkSetupIGP(hw *HW) error {
	phy := &hw.PHY

	err := phy.Op.Reset()
	if err != nil {
		return err
	}

	// Wait 100ms for MAC to configure PHY from NVM settings, to avoid
	// timeout issues when LFS is enabled.
	time.Sleep(100 * time.Millisecond)

	// The NVM settings will configure LPLU in D3 for
	// non-IGP1 PHYs.
	if phy.PHYType == PHYTypeIgp {
		// disable lplu d3 during driver init
		err := phy.Op.SetD3LpluState(false)
		if err != nil {
			return err
		}
	}

	// disable lplu d0 during driver init
	err = phy.Op.SetD0LpluState(false)
	if err != nil {
		return err
	}

	// Configure mdi-mdix settings
	data, err := phy.Op.ReadReg(IGP01E1000_PHY_PORT_CTRL)
	if err != nil {
		return err
	}
	data &^= IGP01E1000_PSCR_AUTO_MDIX

	switch phy.MDIX {
	case 1:
		data &^= IGP01E1000_PSCR_FORCE_MDI_MDIX
	case 2:
		data |= IGP01E1000_PSCR_FORCE_MDI_MDIX
	default:
		data |= IGP01E1000_PSCR_AUTO_MDIX
	}
	err = phy.Op.WriteReg(IGP01E1000_PHY_PORT_CTRL, data)
	if err != nil {
		return err
	}

	// set auto-master slave resolution settings
	if hw.MAC.Autoneg {
		// when autonegotiation advertisement is only 1000Mbps then we
		// should disable SmartSpeed and enable Auto MasterSlave
		// resolution as hardware default.
		if phy.AutonegAdvertised == ADVERTISE_1000_FULL {
			// Disable SmartSpeed
			data, err := phy.Op.ReadReg(IGP01E1000_PHY_PORT_CONFIG)
			if err != nil {
				return err
			}

			data &^= IGP01E1000_PSCFR_SMART_SPEED
			err = phy.Op.WriteReg(IGP01E1000_PHY_PORT_CONFIG, data)
			if err != nil {
				return err
			}

			// Set auto Master/Slave resolution process
			data, err = phy.Op.ReadReg(PHY_1000T_CTRL)
			if err != nil {
				return err
			}

			data &^= CR_1000T_MS_ENABLE
			err = phy.Op.WriteReg(PHY_1000T_CTRL, data)
			if err != nil {
				return err
			}
		}
		return SetMasterSlaveMode(hw)
	}
	return nil
}
