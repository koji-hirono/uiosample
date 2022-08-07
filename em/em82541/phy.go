package em82541

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

func (p *PHY) InitParam() error {
	phy := &p.hw.PHY

	phy.Addr = 1
	phy.AutonegMask = em.AUTONEG_ADVERTISE_SPEED_DEFAULT
	phy.ResetDelayUS = 10000
	phy.PHYType = em.PHYTypeIgp

	// err := e1000_get_phy_id(hw)
	// if err != nil {
	//	return err
	// }

	if phy.ID != em.IGP01E1000_I_PHY_ID {
		return errors.New("Invalid phy id")
	}

	return nil
}

func (p *PHY) Aquire() error {
	// e1000_null_ops_generic
	return nil
}

func (p *PHY) CheckPolarity() error {
	// e1000_check_polarity_igp
	return nil
}

func (p *PHY) CheckResetBlock() error {
	// e1000_null_ops_generic
	return nil
}

func (p *PHY) Commit() error {
	// e1000_null_ops_generic
	return nil
}

func (p *PHY) ForceSpeedDuplex() error {
	// e1000_phy_force_speed_duplex_igp
	return nil
}

func (p *PHY) GetCableLength() error {
	// e1000_get_cable_length_igp_82541
	return nil
}

func (p *PHY) GetCfgDone() error {
	// e1000_get_cfg_done_generic
	return nil
}

func (p *PHY) GetInfo() error {
	// e1000_get_phy_info_igp
	return nil
}

func (p *PHY) SetPage(val uint16) error {
	// e1000_null_set_page
	return nil
}

func (p *PHY) ReadReg(offset uint32) (uint16, error) {
	// e1000_read_phy_reg_igp
	return 0, nil
}

func (p *PHY) ReadRegLocked(offset uint32) (uint16, error) {
	// e1000_null_read_reg
	return 0, nil
}

func (p *PHY) ReadRegPage(offset uint32) (uint16, error) {
	// e1000_null_read_reg
	return 0, nil
}

func (p *PHY) Release() {
	// e1000_null_phy_generic
}

func (p *PHY) Reset() error {
	// e1000_phy_hw_reset_82541

	// err := e1000_phy_hw_reset_generic(hw)
	// if err != nil {
	//	return err
	// }

	// e1000_phy_init_script_82541(hw)

	switch p.hw.MAC.Type {
	case em.MACType82541, em.MACType82547:
		// Configure activity LED after PHY reset
		x := p.hw.RegRead(em.LEDCTL)
		x &= IGP_ACTIVITY_LED_MASK
		x |= IGP_ACTIVITY_LED_ENABLE | IGP_LED3_MODE
		p.hw.RegWrite(em.LEDCTL, x)
	}
	return nil
}

func (p *PHY) SetD0LpluState(e bool) error {
	// e1000_null_lplu_state
	return nil
}

func (p *PHY) SetD3LpluState(e bool) error {
	// e1000_set_d3_lplu_state_82541
	return nil
}

func (p *PHY) WriteReg(offset uint32, val uint16) error {
	// e1000_write_phy_reg_igp
	return nil
}

func (p *PHY) WriteRegLocked(offset uint32, val uint16) error {
	// e1000_null_write_reg
	return nil
}

func (p *PHY) WriteRegPage(offset uint32, val uint16) error {
	// e1000_null_write_reg
	return nil
}

func (p *PHY) PowerUp() {
	// e1000_power_up_phy_copper
}

func (p *PHY) PowerDown() {
	// e1000_power_down_phy_copper_82541
}

func (p *PHY) ReadI2CByte(offset, addr byte) (byte, error) {
	// e1000_read_i2c_byte_null
	return 0, nil
}

func (p *PHY) WriteI2CByte(offset, addr, val byte) error {
	// e1000_write_i2c_byte_null
	return nil
}

func (p *PHY) CfgOnLinkUp() error {
	// e1000_null_ops_generic
	return nil
}
