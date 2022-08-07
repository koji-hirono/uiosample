package ich8lan

type ShadowRAM struct {
	Value    uint16
	Modified bool
}

const SHADOW_RAM_WORDS = 2048

type ULPState int

const (
	ULPStateUnknown ULPState = iota
	ULPStateOff
	ULPStateOn
)

type Device struct {
	KmrnLockLossWorkaroundEnabled bool
	ShadowRAM                     [SHADOW_RAM_WORDS]ShadowRAM
	nvm_k1_enabled                bool
	disable_k1_off                bool
	eee_disable                   bool
	eee_lp_ability                uint16
	ulp_state                     ULPState
	ulp_capability_disabled       bool
	during_suspend_flow           bool
	smbus_disable                 bool
}
