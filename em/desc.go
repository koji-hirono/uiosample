package em

type RxDesc struct {
	Addr    uintptr
	Length  uint16
	Chksum  uint16
	Status  uint8
	Errors  uint8
	Special uint16
}

const SizeofRxDesc = 16

// RxDesc.Status
const (
	// descriptor done
	RxStatusDD uint8 = 1 << iota
	// end of packet
	RxStatusEOP
	// ignore checksum indication
	RxStatusIXSM
	// 802.1Q
	RxStatusVP
	// UDP checksum calculated
	RxStatusUDPCS
	// TCP checksum calculated
	RxStatusTCPCS
	// IPv4 checksum calculated
	RxStatusIPCS
	// passed in-exact filter
	RxStatusPIF
)

type TxDesc struct {
	Addr    uintptr
	Length  uint16
	CSO     uint8
	Command uint8
	Status  uint8
	CSS     uint8
	Special uint16
}

const SizeofTxDesc = 16

// TxDesc.Command
const (
	// end of packet
	TxCommandEOP uint8 = 1 << iota
	// insert FCS
	TxCommandIFCS
	// insert checksum
	TxCommandIC
	// report status
	TxCommandRS
	// reserved
	TxCommandRSV
	// extension
	TxCommandDext
	// VLAN packet enable
	TxCommandVTE
	// interrupt delay enable
	TxCommandIDE
)

// TxDesc.Status
const (
	TxStatusDD uint8 = 1 << 0
	TxStatusEC uint8 = 1 << 1
	TxStatusLC uint8 = 1 << 2
)
