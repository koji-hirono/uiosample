package em82541

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

// s32  (*init_params)(struct e1000_hw *);
func (m *NVM) InitParams() error {
	return nil
}

// s32  (*acquire)(struct e1000_hw *);
func (m *NVM) Acquire() error {
	return nil
}

// s32  (*read)(struct e1000_hw *, u16, u16, u16 *);
func (m *NVM) Read(offset uint16, val []uint16) error {
	return nil
}

// void (*release)(struct e1000_hw *);
func (m *NVM) Release() {
}

// void (*reload)(struct e1000_hw *);
func (m *NVM) Reload() {
	em.NVMReload(m.hw)
}

// s32  (*update)(struct e1000_hw *);
func (m *NVM) Update() error {
	return nil
}

// s32  (*valid_led_default)(struct e1000_hw *, u16 *);
func (m *NVM) ValidLEDDefault() (uint16, error) {
	return 0, nil
}

// s32  (*validate)(struct e1000_hw *);
func (m *NVM) Validate() error {
	return nil
}

// s32  (*write)(struct e1000_hw *, u16, u16, u16 *);
func (m *NVM) Write(offset uint16, val []uint16) error {
	return nil
}
