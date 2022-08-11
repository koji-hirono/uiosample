package em82540

import (
	"uiosample/em"
)

func Init(hw *em.HW) {
	nvm := NewNVM(hw)
	phy := NewPHY(hw)
	mac := NewMAC(hw, nvm, phy)
	hw.MAC.Op = mac
	hw.PHY.Op = phy
	hw.NVM.Op = nvm
}
