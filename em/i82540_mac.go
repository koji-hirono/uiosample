package em

import (
	"time"
)

type I82540MAC struct {
	hw  *HW
	nvm *I82540NVM
	phy *I82540PHY
}

func NewI82540MAC(hw *HW, nvm *I82540NVM, phy *I82540PHY) *I82540MAC {
	m := new(I82540MAC)
	m.hw = hw
	m.nvm = nvm
	m.phy = phy
	return m
}

func (m *I82540MAC) InitParams() error {
	mac := &m.hw.MAC

	switch m.hw.DeviceID {
	case DEV_ID_82545EM_FIBER, DEV_ID_82545GM_FIBER,
		DEV_ID_82546EB_FIBER, DEV_ID_82546GB_FIBER:
		m.hw.PHY.MediaType = MediaTypeFiber
	case DEV_ID_82545GM_SERDES, DEV_ID_82546GB_SERDES:
		m.hw.PHY.MediaType = MediaTypeInternalSerdes
	default:
		m.hw.PHY.MediaType = MediaTypeCopper
	}

	mac.MTARegCount = 128
	mac.RAREntryCount = RAR_ENTRIES
	return nil
}

func (m *I82540MAC) IDLEDInit() error {
	return IDLEDInit(m.hw)
}

func (m *I82540MAC) BlinkLED() error {
	// null
	return nil
}

func (m *I82540MAC) CheckMngMode() bool {
	// null
	return false
}

func (m *I82540MAC) CheckForLink() error {
	switch m.hw.PHY.MediaType {
	case MediaTypeFiber:
		return CheckForFiberLink(m.hw)
	case MediaTypeInternalSerdes:
		return CheckForSerdesLink(m.hw)
	case MediaTypeCopper:
		return CheckForCopperLink(m.hw)
	}
	return nil
}

func (m *I82540MAC) CleanupLED() error {
	return CleanupLED(m.hw)
}

func (m *I82540MAC) ClearHWCounters() {
	ClearHWCounters(m.hw)

	m.hw.RegRead(PRC64)
	m.hw.RegRead(PRC127)
	m.hw.RegRead(PRC255)
	m.hw.RegRead(PRC511)
	m.hw.RegRead(PRC1023)
	m.hw.RegRead(PRC1522)
	m.hw.RegRead(PTC64)
	m.hw.RegRead(PTC127)
	m.hw.RegRead(PTC255)
	m.hw.RegRead(PTC511)
	m.hw.RegRead(PTC1023)
	m.hw.RegRead(PTC1522)

	m.hw.RegRead(ALGNERRC)
	m.hw.RegRead(RXERRC)
	m.hw.RegRead(TNCRS)
	m.hw.RegRead(CEXTERR)
	m.hw.RegRead(TSCTC)
	m.hw.RegRead(TSCTFC)

	m.hw.RegRead(MGTPRC)
	m.hw.RegRead(MGTPDC)
	m.hw.RegRead(MGTPTC)
}

func (m *I82540MAC) ClearVFTA() {
	ClearVFTA(m.hw)
}

func (m *I82540MAC) GetBusInfo() error {
	return GetBusInfoPCI(m.hw)
}

func (m *I82540MAC) SetLANID() {
	SetLANIDMultiPortPCI(m.hw)
}

func (m *I82540MAC) GetLinkUpInfo() (uint16, uint16, error) {
	switch m.hw.PHY.MediaType {
	case MediaTypeCopper:
		return GetSpeedAndDuplexCopper(m.hw)
	default:
		return GetSpeedAndDuplexFiberSerdes(m.hw)
	}
}

func (m *I82540MAC) LEDOn() error {
	return LEDOn(m.hw)
}

func (m *I82540MAC) LEDOff() error {
	return LEDOff(m.hw)
}

func (m *I82540MAC) UpdateMCAddrList(addrs [][6]byte) {
	UpdateMCAddrList(m.hw, addrs)
}

func (m *I82540MAC) ResetHW() error {
	m.hw.RegWrite(IMC, ^uint32(0))

	m.hw.RegWrite(RCTL, 0)
	m.hw.RegWrite(TCTL, TCTL_PSP)
	m.hw.RegWriteFlush()

	time.Sleep(10 * time.Millisecond)

	ctrl := m.hw.RegRead(CTRL)

	switch m.hw.MAC.Type {
	case MACType82545Rev3, MACType82546Rev3:
		m.hw.RegWrite(CTRL_DUP, ctrl|CTRL_RST)
	default:
		m.hw.RegWrite(CTRL, ctrl|CTRL_RST)
	}

	time.Sleep(5 * time.Millisecond)

	manc := m.hw.RegRead(MANC)
	manc &^= MANC_ARP_EN
	m.hw.RegWrite(MANC, manc)

	m.hw.RegWrite(IMC, ^uint32(0))
	m.hw.RegRead(ICR)

	return nil
}

func (m *I82540MAC) InitHW() error {
	err := m.IDLEDInit()
	if err != nil {
		return err
	}

	if m.hw.MAC.Type < MACType82545Rev3 {
		m.hw.RegWrite(VET, 0)
	}

	m.ClearVFTA()

	InitRxAddrs(m.hw, int(m.hw.MAC.RAREntryCount))

	for i := 0; i < int(m.hw.MAC.MTARegCount); i++ {
		m.hw.RegWrite(MTA+(i<<2), 0)
		m.hw.RegWriteFlush()
	}

	if m.hw.MAC.Type < MACType82545Rev3 {
		// e1000_pcix_mmrbc_workaround_generic(hw)
	}

	m.SetupLink()

	txdctl := m.hw.RegRead(TXDCTL(0))
	txdctl &^= TXDCTL_WTHRESH
	txdctl |= TXDCTL_FULL_TX_DESC_WB
	m.hw.RegWrite(TXDCTL(0), txdctl)

	m.ClearHWCounters()

	switch m.hw.DeviceID {
	case DEV_ID_82546GB_QUAD_COPPER, DEV_ID_82546GB_QUAD_COPPER_KSP3:
		x := m.hw.RegRead(CTRL_EXT)
		x |= CTRL_EXT_RO_DIS
		m.hw.RegWrite(CTRL_EXT, x)
	}

	return nil
}

