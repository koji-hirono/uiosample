package em82540

import (
	"errors"

	"uiosample/em"
)

type PHY struct {
	hw *em.HW
}

func NewPHY(hw *em.HW) *PHY {
	p := new(PHY)
	p.hw = hw
	return p
}

func (p *PHY) InitParams() error {
	phy := &p.hw.PHY

	phy.Addr = 1
	phy.AutonegMask = em.AUTONEG_ADVERTISE_SPEED_DEFAULT
	phy.ResetDelayUS = 10000
	phy.PHYType = em.PHYTypeM88

	err := em.GetPHYID(p.hw)
	if err != nil {
		return err
	}

	switch p.hw.MAC.Type {
	case em.MACType82540:
	case em.MACType82545:
	case em.MACType82545Rev3:
	case em.MACType82546:
	case em.MACType82546Rev3:
	default:
		return errors.New("Invalid MAC type")
	}

	if phy.ID != em.M88E1011_I_PHY_ID {
		return errors.New("Invalid phy id")
	}

	return nil
}

func (p *PHY) Acquire() error {
	// null
	return nil
}

func (p *PHY) CheckPolarity() error {
	return em.CheckPolarityM88(p.hw)
}

func (p *PHY) CheckResetBlock() error {
	// null
	return nil
}

func (p *PHY) Commit() error {
	return em.PHYSWReset(p.hw)
}

func (p *PHY) ForceSpeedDuplex() error {
	return em.PHYForceSpeedDuplexM88(p.hw)
}

func (p *PHY) GetCableLength() error {
	return em.GetCableLengthM88(p.hw)
}

func (p *PHY) GetCfgDone() error {
	return em.GetCfgDone(p.hw)
}

func (p *PHY) GetInfo() error {
	return em.GetPHYInfoM88(p.hw)
}

func (p *PHY) SetPage(val uint16) error {
	// null
	return nil
}

func (p *PHY) ReadReg(offset uint32) (uint16, error) {
	return em.ReadPHYRegM88(p.hw, offset)
}

func (p *PHY) ReadRegLocked(offset uint32) (uint16, error) {
	// null
	return 0, nil
}

func (p *PHY) ReadRegPage(offset uint32) (uint16, error) {
	// null
	return 0, nil
}

func (p *PHY) Release() {
	// null
}

func (p *PHY) Reset() error {
	return em.PHYHWReset(p.hw)
}

func (p *PHY) SetD0LpluState(e bool) error {
	// null
	return nil
}

func (p *PHY) SetD3LpluState(e bool) error {
	// null
	return nil
}

func (p *PHY) WriteReg(offset uint32, val uint16) error {
	return em.WritePHYRegM88(p.hw, offset, val)
}

func (p *PHY) WriteRegLocked(offset uint32, val uint16) error {
	// null
	return nil
}

func (p *PHY) WriteRegPage(offset uint32, val uint16) error {
	// null
	return nil
}

func (p *PHY) PowerUp() {
	em.PowerUpPHYCopper(p.hw)
}

func (p *PHY) PowerDown() {
	x := p.hw.RegRead(em.MANC)
	if x&em.MANC_SMBUS_EN != 0 {
		return
	}
	em.PowerDownPHYCopper(p.hw)
}

func (p *PHY) ReadI2CByte(offset, addr byte) (byte, error) {
	// null
	return 0, nil
}

func (p *PHY) WriteI2CByte(offset, addr, val byte) error {
	// null
	return nil
}

func (p *PHY) CfgOnLinkUp() error {
	// null
	return nil
}
