package em

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
	// msec_delay(1)
}
