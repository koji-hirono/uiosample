package em

import (
	"errors"
	"time"
)

const (
	RAR_ENTRIES_82575 = 16
	RAR_ENTRIES_82576 = 24
	RAR_ENTRIES_82580 = 24
	RAR_ENTRIES_I350  = 32
)

type I82575MAC struct {
	hw  *HW
	nvm *I82575NVM
	phy *I82575PHY
}

func NewI82575MAC(hw *HW, nvm *I82575NVM, phy *I82575PHY) *I82575MAC {
	m := new(I82575MAC)
	m.hw = hw
	m.nvm = nvm
	m.phy = phy
	return m
}

func (m *I82575MAC) InitParams() error {
	mac := &m.hw.MAC
	spec := m.hw.Spec.(I82575DeviceSpec)

	// Derives media type
	// e1000_get_media_type_82575(hw)

	// Set MTA register count
	mac.MTARegCount = 128

	// Set UTA register count
	if mac.Type == MACType82575 {
		mac.UTARegCount = 0
	} else {
		mac.UTARegCount = 128
	}

	// Set RAR entry count
	switch mac.Type {
	case MACType82576:
		mac.RAREntryCount = RAR_ENTRIES_82576
	case MACType82580:
		mac.RAREntryCount = RAR_ENTRIES_82580
	case MACTypeI350, MACTypeI354:
		mac.RAREntryCount = RAR_ENTRIES_I350
	default:
		mac.RAREntryCount = RAR_ENTRIES_82575
	}

	// Enable EEE default settings for EEE supported devices
	if mac.Type >= MACTypeI350 {
		spec.EEEDisable = false
	}

	// Allow a single clear of the SW semaphore on I210 and newer
	if mac.Type >= MACTypeI210 {
		spec.ClearSemaphoreOnce = true
	}

	// Set if part includes ASF firmware
	mac.ASFFirmwarePresent = true

	// FWSM register
	mac.HasFWSM = true

	// ARC supported; valid only if manageability features are enabled.
	mac.ArcSubsystemValid = m.hw.RegRead(FWSM)&FWSM_MODE_MASK != 0

	// set lan id for port to determine which phy lock to use
	m.SetLANID()
	return nil
}

func (m *I82575MAC) IDLEDInit() error {
	return IDLEDInit(m.hw)
}

func (m *I82575MAC) BlinkLED() error {
	// e1000_blink_led_generic
	return nil
}

func (m *I82575MAC) CheckMngMode() bool {
	// null
	return false
}

