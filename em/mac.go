package em

import (
	"errors"
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

// VLAN Filter Table (4096 bits)
const VLAN_FILTER_TBL_SIZE = 128

// Word definitions for ID LED Settings
const (
	ID_LED_RESERVED_0000 = 0x0000
	ID_LED_RESERVED_FFFF = 0xFFFF

	ID_LED_DEF1_DEF2 = 0x1
	ID_LED_DEF1_ON2  = 0x2
	ID_LED_DEF1_OFF2 = 0x3
	ID_LED_ON1_DEF2  = 0x4
	ID_LED_ON1_ON2   = 0x5
	ID_LED_ON1_OFF2  = 0x6
	ID_LED_OFF1_DEF2 = 0x7
	ID_LED_OFF1_ON2  = 0x8
	ID_LED_OFF1_OFF2 = 0x9

	ID_LED_DEFAULT = (ID_LED_OFF1_ON2 << 12) |
		(ID_LED_OFF1_OFF2 << 8) |
		(ID_LED_DEF1_DEF2 << 4) |
		ID_LED_DEF1_DEF2
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

	ForcedSpeedduplex uint8

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

// s32 e1000_id_led_init_generic(struct e1000_hw *hw)
func IDLEDInit(hw *HW) error {
	mac := &hw.MAC

	x, err := hw.NVM.Op.ValidLEDDefault()
	if err != nil {
		return err
	}

	mac.LEDCtlDefault = hw.RegRead(LEDCTL)
	mac.LEDCtlMode1 = mac.LEDCtlDefault
	mac.LEDCtlMode2 = mac.LEDCtlDefault

	ledmask := uint16(0x0f)
	ledctlmask := uint32(0xff)
	ledctlon := LEDCTL_MODE_LED_ON
	ledctloff := LEDCTL_MODE_LED_OFF
	for i := 0; i < 4; i++ {
		t := (x >> (i << 2)) & ledmask
		switch t {
		case ID_LED_ON1_DEF2, ID_LED_ON1_ON2, ID_LED_ON1_OFF2:
			mac.LEDCtlMode1 &^= ledctlmask << (i << 3)
			mac.LEDCtlMode1 |= ledctlon << (i << 3)
		case ID_LED_OFF1_DEF2, ID_LED_OFF1_ON2, ID_LED_OFF1_OFF2:
			mac.LEDCtlMode1 &^= ledctlmask << (i << 3)
			mac.LEDCtlMode1 |= ledctloff << (i << 3)
		}
		switch t {
		case ID_LED_DEF1_ON2, ID_LED_ON1_ON2, ID_LED_OFF1_ON2:
			mac.LEDCtlMode2 &^= ledctlmask << (i << 3)
			mac.LEDCtlMode2 |= ledctlon << (i << 3)
		case ID_LED_DEF1_OFF2, ID_LED_ON1_OFF2, ID_LED_OFF1_OFF2:
			mac.LEDCtlMode2 &^= ledctlmask << (i << 3)
			mac.LEDCtlMode2 |= ledctloff << (i << 3)
		}
	}

	return nil
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

// s32 e1000_valid_led_default_generic(struct e1000_hw *hw, u16 *data)
func ValidLEDDefault(hw *HW) (uint16, error) {
	nvm := &hw.NVM

	var x [1]uint16
	err := nvm.Op.Read(NVM_ID_LED_SETTINGS, x[:])
	if err != nil {
		return 0, nil
	}

	if x[0] == uint16(0) || x[0] == ^uint16(0) {
		return ID_LED_DEFAULT, nil
	}

	return x[0], nil
}

func CleanupLED(hw *HW) error {
	hw.RegWrite(LEDCTL, hw.MAC.LEDCtlDefault)
	return nil
}

func ClearVFTA(hw *HW) {
	for offset := 0; offset < VLAN_FILTER_TBL_SIZE; offset++ {
		hw.RegWrite(VFTA+(offset<<2), 0)
		hw.RegWriteFlush()
	}
}

func GetBusInfoPCI(hw *HW) error {
	mac := &hw.MAC
	bus := &hw.Bus
	status := hw.RegRead(STATUS)

	// PCI or PCI-X?
	if status&STATUS_PCIX_MODE != 0 {
		bus.Type = BusTypePCIX
	} else {
		bus.Type = BusTypePCI
	}

	// Bus speed
	if bus.Type == BusTypePCI {
		if status&STATUS_PCI66 != 0 {
			bus.Speed = BusSpeed66
		} else {
			bus.Speed = BusSpeed33
		}
	} else {
		switch status & STATUS_PCIX_SPEED {
		case STATUS_PCIX_SPEED_66:
			bus.Speed = BusSpeed66
		case STATUS_PCIX_SPEED_100:
			bus.Speed = BusSpeed100
		case STATUS_PCIX_SPEED_133:
			bus.Speed = BusSpeed133
		default:
			bus.Speed = BusSpeedReserved
		}
	}

	// Bus width
	if status&STATUS_BUS64 != 0 {
		bus.Width = BusWidth64
	} else {
		bus.Width = BusWidth32
	}

	// Which PCI(-X) function?
	mac.Op.SetLANID()

	return nil
}

func SetLANIDMultiPortPCI(hw *HW) {
	bus := &hw.Bus
	bus.Func = 0
}

func WriteVFTA(hw *HW, offset, val uint32) {
	hw.RegWrite(VFTA+int(offset<<2), val)
	hw.RegWriteFlush()
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

func LEDOn(hw *HW) error {
	switch hw.PHY.MediaType {
	case MediaTypeFiber:
		ctrl := hw.RegRead(CTRL)
		ctrl &^= CTRL_SWDPIN0
		ctrl |= CTRL_SWDPIO0
		hw.RegWrite(CTRL, ctrl)
	case MediaTypeCopper:
		hw.RegWrite(LEDCTL, hw.MAC.LEDCtlMode2)
	}
	return nil
}

func LEDOff(hw *HW) error {
	switch hw.PHY.MediaType {
	case MediaTypeFiber:
		ctrl := hw.RegRead(CTRL)
		ctrl |= CTRL_SWDPIN0
		ctrl |= CTRL_SWDPIO0
		hw.RegWrite(CTRL, ctrl)
	case MediaTypeCopper:
		hw.RegWrite(LEDCTL, hw.MAC.LEDCtlMode1)
	}
	return nil
}

func SetupLED(hw *HW) error {
	switch hw.PHY.MediaType {
	case MediaTypeFiber:
		ledctl := hw.RegRead(LEDCTL)
		hw.MAC.LEDCtlDefault = ledctl
		// Turn off LED0
		ledctl &^= LEDCTL_LED0_IVRT | LEDCTL_LED0_BLINK | LEDCTL_LED0_MODE_MASK
		ledctl |= LEDCTL_MODE_LED_OFF << LEDCTL_LED0_MODE_SHIFT
		hw.RegWrite(LEDCTL, ledctl)
	case MediaTypeCopper:
		hw.RegWrite(LEDCTL, hw.MAC.LEDCtlMode1)
	}
	return nil
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
	return nil
}

func CheckForFiberLink(hw *HW) error {
	return nil
}

func CheckForSerdesLink(hw *HW) error {
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
