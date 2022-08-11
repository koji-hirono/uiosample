package em

import (
	"errors"
	"time"
)

type SerdesLinkState int

const (
	SerdesLinkStateDown SerdesLinkState = iota
	SerdesLinkStateAutonegProgress
	SerdesLinkStateAutonegComplete
	SerdesLinkStateForcedUp
)

const (
	ALT_MAC_ADDRESS_OFFSET_LAN0 = iota * 3
	ALT_MAC_ADDRESS_OFFSET_LAN1
	ALT_MAC_ADDRESS_OFFSET_LAN2
	ALT_MAC_ADDRESS_OFFSET_LAN3
)

const RAR_ENTRIES = 15
const RAH_AV = 0x80000000 // Receive descriptor valid

const (
	HALF_DUPLEX = 1
	FULL_DUPLEX = 2
)

type MACInfo struct {
	Op       MACOp
	Addr     [6]byte
	PermAddr [6]byte

	Type MACType

	CollisionDelta uint32
	LEDCtlDefault  uint32
	LEDCtlMode1    uint32
	LEDCtlMode2    uint32
	MCFilterType   uint32
	TxPacketDelta  uint32
	TxCW           uint32

	CurrentIFSVal uint16
	IFSMaxVal     uint16
	IFSMinVal     uint16
	IFSRatio      uint16
	IFSStepSize   uint16
	MTARegCount   uint16
	UTARegCount   uint16
	MTAShadow     [128]uint32
	RAREntryCount uint16

	ForcedSpeedDuplex uint16

	AdaptiveIFS        bool
	HasFWSM            bool
	ArcSubsystemValid  bool
	ASFFirmwarePresent bool
	Autoneg            bool
	AutonegFailed      bool
	GetLinkStatus      bool
	InIFSMode          bool
	ReportTxEarly      bool
	SerdesLinkState    SerdesLinkState
	SerdesHasLink      bool
	TxPktfiltering     bool
	MaxFrameSize       uint32
}

type MACOp interface {
	// s32  (*init_params)(struct e1000_hw *);
	InitParams() error
	// s32  (*id_led_init)(struct e1000_hw *);
	IDLEDInit() error
	// s32  (*blink_led)(struct e1000_hw *);
	BlinkLED() error
	// bool (*check_mng_mode)(struct e1000_hw *);
	CheckMngMode() bool
	// s32  (*check_for_link)(struct e1000_hw *);
	CheckForLink() error
	// s32  (*cleanup_led)(struct e1000_hw *);
	CleanupLED() error
	// void (*clear_hw_cntrs)(struct e1000_hw *);
	ClearHWCounters()
	// void (*clear_vfta)(struct e1000_hw *);
	ClearVFTA()
	// s32  (*get_bus_info)(struct e1000_hw *);
	GetBusInfo() error
	// void (*set_lan_id)(struct e1000_hw *);
	SetLANID()
	// s32  (*get_link_up_info)(struct e1000_hw *, u16 *, u16 *);
	GetLinkUpInfo() (uint16, uint16, error)
	// s32  (*led_on)(struct e1000_hw *);
	LEDOn() error
	// s32  (*led_off)(struct e1000_hw *);
	LEDOff() error
	// void (*update_mc_addr_list)(struct e1000_hw *, u8 *, u32);
	UpdateMCAddrList([][6]byte)
	// s32  (*reset_hw)(struct e1000_hw *);
	ResetHW() error
	// s32  (*init_hw)(struct e1000_hw *);
	InitHW() error
	// void (*shutdown_serdes)(struct e1000_hw *);
	ShutdownSerdes()
	// void (*power_up_serdes)(struct e1000_hw *);
	PowerUpSerdes()
	// s32  (*setup_link)(struct e1000_hw *);
	SetupLink() error
	// s32  (*setup_physical_interface)(struct e1000_hw *);
	SetupPhysicalInterface() error
	// s32  (*setup_led)(struct e1000_hw *);
	SetupLED() error
	// void (*write_vfta)(struct e1000_hw *, u32, u32);
	WriteVFTA(uint32, uint32)
	// void (*config_collision_dist)(struct e1000_hw *);
	ConfigCollisionDist()
	// int  (*rar_set)(struct e1000_hw *, u8*, u32);
	SetRAR(addr [6]byte, index int) error
	// s32  (*read_mac_addr)(struct e1000_hw *);
	ReadMACAddr() error
	// s32  (*validate_mdi_setting)(struct e1000_hw *);
	ValidateMDISetting() error
	// s32  (*set_obff_timer)(struct e1000_hw *, u32);
	SetOBFFTimer(uint32) error
	// s32  (*acquire_swfw_sync)(struct e1000_hw *, u16);
	AcquireSWFWSync(uint16) error
	// void (*release_swfw_sync)(struct e1000_hw *, u16);
	ReleaseSWFWSync(uint16)
}

