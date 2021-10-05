package e1000

import (
	"log"
	"reflect"
	"time"
	"unsafe"

	"uiosample/hugetlb"
	"uiosample/pci"
)

type Driver struct {
	Dev       *pci.Device
	NumRxDesc int
	NumTxDesc int
	RxBuf     [][]byte
	TxBuf     [][]byte
	RxRing    []RxDesc
	TxRing    []TxDesc
	Mac       []byte
}

func NewDriver(dev *pci.Device, nrxd, ntxd int) *Driver {
	return &Driver{Dev: dev, NumRxDesc: nrxd, NumTxDesc: ntxd}
}

func (d *Driver) RegRead(reg int) uint32 {
	return d.Dev.Ress[0].Read32(reg)
}

func (d *Driver) RegWrite(reg int, val uint32) {
	d.Dev.Ress[0].Write32(reg, val, ^uint32(0))
}

func (d *Driver) RegMaskWrite(reg int, val, mask uint32) {
	d.Dev.Ress[0].Write32(reg, val, mask)
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
	time.Sleep(time.Millisecond * 500)
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
	d.RegMaskWrite(CTRL, CTRL_SLU, CTRL_SLU)
	log.Println("waiting linkup.")
	for {
		status := d.RegRead(STATUS)
		if status&0x2 == 0x2 {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	log.Println("done.")
}

func (d *Driver) InitRx() {
	addr := d.InitRxBuf()

	d.RegWrite(RDLEN, uint32(d.NumRxDesc*SizeofTxDesc))

	d.RegWrite(RDBAL, uint32(addr))
	d.RegWrite(RDBAH, uint32(addr>>32))

	d.RegWrite(RDH, 0)
	d.RegWrite(RDT, uint32(d.NumRxDesc-1))

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
}

func (d *Driver) InitTx() {
	addr := d.InitTxBuf()

	d.RegWrite(TDBAL, uint32(addr))
	d.RegWrite(TDBAH, uint32(addr>>32))

	d.RegWrite(TDLEN, uint32(d.NumTxDesc*SizeofTxDesc))
	d.RegWrite(TDH, 0)
	d.RegWrite(TDT, 0)

	// Enable transmit
	val := TCTL_EN  // Enable
	val |= TCTL_PSP // Pad short packets
	d.RegWrite(TCTL, val)
}

func (d *Driver) InitRxBuf() uintptr {
	size := d.NumRxDesc * SizeofRxDesc
	desc, err := hugetlb.Alloc(size)
	if err != nil {
		log.Println(err)
		return 0
	}
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&d.RxRing))
	hdr.Data = uintptr(unsafe.Pointer(&desc[0]))
	hdr.Cap = d.NumRxDesc
	hdr.Len = d.NumRxDesc

	d.RxBuf = make([][]byte, d.NumRxDesc)

	for i := 0; i < d.NumRxDesc; i++ {
		size := 4096
		buf, err := hugetlb.Alloc(size)
		if err != nil {
			log.Println(err)
			return 0
		}
		d.RxBuf[i] = buf
		virt := uintptr(unsafe.Pointer(&buf[0]))
		phys, err := hugetlb.VirtToPhys(virt)
		if err != nil {
			log.Println(err)
			return 0
		}
		d.RxRing[i].Addr = phys
	}

	virt := uintptr(unsafe.Pointer(&desc[0]))
	phys, err := hugetlb.VirtToPhys(virt)
	if err != nil {
		log.Println(err)
		return 0
	}

	return phys
}

