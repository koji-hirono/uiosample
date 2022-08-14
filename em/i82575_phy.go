package em

import (
	"errors"
	"time"
)

type I82575PHY struct {
	hw *HW
}

func NewI82575PHY(hw *HW) *I82575PHY {
	p := new(I82575PHY)
	p.hw = hw
	return p
}

func (p *I82575PHY) InitParams() error {
	hw := p.hw
	phy := &hw.PHY

	if phy.MediaType != MediaTypeCopper {
		phy.PHYType = PHYTypeNone
		return nil
	}

	phy.AutonegMask = AUTONEG_ADVERTISE_SPEED_DEFAULT
	phy.ResetDelayUS = 100

	ctrl := hw.RegRead(CTRL_EXT)
	if SGMIIActive82575(hw) {
		ctrl |= CTRL_I2C_ENA
	} else {
		ctrl &^= CTRL_I2C_ENA
	}
	hw.RegWrite(CTRL_EXT, ctrl)
	ResetMDIConfig82580(hw)

	// Set phy->phy_addr and phy->id.
	GetPHYID82575(p.hw)

	// Verify phy id and set remaining function pointers
	switch phy.ID {
	case M88E1543_E_PHY_ID, M88E1512_E_PHY_ID, I347AT4_E_PHY_ID, M88E1112_E_PHY_ID, M88E1340M_E_PHY_ID:
		phy.PHYType = PHYTypeM88
	case M88E1111_I_PHY_ID:
		phy.PHYType = PHYTypeM88
	case IGP03E1000_E_PHY_ID, IGP04E1000_E_PHY_ID:
		phy.PHYType = PHYTypeIgp3
	case I82580_I_PHY_ID, I350_I_PHY_ID:
		phy.PHYType = PHYType82580
	case I210_I_PHY_ID:
		phy.PHYType = PHYTypeI210
	case BCM54616_E_PHY_ID:
		phy.PHYType = PHYTypeNone
	default:
		return errors.New("not support")
	}

	// Check if this PHY is configured for media swap.
	switch phy.ID {
	case M88E1112_E_PHY_ID:
		err := p.WriteReg(M88E1112_PAGE_ADDR, 2)
		if err != nil {
			return err
		}
		data, err := p.ReadReg(M88E1112_MAC_CTRL_1)
		if err != nil {
			return err
		}
		data &= M88E1112_MAC_CTRL_1_MODE_MASK
		data >>= M88E1112_MAC_CTRL_1_MODE_SHIFT
		if data == M88E1112_AUTO_COPPER_SGMII || data == M88E1112_AUTO_COPPER_BASEX {
			// TODO
			// hw.MAC.Op.CheckForLink = CheckForLinkMediaSwap
		}
	case M88E1512_E_PHY_ID:
		return InitM88E1512PHY(hw)
	case M88E1543_E_PHY_ID:
		return InitM88E1543PHY(hw)
	}

	return nil
}

func (p *I82575PHY) Acquire() error {
	return AcquirePHYBase(p.hw)
}

func (p *I82575PHY) CheckPolarity() error {
	phy := &p.hw.PHY
	// TODO
	switch phy.ID {
	case M88E1543_E_PHY_ID, M88E1512_E_PHY_ID, I347AT4_E_PHY_ID, M88E1112_E_PHY_ID, M88E1340M_E_PHY_ID:
		return CheckPolarityM88(p.hw)
	case M88E1111_I_PHY_ID:
		return CheckPolarityM88(p.hw)
	case IGP03E1000_E_PHY_ID, IGP04E1000_E_PHY_ID:
		return CheckPolarityIGP(p.hw)
	case I82580_I_PHY_ID, I350_I_PHY_ID:
		return CheckPolarity82577(p.hw)
	case I210_I_PHY_ID:
		return CheckPolarityM88(p.hw)
	default:
		return nil
	}
}

func (p *I82575PHY) CheckResetBlock() error {
	return CheckResetBlock(p.hw)
}

func (p *I82575PHY) Commit() error {
	return PHYSWReset(p.hw)
}

func (p *I82575PHY) ForceSpeedDuplex() error {
	phy := &p.hw.PHY
	// TODO
	switch phy.ID {
	case M88E1543_E_PHY_ID, M88E1512_E_PHY_ID, I347AT4_E_PHY_ID, M88E1112_E_PHY_ID, M88E1340M_E_PHY_ID:
		return PHYForceSpeedDuplexM88(p.hw)
	case M88E1111_I_PHY_ID:
		return PHYForceSpeedDuplexM88(p.hw)
	case IGP03E1000_E_PHY_ID, IGP04E1000_E_PHY_ID:
		return PHYForceSpeedDuplexIGP(p.hw)
	case I82580_I_PHY_ID, I350_I_PHY_ID:
		return PHYForceSpeedDuplex82577(p.hw)
	case I210_I_PHY_ID:
		return PHYForceSpeedDuplexM88(p.hw)
	default:
		return nil
	}
}

