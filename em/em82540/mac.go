package em82540

import (
	"uiosample/em"
)

type MAC struct {
	hw  *em.HW
	nvm *NVM
	phy *PHY
}

func NewMAC(hw *em.HW, nvm *NVM, phy *PHY) *MAC {
	m := new(MAC)
	m.hw = hw
	m.nvm = nvm
	m.phy = phy
	return m
}

func (m *MAC) InitParams() error {
	mac := &m.hw.MAC

	switch m.hw.DeviceID {
	case em.DEV_ID_82545EM_FIBER, em.DEV_ID_82545GM_FIBER,
		em.DEV_ID_82546EB_FIBER, em.DEV_ID_82546GB_FIBER:
		m.hw.PHY.MediaType = em.MediaTypeFiber
	case em.DEV_ID_82545GM_SERDES, em.DEV_ID_82546GB_SERDES:
		m.hw.PHY.MediaType = em.MediaTypeInternalSerdes
	default:
		m.hw.PHY.MediaType = em.MediaTypeCopper
	}

	mac.MTARegCount = 128
	mac.RAREntryCount = em.RAR_ENTRIES
	return nil
}

func (m *MAC) IDLEDInit() error {
	return em.IDLEDInit(m.hw)
}

func (m *MAC) BlinkLED() error {
	// null
	return nil
}

func (m *MAC) CheckMngMode() bool {
	// null
	return false
}

func (m *MAC) CheckForLink() error {
	switch m.hw.PHY.MediaType {
	case em.MediaTypeFiber:
		return em.CheckForFiberLink(m.hw)
	case em.MediaTypeInternalSerdes:
		return em.CheckForSerdesLink(m.hw)
	case em.MediaTypeCopper:
		return em.CheckForCopperLink(m.hw)
	}
	return nil
}

func (m *MAC) CleanupLED() error {
	return em.CleanupLED(m.hw)
}

func (m *MAC) ClearHWCounters() {
	em.ClearHWCounters(m.hw)

	m.hw.RegRead(em.PRC64)
	m.hw.RegRead(em.PRC127)
	m.hw.RegRead(em.PRC255)
	m.hw.RegRead(em.PRC511)
	m.hw.RegRead(em.PRC1023)
	m.hw.RegRead(em.PRC1522)
	m.hw.RegRead(em.PTC64)
	m.hw.RegRead(em.PTC127)
	m.hw.RegRead(em.PTC255)
	m.hw.RegRead(em.PTC511)
	m.hw.RegRead(em.PTC1023)
	m.hw.RegRead(em.PTC1522)

	m.hw.RegRead(em.ALGNERRC)
	m.hw.RegRead(em.RXERRC)
	m.hw.RegRead(em.TNCRS)
	m.hw.RegRead(em.CEXTERR)
	m.hw.RegRead(em.TSCTC)
	m.hw.RegRead(em.TSCTFC)

	m.hw.RegRead(em.MGTPRC)
	m.hw.RegRead(em.MGTPDC)
	m.hw.RegRead(em.MGTPTC)
}

func (m *MAC) ClearVFTA() {
	em.ClearVFTA(m.hw)
}

func (m *MAC) GetBusInfo() error {
	return em.GetBusInfoPCI(m.hw)
}

func (m *MAC) SetLANID() {
	em.SetLANIDMultiPortPCI(m.hw)
}

func (m *MAC) GetLinkUpInfo() (uint16, uint16, error) {
	switch m.hw.PHY.MediaType {
	case em.MediaTypeCopper:
		return em.GetSpeedAndDuplexCopper(m.hw)
	default:
		return em.GetSpeedAndDuplexFiberSerdes(m.hw)
	}
}

func (m *MAC) LEDOn() error {
	return em.LEDOn(m.hw)
}

func (m *MAC) LEDOff() error {
	return em.LEDOff(m.hw)
}

func (m *MAC) UpdateMCAddrList(addrs [][6]byte) {
	em.UpdateMCAddrList(m.hw, addrs)
}

