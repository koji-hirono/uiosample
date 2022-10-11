package e1000

import (
	"log"
	"time"

	"uiosample/ethdev"
	"uiosample/pci"
)

type Driver struct {
	Dev     *pci.Device
	Reg     Reg
	Logger  *log.Logger
	link    *Link
	led     *LED
	counter *ethdev.CounterGroup
	mac     [6]byte
	rxq     [1]RxQueue
	txq     [1]TxQueue
}

func AttachDriver(dev *pci.Device, logger *log.Logger) (*Driver, error) {
	d := new(Driver)
	res, err := dev.GetResource(0)
	if err != nil {
		return nil, err
	}
	d.Dev = dev
	d.Reg = Reg{res: res}
	if logger == nil {
		d.Logger = log.Default()
	} else {
		d.Logger = logger
	}
	d.link = &Link{}
	d.led = &LED{}
	d.counter = NewCounterGroup(d.Reg)
	return d, nil
}

func (d *Driver) Detach() error {
	return d.Close()
}

func (d *Driver) Close() error {
	return nil
}

func (d *Driver) DeviceInfo() (*ethdev.DeviceInfo, error) {
	info := &ethdev.DeviceInfo{}
	info.MaxMACAddrs = 1
	info.MaxRxQueue = 1
	info.MaxTxQueue = 1
	return info, nil
}

func (d *Driver) Configure(nrxq, ntxq int, conf *ethdev.Config) error {
	return nil
}

func (d *Driver) RxQueueSetup(qid, ndesc int, conf *ethdev.RxConfig) error {
	if qid >= len(d.rxq) {
		return nil
	}
	q := &d.rxq[qid]
	q.ID = qid
	q.NumDesc = ndesc
	q.Reg = d.Reg
	addr, err := q.InitBuf()
	if err != nil {
		return err
	}
	q.RingAddr = addr
	return nil
}

func (d *Driver) TxQueueSetup(qid, ndesc int, conf *ethdev.TxConfig) error {
	if qid >= len(d.txq) {
		return nil
	}
	q := &d.txq[qid]
	q.ID = qid
	q.NumDesc = ndesc
	q.Reg = d.Reg
	addr, err := q.InitBuf()
	if err != nil {
		return err
	}
	q.RingAddr = addr
	return nil
}

func (d *Driver) RxQueue(qid int) ethdev.RxQueue {
	if qid >= len(d.rxq) {
		return nil
	}
	return &d.rxq[qid]
}

func (d *Driver) TxQueue(qid int) ethdev.TxQueue {
	if qid >= len(d.txq) {
		return nil
	}
	return &d.txq[qid]
}

func (d *Driver) Start() error {
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
	d.Reg.Read(ICR)
	d.Reg.Write(ICR, ^uint32(0))

	// MAC Addr
	rah0 := d.Reg.Read(RAH0)
	ral0 := d.Reg.Read(RAL0)
	d.mac[0] = byte(ral0)
	d.mac[1] = byte(ral0 >> 8)
	d.mac[2] = byte(ral0 >> 16)
	d.mac[3] = byte(ral0 >> 24)
	d.mac[4] = byte(rah0)
	d.mac[5] = byte(rah0 >> 8)
	d.logf("MAC Address: %x\n", d.mac)

	ctrl := d.Reg.Read(CTRL)
	d.logf("CTRL   : %08x\n", ctrl)
	status := d.Reg.Read(STATUS)
	d.logf("STATUS : %08x\n", status)
	d.logf("  FD   : %x\n", status&0x1)
	d.logf("  LU   : %x\n", (status>>1)&0x1)
	d.logf("  SPEED: %x\n", (status>>6)&0x3)
	d.logf("RCTL   : %08x\n", d.Reg.Read(RCTL))
	d.logf("RDBAL  : %08x\n", d.Reg.Read(RDBAL))
	d.logf("RDBAH  : %08x\n", d.Reg.Read(RDBAH))
	d.logf("RDLEN  : %08x\n", d.Reg.Read(RDLEN))
	d.logf("TCTL   : %08x\n", d.Reg.Read(TCTL))
	d.logf("TDBAL  : %08x\n", d.Reg.Read(TDBAL))
	d.logf("TDBAH  : %08x\n", d.Reg.Read(TDBAH))
	d.logf("TDLEN  : %08x\n", d.Reg.Read(TDLEN))

	return nil
}

func (d *Driver) Stop() error {
	return nil
}

