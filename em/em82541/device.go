package em82541

import (
	"uiosample/em"
)

type FFEConfig int

const (
	FFEConfigEnabled FFEConfig = iota
	FFEConfigActive
	FFEConfigBlocked
)

type DSPConfig int

const (
	DSPConfigDisabled DSPConfig = iota
	DSPConfigEnabled
	DSPConfigActivated
	DSPConfigUndefined = 0xff
)

type Device struct {
	hw            em.HW
	MAC           *MAC
	PHY           *PHY
	NVM           *NVM
	DSPConfig     DSPConfig
	FFEConfig     FFEConfig
	SpdDefault    uint16
	PHYInitScript bool
}

func NewDevice() *Device {
	d := new(Device)
	d.MAC = NewMAC(&d.hw)
	d.PHY = NewPHY(&d.hw)
	d.NVM = NewNVM(&d.hw)
	return d
}
