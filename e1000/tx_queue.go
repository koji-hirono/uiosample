package e1000

import (
	"reflect"
	"unsafe"

	"uiosample/hugetlb"
)

type TxQueue struct {
	ID       int
	NumDesc  int
	RingAddr uintptr
	Reg      Reg
	Buf      [][]byte
	Ring     []TxDesc
	head     int
	tail     int
}

func (q *TxQueue) Start() error {
	return nil
}

func (q *TxQueue) Stop() error {
	return nil
}

func (q *TxQueue) Prepare([][]byte) int {
	return 0
}

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
	q.Reg.Write(TDT, uint32(q.tail))
	return nil
}

func (q *TxQueue) Can() bool {
	tail := (q.tail + 1) % q.NumDesc
	return tail != q.head
}

func (q *TxQueue) Sync() {
	i := q.head
	q.head = int(q.Reg.Read(TDH))
	for i != q.head {
		hugetlb.Free(q.Buf[i])
		q.Buf[i] = nil
		i = (i + 1) % q.NumDesc
	}
}

func (q *TxQueue) DiscardUnsetPackets() {
	i := q.tail
	q.tail = int(q.Reg.Read(TDT))
	q.head = int(q.Reg.Read(TDH))
	q.Reg.Write(TDT, uint32(q.head))
	for i != q.head {
		i = (i - 1) % q.NumDesc
		hugetlb.Free(q.Buf[i])
		q.Ring[i].Addr = ^uintptr(0)
		q.Ring[i].Length = 0
	}
	q.tail = q.head
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
