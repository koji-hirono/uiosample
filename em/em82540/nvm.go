package em82540

import (
	"uiosample/em"
)

type NVM struct {
	hw *em.HW
}

func NewNVM(hw *em.HW) *NVM {
	m := new(NVM)
	m.hw = hw
	return m
}

func (m *NVM) InitParams() error {
	nvm := &m.hw.NVM

	nvm.Type = em.NVMTypeEepromMicrowire
	nvm.DelayUsec = 50
	nvm.OpcodeBits = 3

	switch nvm.Override {
	case em.NVMOverrideMicrowireLarge:
		nvm.AddressBits = 8
		nvm.WordSize = 256
	case em.NVMOverrideMicrowireSmall:
		nvm.AddressBits = 6
		nvm.WordSize = 64
	default:
		eecd := m.hw.RegRead(em.EECD)
		if eecd&em.EECD_SIZE != 0 {
			nvm.AddressBits = 8
			nvm.WordSize = 256
		} else {
			nvm.AddressBits = 6
			nvm.WordSize = 64
		}
	}

	return nil
}

func (m *NVM) Acquire() error {
	return em.AcquireNVM(m.hw)
}

func (m *NVM) Read(offset uint16, val []uint16) error {
	return em.ReadNVMMicrowire(m.hw, offset, val)
}

func (m *NVM) Release() {
	em.ReleaseNVM(m.hw)
}

func (m *NVM) Reload() {
	em.NVMReload(m.hw)
}

func (m *NVM) Update() error {
	return em.UpdateNVMChecksum(m.hw)
}

func (m *NVM) ValidLEDDefault() (uint16, error) {
	return em.ValidLEDDefault(m.hw)
}

func (m *NVM) Validate() error {
	return em.ValidateNVMChecksum(m.hw)
}

func (m *NVM) Write(offset uint16, val []uint16) error {
	return em.WriteNVMMicrowire(m.hw, offset, val)
}
