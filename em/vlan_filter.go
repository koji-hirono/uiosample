package em

// VLAN Filter Table (4096 bits)
const VLAN_FILTER_TBL_SIZE = 128

func ClearVFTA(hw *HW) {
	for offset := 0; offset < VLAN_FILTER_TBL_SIZE; offset++ {
		hw.RegWrite(VFTA+(offset<<2), 0)
		hw.RegWriteFlush()
	}
}

func WriteVFTA(hw *HW, offset, val uint32) {
	hw.RegWrite(VFTA+int(offset<<2), val)
	hw.RegWriteFlush()
}
