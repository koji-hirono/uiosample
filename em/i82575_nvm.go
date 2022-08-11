package em

type I82575NVM struct {
	hw *HW
}

func NewI82575NVM(hw *HW) *I82575NVM {
	m := new(I82575NVM)
	m.hw = hw
	return m
}

func (m *I82575NVM) InitParams() error {
	nvm := &m.hw.NVM

	nvm.Type = NVMTypeEepromMicrowire
	nvm.DelayUsec = 50
	nvm.OpcodeBits = 3

	switch nvm.Override {
	case NVMOverrideMicrowireLarge:
		nvm.AddressBits = 8
		nvm.WordSize = 256
	case NVMOverrideMicrowireSmall:
		nvm.AddressBits = 6
		nvm.WordSize = 64
	default:
		eecd := m.hw.RegRead(EECD)
		if eecd&EECD_SIZE != 0 {
			nvm.AddressBits = 8
			nvm.WordSize = 256
		} else {
			nvm.AddressBits = 6
			nvm.WordSize = 64
		}
	}

	return nil
}

func (m *I82575NVM) Acquire() error {
	return AcquireNVM(m.hw)
}

func (m *I82575NVM) Read(offset uint16, val []uint16) error {
	return ReadNVMMicrowire(m.hw, offset, val)
}

func (m *I82575NVM) Release() {
	ReleaseNVM(m.hw)
}

func (m *I82575NVM) Reload() {
	NVMReload(m.hw)
}

func (m *I82575NVM) Update() error {
	return UpdateNVMChecksum(m.hw)
}

func (m *I82575NVM) ValidLEDDefault() (uint16, error) {
	return ValidLEDDefault(m.hw)
}

func (m *I82575NVM) Validate() error {
	return ValidateNVMChecksum(m.hw)
}

func (m *I82575NVM) Write(offset uint16, val []uint16) error {
	return WriteNVMMicrowire(m.hw, offset, val)
}
