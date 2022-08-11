package em

type I82540NVM struct {
	hw *HW
}

func NewI82540NVM(hw *HW) *I82540NVM {
	m := new(I82540NVM)
	m.hw = hw
	return m
}

func (m *I82540NVM) InitParams() error {
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

func (m *I82540NVM) Acquire() error {
	return AcquireNVM(m.hw)
}

func (m *I82540NVM) Read(offset uint16, val []uint16) error {
	return ReadNVMMicrowire(m.hw, offset, val)
}

func (m *I82540NVM) Release() {
	ReleaseNVM(m.hw)
}

func (m *I82540NVM) Reload() {
	NVMReload(m.hw)
}

func (m *I82540NVM) Update() error {
	return UpdateNVMChecksum(m.hw)
}

func (m *I82540NVM) ValidLEDDefault() (uint16, error) {
	return ValidLEDDefault(m.hw)
}

func (m *I82540NVM) Validate() error {
	return ValidateNVMChecksum(m.hw)
}

func (m *I82540NVM) Write(offset uint16, val []uint16) error {
	return WriteNVMMicrowire(m.hw, offset, val)
}
