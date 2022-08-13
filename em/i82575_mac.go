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

const SW_SYNCH_MB uint16 = 0x0100
const STAT_DEV_RST_SET uint32 = 0x00100000

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
		m.getPCSSpeedAndDuplex()

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
		m.clearVFTAI350()
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
		return m.getPCSSpeedAndDuplex()
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
		return m.resetHW82580()
	} else {
		return m.resetHW82575()
	}
}

func (m *I82575MAC) InitHW() error {
	switch m.hw.MAC.Type {
	case MACTypeI210, MACTypeI211:
		return m.initHWI210()
	default:
		return m.initHW82575()
	}
}

func (m *I82575MAC) ShutdownSerdes() {
	hw := m.hw
	if hw.PHY.MediaType != MediaTypeInternalSerdes && !m.SGMIIActive() {
		return
	}

	if EnableManagePT(hw) {
		return
	}

	// Disable PCS to turn off link
	pcs := hw.RegRead(PCS_CFG0)
	pcs &^= PCS_CFG_PCS_EN
	hw.RegWrite(PCS_CFG0, pcs)

	// shutdown the laser
	ctrl := hw.RegRead(CTRL_EXT)
	ctrl |= CTRL_EXT_SDP3_DATA
	hw.RegWrite(CTRL_EXT, ctrl)

	// flush the write to verify completion
	hw.RegWriteFlush()
	time.Sleep(1 * time.Millisecond)
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
		m.writeVFTAI350(offset, val)
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

func (m *I82575MAC) AcquireSWFWSync(mask uint16) error {
	switch m.hw.MAC.Type {
	case MACTypeI210, MACTypeI211:
		return acquireSWFWSyncI210(m.hw, mask)
	default:
		return acquireSWFWSync82575(m.hw, mask)
	}
}

func (m *I82575MAC) ReleaseSWFWSync(mask uint16) {
	switch m.hw.MAC.Type {
	case MACTypeI210, MACTypeI211:
		releaseSWFWSyncI210(m.hw, mask)
	default:
		releaseSWFWSync82575(m.hw, mask)
	}
}

func (m *I82575MAC) clearVFTAI350() {
	hw := m.hw
	for offset := 0; offset < VLAN_FILTER_TBL_SIZE; offset++ {
		for i := 0; i < 10; i++ {
			hw.RegWrite(VFTA+(offset<<2), 0)
		}
		hw.RegWriteFlush()
	}
}

func (m *I82575MAC) getPCSSpeedAndDuplex() (uint16, uint16, error) {
	hw := m.hw
	mac := &hw.MAC
	// Read the PCS Status register for link state. For non-copper mode,
	// the status register is not accurate. The PCS status register is
	// used instead.
	pcs := hw.RegRead(PCS_LSTAT)

	// The link up bit determines when link is up on autoneg.
	if pcs&PCS_LSTS_LINK_OK == 0 {
		mac.SerdesHasLink = false
		return 0, 0, nil
	}
	mac.SerdesHasLink = true

	// Detect and store PCS speed
	var speed uint16
	if pcs&PCS_LSTS_SPEED_1000 != 0 {
		speed = 1000
	} else if pcs&PCS_LSTS_SPEED_100 != 0 {
		speed = 100
	} else {
		speed = 10
	}

	// Detect and store PCS duplex
	var duplex uint16
	if pcs&PCS_LSTS_DUPLEX_FULL != 0 {
		duplex = FULL_DUPLEX
	} else {
		duplex = HALF_DUPLEX
	}

	// Check if it is an I354 2.5Gb backplane connection.
	if mac.Type == MACTypeI354 {
		status := hw.RegRead(STATUS)
		if status&STATUS_2P5_SKU != 0 && status&STATUS_2P5_SKU_OVER == 0 {
			speed = 2500
			duplex = FULL_DUPLEX
		}
	}
	return speed, duplex, nil
}

func (m *I82575MAC) resetHW82580() error {
	hw := m.hw
	spec := hw.Spec.(*I82575DeviceSpec)

	global_device_reset := spec.GlobalDeviceReset
	spec.GlobalDeviceReset = false

	// 82580 does not reliably do global_device_reset due to hw errata
	if hw.MAC.Type == MACType82580 {
		global_device_reset = false
	}

	// Get current control state.
	ctrl := hw.RegRead(CTRL)

	// Prevent the PCI-E bus from sticking if there is no TLP connection
	// on the last TLP read/write transaction when MAC is reset.
	DisablePCIEMaster(hw)

	hw.RegWrite(IMC, ^uint32(0))
	hw.RegWrite(RCTL, 0)
	hw.RegWrite(TCTL, TCTL_PSP)
	hw.RegWriteFlush()

	time.Sleep(10 * time.Millisecond)

	// BH SW mailbox bit in SW_FW_SYNC
	swmbsw_mask := SW_SYNCH_MB
	// Determine whether or not a global dev reset is requested
	if global_device_reset && m.AcquireSWFWSync(swmbsw_mask) != nil {
		global_device_reset = false
	}

	if global_device_reset && hw.RegRead(STATUS)&STAT_DEV_RST_SET == 0 {
		ctrl |= CTRL_DEV_RST
	} else {
		ctrl |= CTRL_RST
	}
	hw.RegWrite(CTRL, ctrl)

	switch hw.DeviceID {
	case DEV_ID_DH89XXCC_SGMII:
	default:
		hw.RegWriteFlush()
	}

	// Add delay to insure DEV_RST or RST has time to complete
	time.Sleep(5 * time.Millisecond)

	// When auto config read does not complete, do not
	// return with an error. This can happen in situations
	// where there is no eeprom and prevents getting link.
	GetAutoRDDone(hw)

	// clear global device reset status bit
	hw.RegWrite(STATUS, STAT_DEV_RST_SET)

	// Clear any pending interrupt events.
	hw.RegWrite(IMC, ^uint32(0))
	hw.RegRead(ICR)

	m.resetMDIConfig82580()

	// Install any alternate MAC address into RAR0
	err := CheckAltMACAddr(hw)

	// Release semaphore
	if global_device_reset {
		m.ReleaseSWFWSync(swmbsw_mask)
	}

	return err
}

func (m *I82575MAC) resetMDIConfig82580() error {
	hw := m.hw
	if hw.MAC.Type != MACType82580 {
		return nil
	}
	if !m.SGMIIActive() {
		return nil
	}

	var data [1]uint16
	err := hw.NVM.Op.Read(NVM_INIT_CONTROL3_PORT_A+NVM_82580_LAN_FUNC_OFFSET(hw.Bus.Func), data[:])
	if err != nil {
		return err
	}

	mdicnfg := hw.RegRead(MDICNFG)
	if data[0]&NVM_WORD24_EXT_MDIO != 0 {
		mdicnfg |= MDICNFG_EXT_MDIO
	}
	if data[0]&NVM_WORD24_COM_MDIO != 0 {
		mdicnfg |= MDICNFG_COM_MDIO
	}
	hw.RegWrite(MDICNFG, mdicnfg)
	return nil
}

func (m *I82575MAC) resetHW82575() error {
	hw := m.hw
	// Prevent the PCI-E bus from sticking if there is no TLP connection
	// on the last TLP read/write transaction when MAC is reset.
	DisablePCIEMaster(hw)

	// set the completion timeout for interface
	m.setPCIECompletionTimeout()

	hw.RegWrite(IMC, ^uint32(0))

	hw.RegWrite(RCTL, 0)
	hw.RegWrite(TCTL, TCTL_PSP)
	hw.RegWriteFlush()

	time.Sleep(10 * time.Millisecond)

	ctrl := hw.RegRead(CTRL)
	hw.RegWrite(CTRL, ctrl|CTRL_RST)

	// When auto config read does not complete, do not
	// return with an error. This can happen in situations
	// where there is no eeprom and prevents getting link.
	GetAutoRDDone(hw)

	// If EEPROM is not present, run manual init scripts
	if hw.RegRead(EECD)&EECD_PRES == 0 {
		m.resetInitScript()
	}

	// Clear any pending interrupt events.
	hw.RegWrite(IMC, ^uint32(0))
	hw.RegRead(ICR)

	// Install any alternate MAC address into RAR0
	return CheckAltMACAddr(hw)
}

func (m *I82575MAC) setPCIECompletionTimeout() error {
	hw := m.hw
	gcr := hw.RegRead(GCR)
	defer func() {
		// disable completion timeout resend
		gcr &^= GCR_CMPL_TMOUT_RESEND
		hw.RegWrite(GCR, gcr)
	}()

	// only take action if timeout value is defaulted to 0
	if gcr&GCR_CMPL_TMOUT_MASK != 0 {
		return nil
	}

	// if capababilities version is type 1 we can write the
	// timeout of 10ms to 200ms through the GCR register
	if gcr&GCR_CAP_VER2 == 0 {
		gcr |= GCR_CMPL_TMOUT_10ms
		return nil
	}

	// for version 2 capabilities we need to write the config space
	// directly in order to set the completion timeout value for
	// 16ms to 55ms
	devctl, err := ReadPCIECapReg(hw, PCIE_DEVICE_CONTROL2)
	if err != nil {
		return err
	}
	devctl |= PCIE_DEVICE_CONTROL2_16ms
	return WritePCIECapReg(hw, PCIE_DEVICE_CONTROL2, devctl)
}

func (m *I82575MAC) resetInitScript() error {
	hw := m.hw
	if hw.MAC.Type != MACType82575 {
		return nil
	}
	// SerDes configuration via SERDESCTRL
	Write8bitCtrlReg(hw, SCTL, 0x00, 0x0c)
	Write8bitCtrlReg(hw, SCTL, 0x01, 0x78)
	Write8bitCtrlReg(hw, SCTL, 0x1b, 0x23)
	Write8bitCtrlReg(hw, SCTL, 0x23, 0x15)

	// CCM configuration via CCMCTL register
	Write8bitCtrlReg(hw, CCMCTL, 0x14, 0x00)
	Write8bitCtrlReg(hw, CCMCTL, 0x10, 0x00)

	// PCIe lanes configuration
	Write8bitCtrlReg(hw, GIOCTL, 0x00, 0xec)
	Write8bitCtrlReg(hw, GIOCTL, 0x61, 0xdf)
	Write8bitCtrlReg(hw, GIOCTL, 0x34, 0x05)
	Write8bitCtrlReg(hw, GIOCTL, 0x2f, 0x81)

	// PCIe PLL Configuration
	Write8bitCtrlReg(hw, SCCTL, 0x02, 0x47)
	Write8bitCtrlReg(hw, SCCTL, 0x14, 0x00)
	Write8bitCtrlReg(hw, SCCTL, 0x10, 0x00)

	return nil
}

func (m *I82575MAC) initHWI210() error {
	return nil
}

func (m *I82575MAC) initHW82575() error {
	// Initialize identification LED
	// This is not fatal and we should not stop init due to this
	m.IDLEDInit()

	// Disabling VLAN filtering
	m.ClearVFTA()

	err := InitHWBase(m.hw)

	// Set the default MTU size
	spec := m.hw.Spec.(*I82575DeviceSpec)
	spec.MTU = 1500

	// Clear all of the statistics registers (clear on read).  It is
	// important that we do this after we have tried to establish link
	// because the symbol error count will increment wildly if there
	// is no link.
	m.ClearHWCounters()

	return err
}

func (m *I82575MAC) writeVFTAI350(offset, val uint32) {
	hw := m.hw
	for i := 0; i < 10; i++ {
		hw.RegWrite(VFTA+int(offset<<2), val)
	}
	hw.RegWriteFlush()
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

func acquireSWFWSync82575(hw *HW, mask uint16) error {
	swmask := uint32(mask)
	fwmask := uint32(mask) << 16
	timeout := 200
	for i := 0; i < timeout; i++ {
		err := GetHWSemaphore(hw)
		if err != nil {
			return err
		}
		swfw_sync := hw.RegRead(SW_FW_SYNC)
		if swfw_sync&(fwmask|swmask) == 0 {
			swfw_sync |= swmask
			hw.RegWrite(SW_FW_SYNC, swfw_sync)
			PutHWSemaphore(hw)
			return nil
		}
		// Firmware currently using resource (fwmask)
		// or other software thread using resource (swmask)
		PutHWSemaphore(hw)
		time.Sleep(5 * time.Millisecond)
	}
	return errors.New("SW_FW_SYNC timeout")
}

func releaseSWFWSync82575(hw *HW, mask uint16) {
	for {
		err := GetHWSemaphore(hw)
		if err == nil {
			break
		}
	}

	swfw_sync := hw.RegRead(SW_FW_SYNC)
	swfw_sync &^= uint32(mask)
	hw.RegWrite(SW_FW_SYNC, swfw_sync)

	PutHWSemaphore(hw)
}

func acquireSWFWSyncI210(hw *HW, mask uint16) error {
	swmask := uint32(mask)
	fwmask := uint32(mask) << 16
	timeout := 200
	for i := 0; i < timeout; i++ {
		err := GetHWSemaphoreI210(hw)
		if err != nil {
			return err
		}
		swfw_sync := hw.RegRead(SW_FW_SYNC)
		if swfw_sync&(fwmask|swmask) == 0 {
			swfw_sync |= swmask
			hw.RegWrite(SW_FW_SYNC, swfw_sync)
			PutHWSemaphore(hw)
			return nil
		}
		// Firmware currently using resource (fwmask)
		// or other software thread using resource (swmask)
		PutHWSemaphore(hw)
		time.Sleep(5 * time.Millisecond)
	}
	return errors.New("SW_FW_SYNC timeout")
}

func releaseSWFWSyncI210(hw *HW, mask uint16) {
	for {
		err := GetHWSemaphoreI210(hw)
		if err == nil {
			break
		}
	}

	swfw_sync := hw.RegRead(SW_FW_SYNC)
	swfw_sync &^= uint32(mask)
	hw.RegWrite(SW_FW_SYNC, swfw_sync)

	PutHWSemaphore(hw)
}

func GetHWSemaphoreI210(hw *HW) error {
	// Get the SW semaphore
	timeout := int(hw.NVM.WordSize) + 1
	var i int
	for i < timeout {
		swsm := hw.RegRead(SWSM)
		if swsm&SWSM_SMBI == 0 {
			break
		}
		time.Sleep(50 * time.Microsecond)
		i++
	}
	if i == timeout {
		// In rare circumstances, the SW semaphore may already be held
		// unintentionally. Clear the semaphore once before giving up.
		spec := hw.Spec.(*I82575DeviceSpec)
		if spec.ClearSemaphoreOnce {
			spec.ClearSemaphoreOnce = false
			PutHWSemaphore(hw)
			for i = 0; i < timeout; i++ {
				swsm := hw.RegRead(SWSM)
				if swsm&SWSM_SMBI == 0 {
					break
				}
				time.Sleep(50 * time.Microsecond)
			}
		}
		// If we do not have the semaphore here, we have to give up.
		if i == timeout {
			return errors.New("SMBI bit is set")
		}
	}

	// Get the FW semaphore.
	for i = 0; i < timeout; i++ {
		swsm := hw.RegRead(SWSM)
		hw.RegWrite(SWSM, swsm|SWSM_SWESMBI)
		// Semaphore acquired if bit latched
		if hw.RegRead(SWSM)&SWSM_SWESMBI != 0 {
			break
		}
		time.Sleep(50 * time.Microsecond)
	}
	if i == timeout {
		// Release semaphores
		PutHWSemaphore(hw)
		return errors.New("timeout")
	}
	return nil
}
