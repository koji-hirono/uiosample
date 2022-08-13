package em

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

type LED struct {
	mac *MACInfo
}

func NewLED(mac *MACInfo) *LED {
	return &LED{mac: mac}
}

func (l *LED) On() error {
	return l.mac.Op.LEDOn()
}

func (l *LED) Off() error {
	return l.mac.Op.LEDOff()
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

func CleanupLED(hw *HW) error {
	hw.RegWrite(LEDCTL, hw.MAC.LEDCtlDefault)
	return nil
}

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

func BlinkLED(hw *HW) error {
	switch hw.PHY.MediaType {
	case MediaTypeFiber:
		// always blink LED0 for PCI-E fiber
		ledctl := LEDCTL_LED0_BLINK
		ledctl |= LEDCTL_MODE_LED_ON << LEDCTL_LED0_MODE_SHIFT
		hw.RegWrite(LEDCTL, ledctl)
	default:
		// Set the blink bit for each LED that's "on" (0x0E)
		// (or "off" if inverted) in ledctl_mode2.  The blink
		// logic in hardware only works when mode is set to "on"
		// so it must be changed accordingly when the mode is
		// "off" and inverted.
		ledctl := hw.MAC.LEDCtlMode2
		for i := 0; i < 32; i += 8 {
			mode := (hw.MAC.LEDCtlMode2 >> i) & LEDCTL_LED0_MODE_MASK
			def := hw.MAC.LEDCtlDefault >> i
			if (def&LEDCTL_LED0_IVRT == 0 && mode == LEDCTL_MODE_LED_ON) ||
				(def&LEDCTL_LED0_IVRT != 0 && mode == LEDCTL_MODE_LED_OFF) {
				ledctl &^= LEDCTL_LED0_MODE_MASK << i
				ledctl |= LEDCTL_LED0_BLINK
				ledctl |= LEDCTL_MODE_LED_ON << i
			}
		}
		hw.RegWrite(LEDCTL, ledctl)
	}
	return nil
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
