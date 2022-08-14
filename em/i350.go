package em

func ClearVFTAI350(hw *HW) {
	for offset := 0; offset < VLAN_FILTER_TBL_SIZE; offset++ {
		for i := 0; i < 10; i++ {
			hw.RegWrite(VFTA+(offset<<2), 0)
		}
		hw.RegWriteFlush()
	}
}

func WriteVFTAI350(hw *HW, offset, val uint32) {
	for i := 0; i < 10; i++ {
		hw.RegWrite(VFTA+int(offset<<2), val)
	}
	hw.RegWriteFlush()
}

func UpdateNVMChecksumI350(hw *HW) error {
	for i := 0; i < 4; i++ {
		offset := NVM_82580_LAN_FUNC_OFFSET(uint16(i))
		err := UpdateChecksumWithOffset(hw, offset)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidateNVMChecksumI350(hw *HW) error {
	for i := 0; i < 4; i++ {
		offset := NVM_82580_LAN_FUNC_OFFSET(uint16(i))
		err := ValidateChecksumWithOffset(hw, offset)
		if err != nil {
			return err
		}
	}
	return nil
}
