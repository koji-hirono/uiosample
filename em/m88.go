package em

import (
	"errors"
	"time"
)

// M88E1000 Specific Registers
const (
	M88E1000_PHY_SPEC_CTRL     uint32 = 0x10
	M88E1000_PHY_SPEC_STATUS          = 0x11
	M88E1000_EXT_PHY_SPEC_CTRL        = 0x14
	M88E1000_RX_ERR_CNTR              = 0x15
	M88E1000_PHY_EXT_CTRL             = 0x1a
	M88E1000_PHY_PAGE_SELECT          = 0x1d
	M88E1000_PHY_GEN_CONTROL          = 0x1e

	M88E1000_PHY_VCO_REG_BIT8  = 0x100
	M88E1000_PHY_VCO_REG_BIT11 = 0x800
)

// M88E1000 PHY Specific Control Register
const (
	M88E1000_PSCR_POLARITY_REVERSAL uint16 = 0x0002 // 1=Polarity Reverse enabled
	// MDI Crossover Mode bits 6:5 Manual MDI configuration
	M88E1000_PSCR_MDI_MANUAL_MODE  = 0x0000
	M88E1000_PSCR_MDIX_MANUAL_MODE = 0x0020 // Manual MDIX configuration
	// 1000BASE-T: Auto crossover, 100BASE-TX/10BASE-T: MDI Mode
	M88E1000_PSCR_AUTO_X_1000T = 0x0040
	// Auto crossover enabled all speeds
	M88E1000_PSCR_AUTO_X_MODE      = 0x0060
	M88E1000_PSCR_ASSERT_CRS_ON_TX = 0x0800 // 1=Assert CRS on Tx
)

// M88E1000 PHY Specific Status Register
const (
	M88E1000_PSSR_REV_POLARITY uint16 = 0x0002 // 1=Polarity reversed
	M88E1000_PSSR_DOWNSHIFT           = 0x0020 // 1=Downshifted
	M88E1000_PSSR_MDIX                = 0x0040 // 1=MDIX; 0=MDI
)

// 0 = <50M
// 1 = 50-80M
// 2 = 80-110M
// 3 = 110-140M
// 4 = >140M
const (
	M88E1000_PSSR_CABLE_LENGTH      = 0x0380
	M88E1000_PSSR_LINK              = 0x0400 // 1=Link up, 0=Link down
	M88E1000_PSSR_SPD_DPLX_RESOLVED = 0x0800 // 1=Speed & Duplex resolved
	M88E1000_PSSR_DPLX              = 0x2000 // 1=Duplex 0=Half Duplex
	M88E1000_PSSR_SPEED             = 0xC000 // Speed, bits 14:15
	M88E1000_PSSR_100MBS            = 0x4000 // 01=100Mbs
	M88E1000_PSSR_1000MBS           = 0x8000 // 10=1000Mbs

	M88E1000_PSSR_CABLE_LENGTH_SHIFT = 7
)

// Number of times we will attempt to autonegotiate before downshifting if we
// are the master
const (
	M88E1000_EPSCR_MASTER_DOWNSHIFT_MASK = 0x0C00
	M88E1000_EPSCR_MASTER_DOWNSHIFT_1X   = 0x0000
)

// Number of times we will attempt to autonegotiate before downshifting if we
// are the slave
const (
	M88E1000_EPSCR_SLAVE_DOWNSHIFT_MASK = 0x0300
	M88E1000_EPSCR_SLAVE_DOWNSHIFT_1X   = 0x0100
	M88E1000_EPSCR_TX_CLK_25            = 0x0070 // 25  MHz TX_CLK
)

// M88E1112 only registers
const (
	M88E1112_VCT_DSP_DISTANCE     = 0x001A
	M88E1112_AUTO_COPPER_SGMII    = 0x2
	M88E1112_AUTO_COPPER_BASEX    = 0x3
	M88E1112_STATUS_LINK          = 0x0004 // Interface Link Bit
	M88E1112_MAC_CTRL_1           = 0x10
	M88E1112_MAC_CTRL_1_MODE_MASK = 0x0380 // Mode Select

	M88E1112_MAC_CTRL_1_MODE_SHIFT = 7

	M88E1112_PAGE_ADDR = 0x16
	M88E1112_STATUS    = 0x01
)

