package e1000

import (
	"reflect"
	"unsafe"

	"uiosample/hugetlb"
)

type RxQueue struct {
	ID       int
	NumDesc  int
	RingAddr uintptr
	Reg      Reg
	Buf      [][]byte
	Ring     []RxDesc
	head     int
	tail     int
	next     int
}

func (q *RxQueue) Start() error {
	return nil
}

func (q *RxQueue) Stop() error {
	return nil
}

func (q *RxQueue) Count() int {
	return 0
}

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
	length := q.Ring[q.next].Length
	pkt := q.Buf[q.next][:length]
	q.Buf[q.next] = nil
	q.next = (q.next + 1) % q.NumDesc
	return pkt
}

func (q *RxQueue) Can() bool {
	if q.next == q.head {
		return false
	}
	if q.Ring[q.next].Status&RxStatusDD == 0 {
		return false
	}
	return true
}

func (q *RxQueue) CanAddBuf() bool {
	tail := (q.tail + 1) % q.NumDesc
	return tail != q.next
}

func (q *RxQueue) AddBuf(p []byte, phys uintptr) {
	desc := &q.Ring[q.tail]
	desc.Addr = phys
	desc.Status &^= RxStatusDD
	q.Buf[q.tail] = p
	q.tail = (q.tail + 1) % q.NumDesc
}

func (q *RxQueue) FreeAllBuf() {
	for q.tail != q.head {
		q.tail = (q.tail - 1) % q.NumDesc
		desc := &q.Ring[q.tail]
		desc.Addr = ^uintptr(0)
		desc.Status &^= RxStatusDD
		hugetlb.Free(q.Buf[q.tail])
		q.Buf[q.tail] = nil
	}
}

func (q *RxQueue) Sync() {
	head := int(q.Reg.Read(RDH))
	if head < q.NumDesc-1 {
		q.head = head
	} else {
		q.head = q.NumDesc - 1
	}
	q.Reg.Write(RDT, uint32(q.tail))
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