func (d *Driver) Reset() error {
	d.Reg.MaskWrite(CTRL, CTRL_RST, CTRL_RST)
	d.logf("reset...\n")
	// time.Sleep(time.Millisecond * 500)
	for d.Reg.Read(CTRL)&CTRL_RST != 0 {
	}
	d.logf("reset done\n")
	return nil
}

func (d *Driver) SetPromisc(unicast, multicast bool) error {
	x := d.Reg.Read(RCTL)
	if unicast {
		x |= RCTL_UPE
	} else {
		x &^= RCTL_UPE
	}
	if multicast {
		x |= RCTL_MPE
	} else {
		x &^= RCTL_MPE
	}
	d.Reg.Write(RCTL, x)
	return nil
}

func (d *Driver) GetMACAddr() ([6]byte, error) {
	return d.mac, nil
}

func (d *Driver) CounterGroup() *ethdev.CounterGroup {
	return d.counter
}

func (d *Driver) LED() ethdev.LED {
	return d.led
}

func (d *Driver) Link() ethdev.Link {
	return d.link
}

func (d *Driver) IntrDisable() {
	d.Reg.Write(IMC, 0xffffffff)
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
	d.Reg.Write(IMS, val)

	// if TXINT {
	//   set_flags_u32(dev, IMS, IMS_TXDW)
	//   write_u32(dev, TIDV, 1)
	// }

	// if MSIX {
	//   set_flags_u32(dev, IMS, IMS_RXQ0 | IMS_TXQ | IMS_OTHER
	// }
}

func (d *Driver) GlobalConfiguration() {
	// CTRL.FD = 1
	d.Reg.MaskWrite(CTRL, CTRL_FD, CTRL_FD)

	// GCR[22] = 1
	val := uint32(1) << 22
	d.Reg.MaskWrite(GCR, val, val)

	// no flow control
	d.Reg.Write(FCAH, 0)
	d.Reg.Write(FCAL, 0)
	d.Reg.Write(FCT, 0)
	d.Reg.Write(FCTTV, 0)
}

func (d *Driver) LinkUp() {
	v := CTRL_SLU | CTRL_ASDE
	d.Reg.MaskWrite(CTRL, v, v)
	d.logf("waiting linkup.\n")
	for {
		status := d.Reg.Read(STATUS)
		if status&0x2 == 0x2 {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	d.logf("done.\n")
}

func (d *Driver) InitStatRegs() {
	d.Reg.Read(MPC)
	d.Reg.Read(GPRC)
	d.Reg.Read(GPTC)
	d.Reg.Read(GORCL)
	d.Reg.Read(GORCH)
	d.Reg.Read(GOTCL)
	d.Reg.Read(GOTCH)
}

func (d *Driver) InitRx() error {
	rxq := &d.rxq[0]

	d.Reg.Write(RDBAL, uint32(rxq.RingAddr))
	d.Reg.Write(RDBAH, uint32(rxq.RingAddr>>32))

	d.Reg.Write(RDLEN, uint32(rxq.NumDesc*SizeofRxDesc))

	d.Reg.Write(RDH, 0)
	d.Reg.Write(RDT, 0)

	val := RCTL_EN     // Enable
	val |= RCTL_UPE    // Unicast Promiscuous Enable
	val |= RCTL_MPE    // Multicast Promiscuous Enable
	val |= RCTL_BSIZE1 // BSIZE == 11b => 4096 bytes (if BSEX = 1)
	val |= RCTL_BSIZE2
	val |= RCTL_LPE   // Long Packet Enable
	val |= RCTL_BAM   // Broadcast Accept Mode
	val |= RCTL_BSEX  // Buffer Size Extension
	val |= RCTL_SECRC // Strip Ethernet CRC from incoming packet
	d.Reg.Write(RCTL, val)
	return nil
}

func (d *Driver) InitTx() error {
	txq := &d.txq[0]

	d.Reg.Write(TDBAL, uint32(txq.RingAddr))
	d.Reg.Write(TDBAH, uint32(txq.RingAddr>>32))

	d.Reg.Write(TDLEN, uint32(txq.NumDesc*SizeofTxDesc))
	d.Reg.Write(TDH, 0)
	d.Reg.Write(TDT, 0)

	// Enable transmit
	val := TCTL_EN  // Enable
	val |= TCTL_PSP // Pad short packets
	d.Reg.Write(TCTL, val)
	return nil
}

func (d *Driver) logf(format string, v ...interface{}) {
	d.Logger.Printf(format, v...)
}
