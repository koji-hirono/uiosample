package ethdev

type Config struct {
	LinkSpeedCap LinkSpeedCap
	Rx           RxMode
	Tx           TxMode
	VNIC         bool
}

type RingThreshold struct {
	Prefetch  uint8
	Host      uint8
	Writeback uint8
}

type RxMode struct {
	MTU        uint32
	OffloadCap RxOffloadCap
}

type TxMode struct {
	OffloadCap TxOffloadCap
}

type RxConfig struct {
	Threshold  RingThreshold
	OffloadCap RxOffloadCap
}

type TxConfig struct {
	Threshold  RingThreshold
	OffloadCap TxOffloadCap
}