// void e1000_clear_hw_cntrs_base_generic(struct e1000_hw *hw)
func ClearHWCounters(hw *HW) {
	hw.RegRead(CRCERRS)
	hw.RegRead(SYMERRS)
	hw.RegRead(MPC)
	hw.RegRead(SCC)
	hw.RegRead(ECOL)
	hw.RegRead(MCC)
	hw.RegRead(LATECOL)
	hw.RegRead(COLC)
	hw.RegRead(DC)
	hw.RegRead(SEC)
	hw.RegRead(RLEC)
	hw.RegRead(XONRXC)
	hw.RegRead(XONTXC)
	hw.RegRead(XOFFRXC)
	hw.RegRead(XOFFTXC)
	hw.RegRead(FCRUC)
	hw.RegRead(GPRC)
	hw.RegRead(BPRC)
	hw.RegRead(MPRC)
	hw.RegRead(GPTC)
	hw.RegRead(GORCL)
	hw.RegRead(GORCH)
	hw.RegRead(GOTCL)
	hw.RegRead(GOTCH)
	hw.RegRead(RNBC)
	hw.RegRead(RUC)
	hw.RegRead(RFC)
	hw.RegRead(ROC)
	hw.RegRead(RJC)
	hw.RegRead(TORL)
	hw.RegRead(TORH)
	hw.RegRead(TOTL)
	hw.RegRead(TOTH)
	hw.RegRead(TPR)
	hw.RegRead(TPT)
	hw.RegRead(MPTC)
	hw.RegRead(BPTC)
}

func SetRAR(hw *HW, addr [6]byte, index int) error {
	// HW expects these in little endian so we reverse the byte order
	// from network order (big endian) to little endian
	low := uint32(addr[0])
	low |= uint32(addr[1]) << 8
	low |= uint32(addr[2]) << 16
	low |= uint32(addr[3]) << 24

	high := uint32(addr[4])
	high |= uint32(addr[5]) << 8

	// If MAC address zero, no need to set the AV bit
	if low == 0 || high == 0 {
		high |= RAH_AV
	}

	// Some bridges will combine consecutive 32-bit writes into
	// a single burst write, which will malfunction on some parts.
	// The flushes avoid this.
	hw.RegWrite(RAL(index), low)
	hw.RegWriteFlush()
	hw.RegWrite(RAH(index), high)
	hw.RegWriteFlush()
	return nil
}

func ValidateMDISetting(hw *HW) error {
	if !hw.MAC.Autoneg && (hw.PHY.MDIX == 0 || hw.PHY.MDIX == 3) {
		hw.PHY.MDIX = 1
		return errors.New("Invalid MDI setting detected")
	}
	return nil
}

func ConfigCollisionDist(hw *HW) {
	tctl := hw.RegRead(TCTL)
	tctl &^= TCTL_COLD
	tctl |= COLLISION_DISTANCE << COLD_SHIFT
	hw.RegWrite(TCTL, tctl)
	hw.RegWriteFlush()
}