// M88E1543
const (
	M88E1543_PAGE_ADDR     = 0x16 // Page Offset Register
	M88E1543_EEE_CTRL_1    = 0x0
	M88E1543_EEE_CTRL_1_MS = 0x0001 // EEE Master/Slave
	M88E1543_FIBER_CTRL    = 0x0    // Fiber Control Register
)

// M88E1512
const (
	M88E1512_CFG_REG_1 = 0x0010
	M88E1512_CFG_REG_2 = 0x0011
	M88E1512_CFG_REG_3 = 0x0007
	M88E1512_MODE      = 0x0014
)

// M88EC018 Rev 2 specific DownShift settings
const (
	M88EC018_EPSCR_DOWNSHIFT_COUNTER_MASK = 0x0E00
	M88EC018_EPSCR_DOWNSHIFT_COUNTER_5X   = 0x0800
)

func CheckPolarityM88(hw *HW) error {
	phy := &hw.PHY
	data, err := phy.Op.ReadReg(M88E1000_PHY_SPEC_STATUS)
	if err != nil {
		return err
	}
	if data&M88E1000_PSSR_REV_POLARITY != 0 {
		phy.CablePolarity = RevPolarityReversed
	} else {
		phy.CablePolarity = RevPolarityNormal
	}
	return nil
}

func PHYForceSpeedDuplexM88(hw *HW) error {
	phy := &hw.PHY

	// I210 and I211 devices support Auto-Crossover in forced operation.
	if phy.PHYType != PHYTypeI210 {
		// Clear Auto-Crossover to force MDI manually.  M88E1000
		// requires MDI forced whenever speed and duplex are forced.
		data, err := phy.Op.ReadReg(M88E1000_PHY_SPEC_CTRL)
		if err != nil {
			return err
		}

		data &^= M88E1000_PSCR_AUTO_X_MODE
		err = phy.Op.WriteReg(M88E1000_PHY_SPEC_CTRL, data)
		if err != nil {
			return err
		}
	}

	data, err := phy.Op.ReadReg(PHY_CONTROL)
	if err != nil {
		return err
	}

	// e1000_phy_force_speed_duplex_setup(hw, &phy_data);

	err = phy.Op.WriteReg(PHY_CONTROL, data)
	if err != nil {
		return err
	}

	// Reset the phy to commit changes.
	err = phy.Op.Commit()
	if err != nil {
		return err
	}

	if phy.AutonegWaitToComplete {
		link, err := PHYHasLink(hw, PHY_FORCE_LIMIT, 100000)
		if err != nil {
			return err
		}
		if !link {
			reset_dsp := true
			switch phy.ID {
			case I347AT4_E_PHY_ID, M88E1340M_E_PHY_ID, M88E1112_E_PHY_ID, M88E1543_E_PHY_ID, M88E1512_E_PHY_ID, I210_I_PHY_ID:
				reset_dsp = false
			default:
				if phy.PHYType != PHYTypeM88 {
					reset_dsp = false
				}
			}
			if reset_dsp {
				// We didn't get link.
				// Reset the DSP and cross our fingers.
				err := phy.Op.WriteReg(M88E1000_PHY_PAGE_SELECT, 0x001d)
				if err != nil {
					return err
				}
				err = PHYResetDSP(hw)
				if err != nil {
					return err
				}
			}
		}

		// Try once more
		link, err = PHYHasLink(hw, PHY_FORCE_LIMIT, 100000)
		if err != nil {
			return err
		}
	}
	if phy.PHYType != PHYTypeM88 {
		return nil
	}

	if phy.ID == I347AT4_E_PHY_ID ||
		phy.ID == M88E1340M_E_PHY_ID ||
		phy.ID == M88E1112_E_PHY_ID {
		return nil
	}
	if phy.ID == I210_I_PHY_ID {
		return nil
	}
	if phy.ID == M88E1543_E_PHY_ID || phy.ID == M88E1512_E_PHY_ID {
		return nil
	}

	data, err = phy.Op.ReadReg(M88E1000_EXT_PHY_SPEC_CTRL)
	if err != nil {
		return err
	}

	// Resetting the phy means we need to re-force TX_CLK in the
	// Extended PHY Specific Control Register to 25MHz clock from
	// the reset value of 2.5MHz.
	data |= M88E1000_EPSCR_TX_CLK_25
	err = phy.Op.WriteReg(M88E1000_EXT_PHY_SPEC_CTRL, data)
	if err != nil {
		return err
	}

	// In addition, we must re-enable CRS on Tx for both half and full
	// duplex.
	data, err = phy.Op.ReadReg(M88E1000_PHY_SPEC_CTRL)
	if err != nil {
		return err
	}

	data |= M88E1000_PSCR_ASSERT_CRS_ON_TX
	return phy.Op.WriteReg(M88E1000_PHY_SPEC_CTRL, data)
}