func (p *I82575PHY) GetCableLength() error {
	phy := &p.hw.PHY
	// TODO
	switch phy.ID {
	case M88E1543_E_PHY_ID, M88E1512_E_PHY_ID, I347AT4_E_PHY_ID, M88E1112_E_PHY_ID, M88E1340M_E_PHY_ID:
		return GetCableLengthM88gen2(p.hw)
	case M88E1111_I_PHY_ID:
		return GetCableLengthM88(p.hw)
	case IGP03E1000_E_PHY_ID, IGP04E1000_E_PHY_ID:
		return GetCableLengthIGP2(p.hw)
	case I82580_I_PHY_ID, I350_I_PHY_ID:
		return GetCableLength82577(p.hw)
	case I210_I_PHY_ID:
		return GetCableLengthM88(p.hw)
	default:
		return nil
	}
}

func (p *I82575PHY) GetCfgDone() error {
	hw := p.hw
	var mask uint32
	switch hw.Bus.Func {
	case 1:
		mask = NVM_CFG_DONE_PORT_1
	case 2:
		mask = NVM_CFG_DONE_PORT_2
	case 3:
		mask = NVM_CFG_DONE_PORT_3
	default:
		mask = NVM_CFG_DONE_PORT_0
	}
	timeout := PHY_CFG_TIMEOUT
	for timeout > 0 {
		if hw.RegRead(EEMNGCTL)&mask != 0 {
			break
		}
		time.Sleep(1 * time.Millisecond)
		timeout--
	}

	// If EEPROM is not marked present, init the PHY manually
	if hw.RegRead(EECD)&EECD_PRES == 0 && hw.PHY.PHYType == PHYTypeIgp3 {
		PHYInitScriptIGP3(hw)
	}
	return nil
}

func (p *I82575PHY) GetInfo() error {
	phy := &p.hw.PHY
	// TODO
	switch phy.ID {
	case M88E1543_E_PHY_ID, M88E1512_E_PHY_ID, I347AT4_E_PHY_ID, M88E1112_E_PHY_ID, M88E1340M_E_PHY_ID:
		return GetPHYInfoM88(p.hw)
	case M88E1111_I_PHY_ID:
		return GetPHYInfoM88(p.hw)
	case IGP03E1000_E_PHY_ID, IGP04E1000_E_PHY_ID:
		return GetPHYInfoIGP(p.hw)
	case I82580_I_PHY_ID, I350_I_PHY_ID:
		return GetPHYInfo82577(p.hw)
	case I210_I_PHY_ID:
		return GetPHYInfoM88(p.hw)
	default:
		return nil
	}
}

func (p *I82575PHY) SetPage(val uint16) error {
	// null
	return nil
}

func (p *I82575PHY) ReadReg(offset uint32) (uint16, error) {
	// TODO
	if SGMIIActive82575(p.hw) && !SGMIIUsesMDIO82575(p.hw) {
		return ReadPHYRegSGMII82575(p.hw, offset)
	}
	switch p.hw.MAC.Type {
	case MACType82580, MACTypeI350, MACTypeI354:
		return ReadPHYReg82580(p.hw, offset)
	case MACTypeI210, MACTypeI211:
		return ReadPHYRegGS40G(p.hw, offset)
	default:
		return ReadPHYRegIGP(p.hw, offset)
	}
}

func (p *I82575PHY) ReadRegLocked(offset uint32) (uint16, error) {
	// null
	return 0, nil
}

func (p *I82575PHY) ReadRegPage(offset uint32) (uint16, error) {
	// null
	return 0, nil
}

func (p *I82575PHY) Release() {
	ReleasePHYBase(p.hw)
}

func (p *I82575PHY) Reset() error {
	if SGMIIActive82575(p.hw) {
		return PHYHWResetSGMII82575(p.hw)
	} else {
		return PHYHWReset(p.hw)
	}
}

func (p *I82575PHY) SetD0LpluState(active bool) error {
	phy := &p.hw.PHY
	// TODO
	switch phy.ID {
	case M88E1543_E_PHY_ID, M88E1512_E_PHY_ID, I347AT4_E_PHY_ID, M88E1112_E_PHY_ID, M88E1340M_E_PHY_ID:
		return nil
	case M88E1111_I_PHY_ID:
		return nil
	case IGP03E1000_E_PHY_ID, IGP04E1000_E_PHY_ID:
		return SetD0LpluState82575(p.hw, active)
	case I82580_I_PHY_ID, I350_I_PHY_ID:
		return SetD0LpluState82580(p.hw, active)
	case I210_I_PHY_ID:
		return SetD0LpluState82580(p.hw, active)
	default:
		return nil
	}
}

