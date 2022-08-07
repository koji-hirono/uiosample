package em

import (
	"errors"
)

func ReadNVMMicrowire(hw *HW, offset uint16, val []uint16) error {
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

	for i := 0; i < len(val); i++ {
		shiftOutEECbits(hw, NVM_READ_OPCODE_MICROWIRE, nvm.OpcodeBits)
		shiftOutEECbits(hw, offset+uint16(i), nvm.AddressBits)
		val[i] = shiftInEECbits(hw, 16)
		StandbyNVM(hw)
	}

	return nil
}

func WriteNVMMicrowire(hw *HW, offset uint16, val []uint16) error {
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

	shiftOutEECbits(hw, NVM_EWEN_OPCODE_MICROWIRE, uint16(nvm.OpcodeBits+2))
	shiftOutEECbits(hw, 0, uint16(nvm.AddressBits-2))
	StandbyNVM(hw)

	for i := 0; i < len(val); i++ {
		shiftOutEECbits(hw, NVM_WRITE_OPCODE_MICROWIRE, nvm.OpcodeBits)
		shiftOutEECbits(hw, offset+uint16(i), nvm.AddressBits)
		shiftOutEECbits(hw, val[i], 16)
		StandbyNVM(hw)
		var widx int
		for widx < 200 {
			eecd := hw.RegRead(EECD)
			if eecd&EECD_DO != 0 {
				break
			}
			widx++
		}
		if widx == 200 {
			return errors.New("NVM Write did not complete")
		}
		StandbyNVM(hw)

	}

	shiftOutEECbits(hw, NVM_EWDS_OPCODE_MICROWIRE, uint16(nvm.OpcodeBits+2))
	shiftOutEECbits(hw, 0, uint16(nvm.AddressBits-2))

	return nil
}