// Cable length tables
var m88CableLengthTable = [...]uint16{0, 50, 80, 110, 140, 140}

func GetCableLengthM88(hw *HW) error {
	phy := &hw.PHY

	data, err := phy.Op.ReadReg(M88E1000_PHY_SPEC_STATUS)
	if err != nil {
		return err
	}

	index := int(data & M88E1000_PSSR_CABLE_LENGTH)
	index >>= M88E1000_PSSR_CABLE_LENGTH_SHIFT
	if index >= len(m88CableLengthTable) {
		return errors.New("invalid index")
	}

	phy.MinCableLength = m88CableLengthTable[index]
	phy.MaxCableLength = m88CableLengthTable[index+1]
	phy.CableLength = (phy.MinCableLength + phy.MaxCableLength) / 2
	return nil
}

func GetCableLengthM88gen2(hw *HW) error {
	phy := &hw.PHY
	switch phy.ID {
	case I210_I_PHY_ID:
		// Get cable length from PHY Cable Diagnostics Control Reg
		length, err := phy.Op.ReadReg((0x7 << GS40G_PAGE_SHIFT) + (I347AT4_PCDL + phy.Addr))
		if err != nil {
			return err
		}
		// Check if the unit of cable length is meters or cm
		unit, err := phy.Op.ReadReg((0x7 << GS40G_PAGE_SHIFT) + I347AT4_PCDC)
		if err != nil {
			return err
		}
		if unit&I347AT4_PCDC_CABLE_LENGTH_UNIT == 0 {
			length /= 100
		}
		// Populate the phy structure with cable length in meters
		phy.MinCableLength = length
		phy.MaxCableLength = length
		phy.CableLength = length
	case M88E1543_E_PHY_ID, M88E1512_E_PHY_ID, M88E1340M_E_PHY_ID, I347AT4_E_PHY_ID:
		// Remember the original page select and set it to 7
		defpage, err := phy.Op.ReadReg(I347AT4_PAGE_SELECT)
		if err != nil {
			return err
		}
		err = phy.Op.WriteReg(I347AT4_PAGE_SELECT, 0x07)
		if err != nil {
			return err
		}
		// Get cable length from PHY Cable Diagnostics Control Reg
		length, err := phy.Op.ReadReg(I347AT4_PCDL + phy.Addr)
		if err != nil {
			return err
		}
		// Check if the unit of cable length is meters or cm
		unit, err := phy.Op.ReadReg(I347AT4_PCDC)
		if err != nil {
			return err
		}
		if unit&I347AT4_PCDC_CABLE_LENGTH_UNIT == 0 {
			length /= 100
		}
		// Populate the phy structure with cable length in meters
		phy.MinCableLength = length
		phy.MaxCableLength = length
		phy.CableLength = length

		// Reset the page select to its original value
		err = phy.Op.WriteReg(I347AT4_PAGE_SELECT, defpage)
		if err != nil {
			return err
		}
	case M88E1112_E_PHY_ID:
		// Remember the original page select and set it to 5
		defpage, err := phy.Op.ReadReg(I347AT4_PAGE_SELECT)
		if err != nil {
			return err
		}
		err = phy.Op.WriteReg(I347AT4_PAGE_SELECT, 0x05)
		if err != nil {
			return err
		}
		data, err := phy.Op.ReadReg(M88E1112_VCT_DSP_DISTANCE)
		if err != nil {
			return err
		}
		data &= M88E1000_PSSR_CABLE_LENGTH
		data >>= M88E1000_PSSR_CABLE_LENGTH_SHIFT
		if int(data) >= len(m88CableLengthTable) {
			return errors.New("over table size")
		}
		phy.MinCableLength = m88CableLengthTable[data]
		phy.MaxCableLength = m88CableLengthTable[data+1]
		phy.CableLength = (phy.MinCableLength + phy.MaxCableLength) / 2
		// Reset the page select to its original value */
		err = phy.Op.WriteReg(I347AT4_PAGE_SELECT, defpage)
		if err != nil {
			return err
		}
	default:
		return errors.New("not support")
	}
	return nil
}

