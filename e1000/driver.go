package e1000

import (
	"log"
	"reflect"
	"time"
	"unsafe"

	"uiosample/hugetlb"
	"uiosample/pci"
)

type Stat struct {
	MPC  uint64 // Missed Packets Counts
	GPRC uint64 // Good Packets Received Counts
	GPTC uint64 // Good Packest Transmitted Count
	GORC uint64 // Good Octets Received Count
	GOTC uint64 // Good Octets Transmitted Count
}

type Driver struct {
	Dev       *pci.Device
	Logger    *log.Logger
	NumRxDesc int
	NumTxDesc int
	RxBuf     [][]byte
	TxBuf     [][]byte
	RxRing    []RxDesc
	TxRing    []TxDesc
	Mac       []byte
	rdh       int
	rdt       int
	rxnext    int
	tdh       int
	tdt       int
}

func NewDriver(dev *pci.Device, nrxd, ntxd int, logger *log.Logger) *Driver {
	d := new(Driver)
	d.Dev = dev
	if logger == nil {
		d.Logger = log.Default()
	} else {
		d.Logger = logger
	}
	d.NumRxDesc = nrxd
	d.NumTxDesc = ntxd
	return d
}

func (d *Driver) RegRead(reg int) uint32 {
	return d.Dev.Ress[0].Read32(reg)
}

func (d *Driver) RegWrite(reg int, val uint32) {
	d.Dev.Ress[0].Write32(reg, val)
}

func (d *Driver) RegMaskWrite(reg int, val, mask uint32) {
	d.Dev.Ress[0].MaskWrite32(reg, val, mask)
}

func (d *Driver) logf(format string, v ...interface{}) {
	d.Logger.Printf(format, v...)
}

func (d *Driver) IntrDisable() {
	d.RegWrite(IMC, 0xffffffff)
}

func (d *Driver) IntrEnable() {
	// switch {
	// case MSI:
	//   enable_msi()
	// case MSIX:
	//   enable_msix()
	// default:
	//   enable_intx()
	// }

	val := IMS_LSC | IMS_RXT | IMS_RXDMT
	d.RegWrite(IMS, val)

	// if TXINT {
	//   set_flags_u32(dev, IMS, IMS_TXDW)
	//   write_u32(dev, TIDV, 1)
	// }

	// if MSIX {
	//   set_flags_u32(dev, IMS, IMS_RXQ0 | IMS_TXQ | IMS_OTHER
	// }
}

func (d *Driver) Reset() {
	d.RegMaskWrite(CTRL, CTRL_RST, CTRL_RST)
	d.logf("reset...\n")
	// time.Sleep(time.Millisecond * 500)
	for d.RegRead(CTRL)&CTRL_RST != 0 {
	}
	d.logf("reset done\n")
}

func (d *Driver) GlobalConfiguration() {
	// CTRL.FD = 1
	d.RegMaskWrite(CTRL, CTRL_FD, CTRL_FD)

	// GCR[22] = 1
	val := uint32(1) << 22
	d.RegMaskWrite(GCR, val, val)

	// no flow control
	d.RegWrite(FCAH, 0)
	d.RegWrite(FCAL, 0)
	d.RegWrite(FCT, 0)
	d.RegWrite(FCTTV, 0)
}

func (d *Driver) InitStatRegs() {
	d.RegRead(MPC)
	d.RegRead(GPRC)
	d.RegRead(GPTC)
	d.RegRead(GORCL)
	d.RegRead(GORCH)
	d.RegRead(GOTCL)
	d.RegRead(GOTCH)
}

func (d *Driver) LinkUp() {
	v := CTRL_SLU | CTRL_ASDE
	d.RegMaskWrite(CTRL, v, v)
	d.logf("waiting linkup.\n")
	for {
		status := d.RegRead(STATUS)
		if status&0x2 == 0x2 {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	d.logf("done.\n")
}

func (d *Driver) InitRx() error {
	addr, err := d.InitRxBuf()
	if err != nil {
		return err
	}

	d.RegWrite(RDBAL, uint32(addr))
	d.RegWrite(RDBAH, uint32(addr>>32))

	d.RegWrite(RDLEN, uint32(d.NumRxDesc*SizeofTxDesc))

	d.RegWrite(RDH, 0)
	d.RegWrite(RDT, 0)

	val := RCTL_EN     // Enable
	val |= RCTL_UPE    // Unicast Promiscuous Enable
	val |= RCTL_MPE    // Multicast Promiscuous Enable
	val |= RCTL_BSIZE1 // BSIZE == 11b => 4096 bytes (if BSEX = 1)
	val |= RCTL_BSIZE2
	val |= RCTL_LPE   // Long Packet Enable
	val |= RCTL_BAM   // Broadcast Accept Mode
	val |= RCTL_BSEX  // Buffer Size Extension
	val |= RCTL_SECRC // Strip Ethernet CRC from incoming packet
	d.RegWrite(RCTL, val)
	return nil
}

func (d *Driver) InitTx() error {
	addr, err := d.InitTxBuf()
	if err != nil {
		return err
	}

	d.RegWrite(TDBAL, uint32(addr))
	d.RegWrite(TDBAH, uint32(addr>>32))

	d.RegWrite(TDLEN, uint32(d.NumTxDesc*SizeofTxDesc))
	d.RegWrite(TDH, 0)
	d.RegWrite(TDT, 0)

	// Enable transmit
	val := TCTL_EN  // Enable
	val |= TCTL_PSP // Pad short packets
	d.RegWrite(TCTL, val)
	return nil
}

func (d *Driver) InitRxBuf() (uintptr, error) {
	size := d.NumRxDesc * SizeofRxDesc
	desc, phys, err := hugetlb.Alloc(size)
	if err != nil {
		return 0, err
	}
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&d.RxRing))
	hdr.Data = uintptr(unsafe.Pointer(&desc[0]))
	hdr.Cap = d.NumRxDesc
	hdr.Len = d.NumRxDesc

	d.RxBuf = make([][]byte, d.NumRxDesc)

	return phys, nil
}

