package main

type QER struct {
	SEID   uint64
	ID     uint32
	Gate   uint8
	MBR    BitRate
	GBR    BitRate
	CorrID uint32
	RQI    uint8
	QFI    uint8
	PPI    uint8
	PDRIDs []uint16
}

type BitRate struct {
	UL uint64
	DL uint64
}

type QERTable struct {
	s map[uint64]map[uint32]*QER
}

func NewQERTable() *QERTable {
	t := new(QERTable)
	t.s = make(map[uint64]map[uint32]*QER)
	return t
}

func (t *QERTable) Get(seid uint64, id uint32) *QER {
	_, ok := t.s[seid]
	if !ok {
		return nil
	}
	qer, ok := t.s[seid][id]
	if !ok {
		return nil
	}
	return qer
}

func (t *QERTable) Put(seid uint64, id uint32, qer *QER) {
	_, ok := t.s[seid]
	if !ok {
		t.s[seid] = make(map[uint32]*QER)
	}
	t.s[seid][id] = qer
}

func (t *QERTable) Delete(seid uint64, id uint32) {
	_, ok := t.s[seid]
	if !ok {
		return
	}
	delete(t.s[seid], id)
	if len(t.s[seid]) > 0 {
		return
	}
	delete(t.s, seid)
}