func GetPHYInfoM88(hw *HW) error {
	phy := &hw.PHY

	if phy.MediaType != MediaTypeCopper {
		return errors.New("Phy info is only valid for copper media")
	}

	link, err := PHYHasLink(hw, 1, 0)
	if err != nil {
		return err
	}
	if !link {
		return errors.New("Phy info is only valid if link is up")
	}

	data, err := phy.Op.ReadReg(M88E1000_PHY_SPEC_CTRL)
	if err != nil {
		return err
	}

	phy.PolarityCorrection = data&M88E1000_PSCR_POLARITY_REVERSAL != 0

	err = CheckPolarityM88(hw)
	if err != nil {
		return err
	}

	data, err = phy.Op.ReadReg(M88E1000_PHY_SPEC_STATUS)
	if err != nil {
		return err
	}

	phy.IsMDIX = data&M88E1000_PSSR_MDIX != 0

	if data&M88E1000_PSSR_SPEED == M88E1000_PSSR_1000MBS {
		err := phy.Op.GetCableLength()
		if err != nil {
			return err
		}

		data, err := phy.Op.ReadReg(PHY_1000T_STATUS)
		if err != nil {
			return err
		}

		if data&SR_1000T_LOCAL_RX_STATUS != 0 {
			phy.LocalRx = E1000TRxStatusOk
		} else {
			phy.LocalRx = E1000TRxStatusNotOk
		}

		if data&SR_1000T_REMOTE_RX_STATUS != 0 {
			phy.RemoteRx = E1000TRxStatusOk
		} else {
			phy.RemoteRx = E1000TRxStatusNotOk
		}
	} else {
		// Set values to "undefined"
		phy.CableLength = CABLE_LENGTH_UNDEFINED
		phy.LocalRx = E1000TRxStatusUndefined
		phy.RemoteRx = E1000TRxStatusUndefined
	}

	return nil
}

func ReadPHYRegM88(hw *HW, offset uint32) (uint16, error) {
	phy := &hw.PHY
	err := phy.Op.Acquire()
	if err != nil {
		return 0, err
	}
	defer phy.Op.Release()

	return ReadPHYRegMDIC(hw, MAX_PHY_REG_ADDRESS&offset)
}

func WritePHYRegM88(hw *HW, offset uint32, val uint16) error {
	phy := &hw.PHY
	err := phy.Op.Acquire()
	if err != nil {
		return err
	}
	defer phy.Op.Release()

	return WritePHYRegMDIC(hw, MAX_PHY_REG_ADDRESS&offset, val)
}