func (d *Driver) InitTxBuf() uintptr {
	size := d.NumTxDesc * SizeofTxDesc
	desc, err := hugetlb.Alloc(size)
	if err != nil {
		log.Println(err)
		return 0
	}
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&d.TxRing))
	hdr.Data = (uintptr)(unsafe.Pointer(&desc[0]))
	hdr.Cap = d.NumTxDesc
	hdr.Len = d.NumTxDesc

	d.TxBuf = make([][]byte, d.NumTxDesc)

	for i := 0; i < d.NumTxDesc; i++ {
		size := 4096
		buf, err := hugetlb.Alloc(size)
		if err != nil {
			log.Println(err)
			return 0
		}
		d.TxBuf[i] = buf
		virt := uintptr(unsafe.Pointer(&buf[0]))
		phys, err := hugetlb.VirtToPhys(virt)
		if err != nil {
			log.Println(err)
			return 0
		}
		d.TxRing[i].Addr = phys
	}

	virt := uintptr(unsafe.Pointer(&desc[0]))
	phys, err := hugetlb.VirtToPhys(virt)
	if err != nil {
		log.Println(err)
		return 0
	}

	return phys
}

func (d *Driver) Init() {
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
	d.InitRx()

	// 6. Initialize Transmit
	d.InitTx()

	// 7. Enable Interrupts (if not pollmode)
	// enable_interrupt(dev)

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
	log.Printf("MAC Address: %x\n", mac)
	d.Mac = mac

	ctrl := d.RegRead(CTRL)
	log.Printf("CTRL   : %08x\n", ctrl)
	status := d.RegRead(STATUS)
	log.Printf("STATUS : %08x\n", status)
	log.Printf("  FD   : %x\n", status&0x1)
	log.Printf("  LU   : %x\n", (status>>1)&0x1)
	log.Printf("  SPEED: %x\n", (status>>6)&0x3)
	log.Printf("RCTL   : %08x\n", d.RegRead(RCTL))
	log.Printf("RDBAL  : %08x\n", d.RegRead(RDBAL))
	log.Printf("RDBAH  : %08x\n", d.RegRead(RDBAH))
	log.Printf("RDLEN  : %08x\n", d.RegRead(RDLEN))
	log.Printf("TCTL   : %08x\n", d.RegRead(TCTL))
	log.Printf("TDBAL  : %08x\n", d.RegRead(TDBAL))
	log.Printf("TDBAH  : %08x\n", d.RegRead(TDBAH))
	log.Printf("TDLEN  : %08x\n", d.RegRead(TDLEN))
}

func (d *Driver) Tx(pkt []byte) {
	tdt := d.RegRead(TDT)
	tdh := d.RegRead(TDH)
	if tdh == (tdt+1)%uint32(d.NumTxDesc) {
		return
	}

	n := copy(d.TxBuf[tdt], pkt)
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
		log.Printf("Tx status: %x\n", d.TxRing[tdt].Status)
	}
	log.Printf("sned %v bytes\n", n)
}

func (d *Driver) Rx(i int) ([]byte, int) {
	length := d.RxRing[i].Length
	pkt := make([]byte, length)
	copy(pkt, d.RxBuf[i][:length])

	// clear desc
	d.RxRing[i].Status &^= RxStatusDD

	head := d.RegRead(RDH)
	if head != uint32(i) {
		d.RegWrite(RDT, uint32(i))
	}
	i = (i + 1) % d.NumRxDesc
	return pkt, i
}

func (d *Driver) Serve(ch chan []byte) {
	i := 0
	for {
		rdt := d.RegRead(RDT)
		log.Printf("RDT: %v\n", rdt)
		for j := 0; j < d.NumRxDesc; j++ {
			log.Printf("RxDesc[%v]: %#+v\n", j, d.RxRing[j])
		}
		if d.RxRing[i].Status&RxStatusDD != 0 {
			pkt, next := d.Rx(i)
			i = next
			ch <- pkt
		}
		/*
			log.Printf("MPC  : %v\n", d.RegRead(MPC))
			log.Printf("GPRC : %v\n", d.RegRead(GPRC))
			log.Printf("GPTC : %v\n", d.RegRead(GPTC))
			log.Printf("GORCL: %v\n", d.RegRead(GORCL))
			log.Printf("GORCH: %v\n", d.RegRead(GORCH))
			log.Printf("GOTCL: %v\n", d.RegRead(GOTCL))
			log.Printf("GOTCH: %v\n", d.RegRead(GOTCH))
			log.Printf("RxBuf(%v): %x\n", i, d.RxBuf[i][:60])
		*/

		time.Sleep(time.Millisecond * 500)
	}
}