func (d *Driver) InitTxBuf() (uintptr, error) {
	size := d.NumTxDesc * SizeofTxDesc
	desc, phys, err := hugetlb.Alloc(size)
	if err != nil {
		return 0, err
	}
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&d.TxRing))
	hdr.Data = (uintptr)(unsafe.Pointer(&desc[0]))
	hdr.Cap = d.NumTxDesc
	hdr.Len = d.NumTxDesc

	d.TxBuf = make([][]byte, d.NumTxDesc)

	return phys, nil
}

func (d *Driver) Init() error {
	// 1. Disable Interrupts
	d.IntrDisable()

	// 2. Global reset & general configuration
	d.Reset()
	d.IntrDisable()
	d.GlobalConfiguration()

	// 3. Setup the PHY and the link
	d.LinkUp()

	// 4. Initialize statistical counters
	d.InitStatRegs()

	// 5. Initialize Receive
	err := d.InitRx()
	if err != nil {
		return err
	}

	// 6. Initialize Transmit
	err = d.InitTx()
	if err != nil {
		return err
	}

	// 7. Enable Interrupts (if not pollmode)
	// d.IntrEnable()

	// clear pending intrs
	d.RegRead(ICR)
	d.RegWrite(ICR, ^uint32(0))

	// MAC Addr
	rah0 := d.RegRead(RAH0)
	ral0 := d.RegRead(RAL0)
	mac := make([]byte, 6)
	mac[0] = byte(ral0)
	mac[1] = byte(ral0 >> 8)
	mac[2] = byte(ral0 >> 16)
	mac[3] = byte(ral0 >> 24)
	mac[4] = byte(rah0)
	mac[5] = byte(rah0 >> 8)
	d.logf("MAC Address: %x\n", mac)
	d.Mac = mac

	ctrl := d.RegRead(CTRL)
	d.logf("CTRL   : %08x\n", ctrl)
	status := d.RegRead(STATUS)
	d.logf("STATUS : %08x\n", status)
	d.logf("  FD   : %x\n", status&0x1)
	d.logf("  LU   : %x\n", (status>>1)&0x1)
	d.logf("  SPEED: %x\n", (status>>6)&0x3)
	d.logf("RCTL   : %08x\n", d.RegRead(RCTL))
	d.logf("RDBAL  : %08x\n", d.RegRead(RDBAL))
	d.logf("RDBAH  : %08x\n", d.RegRead(RDBAH))
	d.logf("RDLEN  : %08x\n", d.RegRead(RDLEN))
	d.logf("TCTL   : %08x\n", d.RegRead(TCTL))
	d.logf("TDBAL  : %08x\n", d.RegRead(TDBAL))
	d.logf("TDBAH  : %08x\n", d.RegRead(TDBAH))
	d.logf("TDLEN  : %08x\n", d.RegRead(TDLEN))
	return err
}

func (d *Driver) Tx(pkt []byte) error {
	phys, err := hugetlb.PhysAddr(pkt)
	if err != nil {
		return err
	}
	d.TxBuf[d.tdt] = pkt
	desc := &d.TxRing[d.tdt]
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
	d.tdt = (d.tdt + 1) % d.NumTxDesc
	d.RegWrite(TDT, uint32(d.tdt))
	return nil
}

func (d *Driver) CanTx() bool {
	tdt := (d.tdt + 1) % d.NumTxDesc
	return tdt != d.tdh
}

func (d *Driver) SyncTx() {
	i := d.tdh
	d.tdh = int(d.RegRead(TDH))
	for i != d.tdh {
		hugetlb.Free(d.TxBuf[i])
		d.TxBuf[i] = nil
		i = (i + 1) % d.NumTxDesc
	}
}

