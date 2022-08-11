package em

import (
	"reflect"
	"unsafe"

	"uiosample/ethdev"
	"uiosample/hugetlb"
)

type RxQueue struct {
	ID        int
	NumDesc   int
	RingAddr  uintptr
	Threshold ethdev.RingThreshold
	Reg       Reg
	Buf       [][]byte
	Ring      []RxDesc
	rdh       int
	rdt       int
	rxnext    int
}

const RxQueueOffloadCap = ethdev.RxOffloadCapVLANStrip |
	ethdev.RxOffloadCapVLANFilter |
	ethdev.RxOffloadCapIPv4Checksum |
	ethdev.RxOffloadCapUDPChecksum |
	ethdev.RxOffloadCapTCPChecksum |
	ethdev.RxOffloadCapKeepCRC |
	ethdev.RxOffloadCapScatter

func (q *RxQueue) Start() error {
	return nil
}

func (q *RxQueue) Stop() error {
	return nil
}

// uint32_t eth_em_rx_queue_count(void *rx_queue)
func (q *RxQueue) Count() int {
	return 0
}

// int eth_em_rx_descriptor_status(void *rx_queue, uint16_t offset)
func (q *RxQueue) Status(offset uint16) int {
	return 0
}

// uint16_t eth_em_recv_pkts(void *rx_queue, struct rte_mbuf **rx_pkts,
//
//	uint16_t nb_pkts)
func (q *RxQueue) Do(pkts [][]byte) int {
	q.Sync()
	n := 0
	for i := 0; i < len(pkts); i++ {
		if !q.Can() {
			break
		}
		pkt := q.Rx()
		pkts[i] = pkt
		n++
	}
	for q.CanAddBuf() {
		p, phys, err := hugetlb.Alloc(2048)
		if err != nil {
			break
		}
		q.AddBuf(p, phys)
	}
	return n
}

func (q *RxQueue) Rx() []byte {
	length := q.Ring[q.rxnext].Length
	pkt := q.Buf[q.rxnext][:length]
	q.Buf[q.rxnext] = nil
	q.rxnext = (q.rxnext + 1) % q.NumDesc
	return pkt
}

func (q *RxQueue) Can() bool {
	if q.rxnext == q.rdh {
		return false
	}
	if q.Ring[q.rxnext].Status&RxStatusDD == 0 {
		return false
	}
	return true
}

func (q *RxQueue) CanAddBuf() bool {
	rdt := (q.rdt + 1) % q.NumDesc
	return rdt != q.rxnext
}

func (q *RxQueue) AddBuf(p []byte, phys uintptr) {
	desc := &q.Ring[q.rdt]
	desc.Addr = phys
	desc.Status &^= RxStatusDD
	q.Buf[q.rdt] = p
	q.rdt = (q.rdt + 1) % q.NumDesc
}

func (q *RxQueue) FreeAllBuf() {
	for q.rdt != q.rdh {
		q.rdt = (q.rdt - 1) % q.NumDesc
		desc := &q.Ring[q.rdt]
		desc.Addr = ^uintptr(0)
		desc.Status &^= RxStatusDD
		hugetlb.Free(q.Buf[q.rdt])
		q.Buf[q.rdt] = nil
	}
}

func (q *RxQueue) Sync() {
	rdh := int(q.Reg.Read(RDH(q.ID)))
	if rdh < q.NumDesc-1 {
		q.rdh = rdh
	} else {
		q.rdh = q.NumDesc - 1
	}
	q.Reg.Write(RDT(q.ID), uint32(q.rdt))
}

func (q *RxQueue) InitBuf() (uintptr, error) {
	size := q.NumDesc * SizeofRxDesc
	desc, phys, err := hugetlb.Alloc(size)
	if err != nil {
		return 0, err
	}
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&q.Ring))
	hdr.Data = uintptr(unsafe.Pointer(&desc[0]))
	hdr.Cap = q.NumDesc
	hdr.Len = q.NumDesc

	q.Buf = make([][]byte, q.NumDesc)

	return phys, nil
}
