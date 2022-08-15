package em

import (
	"time"
)

func CopperLinkSetup82577(hw *HW) error {
	phy := &hw.PHY
	if phy.PHYType == PHYType82580 {
		err := phy.Op.Reset()
		if err != nil {
			return err
		}
	}

	// Enable CRS on Tx. This must be set for half-duplex operation.
	data, err := phy.Op.ReadReg(I82577_CFG_REG)
	if err != nil {
		return err
	}
	data |= I82577_CFG_ASSERT_CRS_ON_TX

	// Enable downshift
	data |= I82577_CFG_ENABLE_DOWNSHIFT

	err = phy.Op.WriteReg(I82577_CFG_REG, data)
	if err != nil {
		return err
	}

	// Set MDI/MDIX mode
	data, err = phy.Op.ReadReg(I82577_PHY_CTRL_2)
	if err != nil {
		return err
	}
	data &^= I82577_PHY_CTRL2_MDIX_CFG_MASK

	// Options:
	//   0 - Auto (default)
	//   1 - MDI mode
	//   2 - MDI-X mode
	switch phy.MDIX {
	case 1:
	case 2:
		data |= I82577_PHY_CTRL2_MANUAL_MDIX
	default:
		data |= I82577_PHY_CTRL2_AUTO_MDI_MDIX
	}
	err = phy.Op.WriteReg(I82577_PHY_CTRL_2, data)
	if err != nil {
		return err
	}

	return SetMasterSlaveMode(hw)
}

func CheckPolarity82577(hw *HW) error {
	phy := &hw.PHY
	data, err := phy.Op.ReadReg(I82577_PHY_STATUS_2)
	if err != nil {
		return err
	}
	if data&I82577_PHY_STATUS2_REV_POLARITY != 0 {
		phy.CablePolarity = RevPolarityReversed
	} else {
		phy.CablePolarity = RevPolarityNormal
	}
	return nil
}

func PHYForceSpeedDuplex82577(hw *HW) error {
	phy := &hw.PHY

	data, err := phy.Op.ReadReg(PHY_CONTROL)
	if err != nil {
		return err
	}

	data = PHYForceSpeedDuplexSetup(hw, data)

	err = phy.Op.WriteReg(PHY_CONTROL, data)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Microsecond)

	if phy.AutonegWaitToComplete {
		_, err := PHYHasLink(hw, PHY_FORCE_LIMIT, 100000)
		if err != nil {
			return err
		}
		// Try once more
		_, err = PHYHasLink(hw, PHY_FORCE_LIMIT, 100000)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetCableLength82577(hw *HW) error {
	return nil
}

func GetPHYInfo82577(hw *HW) error {
	return nil
}
