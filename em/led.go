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
