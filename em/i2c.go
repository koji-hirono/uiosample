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
