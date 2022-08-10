package em

import (
	"time"
)

func PowerUpPHYCopper(hw *HW) {
	phy := &hw.PHY
	x, err := phy.Op.ReadReg(PHY_CONTROL)
	if err != nil {
		return
	}
	x &^= MII_CR_POWER_DOWN
	phy.Op.WriteReg(PHY_CONTROL, x)
}

func PowerDownPHYCopper(hw *HW) {
	phy := &hw.PHY
	x, err := phy.Op.ReadReg(PHY_CONTROL)
	if err != nil {
		return
	}
	x |= MII_CR_POWER_DOWN
	phy.Op.WriteReg(PHY_CONTROL, x)
	time.Sleep(1 * time.Millisecond)
}

func SetupCopperLink(hw *HW) error {
	if hw.MAC.Autoneg {
		// Setup autoneg and flow control advertisement and perform
		// autonegotiation.
		err := CopperLinkAutoneg(hw)
		if err != nil {
			return err
		}
	} else {
		// PHY will be set to 10H, 10F, 100H or 100F
		// depending on user settings.
		err := hw.PHY.Op.ForceSpeedDuplex()
		if err != nil {
			return err
		}
	}

	// Check link status. Wait up to 100 microseconds for link to become
	// valid.
	link, err := PHYHasLink(hw, COPPER_LINK_UP_LIMIT, 10)
	if err != nil {
		return err
	}

	if link {
		hw.MAC.Op.ConfigCollisionDist()
		err := ConfigFCAfterLinkUp(hw)
		if err != nil {
			return err
		}
	}
	return nil
}

func CopperLinkAutoneg(hw *HW) error {
	phy := &hw.PHY

	// Perform some bounds checking on the autoneg advertisement
	// parameter.
	phy.AutonegAdvertised &= phy.AutonegMask

	// If autoneg_advertised is zero, we assume it was not defaulted
	// by the calling code so we set to advertise full capability.
	if phy.AutonegAdvertised == 0 {
		phy.AutonegAdvertised = phy.AutonegMask
	}

	err := PHYSetupAutoneg(hw)
	if err != nil {
		return err
	}

	// Restart auto-negotiation by setting the Auto Neg Enable bit and
	// the Auto Neg Restart bit in the PHY control register.
	ctrl, err := phy.Op.ReadReg(PHY_CONTROL)
	if err != nil {
		return err
	}

	ctrl |= MII_CR_AUTO_NEG_EN | MII_CR_RESTART_AUTO_NEG
	err = phy.Op.WriteReg(PHY_CONTROL, ctrl)
	if err != nil {
		return err
	}

	// Does the user want to wait for Auto-Neg to complete here, or
	// check at a later time (for example, callback routine).
	if phy.AutonegWaitToComplete {
		err := WaitAutoneg(hw)
		if err != nil {
			return err
		}
	}

	hw.MAC.GetLinkStatus = true
	return nil
}