func (m *I82575MAC) CheckForLink() error {
	switch m.hw.PHY.MediaType {
	case MediaTypeCopper:
		return CheckForCopperLink(m.hw)
	default:
		// TODO:
		// e1000_get_pcs_speed_and_duplex_82575(hw, &speed, &duplex)

		// Use this flag to determine if link needs to be checked or
		// not.  If we have link clear the flag so that we do not
		// continue to check for link.
		m.hw.MAC.GetLinkStatus = !m.hw.MAC.SerdesHasLink

		// Configure Flow Control now that Auto-Neg has completed.
		// First, we need to restore the desired flow control
		// settings because we may have had to re-autoneg with a
		// different link partner.
		err := ConfigFCAfterLinkUp(m.hw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *I82575MAC) CleanupLED() error {
	return CleanupLED(m.hw)
}

func (m *I82575MAC) ClearHWCounters() {
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

	m.hw.RegRead(IAC)
	m.hw.RegRead(ICRXOC)

	m.hw.RegRead(ICRXPTC)
	m.hw.RegRead(ICRXATC)
	m.hw.RegRead(ICTXPTC)
	m.hw.RegRead(ICTXATC)
	m.hw.RegRead(ICTXQEC)
	m.hw.RegRead(ICTXQMTC)
	m.hw.RegRead(ICRXDMTC)

	m.hw.RegRead(CBTMPC)
	m.hw.RegRead(HTDPMC)
	m.hw.RegRead(CBRMPC)
	m.hw.RegRead(RPTHC)
	m.hw.RegRead(HGPTC)
	m.hw.RegRead(HTCBDPC)
	m.hw.RegRead(HGORCL)
	m.hw.RegRead(HGORCH)
	m.hw.RegRead(HGOTCL)
	m.hw.RegRead(HGOTCH)
	m.hw.RegRead(LENERRS)

	// This register should not be read in copper configurations
	if m.hw.PHY.MediaType == MediaTypeInternalSerdes || m.SGMIIActive() {
		m.hw.RegRead(SCVPC)
	}
}

func (m *I82575MAC) ClearVFTA() {
	switch m.hw.MAC.Type {
	case MACTypeI350, MACTypeI354:
		//ClearVFTAI350(m.hw)
	default:
		ClearVFTA(m.hw)
	}
}

func (m *I82575MAC) GetBusInfo() error {
	return GetBusInfoPCIE(m.hw)
}

func (m *I82575MAC) SetLANID() {
	SetLANIDMultiPortPCI(m.hw)
}

func (m *I82575MAC) GetLinkUpInfo() (uint16, uint16, error) {
	switch m.hw.PHY.MediaType {
	case MediaTypeCopper:
		return GetSpeedAndDuplexCopper(m.hw)
	default:
		// e1000_get_pcs_speed_and_duplex_82575
		return 0, 0, nil
	}
}

func (m *I82575MAC) LEDOn() error {
	return LEDOn(m.hw)
}

func (m *I82575MAC) LEDOff() error {
	return LEDOff(m.hw)
}

func (m *I82575MAC) UpdateMCAddrList(addrs [][6]byte) {
	UpdateMCAddrList(m.hw, addrs)
}

func (m *I82575MAC) ResetHW() error {
	if m.hw.MAC.Type >= MACType82580 {
		// e1000_reset_hw_82580
		return nil
	} else {
		// e1000_reset_hw_82575
		return nil
	}
}

func (m *I82575MAC) InitHW() error {
	switch m.hw.MAC.Type {
	case MACTypeI210, MACTypeI211:
		// e1000_init_hw_i210
		return nil
	default:
		// e1000_init_hw_82575
		return nil
	}
}

func (m *I82575MAC) ShutdownSerdes() {
	// e1000_shutdown_serdes_link_82575
}

func (m *I82575MAC) PowerUpSerdes() {
	// e1000_power_up_serdes_link_82575
}

func (m *I82575MAC) SetupLink() error {
	return SetupLink(m.hw)
}

func (m *I82575MAC) SetupPhysicalInterface() error {
	switch m.hw.PHY.MediaType {
	case MediaTypeCopper:
		return m.setupCopperLink()
	default:
		return m.setupSerdesLink()
	}
}

func (m *I82575MAC) SetupLED() error {
	return SetupLED(m.hw)
}

func (m *I82575MAC) WriteVFTA(offset, val uint32) {
	switch m.hw.MAC.Type {
	case MACTypeI350, MACTypeI354:
		//WriteVFTAI350(m.hw, offset, val)
	default:
		WriteVFTA(m.hw, offset, val)
	}
}

func (m *I82575MAC) ConfigCollisionDist() {
	// e1000_config_collision_dist_82575
}

func (m *I82575MAC) SetRAR(addr [6]byte, index int) error {
	return SetRAR(m.hw, addr, index)
}

func (m *I82575MAC) ReadMACAddr() error {
	// If there's an alternate MAC address place it in RAR0
	// so that it will override the Si installed default perm
	// address.
	err := CheckAltMACAddr(m.hw)
	if err != nil {
		return err
	}
	return ReadMACAddr(m.hw)
}

func (m *I82575MAC) ValidateMDISetting() error {
	if m.hw.MAC.Type >= MACType82580 {
		// null
		return nil
	}
	return ValidateMDISetting(m.hw)
}

func (m *I82575MAC) SetOBFFTimer(uint32) error {
	return nil
}

func (m *I82575MAC) AcquireSWFWSync(uint16) error {
	switch m.hw.MAC.Type {
	case MACTypeI210, MACTypeI211:
		// e1000_acquire_swfw_sync_i210
		return nil
	default:
		// e1000_acquire_swfw_sync_82575
		return nil
	}
}

func (m *I82575MAC) ReleaseSWFWSync(uint16) {
	switch m.hw.MAC.Type {
	case MACTypeI210, MACTypeI211:
		// e1000_release_swfw_sync_i210
	default:
		// e1000_release_swfw_sync_82575
	}
}

func (m *I82575MAC) SGMIIActive() bool {
	spec := m.hw.Spec.(*I82575DeviceSpec)
	return spec.SGMIIActive
}

func (m *I82575MAC) setupSerdesLink() error {
	return nil
}

func (m *I82575MAC) setupCopperLink() error {
	ctrl := m.hw.RegRead(CTRL)
	ctrl |= CTRL_SLU
	ctrl &^= CTRL_FRCSPD | CTRL_FRCDPX
	m.hw.RegWrite(CTRL, ctrl)

	switch m.hw.MAC.Type {
	case MACType82580, MACTypeI350, MACTypeI210, MACTypeI211:
		phpm := m.hw.RegRead(I82580_PHY_POWER_MGMT)
		phpm &^= I82580_PM_GO_LINKD
		m.hw.RegWrite(I82580_PHY_POWER_MGMT, phpm)
	}

	err := m.setupSerdesLink()
	if err != nil {
		return err
	}

	if m.SGMIIActive() {
		// allow time for SFP cage time to power up phy
		time.Sleep(300 * time.Millisecond)

		err := m.phy.Reset()
		if err != nil {
			return err
		}
	}

	switch m.hw.PHY.PHYType {
	case PHYTypeI210, PHYTypeM88:
		switch m.hw.PHY.ID {
		case I347AT4_E_PHY_ID, M88E1112_E_PHY_ID, M88E1340M_E_PHY_ID, M88E1543_E_PHY_ID, M88E1512_E_PHY_ID, I210_I_PHY_ID:
			//err := CopperLinkSetupM88gen2(m.hw)
			//if err != nil {
			//	return err
			//}
		default:
			//err := CopperLinkSetupM88(m.hw)
			//if err != nil {
			//	return err
			//}
		}
	case PHYTypeIgp3:
		//err := CopperLinkSetupIgp(m.hw)
		//if err != nil {
		//	return err
		//}
	case PHYType82580:
		//err := CopperLinkSetup82577(m.hw)
		//if err != nil {
		//	return err
		//}
	default:
		return errors.New("invalid phy type")
	}

	return SetupCopperLink(m.hw)
}
