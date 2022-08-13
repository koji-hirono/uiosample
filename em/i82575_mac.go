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
	m.getMediaType()

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
	return BlinkLED(m.hw)
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
			err := CopperLinkSetupM88gen2(m.hw)
			if err != nil {
				return err
			}
		default:
			err := CopperLinkSetupM88(m.hw)
			if err != nil {
				return err
			}
		}
	case PHYTypeIgp3:
		err := CopperLinkSetupIgp(m.hw)
		if err != nil {
			return err
		}
	case PHYType82580:
		err := CopperLinkSetup82577(m.hw)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid phy type")
	}

	return SetupCopperLink(m.hw)
}

func (m *I82575MAC) getMediaType() error {
	hw := m.hw
	phy := &hw.PHY
	spec := hw.Spec.(I82575DeviceSpec)

	// Set internal phy as default
	spec.SGMIIActive = false
	spec.ModulePlugged = false

	// Get CSR setting
	ctrl := hw.RegRead(CTRL_EXT)

	// extract link mode setting
	linkmode := ctrl & CTRL_EXT_LINK_MODE_MASK

	switch linkmode {
	case CTRL_EXT_LINK_MODE_1000BASE_KX:
		phy.MediaType = MediaTypeInternalSerdes
	case CTRL_EXT_LINK_MODE_GMII:
		phy.MediaType = MediaTypeCopper
	case CTRL_EXT_LINK_MODE_SGMII:
		// Get phy control interface type set (MDIO vs. I2C)
		if m.SGMIIUsesMDIO() {
			phy.MediaType = MediaTypeCopper
			spec.SGMIIActive = true
			break
		}
		// Fall through for I2C based SGMII
		fallthrough
	case CTRL_EXT_LINK_MODE_PCIE_SERDES:
		// read media type from SFP EEPROM
		err := m.setSFPMediaType()
		if err != nil || phy.MediaType == MediaTypeUnknown {
			// If media type was not identified then return media
			// type defined by the CTRL_EXT settings.
			phy.MediaType = MediaTypeInternalSerdes
			if linkmode == CTRL_EXT_LINK_MODE_SGMII {
				phy.MediaType = MediaTypeCopper
				spec.SGMIIActive = true
			}
			break
		}
		// do not change link mode for 100BaseFX
		if spec.Ethflags&SFPFlags100BaseFX != 0 {
			break
		}

		// change current link mode setting
		ctrl &^= CTRL_EXT_LINK_MODE_MASK
		if phy.MediaType == MediaTypeCopper {
			ctrl |= CTRL_EXT_LINK_MODE_SGMII
		} else {
			ctrl |= CTRL_EXT_LINK_MODE_PCIE_SERDES
		}
		hw.RegWrite(CTRL_EXT, ctrl)
	}

	return nil
}

func (m *I82575MAC) SGMIIUsesMDIO() bool {
	hw := m.hw
	switch hw.MAC.Type {
	case MACType82575, MACType82576:
		x := hw.RegRead(MDIC)
		return x&MDIC_DEST != 0
	case MACType82580, MACTypeI350, MACTypeI354, MACTypeI210, MACTypeI211:
		x := hw.RegRead(MDICNFG)
		return x&MDICNFG_EXT_MDIO != 0
	default:
		return false
	}
}

func (m *I82575MAC) setSFPMediaType() error {
	hw := m.hw
	phy := &hw.PHY
	spec := hw.Spec.(I82575DeviceSpec)

	// Turn I2C interface ON and power on sfp cage
	ctrl := hw.RegRead(CTRL_EXT)
	ctrl &^= CTRL_EXT_SDP3_DATA
	hw.RegWrite(CTRL_EXT, ctrl|CTRL_I2C_ENA)
	hw.RegWriteFlush()
	defer hw.RegWrite(CTRL_EXT, ctrl)

	// Read SFP module data
	var tranceiver_type uint8
	timeout := 3
	for timeout > 0 {
		data, err := ReadSFPDataByte(hw, I2CCMD_SFP_DATA_ADDR(SFF_IDENTIFIER_OFFSET))
		if err == nil {
			tranceiver_type = data
			break
		}
		time.Sleep(100 * time.Millisecond)
		timeout--
	}
	if timeout == 0 {
		return errors.New("timeout")
	}

	flags, err := ReadSFPDataByte(hw, I2CCMD_SFP_DATA_ADDR(SFF_ETH_FLAGS_OFFSET))
	if err != nil {
		return err
	}
	spec.Ethflags = SFPFlags(flags)

	// Check if there is some SFP module plugged and powered
	if tranceiver_type == SFF_IDENTIFIER_SFP || tranceiver_type == SFF_IDENTIFIER_SFF {
		spec.ModulePlugged = true
		if spec.Ethflags&SFPFlags1000BaseLX != 0 || spec.Ethflags&SFPFlags1000BaseSX != 0 {
			phy.MediaType = MediaTypeInternalSerdes
		} else if spec.Ethflags&SFPFlags100BaseFX != 0 {
			spec.SGMIIActive = true
			phy.MediaType = MediaTypeInternalSerdes
		} else if spec.Ethflags&SFPFlags1000BaseT != 0 {
			spec.SGMIIActive = true
			phy.MediaType = MediaTypeCopper
		} else {
			phy.MediaType = MediaTypeUnknown
		}
	} else {
		phy.MediaType = MediaTypeUnknown
	}
	return nil
}