func SetupLink(hw *HW) error {
	// In the case of the phy reset being blocked, we already have a link.
	// We do not need to set it up again.
	if hw.PHY.Op.CheckResetBlock() != nil {
		return nil
	}

	// If requested flow control is set to default, set flow control
	// based on the EEPROM flow control settings.
	if hw.FC.RequestedMode == FCModeDefault {
		err := SetDefaultFC(hw)
		if err != nil {
			return err
		}
	}

	// Save off the requested flow control mode for use later.  Depending
	// on the link partner's capabilities, we may or may not use this mode.
	hw.FC.CurrentMode = hw.FC.RequestedMode

	// Call the necessary media_type subroutine to configure the link.
	err := hw.MAC.Op.SetupPhysicalInterface()
	if err != nil {
		return err
	}

	// Initialize the flow control address, type, and PAUSE timer
	// registers to their default values.  This is done even if flow
	// control is disabled, because it does not hurt anything to
	// initialize these registers.
	hw.RegWrite(FCT, FLOW_CONTROL_TYPE)
	hw.RegWrite(FCAH, FLOW_CONTROL_ADDRESS_HIGH)
	hw.RegWrite(FCAL, FLOW_CONTROL_ADDRESS_LOW)

	hw.RegWrite(FCTTV, uint32(hw.FC.PauseTime))

	return SetFCWatermarks(hw)
}

func UpdateMCAddrList(hw *HW, addrs [][6]byte) {
	// clear mta_shadow
	for i := 0; i < len(hw.MAC.MTAShadow); i++ {
		hw.MAC.MTAShadow[i] = 0
	}

	// update mta_shadow from mc_addr_list
	for _, addr := range addrs {
		value := HashMCAddr(hw, addr)

		reg := (value >> 5) & uint32(hw.MAC.MTARegCount-1)
		bit := value & 0x1f

		hw.MAC.MTAShadow[reg] |= 1 << bit
	}

	// replace the entire MTA table
	for i := int(hw.MAC.MTARegCount) - 1; i >= 0; i-- {
		hw.RegWrite(MTA+(i<<2), hw.MAC.MTAShadow[i])
	}
	hw.RegWriteFlush()
}

func HashMCAddr(hw *HW, addr [6]byte) uint32 {
	// Register count multiplied by bits per register
	mask := (uint32(hw.MAC.MTARegCount) * 32) - 1

	// For a mc_filter_type of 0, bit_shift is the number of left-shifts
	// where 0xFF would still fall within the hash mask.
	shift := 0
	for mask>>shift != 0xff {
		shift++
	}

	switch hw.MAC.MCFilterType {
	case 1:
		shift += 1
	case 2:
		shift += 2
	case 3:
		shift += 4
	}

	value := uint32(addr[4]) >> (8 - shift)
	value |= uint32(addr[5]) << shift

	return value & mask
}

func CheckForCopperLink(hw *HW) error {
	mac := &hw.MAC

	// We only want to go out to the PHY registers to see if Auto-Neg
	// has completed and/or if our link status has changed.  The
	// get_link_status flag is set upon receiving a Link Status
	// Change or Rx Sequence Error interrupt.
	if !mac.GetLinkStatus {
		return nil
	}

	// First we want to see if the MII Status Register reports
	// link.  If so, then we want to get the current speed/duplex
	// of the PHY.
	link, err := PHYHasLink(hw, 1, 0)
	if err != nil {
		return err
	}
	if !link {
		// No link detected
		return nil
	}

	mac.GetLinkStatus = false

	// Check if there was DownShift, must be checked
	// immediately after link-up
	CheckDownshift(hw)

	// If we are forcing speed/duplex, then we simply return since
	// we have already determined whether we have link or not.
	if !mac.Autoneg {
		return errors.New("illegal config")
	}

	// Auto-Neg is enabled.  Auto Speed Detection takes care
	// of MAC speed/duplex configuration.  So we only need to
	// configure Collision Distance in the MAC.
	mac.Op.ConfigCollisionDist()

	// Configure Flow Control now that Auto-Neg has completed.
	// First, we need to restore the desired flow control
	// settings because we may have had to re-autoneg with a
	// different link partner.
	return ConfigFCAfterLinkUp(hw)
}

