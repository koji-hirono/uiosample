package ethdev

type Config struct {
	LinkSpeedCap LinkSpeedCap
	Rx           RxMode
	Tx           TxMode
}

type RxMode struct {
	MTU        uint32
	OffloadCap RxOffloadCap
}

type TxMode struct {
	OffloadCap TxOffloadCap
}

type RxConfig struct {
	OffloadCap RxOffloadCap
}

type TxConfig struct {
	OffloadCap TxOffloadCap
}