func (d *Driver) ServeTx(ch chan []byte) {
	for {
		for len(ch) > 0 && d.CanTx() {
			select {
			case pkt, ok := <-ch:
				if !ok {
					return
				}
				d.Tx(pkt)
			default:
				break
			}
		}
		d.SyncTx()
	}
}

func (d *Driver) DiscardUnsetPackets() {
	i := d.tdt
	d.tdt = int(d.RegRead(TDT))
	d.tdh = int(d.RegRead(TDH))
	d.RegWrite(TDT, uint32(d.tdh))
	for i != d.tdh {
		i = (i - 1) % d.NumTxDesc
		hugetlb.Free(d.TxBuf[i])
		d.TxRing[i].Addr = ^uintptr(0)
		d.TxRing[i].Length = 0
	}
	d.tdt = d.tdh
}

func (d *Driver) OldTx(pkt []byte) int {
	tdt := d.RegRead(TDT)
	tdh := d.RegRead(TDH)
	if tdh == (tdt+1)%uint32(d.NumTxDesc) {
		return 0
	}

	n := len(pkt)
	d.TxBuf[tdt] = pkt
	phys, err := hugetlb.PhysAddr(pkt)
	if err != nil {
		return 0
	}
	d.TxRing[tdt].Addr = phys

	d.TxRing[tdt].Length = uint16(n)
	cmd := TxCommandEOP
	cmd |= TxCommandIFCS
	cmd |= TxCommandRS
	// cmd |= TxCommandIDE
	d.TxRing[tdt].Command = cmd
	d.TxRing[tdt].CSO = 0
	d.TxRing[tdt].Status = 0
	d.TxRing[tdt].CSS = 0
	d.TxRing[tdt].Special = 0

	d.RegWrite(TDT, (tdt+1)%uint32(d.NumTxDesc))

	for d.TxRing[tdt].Status == 0 {
		// d.logf("Tx status: %x\n", d.TxRing[tdt].Status)
	}

	// clear
	hugetlb.Free(pkt)
	d.TxBuf[tdt] = nil
	d.TxRing[tdt].Addr = 0
	return n
}

func (d *Driver) Rx() []byte {
	length := d.RxRing[d.rxnext].Length
	pkt := d.RxBuf[d.rxnext][:length]
	d.RxBuf[d.rxnext] = nil
	d.rxnext = (d.rxnext + 1) % d.NumRxDesc
	return pkt
}

func (d *Driver) CanRx() bool {
	if d.rxnext == d.rdh {
		return false
	}
	if d.RxRing[d.rxnext].Status&RxStatusDD == 0 {
		return false
	}
	return true
}

func (d *Driver) CanAddRxBuf() bool {
	rdt := (d.rdt + 1) % d.NumRxDesc
	return rdt != d.rxnext
}

func (d *Driver) AddRxBuf(p []byte, phys uintptr) {
	desc := &d.RxRing[d.rdt]
	desc.Addr = phys
	desc.Status &^= RxStatusDD
	d.RxBuf[d.rdt] = p
	d.rdt = (d.rdt + 1) % d.NumRxDesc
}

func (d *Driver) FreeAllRxBuf() {
	for d.rdt != d.rdh {
		d.rdt = (d.rdt - 1) % d.NumRxDesc
		desc := &d.RxRing[d.rdt]
		desc.Addr = ^uintptr(0)
		desc.Status &^= RxStatusDD
		hugetlb.Free(d.RxBuf[d.rdt])
		d.RxBuf[d.rdt] = nil
	}
}

func (d *Driver) SyncRx() {
	rdh := int(d.RegRead(RDH))
	if rdh < d.NumRxDesc-1 {
		d.rdh = rdh
	} else {
		d.rdh = d.NumRxDesc - 1
	}
	d.RegWrite(RDT, uint32(d.rdt))
}

func (d *Driver) Serve(ch chan []byte) {
	n := 32
	for {
		d.SyncRx()
		for i := 0; i < n; i++ {
			if !d.CanRx() {
				break
			}
			pkt := d.Rx()
			ch <- pkt
		}
		for d.CanAddRxBuf() {
			p, phys, err := hugetlb.Alloc(2048)
			if err != nil {
				d.logf("alloc failed %v\n", err)
				break
			}
			d.AddRxBuf(p, phys)
		}
	}
}

func (d *Driver) UpdateStat(stat *Stat) {
	stat.MPC += uint64(d.RegRead(MPC))
	stat.GPRC += uint64(d.RegRead(GPRC))
	stat.GPTC += uint64(d.RegRead(GPTC))
	stat.GORC += uint64(d.RegRead(GORCL))
	stat.GORC += uint64(d.RegRead(GORCH)) << 32
	stat.GOTC += uint64(d.RegRead(GOTCL))
	stat.GOTC += uint64(d.RegRead(GOTCH)) << 32
}