func CopperLinkSetupM88(hw *HW) error {
	phy := &hw.PHY

	// Enable CRS on Tx. This must be set for half-duplex operation.
	data, err := phy.Op.ReadReg(M88E1000_PHY_SPEC_CTRL)
	if err != nil {
		return err
	}

	// For BM PHY this bit is downshift enable
	if phy.PHYType != PHYTypeBm {
		data |= M88E1000_PSCR_ASSERT_CRS_ON_TX
	}

	// Options:
	// MDI/MDI-X = 0 (default)
	// 0 - Auto for all speeds
	// 1 - MDI mode
	// 2 - MDI-X mode
	// 3 - Auto for 1000Base-T only (MDI-X for 10/100Base-T modes)
	data &^= M88E1000_PSCR_AUTO_X_MODE
	switch phy.MDIX {
	case 1:
		data |= M88E1000_PSCR_MDI_MANUAL_MODE
	case 2:
		data |= M88E1000_PSCR_MDIX_MANUAL_MODE
	case 3:
		data |= M88E1000_PSCR_AUTO_X_1000T
	default:
		data |= M88E1000_PSCR_AUTO_X_MODE
	}

	// Options:
	// disable_polarity_correction = 0 (default)
	//     Automatic Correction for Reversed Cable Polarity
	// 0 - Disabled
	// 1 - Enabled
	data &^= M88E1000_PSCR_POLARITY_REVERSAL
	if phy.DisablePolarityCorrection {
		data |= M88E1000_PSCR_POLARITY_REVERSAL
	}

	// Enable downshift on BM (disabled by default)
	if phy.PHYType == PHYTypeBm {
		// For 82574/82583, first disable then enable downshift
		if phy.ID == BME1000_E_PHY_ID_R2 {
			data &^= BME1000_PSCR_ENABLE_DOWNSHIFT
			err := phy.Op.WriteReg(M88E1000_PHY_SPEC_CTRL, data)
			if err != nil {
				return err
			}
			// Commit the changes.
			err = phy.Op.Commit()
			if err != nil {
				return err
			}
		}

		data |= BME1000_PSCR_ENABLE_DOWNSHIFT
	}

	err = phy.Op.WriteReg(M88E1000_PHY_SPEC_CTRL, data)
	if err != nil {
		return err
	}

	if phy.PHYType == PHYTypeM88 && phy.Revision < 4 &&
		phy.ID != BME1000_E_PHY_ID_R2 {
		// Force TX_CLK in the Extended PHY Specific Control Register
		// to 25MHz clock.
		data, err := phy.Op.ReadReg(M88E1000_EXT_PHY_SPEC_CTRL)
		if err != nil {
			return err
		}
		data |= M88E1000_EPSCR_TX_CLK_25

		if phy.Revision == 2 && phy.ID == M88E1111_I_PHY_ID {
			// 82573L PHY - set the downshift counter to 5x.
			data &^= M88EC018_EPSCR_DOWNSHIFT_COUNTER_MASK
			data |= M88EC018_EPSCR_DOWNSHIFT_COUNTER_5X
		} else {
			// Configure Master and Slave downshift values
			data &^= M88E1000_EPSCR_MASTER_DOWNSHIFT_MASK | M88E1000_EPSCR_SLAVE_DOWNSHIFT_MASK
			data |= M88E1000_EPSCR_MASTER_DOWNSHIFT_1X
			data |= M88E1000_EPSCR_SLAVE_DOWNSHIFT_1X
		}
		err = phy.Op.WriteReg(M88E1000_EXT_PHY_SPEC_CTRL, data)
		if err != nil {
			return err
		}
	}

	if phy.PHYType == PHYTypeBm && phy.ID == BME1000_E_PHY_ID_R2 {
		// Set PHY page 0, register 29 to 0x0003
		err := phy.Op.WriteReg(29, 0x0003)
		if err != nil {
			return err
		}

		// Set PHY page 0, register 30 to 0x0000
		err = phy.Op.WriteReg(30, 0x0000)
		if err != nil {
			return err
		}
	}

	// Commit the changes.
	err = phy.Op.Commit()
	if err != nil {
		return err
	}

	if phy.PHYType == PHYType82578 {
		data, err := phy.Op.ReadReg(M88E1000_EXT_PHY_SPEC_CTRL)
		if err != nil {
			return err
		}

		// 82578 PHY - set the downshift count to 1x.
		data |= I82578_EPSCR_DOWNSHIFT_ENABLE
		data &^= I82578_EPSCR_DOWNSHIFT_COUNTER_MASK
		err = phy.Op.WriteReg(M88E1000_EXT_PHY_SPEC_CTRL, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func CopperLinkSetupM88gen2(hw *HW) error {
	phy := &hw.PHY
	// Enable CRS on Tx. This must be set for half-duplex operation.
	data, err := phy.Op.ReadReg(M88E1000_PHY_SPEC_CTRL)
	if err != nil {
		return err
	}

	// Options:
	//   MDI/MDI-X = 0 (default)
	//   0 - Auto for all speeds
	//   1 - MDI mode
	//   2 - MDI-X mode
	//   3 - Auto for 1000Base-T only (MDI-X for 10/100Base-T modes)
	data &^= M88E1000_PSCR_AUTO_X_MODE

	switch phy.MDIX {
	case 1:
		data |= M88E1000_PSCR_MDI_MANUAL_MODE
	case 2:
		data |= M88E1000_PSCR_MDIX_MANUAL_MODE
	case 3:
		// M88E1112 does not support this mode)
		if phy.ID != M88E1112_E_PHY_ID {
			data |= M88E1000_PSCR_AUTO_X_1000T
		}
		data |= M88E1000_PSCR_AUTO_X_MODE
	default:
		data |= M88E1000_PSCR_AUTO_X_MODE
	}

	// Options:
	//   disable_polarity_correction = 0 (default)
	//       Automatic Correction for Reversed Cable Polarity
	//   0 - Disabled
	//   1 - Enabled
	data &^= M88E1000_PSCR_POLARITY_REVERSAL
	if phy.DisablePolarityCorrection {
		data |= M88E1000_PSCR_POLARITY_REVERSAL
	}

	// Enable downshift and setting it to X6
	if phy.ID == M88E1543_E_PHY_ID {
		data &^= I347AT4_PSCR_DOWNSHIFT_ENABLE
		err := phy.Op.WriteReg(M88E1000_PHY_SPEC_CTRL, data)
		if err != nil {
			return err
		}

		err = phy.Op.Commit()
		if err != nil {
			return err
		}
	}

	data &^= I347AT4_PSCR_DOWNSHIFT_MASK
	data |= I347AT4_PSCR_DOWNSHIFT_6X
	data |= I347AT4_PSCR_DOWNSHIFT_ENABLE

	err = phy.Op.WriteReg(M88E1000_PHY_SPEC_CTRL, data)
	if err != nil {
		return err
	}

	// Commit the changes.
	err = phy.Op.Commit()
	if err != nil {
		return err
	}

	return SetMasterSlaveMode(hw)
}

// s32 e1000_initialize_M88E1512_phy(struct e1000_hw *hw)
func InitM88E1512PHY(hw *HW) error {
	phy := &hw.PHY
	// Check if this is correct PHY.
	if phy.ID != M88E1512_E_PHY_ID {
		return nil
	}
	// Switch to PHY page 0xFF.
	err := phy.Op.WriteReg(M88E1543_PAGE_ADDR, 0x00ff)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_2, 0x214b)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_1, 0x2144)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_2, 0x0c28)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_1, 0x2146)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_2, 0xb233)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_1, 0x214d)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_2, 0xcc0c)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_1, 0x2159)
	if err != nil {
		return err
	}

	// Switch to PHY page 0xFB.
	err = phy.Op.WriteReg(M88E1543_PAGE_ADDR, 0x00fb)
	if err != nil {
		return err
	}
	err = phy.Op.WriteReg(M88E1512_CFG_REG_3, 0x000d)
	if err != nil {
		return err
	}

	// Switch to PHY page 0x12.
	err = phy.Op.WriteReg(M88E1543_PAGE_ADDR, 0x12)
	if err != nil {
		return err
	}

	// Change mode to SGMII-to-Copper
	err = phy.Op.WriteReg(M88E1512_MODE, 0x8001)
	if err != nil {
		return err
	}

	// Return the PHY to page 0.
	err = phy.Op.WriteReg(M88E1543_PAGE_ADDR, 0)
	if err != nil {
		return err
	}

	err = phy.Op.Commit()
	if err != nil {
		return err
	}

	time.Sleep(1000 * time.Millisecond)
	return nil
}