func CheckForFiberLink(hw *HW) error {
	mac := &hw.MAC

	ctrl := hw.RegRead(CTRL)
	status := hw.RegRead(STATUS)
	rxcw := hw.RegRead(RXCW)

	// If we don't have link (auto-negotiation failed or link partner
	// cannot auto-negotiate), the cable is plugged in (we have signal),
	// and our link partner is not trying to auto-negotiate with us (we
	// are receiving idles or data), we need to force link up. We also
	// need to give auto-negotiation time to complete, in case the cable
	// was just plugged in. The autoneg_failed flag does this.
	// (ctrl & E1000_CTRL_SWDPIN1) == 1 == have signal
	if ctrl&CTRL_SWDPIN1 != 0 && status&STATUS_LU == 0 && rxcw&RXCW_C == 0 {
		if !mac.AutonegFailed {
			mac.AutonegFailed = true
			return nil
		}

		// Disable auto-negotiation in the TXCW register
		hw.RegWrite(TXCW, mac.TxCW & ^TXCW_ANE)

		// Force link-up and also force full-duplex.
		ctrl = hw.RegRead(CTRL)
		ctrl |= CTRL_SLU | CTRL_FD
		hw.RegWrite(CTRL, ctrl)

		// Configure Flow Control after forcing link up. */
		err := ConfigFCAfterLinkUp(hw)
		if err != nil {
			return err
		}
	} else if ctrl&CTRL_SLU != 0 && rxcw&RXCW_C != 0 {
		// If we are forcing link and we are receiving /C/ ordered
		// sets, re-enable auto-negotiation in the TXCW register
		// and disable forced link in the Device Control register
		// in an attempt to auto-negotiate with our link partner.
		hw.RegWrite(TXCW, mac.TxCW)
		hw.RegWrite(CTRL, ctrl & ^CTRL_SLU)

		mac.SerdesHasLink = true
	}

	return nil
}

func CheckForSerdesLink(hw *HW) error {
	mac := &hw.MAC

	ctrl := hw.RegRead(CTRL)
	status := hw.RegRead(STATUS)
	rxcw := hw.RegRead(RXCW)

	// If we don't have link (auto-negotiation failed or link partner
	// cannot auto-negotiate), and our link partner is not trying to
	// auto-negotiate with us (we are receiving idles or data),
	// we need to force link up. We also need to give auto-negotiation
	// time to complete.
	// (ctrl & E1000_CTRL_SWDPIN1) == 1 == have signal
	if status&STATUS_LU == 0 && rxcw&RXCW_C == 0 {
		if !mac.AutonegFailed {
			mac.AutonegFailed = true
			return nil
		}

		// Disable auto-negotiation in the TXCW register
		hw.RegWrite(TXCW, mac.TxCW & ^TXCW_ANE)

		// Force link-up and also force full-duplex.
		ctrl = hw.RegRead(CTRL)
		ctrl |= CTRL_SLU | CTRL_FD
		hw.RegWrite(CTRL, ctrl)

		// Configure Flow Control after forcing link up.
		err := ConfigFCAfterLinkUp(hw)
		if err != nil {
			return err
		}
	} else if ctrl&CTRL_SLU != 0 && rxcw&RXCW_C != 0 {
		// If we are forcing link and we are receiving /C/ ordered
		// sets, re-enable auto-negotiation in the TXCW register
		// and disable forced link in the Device Control register
		// in an attempt to auto-negotiate with our link partner.
		hw.RegWrite(TXCW, mac.TxCW)
		hw.RegWrite(CTRL, ctrl & ^CTRL_SLU)

		mac.SerdesHasLink = true
	} else if hw.RegRead(TXCW)&TXCW_ANE == 0 {
		// If we force link for non-auto-negotiation switch, check
		// link status based on MAC synchronization for internal
		// serdes media type.
		// SYNCH bit and IV bit are sticky.
		time.Sleep(10 * time.Microsecond)
		rxcw := hw.RegRead(RXCW)
		if rxcw&RXCW_SYNCH != 0 {
			if rxcw&RXCW_IV == 0 {
				mac.SerdesHasLink = true
			}
		} else {
			mac.SerdesHasLink = false
		}
	}

	if hw.RegRead(TXCW)&TXCW_ANE != 0 {
		status = hw.RegRead(STATUS)
		if status&STATUS_LU != 0 {
			// SYNCH bit and IV bit are sticky, so reread rxcw.
			time.Sleep(10 * time.Microsecond)
			rxcw := hw.RegRead(RXCW)
			if rxcw&RXCW_SYNCH != 0 {
				if rxcw&RXCW_IV == 0 {
					mac.SerdesHasLink = true
				} else {
					mac.SerdesHasLink = false
				}
			} else {
				mac.SerdesHasLink = false
			}
		} else {
			mac.SerdesHasLink = false
		}
	}

	return nil
}

