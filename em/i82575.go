package em

// SFP modules ID memory locations
const (
	SFF_IDENTIFIER_OFFSET = 0x00
	SFF_IDENTIFIER_SFF    = 0x02
	SFF_IDENTIFIER_SFP    = 0x03

	SFF_ETH_FLAGS_OFFSET = 0x06
)

// Flags for SFP modules compatible with ETH up to 1Gb
type SFPFlags uint8

const (
	SFPFlags1000BaseSX SFPFlags = 1 << iota
	SFPFlags1000BaseLX
	SFPFlags1000BaseCX
	SFPFlags1000BaseT
	SFPFlags100BaseLX
	SFPFlags100BaseFX
	SFPFlags10BaseBX10
	SFPFlags10BasePX
)

type I82575DeviceSpec struct {
	SGMIIActive        bool
	GlobalDeviceReset  bool
	EEEDisable         bool
	ModulePlugged      bool
	ClearSemaphoreOnce bool
	MTU                uint32
	Ethflags           SFPFlags
	MediaPort          uint8
	MediaChanged       bool
}

func I82575Init(hw *HW) {
	hw.Spec = &I82575DeviceSpec{}
	nvm := NewI82575NVM(hw)
	phy := NewI82575PHY(hw)
	mac := NewI82575MAC(hw, nvm, phy)
	hw.MAC.Op = mac
	hw.PHY.Op = phy
	hw.NVM.Op = nvm
}
