package em

import (
	"errors"
	"log"

	"uiosample/ethdev"
	"uiosample/pci"
)

const ETHER_TYPE_VLAN = 0x8100
const ETHER_MAX_LEN = 1518
const MAX_BUF_SIZE = 2048

const FC_PAUSE_TIME = 0x0680
const FCSetting = FCModeFull

type Driver struct {
	Dev     *pci.Device
	Logger  *log.Logger
	Config  *ethdev.Config
	link    *Link
	led     *LED
	counter *ethdev.CounterGroup
	HW      *HW
	Reg     Reg
	MAC     [][6]byte
	nrxq    int
	ntxq    int
	rxq     [1]RxQueue
	txq     [1]TxQueue
}

// int eth_em_dev_init(struct rte_eth_dev *eth_dev)
func AttachDriver(dev *pci.Device, logger *log.Logger) (*Driver, error) {
	d := new(Driver)
	bar0, err := dev.GetResource(0)
	if err != nil {
		return nil, err
	}
	d.Reg = Reg{res: bar0}
	d.Dev = dev

	if logger == nil {
		d.Logger = log.Default()
	} else {
		d.Logger = logger
	}

	v, err := dev.Config.Read16(2)
	if err != nil {
		return nil, err
	}
	deviceid := DeviceID(v)
	var bar1 pci.Resource
	if deviceid.IsICH8() {
		res, err := dev.GetResource(1)
		if err != nil {
			return nil, err
		}
		bar1 = res
	}
	hw, err := NewHW(deviceid, bar0, bar1)
	if err != nil {
		return nil, err
	}
	d.HW = hw
	/*
		err = SetupInitFuncs(hw, true)
		if err != nil {
			return nil, err
		}
		err = d.HWInit()
		if err != nil {
			return nil, err
		}
		d.link = NewLink(hw)
		d.led = NewLED(&hw.MAC)
		d.MAC = [][6]byte{d.HW.MAC.Addr}
	*/
	d.link = NewLink(d.HW)
	d.led = NewLED(&d.HW.MAC)
	d.counter = NewCounterGroup(d.HW)

	return d, nil
}

// int eth_em_dev_uninit(struct rte_eth_dev *)
func (d *Driver) Detach() {
	d.Close()
}

// int eth_em_infos_get(struct rte_eth_dev *dev, struct rte_eth_dev_info *dev_info)
func (d *Driver) DeviceInfo() (*ethdev.DeviceInfo, error) {
	info := &ethdev.DeviceInfo{}
	info.MaxMACAddrs = int(d.HW.MAC.RAREntryCount)

	info.MaxRxQueue = 1
	info.MaxTxQueue = 1

	info.RxQueueOffloadCap = RxQueueOffloadCap
	info.RxOffloadCap = RxQueueOffloadCap
	info.TxQueueOffloadCap = TxQueueOffloadCap
	info.TxOffloadCap = TxQueueOffloadCap

	info.LinkSpeedCap = ethdev.LinkSpeedCap10MHalf
	info.LinkSpeedCap |= ethdev.LinkSpeedCap10M
	info.LinkSpeedCap |= ethdev.LinkSpeedCap100MHalf
	info.LinkSpeedCap |= ethdev.LinkSpeedCap100M
	info.LinkSpeedCap |= ethdev.LinkSpeedCap1G

	return info, nil
}

func (d *Driver) Configure(nrxq, ntxq int, conf *ethdev.Config) error {
	d.nrxq = nrxq
	d.ntxq = ntxq
	d.Config = conf
	d.link.conf = conf
	d.HW.VNIC = conf.VNIC
	err := SetupInitFuncs(d.HW, true)
	if err != nil {
		return err
	}
	err = d.HWInit()
	if err != nil {
		return err
	}
	d.MAC = [][6]byte{d.HW.MAC.Addr}
	return nil
}

// int eth_em_rx_queue_setup(struct rte_eth_dev *dev,
//
//	uint16_t queue_idx,
//	uint16_t nb_desc,
//	unsigned int socket_id,
//	const struct rte_eth_rxconf *rx_conf,
//	struct rte_mempool *mp)
func (d *Driver) RxQueueSetup(qid, ndesc int, conf *ethdev.RxConfig) error {
	if qid >= len(d.rxq) {
		return nil
	}
	q := &d.rxq[qid]
	q.ID = qid
	q.NumDesc = ndesc
	q.Threshold = conf.Threshold
	q.Reg = d.Reg
	addr, err := q.InitBuf()
	if err != nil {
		return err
	}
	q.RingAddr = addr
	return nil
}

