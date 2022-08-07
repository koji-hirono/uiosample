package em82541

import (
	"uiosample/em"
)

type MAC struct {
	hw *em.HW
}

func NewMAC(hw *em.HW) *MAC {
	m := new(MAC)
	m.hw = hw
	return m
}

// s32  (*init_params)(struct e1000_hw *);
func (m *MAC) InitParam() error {
	mac := &m.hw.MAC
	m.hw.PHY.MediaType = em.MediaTypeCopper
	mac.MTARegCount = 128
	mac.RAREntryCount = em.RAR_ENTRIES
	mac.ASFFirmwarePresent = true
	return nil
}

// s32  (*id_led_init)(struct e1000_hw *);
func (m *MAC) IDLEDInit() error {
	return nil
}

// s32  (*blink_led)(struct e1000_hw *);
func (m *MAC) BlinkLED() error {
	return nil
}

// bool (*check_mng_mode)(struct e1000_hw *);
func (m *MAC) CheckMngMode() bool {
	return false
}

// s32  (*check_for_link)(struct e1000_hw *);
func (m *MAC) CheckForLink() error {
	return nil
}

// s32  (*cleanup_led)(struct e1000_hw *);
func (m *MAC) CleanupLED() error {
	return nil
}

// void (*clear_hw_cntrs)(struct e1000_hw *);
func (m *MAC) ClearHWCounters() {
}

// void (*clear_vfta)(struct e1000_hw *);
func (m *MAC) ClearVFTA() {
}

// s32  (*get_bus_info)(struct e1000_hw *);
func (m *MAC) GetBusInfo() error {
	return nil
}

// void (*set_lan_id)(struct e1000_hw *);
func (m *MAC) SetLANID() {
}

// s32  (*get_link_up_info)(struct e1000_hw *, u16 *, u16 *);
func (m *MAC) GetLinkUpInfo() (uint16, uint16, error) {
	return 0, 0, nil
}

// s32  (*led_on)(struct e1000_hw *);
func (m *MAC) LEDOn() error {
	return nil
}

// s32  (*led_off)(struct e1000_hw *);
func (m *MAC) LEDOff() error {
	return nil
}

// void (*update_mc_addr_list)(struct e1000_hw *, u8 *, u32);
func (m *MAC) UpdateMACAddrList(addr [6]byte, index int) {
}

// s32  (*reset_hw)(struct e1000_hw *);
func (m *MAC) ResetHW() error {
	return nil
}

// s32  (*init_hw)(struct e1000_hw *);
func (m *MAC) InitHW() error {
	return nil
}

// void (*shutdown_serdes)(struct e1000_hw *);
func (m *MAC) ShutdownSerdes() {
}

// void (*power_up_serdes)(struct e1000_hw *);
func (m *MAC) PowerUpSerdes() {
}

// s32  (*setup_link)(struct e1000_hw *);
func (m *MAC) SetupLink() error {
	return nil
}

// s32  (*setup_physical_interface)(struct e1000_hw *);
func (m *MAC) SetupPhysicalInterface() error {
	return nil
}

// s32  (*setup_led)(struct e1000_hw *);
func (m *MAC) SetupLED() error {
	return nil
}

// void (*write_vfta)(struct e1000_hw *, u32, u32);
func (m *MAC) WriteVFTA(offset, val uint32) {
}

// void (*config_collision_dist)(struct e1000_hw *);
func (m *MAC) ConfigCollisionDist() {
}

// int  (*rar_set)(struct e1000_hw *, u8*, u32);
func (m *MAC) SetRAR(addr [6]byte, index int) error {
	return nil
}

// s32  (*read_mac_addr)(struct e1000_hw *);
func (m *MAC) ReadMACAddr() error {
	return nil
}

// s32  (*validate_mdi_setting)(struct e1000_hw *);
func (m *MAC) ValidateMDISetting() error {
	return nil
}

// s32  (*acquire_swfw_sync)(struct e1000_hw *, u16);
func (m *MAC) AcquireSWFWSync(uint16) error {
	return nil
}

// void (*release_swfw_sync)(struct e1000_hw *, u16);
func (m *MAC) ReleaseSWFWSync(uint16) {
}
