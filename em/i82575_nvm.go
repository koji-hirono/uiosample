package em

import (
	"errors"
)

const ID_LED_DEFAULT_82575_SERDES = ID_LED_DEF1_DEF2<<12 |
	ID_LED_DEF1_DEF2<<8 |
	ID_LED_DEF1_DEF2<<4 |
	ID_LED_OFF1_ON2

type I82575NVM struct {
	hw *HW
}

func NewI82575NVM(hw *HW) *I82575NVM {
	m := new(I82575NVM)
	m.hw = hw
	return m
}

func (m *I82575NVM) InitParams() error {
	hw := m.hw
	nvm := &hw.NVM
	eecd := hw.RegRead(EECD)
	size := uint16(eecd&EECD_SIZE_EX_MASK) >> EECD_SIZE_EX_SHIFT

	// Added to a constant, "size" becomes the left-shift value
	// for setting word_size.
	size += NVM_WORD_SIZE_BASE_SHIFT

	// Just in case size is out of range, cap it to the largest
	// EEPROM size supported
	if size > 15 {
		size = 15
	}

	nvm.WordSize = 1 << size
	if hw.MAC.Type < MACTypeI210 {
		nvm.OpcodeBits = 8
		nvm.DelayUsec = 1
		switch nvm.Override {
		case NVMOverrideSpiLarge:
			nvm.PageSize = 32
			nvm.AddressBits = 16
		case NVMOverrideSpiSmall:
			nvm.PageSize = 8
			nvm.AddressBits = 8
		default:
			if eecd&EECD_ADDR_BITS != 0 {
				nvm.PageSize = 32
				nvm.AddressBits = 16
			} else {
				nvm.PageSize = 8
				nvm.AddressBits = 8
			}
		}
		if nvm.WordSize == 1<<15 {
			nvm.PageSize = 128
		}
		nvm.Type = NVMTypeEepromSpi
	} else {
		nvm.Type = NVMTypeFlashHw
	}

	return nil
}

func (m *I82575NVM) Acquire() error {
	hw := m.hw

	err := AcquireSWFWSync82575(hw, SWFW_EEP_SM)
	if err != nil {
		return err
	}
	defer ReleaseSWFWSync82575(hw, SWFW_EEP_SM)

	// Check if there is some access
	// error this access may hook on
	if hw.MAC.Type == MACTypeI350 {
		eecd := hw.RegRead(EECD)
		if eecd&(EECD_BLOCKED|EECD_ABORT|EECD_TIMEOUT) != 0 {
			// Clear all access error flags
			hw.RegWrite(EECD, eecd|EECD_ERROR_CLR)
		}
	}

	if hw.MAC.Type == MACType82580 {
		eecd := hw.RegRead(EECD)
		if eecd&EECD_BLOCKED != 0 {
			// Clear access error flag
			hw.RegWrite(EECD, eecd|EECD_BLOCKED)
		}
	}

	return AcquireNVM(hw)
}

func (m *I82575NVM) Read(offset uint16, val []uint16) error {
	nvm := &m.hw.NVM
	if nvm.WordSize < 1<<15 {
		return ReadNVMEERD(m.hw, offset, val)
	} else {
		return ReadNVMSpi(m.hw, offset, val)
	}
}

func (m *I82575NVM) Release() {
	ReleaseNVM(m.hw)
	ReleaseSWFWSync82575(m.hw, SWFW_EEP_SM)
}

func (m *I82575NVM) Reload() {
	NVMReload(m.hw)
}

func (m *I82575NVM) Update() error {
	switch m.hw.MAC.Type {
	case MACType82580:
		return UpdateNVMChecksum82580(m.hw)
	case MACTypeI350:
		return UpdateNVMChecksumI350(m.hw)
	default:
		return UpdateNVMChecksum(m.hw)
	}
}

func (m *I82575NVM) ValidLEDDefault() (uint16, error) {
	phy := &m.hw.PHY
	var data [1]uint16
	err := m.Read(NVM_ID_LED_SETTINGS, data[:])
	if err != nil {
		return 0, err
	}
	if data[0] == 0 || data[0] == 0xffff {
		switch phy.MediaType {
		case MediaTypeInternalSerdes:
			return ID_LED_DEFAULT_82575_SERDES, nil
		default:
			return ID_LED_DEFAULT, nil
		}
	}
	return data[0], nil
}

func (m *I82575NVM) Validate() error {
	switch m.hw.MAC.Type {
	case MACType82580:
		return ValidateNVMChecksum82580(m.hw)
	case MACTypeI350:
		return ValidateNVMChecksumI350(m.hw)
	default:
		return ValidateNVMChecksum(m.hw)
	}
}

func (m *I82575NVM) Write(offset uint16, val []uint16) error {
	return WriteNVMSpi(m.hw, offset, val)
}

func UpdateChecksumWithOffset(hw *HW, offset uint16) error {
	nvm := &hw.NVM
	var checksum uint16
	var data [1]uint16
	for i := offset; i < NVM_CHECKSUM_REG+offset; i++ {
		err := nvm.Op.Read(i, data[:])
		if err != nil {
			return err
		}
		checksum += data[0]
	}
	data[0] = NVM_SUM - checksum
	return nvm.Op.Write(NVM_CHECKSUM_REG+offset, data[:])
}

func ValidateChecksumWithOffset(hw *HW, offset uint16) error {
	nvm := &hw.NVM
	var checksum uint16
	for i := offset; i < NVM_CHECKSUM_REG+offset+1; i++ {
		var data [1]uint16
		err := nvm.Op.Read(i, data[:])
		if err != nil {
			return err
		}
		checksum += data[0]
	}
	if checksum != NVM_SUM {
		return errors.New("NVM Checksum Invalid")
	}
	return nil
}