func GetSpeedAndDuplexCopper(hw *HW) (uint16, uint16, error) {
	status := hw.RegRead(STATUS)
	var speed uint16
	if status&STATUS_SPEED_1000 != 0 {
		speed = 1000
	} else if status&STATUS_SPEED_100 != 0 {
		speed = 100
	} else {
		speed = 10
	}

	var duplex uint16
	if status&STATUS_FD != 0 {
		duplex = FULL_DUPLEX
	} else {
		duplex = HALF_DUPLEX
	}

	return speed, duplex, nil
}

func GetSpeedAndDuplexFiberSerdes(hw *HW) (uint16, uint16, error) {
	return 1000, FULL_DUPLEX, nil
}

func CheckAltMACAddr(hw *HW) error {
	var data [1]uint16
	err := hw.NVM.Op.Read(NVM_COMPAT, data[:])
	if err != nil {
		return err
	}

	// not supported on older hardware or 82573
	if hw.MAC.Type < MACType82571 || hw.MAC.Type == MACType82573 {
		return nil
	}

	// Alternate MAC address is handled by the option ROM for 82580
	// and newer. SW support not required.
	if hw.MAC.Type >= MACType82580 {
		return nil
	}

	err = hw.NVM.Op.Read(NVM_ALT_MAC_ADDR_PTR, data[:])
	if err != nil {
		return err
	}
	offset := data[0]
	if offset == 0xffff || offset == 0 {
		return nil
	}
	switch hw.Bus.Func {
	case 1:
		offset += ALT_MAC_ADDRESS_OFFSET_LAN1
	case 2:
		offset += ALT_MAC_ADDRESS_OFFSET_LAN2
	case 3:
		offset += ALT_MAC_ADDRESS_OFFSET_LAN3
	}

	var addr [6]byte
	for i := 0; i < 6; i += 2 {
		reg := offset + uint16(i>>1)
		err := hw.NVM.Op.Read(reg, data[:])
		if err != nil {
			return err
		}
		addr[i] = byte(data[0])
		addr[i+1] = byte(data[0] >> 8)
	}

	// if multicast bit is set, the alternate address will not be used
	if addr[0]&0x01 != 0 {
		return nil
	}

	// We have a valid alternate MAC address, and we want to treat it the
	// same as the normal permanent MAC address stored by the HW into the
	// RAR. Do this by mapping this address into RAR0.
	hw.MAC.Op.SetRAR(addr, 0)

	return nil
}
