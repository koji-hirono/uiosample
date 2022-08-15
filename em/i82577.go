package em

import (
	"errors"
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
	phy := &hw.PHY
	data, err := phy.Op.ReadReg(I82577_PHY_DIAG_STATUS)
	if err != nil {
		return err
	}

	data &= I82577_DSTATUS_CABLE_LENGTH
	data >>= I82577_DSTATUS_CABLE_LENGTH_SHIFT
	if data == CABLE_LENGTH_UNDEFINED {
		return errors.New("undefined length")
	}

	phy.CableLength = data

	return nil
}

func GetPHYInfo82577(hw *HW) error {
	phy := &hw.PHY
	link, err := PHYHasLink(hw, 1, 0)
	if err != nil {
		return err
	}
	if !link {
		return errors.New("link down")
	}

	phy.PolarityCorrection = true

	err = CheckPolarity82577(hw)
	if err != nil {
		return err
	}

	data, err := phy.Op.ReadReg(I82577_PHY_STATUS_2)
	if err != nil {
		return err
	}

	phy.IsMDIX = data&I82577_PHY_STATUS2_MDIX != 0

	if data&I82577_PHY_STATUS2_SPEED_MASK == I82577_PHY_STATUS2_SPEED_1000MBPS {
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
		phy.CableLength = CABLE_LENGTH_UNDEFINED
		phy.LocalRx = E1000TRxStatusUndefined
		phy.RemoteRx = E1000TRxStatusUndefined
	}

	return nil
}