func (m *MAC) ResetHW() error {
	m.hw.RegWrite(em.IMC, ^uint32(0))

	m.hw.RegWrite(em.RCTL, 0)
	m.hw.RegWrite(em.TCTL, em.TCTL_PSP)
	m.hw.RegWriteFlush()

	// msec_delay(10)

	ctrl := m.hw.RegRead(em.CTRL)

	switch m.hw.MAC.Type {
	case em.MACType82545Rev3, em.MACType82546Rev3:
		m.hw.RegWrite(em.CTRL_DUP, ctrl|em.CTRL_RST)
	default:
		m.hw.RegWrite(em.CTRL, ctrl|em.CTRL_RST)
	}

	// msec_delay(5)

	manc := m.hw.RegRead(em.MANC)
	manc &^= em.MANC_ARP_EN
	m.hw.RegWrite(em.MANC, manc)

	m.hw.RegWrite(em.IMC, ^uint32(0))
	m.hw.RegRead(em.ICR)

	return nil
}

func (m *MAC) InitHW() error {
	err := m.IDLEDInit()
	if err != nil {
		return err
	}

	if m.hw.MAC.Type < em.MACType82545Rev3 {
		m.hw.RegWrite(em.VET, 0)
	}

	m.ClearVFTA()

	// e1000_init_rx_addrs_generic(hw, mac->rar_entry_count)

	for i := 0; i < int(m.hw.MAC.MTARegCount); i++ {
		m.hw.RegWrite(em.MTA+(i<<2), 0)
		m.hw.RegWriteFlush()
	}

	if m.hw.MAC.Type < em.MACType82545Rev3 {
		// e1000_pcix_mmrbc_workaround_generic(hw)
	}

	m.SetupLink()

	txdctl := m.hw.RegRead(em.TXDCTL(0))
	txdctl &^= em.TXDCTL_WTHRESH
	txdctl |= em.TXDCTL_FULL_TX_DESC_WB
	m.hw.RegWrite(em.TXDCTL(0), txdctl)

	m.ClearHWCounters()

	switch m.hw.DeviceID {
	case em.DEV_ID_82546GB_QUAD_COPPER, em.DEV_ID_82546GB_QUAD_COPPER_KSP3:
		x := m.hw.RegRead(em.CTRL_EXT)
		x |= em.CTRL_EXT_RO_DIS
		m.hw.RegWrite(em.CTRL_EXT, x)
	}

	return nil
}

func (m *MAC) ShutdownSerdes() {
}

func (m *MAC) PowerUpSerdes() {
}

func (m *MAC) SetupLink() error {
	return em.SetupLink(m.hw)
}

func (m *MAC) SetupPhysicalInterface() error {
	switch m.hw.PHY.MediaType {
	case em.MediaTypeCopper:
		return m.setupCopperLink()
	default:
		return m.setupFiberSerdesLink()
	}
}

func (m *MAC) SetupLED() error {
	return em.SetupLED(m.hw)
}

func (m *MAC) WriteVFTA(offset, val uint32) {
	em.WriteVFTA(m.hw, offset, val)
}

func (m *MAC) ConfigCollisionDist() {
	em.ConfigCollisionDist(m.hw)
}

func (m *MAC) SetRAR(addr [6]byte, index int) error {
	return em.SetRAR(m.hw, addr, index)
}

func (m *MAC) ReadMACAddr() error {
	for i := 0; i < 6; i += 2 {
		offset := uint16(i >> 1)
		var x [1]uint16
		err := m.nvm.Read(offset, x[:])
		if err != nil {
			return err
		}
		m.hw.MAC.PermAddr[i] = byte(x[0])
		m.hw.MAC.PermAddr[i+1] = byte(x[0] >> 8)
	}
	if m.hw.Bus.Func == 1 {
		m.hw.MAC.PermAddr[5] ^= 1
	}

	copy(m.hw.MAC.Addr[:], m.hw.MAC.PermAddr[:])
	return nil
}

func (m *MAC) ValidateMDISetting() error {
	return em.ValidateMDISetting(m.hw)
}

func (m *MAC) SetOBFFTimer(uint32) error {
	return nil
}

func (m *MAC) AcquireSWFWSync(uint16) error {
	return nil
}

