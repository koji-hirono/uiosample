package em

import (
	"reflect"
	"unsafe"

	"uiosample/ethdev"
	"uiosample/hugetlb"
)

type TxQueue struct {
	ID        int
	NumDesc   int
	RingAddr  uintptr
	Threshold ethdev.RingThreshold
	Reg       Reg
	Buf       [][]byte
	Ring      []TxDesc
	head      int
	tail      int
}

const TxQueueOffloadCap = ethdev.TxOffloadCapMultiSegs |
	ethdev.TxOffloadCapVLANInsert |
	ethdev.TxOffloadCapIPv4Checksum |
	ethdev.TxOffloadCapUDPChecksum |
	ethdev.TxOffloadCapTCPChecksum

func (q *TxQueue) Start() error {
	return nil
}

func (q *TxQueue) Stop() error {
	return nil
}

// int eth_em_tx_descriptor_status(void *tx_queue, uint16_t offset)
func (q *TxQueue) Status(offset uint16) int {
	return 0
}

// uint16_t eth_em_xmit_pkts(void *tx_queue, struct rte_mbuf **tx_pkts,
//
//	uint16_t nb_pkts)
func (q *TxQueue) Do(pkts [][]byte) int {
	n := 0
	for i := 0; i < len(pkts); i++ {
		if !q.Can() {
			break
		}
		q.Tx(pkts[i])
		n++
	}
	q.Sync()
	return n
}

// uint16_t eth_em_prep_pkts(void *tx_queue, struct rte_mbuf **tx_pkts,
//
//	uint16_t nb_pkts)
func (q *TxQueue) Prepare([][]byte) int {
	return 0
}

func (q *TxQueue) InitBuf() (uintptr, error) {
	size := q.NumDesc * SizeofTxDesc
	desc, phys, err := hugetlb.Alloc(size)
	if err != nil {
		return 0, err
	}
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&q.Ring))
	hdr.Data = (uintptr)(unsafe.Pointer(&desc[0]))
	hdr.Cap = q.NumDesc
	hdr.Len = q.NumDesc

	q.Buf = make([][]byte, q.NumDesc)

	return phys, nil
}

func (q *TxQueue) Tx(pkt []byte) error {
	phys, err := hugetlb.PhysAddr(pkt)
	if err != nil {
		return err
	}
	q.Buf[q.tail] = pkt
	desc := &q.Ring[q.tail]
	desc.Addr = phys
	desc.Length = uint16(len(pkt))
	cmd := TxCommandEOP
	cmd |= TxCommandIFCS
	cmd |= TxCommandRS
	// cmd |= TxCommandIDE
	desc.Command = cmd
	desc.CSO = 0
	desc.Status = 0
	desc.CSS = 0
	desc.Special = 0
	q.tail = (q.tail + 1) % q.NumDesc
	q.Reg.Write(TDT(q.ID), uint32(q.tail))
	return nil
}

func (q *TxQueue) Can() bool {
	tail := (q.tail + 1) % q.NumDesc
	return tail != q.head
}

func (q *TxQueue) Sync() {
	i := q.head
	q.head = int(q.Reg.Read(TDH(q.ID)))
	for i != q.head {
		hugetlb.Free(q.Buf[i])
		q.Buf[i] = nil
		i = (i + 1) % q.NumDesc
	}
}

func (q *TxQueue) DiscardUnsetPackets() {
	i := q.tail
	q.tail = int(q.Reg.Read(TDT(q.ID)))
	q.head = int(q.Reg.Read(TDH(q.ID)))
	q.Reg.Write(TDT(q.ID), uint32(q.head))
	for i != q.head {
		i = (i - 1) % q.NumDesc
		hugetlb.Free(q.Buf[i])
		q.Ring[i].Addr = ^uintptr(0)
		q.Ring[i].Length = 0
	}
	q.tail = q.head
}
