package em

type FCMode int

const (
	FCModeNone FCMode = iota
	FCModeRxPause
	FCModeTxPause
	FCModeFull
	FCModeDefault = 0xff
)

type FCInfo struct {
	HighWater     uint32
	LowWater      uint32
	PauseTime     uint16
	RefreshTime   uint16
	SendXON       bool
	StrictIEEE    bool
	CurrentMode   FCMode
	RequestedMode FCMode
}

// Flow Control Constants
const (
	FLOW_CONTROL_ADDRESS_LOW  = 0x00C28001
	FLOW_CONTROL_ADDRESS_HIGH = 0x00000100
	FLOW_CONTROL_TYPE         = 0x8808
)

func SetDefaultFC(hw *HW) error {
	// Read and store word 0x0F of the EEPROM. This word contains bits
	// that determine the hardware's default PAUSE (flow control) mode,
	// a bit that determines whether the HW defaults to enabling or
	// disabling auto-negotiation, and the direction of the
	// SW defined pins. If there is no SW over-ride of the flow
	// control setting, then the variable hw->fc will
	// be initialized based on a value in the EEPROM.
	var data [1]uint16
	if hw.MAC.Type == MACTypeI350 {
		offset := NVM_82580_LAN_FUNC_OFFSET(hw.Bus.Func)
		err := hw.NVM.Op.Read(NVM_INIT_CONTROL2_REG+offset, data[:])
		if err != nil {
			return err
		}
	} else {
		err := hw.NVM.Op.Read(NVM_INIT_CONTROL2_REG, data[:])
		if err != nil {
			return err
		}
	}

	if data[0]&NVM_WORD0F_PAUSE_MASK == 0 {
		hw.FC.RequestedMode = FCModeNone
	} else if data[0]&NVM_WORD0F_PAUSE_MASK == NVM_WORD0F_ASM_DIR {
		hw.FC.RequestedMode = FCModeTxPause
	} else {
		hw.FC.RequestedMode = FCModeFull
	}

	return nil
}

func SetFCWatermarks(hw *HW) error {
	var fcrtl uint32
	var fcrth uint32
	// Set the flow control receive threshold registers.  Normally,
	// these registers will be set to a default threshold that may be
	// adjusted later by the driver's runtime code.  However, if the
	// ability to transmit pause frames is not enabled, then these
	// registers will be set to 0.
	if hw.FC.CurrentMode&FCModeTxPause != 0 {
		// We need to set up the Receive Threshold high and low water
		// marks as well as (optionally) enabling the transmission of
		// XON frames.
		fcrtl = hw.FC.LowWater
		if hw.FC.SendXON {
			fcrtl |= FCRTL_XONE
		}
		fcrth = hw.FC.HighWater
	}
	hw.RegWrite(FCRTL, fcrtl)
	hw.RegWrite(FCRTH, fcrth)
	return nil
}
