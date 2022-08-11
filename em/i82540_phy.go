package em

import (
	"errors"
)

type I82540PHY struct {
	hw *HW
}

func NewI82540PHY(hw *HW) *I82540PHY {
	p := new(I82540PHY)
	p.hw = hw
	return p
}

func (p *I82540PHY) InitParams() error {
	phy := &p.hw.PHY

	phy.Addr = 1
	phy.AutonegMask = AUTONEG_ADVERTISE_SPEED_DEFAULT
	phy.ResetDelayUS = 10000
	phy.PHYType = PHYTypeM88

	err := GetPHYID(p.hw)
	if err != nil {
		return err
	}

	switch p.hw.MAC.Type {
	case MACType82540:
	case MACType82545:
	case MACType82545Rev3:
	case MACType82546:
	case MACType82546Rev3:
	default:
		return errors.New("Invalid MAC type")
	}

	if phy.ID != M88E1011_I_PHY_ID {
		return errors.New("Invalid phy id")
	}

	return nil
}

func (p *I82540PHY) Acquire() error {
	// null
	return nil
}

func (p *I82540PHY) CheckPolarity() error {
	return CheckPolarityM88(p.hw)
}

func (p *I82540PHY) CheckResetBlock() error {
	// null
	return nil
}

func (p *I82540PHY) Commit() error {
	return PHYSWReset(p.hw)
}

func (p *I82540PHY) ForceSpeedDuplex() error {
	return PHYForceSpeedDuplexM88(p.hw)
}

func (p *I82540PHY) GetCableLength() error {
	return GetCableLengthM88(p.hw)
}

func (p *I82540PHY) GetCfgDone() error {
	return GetCfgDone(p.hw)
}

func (p *I82540PHY) GetInfo() error {
	return GetPHYInfoM88(p.hw)
}

func (p *I82540PHY) SetPage(val uint16) error {
	// null
	return nil
}

func (p *I82540PHY) ReadReg(offset uint32) (uint16, error) {
	return ReadPHYRegM88(p.hw, offset)
}

func (p *I82540PHY) ReadRegLocked(offset uint32) (uint16, error) {
	// null
	return 0, nil
}

func (p *I82540PHY) ReadRegPage(offset uint32) (uint16, error) {
	// null
	return 0, nil
}

func (p *I82540PHY) Release() {
	// null
}

func (p *I82540PHY) Reset() error {
	return PHYHWReset(p.hw)
}

func (p *I82540PHY) SetD0LpluState(e bool) error {
	// null
	return nil
}

func (p *I82540PHY) SetD3LpluState(e bool) error {
	// null
	return nil
}

func (p *I82540PHY) WriteReg(offset uint32, val uint16) error {
	return WritePHYRegM88(p.hw, offset, val)
}

func (p *I82540PHY) WriteRegLocked(offset uint32, val uint16) error {
	// null
	return nil
}

func (p *I82540PHY) WriteRegPage(offset uint32, val uint16) error {
	// null
	return nil
}

func (p *I82540PHY) PowerUp() {
	PowerUpPHYCopper(p.hw)
}

func (p *I82540PHY) PowerDown() {
	x := p.hw.RegRead(MANC)
	if x&MANC_SMBUS_EN != 0 {
		return
	}
	PowerDownPHYCopper(p.hw)
}

func (p *I82540PHY) ReadI2CByte(offset, addr byte) (byte, error) {
	// null
	return 0, nil
}

func (p *I82540PHY) WriteI2CByte(offset, addr, val byte) error {
	// null
	return nil
}

func (p *I82540PHY) CfgOnLinkUp() error {
	// null
	return nil
}
