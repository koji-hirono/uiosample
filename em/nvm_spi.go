package em

import (
	"time"
)

func ReadNVMSpi(hw *HW, offset uint16, val []uint16) error {
	nvm := &hw.NVM

	err := nvm.Op.Acquire()
	if err != nil {
		return err
	}
	defer nvm.Op.Release()

	err = ReadyNVMEeprom(hw)
	if err != nil {
		return err
	}
	StandbyNVM(hw)

	opcode := NVM_READ_OPCODE_SPI
	if nvm.AddressBits == 8 && offset >= 128 {
		opcode |= NVM_A8_OPCODE_SPI
	}

	// Send the READ command (opcode + addr)
	shiftOutEECbits(hw, opcode, nvm.OpcodeBits)
	shiftOutEECbits(hw, uint16(offset*2), nvm.AddressBits)

	// Read the data.  SPI NVMs increment the address with each byte
	// read and will roll over if reading beyond the end.  This allows
	// us to read the whole NVM from any offset
	for i := 0; i < len(val); i++ {
		word := shiftInEECbits(hw, 16)
		val[i] = word>>8 | word<<8
	}

	return nil
}

func WriteNVMSpi(hw *HW, offset uint16, val []uint16) error {
	nvm := &hw.NVM
	var widx int
	for widx < len(val) {
		err := nvm.Op.Acquire()
		if err != nil {
			return err
		}

		err = ReadyNVMEeprom(hw)
		if err != nil {
			nvm.Op.Release()
			return err
		}

		StandbyNVM(hw)

		// Send the WRITE ENABLE command (8 bit opcode)
		shiftOutEECbits(hw, NVM_WREN_OPCODE_SPI, nvm.OpcodeBits)

		StandbyNVM(hw)

		// Some SPI eeproms use the 8th address bit embedded in the
		// opcode
		opcode := NVM_WRITE_OPCODE_SPI
		if nvm.AddressBits == 8 && offset >= 128 {
			opcode |= NVM_A8_OPCODE_SPI
		}
		// Send the Write command (8-bit opcode + addr)
		shiftOutEECbits(hw, opcode, nvm.OpcodeBits)
		shiftOutEECbits(hw, (offset+uint16(widx))*2, nvm.AddressBits)

		// Loop to allow for up to whole page write of eeprom
		for widx < len(val) {
			word := val[widx]
			word = word>>8 | word<<8
			shiftOutEECbits(hw, word, 16)
			widx++

			if ((offset+uint16(widx))*2)%nvm.PageSize == 0 {
				StandbyNVM(hw)
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		nvm.Op.Release()
	}

	return nil
}
