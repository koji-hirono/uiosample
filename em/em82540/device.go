package em82540

import (
	"uiosample/em"
)

type Device struct {
	hw  em.HW
	MAC *MAC
	PHY *PHY
	NVM *NVM
}

func NewDevice() *Device {
	d := new(Device)
	d.NVM = NewNVM(&d.hw)
	d.PHY = NewPHY(&d.hw)
	d.MAC = NewMAC(&d.hw, d.NVM, d.PHY)
	d.hw.MAC.Op = d.MAC
	d.hw.PHY.Op = d.PHY
	d.hw.NVM.Op = d.NVM
	return d
}