// int eth_em_tx_queue_setup(struct rte_eth_dev *dev,
//
//	uint16_t queue_idx,
//	uint16_t nb_desc,
//	unsigned int socket_id,
//	const struct rte_eth_txconf *tx_conf)
func (d *Driver) TxQueueSetup(qid, ndesc int, conf *ethdev.TxConfig) error {
	if qid >= len(d.txq) {
		return nil
	}
	q := &d.txq[qid]
	q.ID = qid
	q.NumDesc = ndesc
	q.Threshold = conf.Threshold
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

// int eth_em_start(struct rte_eth_dev *dev)
func (d *Driver) Start() error {
	hw := d.HW
	mac := &d.HW.MAC
	phy := &d.HW.PHY

	d.Stop()

	hw.PHY.Op.PowerUp()
	hw.MAC.Op.SetupLink()

	// Set default PBA value
	SetPBA(hw)

	// Put the address into the Receive Address Array
	hw.MAC.Op.SetRAR(mac.Addr, 0)

	// With the 82571 adapter, RAR[0] may be overwritten
	// when the other port is reset, we make a duplicate
	// in RAR[14] for that eventuality, this assures
	// the interface continues to function.
	if mac.Type == MACType82571 {
		// e1000_set_laa_state_82571(hw, TRUE)
		hw.MAC.Op.SetRAR(mac.Addr, RAR_ENTRIES-1)
	}

	// Initialize the hardware
	err := d.HardwareInit()
	if err != nil {
		return err
	}

	hw.RegWrite(VET, ETHER_TYPE_VLAN)

	// Configure for OS presence
	d.InitManageAbility()

	d.TxInit()
	err = d.RxInit()
	if err != nil {
		d.ClearQueues()
		return err
	}

	ClearHWCounters(hw)

	// VLANのオフロードを設定する。
	// mask := RTE_ETH_VLAN_STRIP_MASK | RTE_ETH_VLAN_FILTER_MASK |
	//         RTE_ETH_VLAN_EXTEND_MASK
	// err := eth_em_vlan_offload_set(dev, mask)
	// if err != nil {
	//	d.ClearQueues()
	//	return err
	// }

	// Set Interrupt Throttling Rate to maximum allowed value.
	hw.RegWrite(ITR, 0xffff)

	// Setup link speed and duplex
	speedcap := d.Config.LinkSpeedCap
	if speedcap == ethdev.LinkSpeedCapAutoneg {
		phy.AutonegAdvertised = ALL_SPEED_DUPLEX
		mac.Autoneg = true
	} else {
		num_speeds := 0
		autoneg := speedcap&ethdev.LinkSpeedCapFixed == 0

		// Reset
		phy.AutonegAdvertised = 0

		if speedcap&ethdev.LinkSpeedCap10MHalf != 0 {
			phy.AutonegAdvertised |= ADVERTISE_10_HALF
			num_speeds++
		}
		if speedcap&ethdev.LinkSpeedCap10M != 0 {
			phy.AutonegAdvertised |= ADVERTISE_10_FULL
			num_speeds++
		}
		if speedcap&ethdev.LinkSpeedCap100MHalf != 0 {
			phy.AutonegAdvertised |= ADVERTISE_100_HALF
			num_speeds++
		}
		if speedcap&ethdev.LinkSpeedCap100M != 0 {
			phy.AutonegAdvertised |= ADVERTISE_100_FULL
			num_speeds++
		}
		if speedcap&ethdev.LinkSpeedCap1G != 0 {
			phy.AutonegAdvertised |= ADVERTISE_1000_FULL
			num_speeds++
		}
		if num_speeds == 0 || (!autoneg && num_speeds > 1) {
			d.ClearQueues()
			return errors.New("invalid advertised speeds")
		}
		// Set/reset the mac.autoneg based on the link speed,
		// fixed or not
		if !autoneg {
			mac.Autoneg = false
			mac.ForcedSpeedDuplex = phy.AutonegAdvertised
		} else {
			mac.Autoneg = true
		}
	}

	hw.MAC.Op.SetupLink()

	d.RxTxControl(true)

	d.link.UpdateLink(false)
	return nil
}

// int eth_em_stop(struct rte_eth_dev *dev)
func (d *Driver) Stop() error {
	d.RxTxControl(false)

	d.HW.MAC.Op.ResetHW()

	switch d.HW.MAC.Type {
	case MACTypePch_spt, MACTypePch_cnp:
		// Flush desc rings for i219
		// em_flush_desc_rings(dev)
	}

	if d.HW.MAC.Type >= MACType82544 {
		d.HW.RegWrite(WUC, 0)
	}

	// Power down the phy. Needed to make the link go down
	d.HW.PHY.Op.PowerDown()

	d.ClearQueues()
	return nil
}

// int eth_em_close(struct rte_eth_dev *dev)
func (d *Driver) Close() error {
	d.Stop()

	d.FreeQueues()

	d.HW.PHY.Op.Reset()

	d.ReleaseManageAbility()

	d.HWControlRelease()
	return nil
}

func (d *Driver) Reset() error {
	// not support
	return nil
}

// int eth_em_promiscuous_enable(struct rte_eth_dev *dev)
// int eth_em_promiscuous_disable(struct rte_eth_dev *dev)
// int eth_em_allmulticast_enable(struct rte_eth_dev *dev)
// int eth_em_allmulticast_disable(struct rte_eth_dev *dev)
func (d *Driver) SetPromisc(unicast, multicast bool) error {
	x := d.HW.RegRead(RCTL)
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
	d.HW.RegWrite(RCTL, x)
	return nil
}

func (d *Driver) GetMACAddr() ([6]byte, error) {
	return d.HW.MAC.Addr, nil
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

func PMD_ROUNDUP(x, y uint32) uint32 {
	return (x + y - 1) / y * y
}

func (d *Driver) GetRxBufferSize() uint32 {
	hw := d.HW
	pba := hw.RegRead(PBA)
	return (pba & 0xffff) << 10
}

// int em_hardware_init(struct e1000_hw *hw)
func (d *Driver) HardwareInit() error {
	hw := d.HW
	// Issue a global reset
	hw.MAC.Op.ResetHW()

	// Let the firmware know the OS is in control
	d.HWControlAcquire()

	// These parameters control the automatic generation (Tx) and
	// response (Rx) to Ethernet PAUSE frames.
	// - High water mark should allow for at least two standard size (1518)
	//   frames to be received after sending an XOFF.
	// - Low water mark works best when it is very near the high water mark.
	//   This allows the receiver to restart by sending XON when it has
	//   drained a bit. Here we use an arbitrary value of 1500 which will
	//   restart after one full frame is pulled from the buffer. There
	//   could be several smaller frames in the buffer and if so they will
	//   not trigger the XON until their total number reduces the buffer
	//   by 1500.
	// - The pause time is fairly large at 1000 x 512ns = 512 usec.

	size := d.GetRxBufferSize()
	hw.FC.HighWater = size - PMD_ROUNDUP(ETHER_MAX_LEN*2, 1024)
	hw.FC.LowWater = hw.FC.HighWater - 1500

	if hw.MAC.Type == MACType80003es2lan {
		hw.FC.PauseTime = ^uint16(0)
	} else {
		hw.FC.PauseTime = FC_PAUSE_TIME
	}

	hw.FC.SendXON = true

	// Set Flow control, use the tunable location if sane
	if FCSetting <= FCModeFull {
		hw.FC.RequestedMode = FCSetting
	} else {
		hw.FC.RequestedMode = FCModeNone
	}

	// Workaround: no TX flow ctrl for PCH
	if hw.MAC.Type == MACTypePchlan {
		hw.FC.RequestedMode = FCModeRxPause
	}

	// Override - settings for PCH2LAN, ya its magic :)
	if hw.MAC.Type == MACTypePch2lan {
		hw.FC.HighWater = 0x5c20
		hw.FC.LowWater = 0x5048
		hw.FC.PauseTime = 0x0650
		hw.FC.RefreshTime = 0x0400
	} else if hw.MAC.Type == MACTypePch_lpt ||
		hw.MAC.Type == MACTypePch_spt ||
		hw.MAC.Type == MACTypePch_cnp {
		hw.FC.RequestedMode = FCModeFull
	}

	err := hw.MAC.Op.InitHW()
	if err != nil {
		return err
	}
	hw.MAC.Op.CheckForLink()
	return nil
}

// void em_init_manageability(struct e1000_hw *hw)
func (d *Driver) InitManageAbility() {
	hw := d.HW
	if !EnableManagePT(hw) {
		return
	}
	manc2h := hw.RegRead(MANC2H)
	manc := hw.RegRead(MANC)

	// disable hardware interception of ARP
	manc &^= MANC_ARP_EN

	// enable receiving management packets to the host
	manc |= MANC_EN_MNG2HOST
	manc2h |= 1 << 5 // Mng Port 623
	manc2h |= 1 << 6 // Mng Port 664
	hw.RegWrite(MANC2H, manc2h)
	hw.RegWrite(MANC, manc)
}

// void em_release_manageability(struct e1000_hw *hw)
func (d *Driver) ReleaseManageAbility() {
	hw := d.HW
	if !EnableManagePT(hw) {
		return
	}
	manc := hw.RegRead(MANC)

	// re-enable hardware interception of ARP
	manc |= MANC_ARP_EN
	manc &^= MANC_EN_MNG2HOST

	hw.RegWrite(MANC, manc)
}

// void em_hw_control_acquire(struct e1000_hw *hw)
func (d *Driver) HWControlAcquire() {
	hw := d.HW
	// Let firmware know the driver has taken over
	if hw.MAC.Type == MACType82573 {
		swsm := hw.RegRead(SWSM)
		hw.RegWrite(SWSM, swsm|SWSM_DRV_LOAD)
	} else {
		ctrl := hw.RegRead(CTRL_EXT)
		hw.RegWrite(CTRL_EXT, ctrl|CTRL_EXT_DRV_LOAD)
	}
}

// void em_hw_control_release(struct e1000_hw *hw)
func (d *Driver) HWControlRelease() {
	hw := d.HW
	// Let firmware taken over control of h/w
	if hw.MAC.Type == MACType82573 {
		swsm := hw.RegRead(SWSM)
		hw.RegWrite(SWSM, swsm&^SWSM_DRV_LOAD)
	} else {
		ctrl := hw.RegRead(CTRL_EXT)
		hw.RegWrite(CTRL_EXT, ctrl&^CTRL_EXT_DRV_LOAD)
	}
}

// int em_hw_init(struct e1000_hw *hw)
func (d *Driver) HWInit() error {
	hw := d.HW
	err := hw.MAC.Op.InitParams()
	if err != nil {
		return err
	}
	err = hw.NVM.Op.InitParams()
	if err != nil {
		return err
	}
	err = hw.PHY.Op.InitParams()
	if err != nil {
		return err
	}
	hw.MAC.Op.GetBusInfo()

	hw.MAC.Autoneg = true
	hw.PHY.AutonegWaitToComplete = false
	hw.PHY.AutonegAdvertised = ALL_SPEED_DUPLEX

	// TODO:
	// e1000_init_script_state_82541(hw, TRUE)
	// e1000_set_tbi_compatibility_82543(hw, TRUE)

	// Copper options
	if hw.PHY.MediaType == MediaTypeCopper {
		hw.PHY.MDIX = 0
		hw.PHY.DisablePolarityCorrection = false
		hw.PHY.MSType = MSTypeHwDefault
	}

	// Start from a known state, this is important in reading the nvm
	// and mac from that.
	hw.MAC.Op.ResetHW()

	// Make sure we have a good EEPROM before we read from it
	if hw.NVM.Op.Validate() != nil {
		// Some PCI-E parts fail the first check due to
		// the link being in sleep state, call it again,
		// if it fails a second time its a real issue.
		err := hw.NVM.Op.Validate()
		if err != nil {
			d.HWControlRelease()
			return err
		}
	}

	// Read the permanent MAC address out of the EEPROM
	err = hw.MAC.Op.ReadMACAddr()
	if err != nil {
		d.HWControlRelease()
		return err
	}

	// Now initialize the hardware
	err = d.HardwareInit()
	if err != nil {
		d.HWControlRelease()
		return err
	}

	hw.MAC.GetLinkStatus = true

	// Indicate SOL/IDER usage
	hw.PHY.Op.CheckResetBlock()

	return nil
}

// void eth_em_rxtx_control(struct rte_eth_dev *dev, bool enable)
func (d *Driver) RxTxControl(enable bool) {
	hw := d.HW
	tctl := hw.RegRead(TCTL)
	rctl := hw.RegRead(RCTL)
	if enable {
		// enable Tx/Rx
		tctl |= TCTL_EN
		rctl |= RCTL_EN
	} else {
		// disable Tx/Rx
		tctl &^= TCTL_EN
		rctl &^= RCTL_EN
	}
	hw.RegWrite(TCTL, tctl)
	hw.RegWrite(RCTL, rctl)
	hw.RegWriteFlush()
}

// int eth_em_rx_init(struct rte_eth_dev *dev)
func (d *Driver) RxInit() error {
	hw := d.HW
	rxmode := &d.Config.Rx

	// Make sure receives are disabled while setting
	// up the descriptor ring.
	rctl := hw.RegRead(RCTL)
	hw.RegWrite(RCTL, rctl&^RCTL_EN)

	rfctl := hw.RegRead(RFCTL)

	// Disable extended descriptor type.
	rfctl &^= RFCTL_EXTEN
	// Disable accelerated acknowledge
	if hw.MAC.Type == MACType82574 {
		rfctl |= RFCTL_ACK_DIS
	}

	hw.RegWrite(RFCTL, rfctl)

	// XXX TEMPORARY WORKAROUND: on some systems with 82573
	// long latencies are observed, like Lenovo X60. This
	// change eliminates the problem, but since having positive
	// values in RDTR is a known source of problems on other
	// platforms another solution is being sought.
	if hw.MAC.Type == MACType82573 {
		hw.RegWrite(RDTR, 0x20)
	}

	// Determine RX bufsize.
	bsize, ok := d.rctlBsize(MAX_BUF_SIZE)
	if !ok {
		return errors.New("not found")
	}
	rctl |= bsize

	// Configure and enable each RX queue.
	for i := 0; i < d.nrxq; i++ {
		rxq := &d.rxq[i]

		// Allocate buffers for descriptor rings and setup queue
		//err := alloc_rx_queue_mbufs(rxq)
		//if err != nil {
		//	return err
		//}

		addr := rxq.RingAddr
		hw.RegWrite(RDLEN(i), uint32(rxq.NumDesc)*SizeofRxDesc)
		hw.RegWrite(RDBAH(i), uint32(addr>>32))
		hw.RegWrite(RDBAL(i), uint32(addr))

		hw.RegWrite(RDH(i), 0)
		//hw.RegWrite(RDT(i), uint32(rxq.NumDesc-1))
		hw.RegWrite(RDT(i), 0)

		rxdctl := hw.RegRead(RXDCTL(0))
		rxdctl &= 0xfe000000
		rxdctl |= uint32(rxq.Threshold.Prefetch & 0x3f)
		rxdctl |= uint32(rxq.Threshold.Host&0x3f) << 8
		rxdctl |= uint32(rxq.Threshold.Writeback&0x3f) << 16
		rxdctl |= RXDCTL_GRAN
		hw.RegWrite(RXDCTL(i), rxdctl)
	}

	// Setup the Checksum Register.
	// Receive Full-Packet Checksum Offload is mutually exclusive with RSS.
	rxcsum := hw.RegRead(RXCSUM)
	if rxmode.OffloadCap&ethdev.RxOffloadCapChecksum != 0 {
		rxcsum |= RXCSUM_IPOFL
	} else {
		rxcsum &^= RXCSUM_IPOFL
	}
	hw.RegWrite(RXCSUM, rxcsum)

	// Setup the Receive Control Register.
	if rxmode.OffloadCap&ethdev.RxOffloadCapKeepCRC != 0 {
		rctl &^= RCTL_SECRC // Do not Strip Ethernet CRC.
	} else {
		rctl |= RCTL_SECRC // Strip Ethernet CRC.
	}

	rctl &^= 3 << RCTL_MO_SHIFT
	rctl |= RCTL_EN
	rctl |= RCTL_BAM
	rctl |= RCTL_LBM_NO
	rctl |= RCTL_RDMTS_HALF
	rctl |= hw.MAC.MCFilterType << RCTL_MO_SHIFT

	// Make sure VLAN Filters are off.
	rctl &^= RCTL_VFE
	// Don't store bad packets.
	rctl &^= RCTL_SBP
	// Legacy descriptor type.
	rctl &^= RCTL_DTYP_MASK

	// Enable Receives.
	hw.RegWrite(RCTL, rctl)
	return nil
}

func (d *Driver) rctlBsize(size uint32) (uint32, bool) {
	switch {
	case size > 16384:
		return 0, false
	case size > 8192:
		return RCTL_SZ_16384 | RCTL_BSEX, true
	case size > 4096:
		return RCTL_SZ_8192 | RCTL_BSEX, true
	case size > 2048:
		return RCTL_SZ_4096 | RCTL_BSEX, true
	case size > 1024:
		return RCTL_SZ_2048, true
	case size > 512:
		return RCTL_SZ_1024, true
	case size > 256:
		return RCTL_SZ_512, true
	default:
		return RCTL_SZ_256, true
	}
}

// void eth_em_tx_init(struct rte_eth_dev *dev)
func (d *Driver) TxInit() error {
	hw := d.HW
	for i := 0; i < d.ntxq; i++ {
		txq := &d.txq[i]
		addr := txq.RingAddr
		hw.RegWrite(TDLEN(i), uint32(txq.NumDesc)*SizeofTxDesc)
		hw.RegWrite(TDBAH(i), uint32(addr>>32))
		hw.RegWrite(TDBAL(i), uint32(addr))

		// Setup the HW Tx Head and Tail descriptor pointers.
		hw.RegWrite(TDT(i), 0)
		hw.RegWrite(TDH(i), 0)

		// Setup Transmit threshold registers.
		txdctl := hw.RegRead(TXDCTL(i))
		// bit 22 is reserved, on some models should always be 0,
		// on others  - always 1.
		txdctl &= TXDCTL_COUNT_DESC
		txdctl |= uint32(txq.Threshold.Prefetch & 0x3f)
		txdctl |= uint32(txq.Threshold.Host&0x3f) << 8
		txdctl |= uint32(txq.Threshold.Writeback&0x3f) << 16
		txdctl |= TXDCTL_GRAN
		hw.RegWrite(TXDCTL(i), txdctl)
	}

	// Program the Transmit Control Register.
	tctl := hw.RegRead(TCTL)
	tctl &^= TCTL_CT
	tctl |= TCTL_PSP
	tctl |= TCTL_RTLC
	tctl |= TCTL_EN
	tctl |= COLLISION_THRESHOLD << CT_SHIFT

	// SPT and CNP Si errata workaround to avoid data corruption
	if hw.MAC.Type == MACTypePch_spt {
		iosfpc := hw.RegRead(IOSFPC)
		iosfpc |= RCTL_RDMTS_HEX
		hw.RegWrite(IOSFPC, iosfpc)

		// Dropping the number of outstanding requests from
		// 3 to 2 in order to avoid a buffer overrun.
		tarc := hw.RegRead(TARC(0))
		tarc &^= TARC0_CB_MULTIQ_3_REQ
		tarc |= TARC0_CB_MULTIQ_2_REQ
		hw.RegWrite(TARC(0), tarc)
	}

	// This write will effectively turn on the transmit unit.
	hw.RegWrite(TCTL, tctl)
	return nil
}

// void em_dev_clear_queues(struct rte_eth_dev *dev)
func (d *Driver) ClearQueues() {
	/*
		for _, txq := range d.txq {
			em_tx_queue_release_mbufs(txq)
			em_reset_tx_queue(txq)
		}
		for _, rxq := range d.rxq {
			em_rx_queue_release_mbufs(rxq)
			em_reset_rx_queue(rxq)
		}
	*/
}

// void em_dev_free_queues(struct rte_eth_dev *dev)
func (d *Driver) FreeQueues() {
	/*
		for _, rxq := range d.rxq {
			eth_em_rx_queue_release(dev, i)
		}
		d.nrxq = 0
		for _, txq := range d.txq {
			eth_em_tx_queue_release(dev, i)
		}
		d.ntxq = 0
	*/
}
