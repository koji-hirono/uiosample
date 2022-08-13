package ethdev

type Counter interface {
	Value() uint64
	Clear()
}

type CounterGroup struct {
	RxPackets Counter
	TxPackets Counter
	RxOctets  Counter
	TxOctets  Counter
	RxMissed  Counter
	RxErrors  Counter
	TxErrors  Counter
	Ext       map[string]Counter
}
