package em

const (
	PBA_8K  uint32 = 0x0008
	PBA_10K uint32 = 0x000a
	PBA_12K uint32 = 0x000c
	PBA_14K uint32 = 0x000e
	PBA_16K uint32 = 0x0010
	PBA_18K uint32 = 0x0012
	PBA_20K uint32 = 0x0014
	PBA_22K uint32 = 0x0016
	PBA_24K uint32 = 0x0018
	PBA_26K uint32 = 0x001a
	PBA_30K uint32 = 0x001e
	PBA_32K uint32 = 0x0020
	PBA_34K uint32 = 0x0022
	PBA_35K uint32 = 0x0023
	PBA_38K uint32 = 0x0026
	PBA_40K uint32 = 0x0028
	PBA_48K uint32 = 0x0030
	PBA_64K uint32 = 0x0040
)

func SetPBA(hw *HW) {
	var pba uint32
	switch hw.MAC.Type {
	case MACType82547, MACType82547Rev2:
		// 82547: Total Packet Buffer is 40K
		// 22K for Rx, 18K for Tx
		pba = PBA_22K
	case MACType82571, MACType82572, MACType80003es2lan:
		// 32K for Rx, 16K for Tx
		pba = PBA_32K
	case MACType82573:
		// 82573: Total Packet Buffer is 32K
		// 12K for Rx, 20K for Tx
		pba = PBA_12K
	case MACType82574, MACType82583:
		// 20K for Rx, 20K for Tx
		pba = PBA_20K
	case MACTypeIch8lan:
		pba = PBA_8K
	case MACTypeIch9lan, MACTypeIch10lan:
		pba = PBA_10K
	case MACTypePchlan:
	case MACTypePch2lan:
	case MACTypePch_lpt:
	case MACTypePch_spt:
	case MACTypePch_cnp:
		pba = PBA_26K
	default:
		// 40K for Rx, 24K for Tx
		pba = PBA_40K
	}
	hw.RegWrite(PBA, pba)
}