func (p *I82575PHY) SetD3LpluState(active bool) error {
	phy := &p.hw.PHY
	// TODO
	switch phy.ID {
	case M88E1543_E_PHY_ID, M88E1512_E_PHY_ID, I347AT4_E_PHY_ID, M88E1112_E_PHY_ID, M88E1340M_E_PHY_ID:
		return nil
	case M88E1111_I_PHY_ID:
		return nil
	case IGP03E1000_E_PHY_ID, IGP04E1000_E_PHY_ID:
		return SetD3LpluState(p.hw, active)
	case I82580_I_PHY_ID, I350_I_PHY_ID:
		return SetD3LpluState82580(p.hw, active)
	case I210_I_PHY_ID:
		return SetD3LpluState82580(p.hw, active)
	default:
		return nil
	}
}

func (p *I82575PHY) WriteReg(offset uint32, val uint16) error {
	// TODO
	if SGMIIActive82575(p.hw) && !SGMIIUsesMDIO82575(p.hw) {
		return WritePHYRegSGMII82575(p.hw, offset, val)
	}
	switch p.hw.MAC.Type {
	case MACType82580, MACTypeI350, MACTypeI354:
		return WritePHYReg82580(p.hw, offset, val)
	case MACTypeI210, MACTypeI211:
		return WritePHYRegGS40G(p.hw, offset, val)
	default:
		return WritePHYRegIGP(p.hw, offset, val)
	}
}

func (p *I82575PHY) WriteRegLocked(offset uint32, val uint16) error {
	// null
	return nil
}

func (p *I82575PHY) WriteRegPage(offset uint32, val uint16) error {
	// null
	return nil
}

func (p *I82575PHY) PowerUp() {
	PowerUpPHYCopper(p.hw)
}

func (p *I82575PHY) PowerDown() {
	PowerDownPHYCopperBase(p.hw)
}

func (p *I82575PHY) ReadI2CByte(offset, addr byte) (byte, error) {
	return ReadI2CByte(p.hw, offset, addr)
}

func (p *I82575PHY) WriteI2CByte(offset, addr, val byte) error {
	return WriteI2CByte(p.hw, offset, addr, val)
}

func (p *I82575PHY) CfgOnLinkUp() error {
	// null
	return nil
}

func GetPHYID82575(hw *HW) error {
	phy := &hw.PHY
	// some i354 devices need an extra read for phy id
	if hw.MAC.Type == MACTypeI354 {
		GetPHYID(hw)
	}

	// For SGMII PHYs, we try the list of possible addresses until
	// we find one that works.  For non-SGMII PHYs
	// (e.g. integrated copper PHYs), an address of 1 should
	// work.  The result of this function should mean phy->phy_addr
	// and phy->id are set correctly.
	if !SGMIIActive82575(hw) {
		phy.Addr = 1
		return GetPHYID(hw)
	}
	if SGMIIUsesMDIO82575(hw) {
		switch hw.MAC.Type {
		case MACType82575, MACType82576:
			mdic := hw.RegRead(MDIC)
			mdic &= MDIC_PHY_MASK
			phy.Addr = mdic >> MDIC_PHY_SHIFT
		case MACType82580:
		case MACTypeI350:
		case MACTypeI354:
		case MACTypeI210:
		case MACTypeI211:
			mdic := hw.RegRead(MDICNFG)
			mdic &= MDICNFG_PHY_MASK
			phy.Addr = mdic >> MDICNFG_PHY_SHIFT
		default:
			return errors.New("not support")
		}
		return GetPHYID(hw)
	}

	// Power on sgmii phy if it is disabled
	ctrl := hw.RegRead(CTRL_EXT)
	hw.RegWrite(CTRL_EXT, ctrl&^CTRL_EXT_SDP3_DATA)
	hw.RegWriteFlush()
	time.Sleep(300 * time.Millisecond)
	// restore previous sfp cage power state
	defer hw.RegWrite(CTRL_EXT, ctrl)

	// The address field in the I2CCMD register is 3 bits and 0 is invalid.
	// Therefore, we need to test 1-7
	for phy.Addr = 1; phy.Addr < 8; phy.Addr++ {
		phyid, err := ReadPHYRegSGMII82575(hw, PHY_ID1)
		if err == nil {
			// At the time of this writing, The M88 part is
			// the only supported SGMII PHY product.
			if phyid == M88_VENDOR {
				break
			}
		}
	}

	// A valid PHY type couldn't be found.
	if phy.Addr == 8 {
		phy.Addr = 0
		return errors.New("not found")
	}
	return GetPHYID(hw)
}

func ReadPHYRegSGMII82575(hw *HW, offset uint32) (uint16, error) {
	phy := &hw.PHY
	err := phy.Op.Acquire()
	if err != nil {
		return 0, err
	}
	defer phy.Op.Release()
	return ReadPHYRegI2C(hw, offset)
}

func WritePHYRegSGMII82575(hw *HW, offset uint32, data uint16) error {
	phy := &hw.PHY
	err := phy.Op.Acquire()
	if err != nil {
		return err
	}
	defer phy.Op.Release()
	return WritePHYRegI2C(hw, offset, data)
}
