package em

type MBXStats struct {
	MsgsTx uint32
	MsgsRx uint32

	Acks uint32
	Reqs uint32
	Rsts uint32
}

type MBXInfo struct {
	Op        MBXOp
	Stats     MBXStats
	Timeout   uint32
	UsecDelay uint32
	Size      uint16
}

type MBXOp interface {
	InitParams() error
	Read([]uint32, uint16) error
	Write([]uint32, uint16) error
	ReadPosted([]uint32, uint16) error
	WritePosted([]uint32, uint16) error
	CheckForMsg(uint16) error
	CheckForACK(uint16) error
	CheckForRST(uint16) error
}

/*
   s32 (*init_params)(struct e1000_hw *hw);
   s32 (*read)(struct e1000_hw *, u32 *, u16,  u16);
   s32 (*write)(struct e1000_hw *, u32 *, u16, u16);
   s32 (*read_posted)(struct e1000_hw *, u32 *, u16,  u16);
   s32 (*write_posted)(struct e1000_hw *, u32 *, u16, u16);
   s32 (*check_for_msg)(struct e1000_hw *, u16);
   s32 (*check_for_ack)(struct e1000_hw *, u16);
   s32 (*check_for_rst)(struct e1000_hw *, u16);
*/
