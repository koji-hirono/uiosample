package ethdev

type Port interface {
	DeviceInfo() (*DeviceInfo, error)
	Close() error
	Configure(rxd, txd int, conf *Config) error
	RxQueueSetup(qid, rxd int, conf *RxConfig) error
	TxQueueSetup(qid, txd int, conf *TxConfig) error
	RxQueue(qid int) RxQueue
	TxQueue(qid int) TxQueue
	Start() error
	Stop() error
	Reset() error
	SetPromisc(unicast, multicast bool) error
	GetMACAddr() ([6]byte, error)
	CounterGroup() *CounterGroup
	LED() LED
	Link() Link
}
