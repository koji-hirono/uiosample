package em82575

/* Flags for SFP modules compatible with ETH up to 1Gb */
/*
struct sfp_e1000_flags {
        u8 e1000_base_sx:1;
        u8 e1000_base_lx:1;
        u8 e1000_base_cx:1;
        u8 e1000_base_t:1;
        u8 e100_base_lx:1;
        u8 e100_base_fx:1;
        u8 e10_base_bx10:1;
        u8 e10_base_px:1;
};
*/
type SFPFlags uint8

type Device struct {
	sgmii_active         bool
	global_device_reset  bool
	eee_disable          bool
	module_plugged       bool
	clear_semaphore_once bool
	mtu                  uint32
	eth_flags            SFPFlags
	media_port           uint8
	media_changed        bool
}
