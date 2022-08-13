package ethdev

type Counter interface {
	Value() uint64
	Clear()
}

type SumCounter struct {
	counters []Counter
	val      uint64
}

func NewSumCounter(counters ...Counter) *SumCounter {
	return &SumCounter{counters: counters}
}

func (c *SumCounter) Value() uint64 {
	for _, cc := range c.counters {
		c.val += cc.Value()
	}
	return c.val
}

func (c *SumCounter) Clear() {
	for _, cc := range c.counters {
		cc.Clear()
	}
	c.val = 0
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
