package em

func I82540Init(hw *HW) {
	nvm := NewI82540NVM(hw)
	phy := NewI82540PHY(hw)
	mac := NewI82540MAC(hw, nvm, phy)
	hw.MAC.Op = mac
	hw.PHY.Op = phy
	hw.NVM.Op = nvm
}
