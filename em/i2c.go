package em

import (
	"errors"
	"time"
)

const (
	I2CCMD_PHY_TIMEOUT = 200
)

func I2CCMD_SFP_DATA_ADDR(a uint16) uint16 {
	return 0x0000 + a
}

func I2CCMD_SFP_DIAG_ADDR(a uint16) uint16 {
	return 0x0100 + a
}

func ReadSFPDataByte(hw *HW, offset uint16) (uint8, error) {
	// Set up Op-code, EEPROM Address,in the I2CCMD
	// register. The MAC will take care of interfacing with the
	// EEPROM to retrieve the desired data.
	i2ccmd := uint32(offset)<<I2CCMD_REG_ADDR_SHIFT | I2CCMD_OPCODE_READ

	hw.RegWrite(I2CCMD, i2ccmd)

	// Poll the ready bit to see if the I2C read completed
	for i := 0; i < I2CCMD_PHY_TIMEOUT; i++ {
		time.Sleep(50 * time.Microsecond)
		data := hw.RegRead(I2CCMD)
		if data&I2CCMD_READY == 0 {
			continue
		}
		if data&I2CCMD_ERROR != 0 {
			return 0, errors.New("I2CCMD Error bit set")
		}
		return uint8(data), nil
	}

	return 0, errors.New("I2CCMD Read did not complete")
}

func ReadPHYRegI2C(hw *HW, offset uint32) (uint16, error) {
	phy := &hw.PHY

	// Set up Op-code, Phy Address, and register address in the I2CCMD
	// register.  The MAC will take care of interfacing with the
	// PHY to retrieve the desired data.
	i2ccmd := offset << I2CCMD_REG_ADDR_SHIFT
	i2ccmd |= phy.Addr << I2CCMD_PHY_ADDR_SHIFT
	i2ccmd |= I2CCMD_OPCODE_READ
	hw.RegWrite(I2CCMD, i2ccmd)

	// Poll the ready bit to see if the I2C read completed
	for i := 0; i < I2CCMD_PHY_TIMEOUT; i++ {
		time.Sleep(50 * time.Microsecond)
		i2ccmd := hw.RegRead(I2CCMD)
		if i2ccmd&I2CCMD_READY == 0 {
			continue
		}
		if i2ccmd&I2CCMD_ERROR != 0 {
			return 0, errors.New("I2CCMD Error bit set")
		}
		// Need to byte-swap the 16-bit value.
		data := uint16(i2ccmd)
		data = (data >> 8) | (data << 8)
		return data, nil
	}

	return 0, errors.New("I2CCMD Read did not complete")
}

func WritePHYRegI2C(hw *HW, offset uint32, data uint16) error {
	phy := &hw.PHY

	// Swap the data bytes for the I2C interface
	data = (data >> 8) | (data << 8)

	// Set up Op-code, Phy Address, and register address in the I2CCMD
	// register.  The MAC will take care of interfacing with the
	// PHY to retrieve the desired data.
	i2ccmd := offset << I2CCMD_REG_ADDR_SHIFT
	i2ccmd |= phy.Addr << I2CCMD_PHY_ADDR_SHIFT
	i2ccmd |= I2CCMD_OPCODE_WRITE
	i2ccmd |= uint32(data)
	hw.RegWrite(I2CCMD, i2ccmd)

	// Poll the ready bit to see if the I2C read completed
	for i := 0; i < I2CCMD_PHY_TIMEOUT; i++ {
		time.Sleep(50 * time.Microsecond)
		i2ccmd := hw.RegRead(I2CCMD)
		if i2ccmd&I2CCMD_READY == 0 {
			continue
		}
		if i2ccmd&I2CCMD_ERROR != 0 {
			return errors.New("I2CCMD Error bit set")
		}
		return nil
	}

	return errors.New("I2CCMD Write did not complete")
}

func ReadI2CByte(hw *HW, offset, addr byte) (byte, error) {
	return 0, nil
}

func WriteI2CByte(hw *HW, offset, addr, data byte) error {
	return nil
}