func (m *MAC) ReleaseSWFWSync(uint16) {
}

func (m *MAC) setupCopperLink() error {
	ctrl := m.hw.RegRead(em.CTRL)
	ctrl |= em.CTRL_SLU
	ctrl &^= em.CTRL_FRCSPD | em.CTRL_FRCDPX
	m.hw.RegWrite(em.CTRL, ctrl)

	err := m.setPHYMode()
	if err != nil {
		return err
	}

	switch m.hw.MAC.Type {
	case em.MACType82545Rev3, em.MACType82546Rev3:
		x, err := m.phy.ReadReg(em.M88E1000_PHY_SPEC_CTRL)
		if err != nil {
			return err
		}
		x |= 0x00000008
		err = m.phy.WriteReg(em.M88E1000_PHY_SPEC_CTRL, x)
		if err != nil {
			return err
		}
	}

	err = em.CopperLinkSetupM88(m.hw)
	if err != nil {
		return err
	}

	return em.SetupCopperLink(m.hw)
}

func (m *MAC) setPHYMode() error {
	if m.hw.MAC.Type != em.MACType82545Rev3 {
		return nil
	}

	var x [1]uint16
	err := m.nvm.Read(em.NVM_PHY_CLASS_WORD, x[:])
	if err != nil {
		return err
	}
	if x[0] == em.NVM_RESERVED_WORD {
		return nil
	}
	if x[0]&em.NVM_PHY_CLASS_A == 0 {
		return nil
	}

	err = m.phy.WriteReg(em.M88E1000_PHY_PAGE_SELECT, 0x000b)
	if err != nil {
		return err
	}
	return m.phy.WriteReg(em.M88E1000_PHY_GEN_CONTROL, 0x8104)
}

func (m *MAC) setupFiberSerdesLink() error {
	switch m.hw.MAC.Type {
	case em.MACType82545Rev3, em.MACType82546Rev3:
		if m.hw.PHY.MediaType == em.MediaTypeInternalSerdes {
			err := m.adjustSerdesAmplitude()
			if err != nil {
				return err
			}
		}
		err := m.setVCOSpeed()
		if err != nil {
			return err
		}
	}

	return em.SetupFiberSerdesLink(m.hw)
}

func (m *MAC) adjustSerdesAmplitude() error {
	var x [1]uint16
	err := m.nvm.Read(em.NVM_SERDES_AMPLITUDE, x[:])
	if err != nil {
		return err
	}
	if x[0] != em.NVM_RESERVED_WORD {
		x[0] &= em.NVM_SERDES_AMPLITUDE_MASK
		return m.phy.WriteReg(em.M88E1000_PHY_EXT_CTRL, x[0])
	}
	return nil
}

func (m *MAC) setVCOSpeed() error {
	// Set PHY register 30, page 5, bit 8 to 0
	defaultPage, err := m.phy.ReadReg(em.M88E1000_PHY_PAGE_SELECT)
	if err != nil {
		return err
	}

	err = m.phy.WriteReg(em.M88E1000_PHY_PAGE_SELECT, 0x0005)
	if err != nil {
		return err
	}

	x, err := m.phy.ReadReg(em.M88E1000_PHY_GEN_CONTROL)
	if err != nil {
		return err
	}

	x &^= em.M88E1000_PHY_VCO_REG_BIT8
	err = m.phy.WriteReg(em.M88E1000_PHY_GEN_CONTROL, x)
	if err != nil {
		return err
	}

	// Set PHY register 30, page 4, bit 11 to 1
	err = m.phy.WriteReg(em.M88E1000_PHY_PAGE_SELECT, 0x0004)
	if err != nil {
		return err
	}

	x, err = m.phy.ReadReg(em.M88E1000_PHY_GEN_CONTROL)
	if err != nil {
		return err
	}

	x |= em.M88E1000_PHY_VCO_REG_BIT11
	err = m.phy.WriteReg(em.M88E1000_PHY_GEN_CONTROL, x)
	if err != nil {
		return err
	}

	return m.phy.WriteReg(em.M88E1000_PHY_PAGE_SELECT, defaultPage)
}
