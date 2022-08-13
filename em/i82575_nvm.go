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

	//err := e1000_acquire_swfw_sync_82575(hw, SWFW_EEP_SM)
	//if err != nil {
	//	return err
	//}
	//defer e1000_release_swfw_sync_82575(hw, SWFW_EEP_SM)

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
	// e1000_release_swfw_sync_82575(m.hw, SWFW_EEP_SM)
}

func (m *I82575NVM) Reload() {
	NVMReload(m.hw)
}

func (m *I82575NVM) Update() error {
	switch m.hw.MAC.Type {
	case MACType82580:
		return m.updateChecksum82580()
	case MACTypeI350:
		return m.updateChecksumI350()
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
		return m.validateChecksum82580()
	case MACTypeI350:
		return m.validateChecksumI350()
	default:
		return ValidateNVMChecksum(m.hw)
	}
}

func (m *I82575NVM) Write(offset uint16, val []uint16) error {
	return WriteNVMSpi(m.hw, offset, val)
}

func (m *I82575NVM) updateChecksum82580() error {
	var data [1]uint16
	err := m.Read(NVM_COMPATIBILITY_REG_3, data[:])
	if err != nil {
		return err
	}
	if data[0]&NVM_COMPATIBILITY_BIT_MASK == 0 {
		// set compatibility bit to validate checksums appropriately */
		data[0] |= NVM_COMPATIBILITY_BIT_MASK
		err := m.Write(NVM_COMPATIBILITY_REG_3, data[:])
		if err != nil {
			return err
		}
	}
	for i := 0; i < 4; i++ {
		offset := NVM_82580_LAN_FUNC_OFFSET(uint16(i))
		err := m.updateChecksumWithOffset(offset)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *I82575NVM) updateChecksumI350() error {
	for i := 0; i < 4; i++ {
		offset := NVM_82580_LAN_FUNC_OFFSET(uint16(i))
		err := m.updateChecksumWithOffset(offset)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *I82575NVM) updateChecksumWithOffset(offset uint16) error {
	var checksum uint16
	var data [1]uint16
	for i := offset; i < NVM_CHECKSUM_REG+offset; i++ {
		err := m.Read(i, data[:])
		if err != nil {
			return err
		}
		checksum += data[0]
	}
	data[0] = NVM_SUM - checksum
	return m.Write(NVM_CHECKSUM_REG+offset, data[:])
}

func (m *I82575NVM) validateChecksum82580() error {
	var data [1]uint16
	err := m.Read(NVM_COMPATIBILITY_REG_3, data[:])
	if err != nil {
		return err
	}
	n := 1
	if data[0]&NVM_COMPATIBILITY_BIT_MASK != 0 {
		// if chekcsums compatibility bit is set validate checksums
		// for all 4 ports.
		n = 4
	}
	for i := 0; i < n; i++ {
		offset := NVM_82580_LAN_FUNC_OFFSET(uint16(i))
		err := m.validateChecksumWithOffset(offset)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *I82575NVM) validateChecksumI350() error {
	for i := 0; i < 4; i++ {
		offset := NVM_82580_LAN_FUNC_OFFSET(uint16(i))
		err := m.validateChecksumWithOffset(offset)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *I82575NVM) validateChecksumWithOffset(offset uint16) error {
	var checksum uint16
	for i := offset; i < NVM_CHECKSUM_REG+offset+1; i++ {
		var data [1]uint16
		err := m.Read(i, data[:])
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
