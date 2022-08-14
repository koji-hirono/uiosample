package em

// GS40G - I210 PHY defines
const (
	GS40G_PAGE_SELECT  = 0x16
	GS40G_PAGE_SHIFT   = 16
	GS40G_OFFSET_MASK  = 0xFFFF
	GS40G_PAGE_2       = 0x20000
	GS40G_MAC_REG2     = 0x15
	GS40G_MAC_LB       = 0x4140
	GS40G_MAC_SPEED_1G = 0x0006
	GS40G_COPPER_SPEC  = 0x0010
)

func ReadPHYRegGS40G(hw *HW, offset uint32) (uint16, error) {
	phy := &hw.PHY
	err := phy.Op.Acquire()
	if err != nil {
		return 0, err
	}
	defer phy.Op.Release()

	page := offset >> GS40G_PAGE_SHIFT
	offset = offset & GS40G_OFFSET_MASK
	err = WritePHYRegMDIC(hw, GS40G_PAGE_SELECT, uint16(page))
	if err != nil {
		return 0, err
	}
	return ReadPHYRegMDIC(hw, offset)
}

func WritePHYRegGS40G(hw *HW, offset uint32, data uint16) error {
	phy := &hw.PHY
	err := phy.Op.Acquire()
	if err != nil {
		return err
	}
	defer phy.Op.Release()

	page := offset >> GS40G_PAGE_SHIFT
	offset = offset & GS40G_OFFSET_MASK
	err = WritePHYRegMDIC(hw, GS40G_PAGE_SELECT, uint16(page))
	if err != nil {
		return err
	}
	return WritePHYRegMDIC(hw, offset, data)
}
