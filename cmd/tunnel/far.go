package main

import (
	"net"
)

type ApplyAction uint8

const (
	ApplyActionDROP ApplyAction = 1 << iota
	ApplyActionFORW
	ApplyActionBUFF
	ApplyActionNOCP
)

type FAR struct {
	SEID   uint64
	ID     uint32
	Action ApplyAction
	Param  *ForwardParam
	PDRIDs []uint16
	BARID  *uint8
}

type ForwardParam struct {
	Creation *HeaderCreation
	Policy   string
}

type HeaderCreation struct {
	Desc     uint16
	TEID     uint32
	PeerAddr net.IP
	Port     uint16
}

type FARTable struct {
	s map[uint64]map[uint32]*FAR
}

func NewFARTable() *FARTable {
	t := new(FARTable)
	t.s = make(map[uint64]map[uint32]*FAR)
	return t
}

func (t *FARTable) Get(seid uint64, id uint32) *FAR {
	_, ok := t.s[seid]
	if !ok {
		return nil
	}
	far, ok := t.s[seid][id]
	if !ok {
		return nil
	}
	return far
}

func (t *FARTable) Put(seid uint64, id uint32, far *FAR) {
	_, ok := t.s[seid]
	if !ok {
		t.s[seid] = make(map[uint32]*FAR)
	}
	t.s[seid][id] = far
}

func (t *FARTable) Delete(seid uint64, id uint32) {
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