// s32 e1000_initialize_M88E1543_phy(struct e1000_hw *hw)
func InitM88E1543PHY(hw *HW) error {
	phy := &hw.PHY
	// Check if this is correct PHY.
	if phy.ID != M88E1543_E_PHY_ID {
		return nil
	}

	// Switch to PHY page 0xFF.
	err := phy.Op.WriteReg(M88E1543_PAGE_ADDR, 0x00ff)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_2, 0x214b)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_1, 0x2144)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_2, 0x0c28)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_1, 0x2146)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_2, 0xb233)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_1, 0x214d)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_2, 0xdc0c)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_1, 0x2159)
	if err != nil {
		return err
	}

	// Switch to PHY page 0xFB.
	err = phy.Op.WriteReg(M88E1543_PAGE_ADDR, 0x00fb)
	if err != nil {
		return err
	}

	err = phy.Op.WriteReg(M88E1512_CFG_REG_3, 0xc00d)
	if err != nil {
		return err
	}

	// Switch to PHY page 0x12.
	err = phy.Op.WriteReg(M88E1543_PAGE_ADDR, 0x12)
	if err != nil {
		return err
	}

	// Change mode to SGMII-to-Copper
	err = phy.Op.WriteReg(M88E1512_MODE, 0x8001)
	if err != nil {
		return err
	}

	// Switch to PHY page 1.
	err = phy.Op.WriteReg(M88E1543_PAGE_ADDR, 0x1)
	if err != nil {
		return err
	}

	// Change mode to 1000BASE-X/SGMII and autoneg enable; reset
	err = phy.Op.WriteReg(M88E1543_FIBER_CTRL, 0x9140)
	if err != nil {
		return err
	}

	// Return the PHY to page 0.
	err = phy.Op.WriteReg(M88E1543_PAGE_ADDR, 0)
	if err != nil {
		return err
	}

	err = phy.Op.Commit()
	if err != nil {
		return err
	}

	time.Sleep(1000 * time.Millisecond)
	return nil
}