func (m *I82540MAC) ShutdownSerdes() {
}

func (m *I82540MAC) PowerUpSerdes() {
}

func (m *I82540MAC) SetupLink() error {
	return SetupLink(m.hw)
}

func (m *I82540MAC) SetupPhysicalInterface() error {
	switch m.hw.PHY.MediaType {
	case MediaTypeCopper:
		return m.setupCopperLink()
	default:
		return m.setupFiberSerdesLink()
	}
}

func (m *I82540MAC) SetupLED() error {
	return SetupLED(m.hw)
}

func (m *I82540MAC) WriteVFTA(offset, val uint32) {
	WriteVFTA(m.hw, offset, val)
}

func (m *I82540MAC) ConfigCollisionDist() {
	ConfigCollisionDist(m.hw)
}

func (m *I82540MAC) SetRAR(addr [6]byte, index int) error {
	return SetRAR(m.hw, addr, index)
}

func (m *I82540MAC) ReadMACAddr() error {
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

func (m *I82540MAC) ValidateMDISetting() error {
	return ValidateMDISetting(m.hw)
}

func (m *I82540MAC) SetOBFFTimer(uint32) error {
	return nil
}

func (m *I82540MAC) AcquireSWFWSync(uint16) error {
	return nil
}

func (m *I82540MAC) ReleaseSWFWSync(uint16) {
}

func (m *I82540MAC) setupCopperLink() error {
	ctrl := m.hw.RegRead(CTRL)
	ctrl |= CTRL_SLU
	ctrl &^= CTRL_FRCSPD | CTRL_FRCDPX
	m.hw.RegWrite(CTRL, ctrl)

	err := m.setPHYMode()
	if err != nil {
		return err
	}

	switch m.hw.MAC.Type {
	case MACType82545Rev3, MACType82546Rev3:
		x, err := m.phy.ReadReg(M88E1000_PHY_SPEC_CTRL)
		if err != nil {
			return err
		}
		x |= 0x00000008
		err = m.phy.WriteReg(M88E1000_PHY_SPEC_CTRL, x)
		if err != nil {
			return err
		}
	}

	err = CopperLinkSetupM88(m.hw)
	if err != nil {
		return err
	}

	return SetupCopperLink(m.hw)
}

func (m *I82540MAC) setPHYMode() error {
	if m.hw.MAC.Type != MACType82545Rev3 {
		return nil
	}

	var x [1]uint16
	err := m.nvm.Read(NVM_PHY_CLASS_WORD, x[:])
	if err != nil {
		return err
	}
	if x[0] == NVM_RESERVED_WORD {
		return nil
	}
	if x[0]&NVM_PHY_CLASS_A == 0 {
		return nil
	}

	err = m.phy.WriteReg(M88E1000_PHY_PAGE_SELECT, 0x000b)
	if err != nil {
		return err
	}
	return m.phy.WriteReg(M88E1000_PHY_GEN_CONTROL, 0x8104)
}

func (m *I82540MAC) setupFiberSerdesLink() error {
	switch m.hw.MAC.Type {
	case MACType82545Rev3, MACType82546Rev3:
		if m.hw.PHY.MediaType == MediaTypeInternalSerdes {
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

	return SetupFiberSerdesLink(m.hw)
}

func (m *I82540MAC) adjustSerdesAmplitude() error {
	var x [1]uint16
	err := m.nvm.Read(NVM_SERDES_AMPLITUDE, x[:])
	if err != nil {
		return err
	}
	if x[0] != NVM_RESERVED_WORD {
		x[0] &= NVM_SERDES_AMPLITUDE_MASK
		return m.phy.WriteReg(M88E1000_PHY_EXT_CTRL, x[0])
	}
	return nil
}

func (m *I82540MAC) setVCOSpeed() error {
	// Set PHY register 30, page 5, bit 8 to 0
	defaultPage, err := m.phy.ReadReg(M88E1000_PHY_PAGE_SELECT)
	if err != nil {
		return err
	}

	err = m.phy.WriteReg(M88E1000_PHY_PAGE_SELECT, 0x0005)
	if err != nil {
		return err
	}

	x, err := m.phy.ReadReg(M88E1000_PHY_GEN_CONTROL)
	if err != nil {
		return err
	}

	x &^= M88E1000_PHY_VCO_REG_BIT8
	err = m.phy.WriteReg(M88E1000_PHY_GEN_CONTROL, x)
	if err != nil {
		return err
	}

	// Set PHY register 30, page 4, bit 11 to 1
	err = m.phy.WriteReg(M88E1000_PHY_PAGE_SELECT, 0x0004)
	if err != nil {
		return err
	}

	x, err = m.phy.ReadReg(M88E1000_PHY_GEN_CONTROL)
	if err != nil {
		return err
	}

	x |= M88E1000_PHY_VCO_REG_BIT11
	err = m.phy.WriteReg(M88E1000_PHY_GEN_CONTROL, x)
	if err != nil {
		return err
	}

	return m.phy.WriteReg(M88E1000_PHY_PAGE_SELECT, defaultPage)
}
