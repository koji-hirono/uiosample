package em

import (
	"errors"
	"log"

	"uiosample/ethdev"
	"uiosample/pci"
)

const ETHER_TYPE_VLAN = 0x8100

type Driver struct {
	Dev    *pci.Device
	Logger *log.Logger
	Config *ethdev.Config
	link   *Link
	led    *LED
	HW     *HW
	MAC    [][6]byte
	rxq    [1]RxQueue
	txq    [1]TxQueue
}

// int eth_em_dev_init(struct rte_eth_dev *eth_dev)
func AttachDriver(dev *pci.Device, logger *log.Logger) (*Driver, error) {
	d := new(Driver)
	bar0, err := dev.GetResource(0)
	if err != nil {
		return nil, err
	}
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
	err = SetupInitFuncs(hw, true)
	if err != nil {
		return nil, err
	}
	err = d.HWInit()
	if err != nil {
		return nil, err
	}
	d.HW = hw
	d.link = NewLink(hw)
	d.led = NewLED(&hw.MAC)
	d.MAC = [][6]byte{d.HW.MAC.Addr}

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

func (d *Driver) Configure(rxd, txd int, conf *ethdev.Config) error {
	d.Config = conf
	d.link.conf = conf
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
func (d *Driver) Stop() {
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
}

// int eth_em_close(struct rte_eth_dev *dev)
func (d *Driver) Close() {
	d.Stop()

	d.FreeQueues()

	d.HW.PHY.Op.Reset()

	d.ReleaseManageAbility()

	d.HWControlRelease()
}

func (d *Driver) Reset() error {
	// not support
	return nil
}

// int eth_em_promiscuous_enable(struct rte_eth_dev *dev)
// int eth_em_promiscuous_disable(struct rte_eth_dev *dev)
func (d *Driver) SetPromisc(enable bool) {
	x := d.HW.RegRead(RCTL)
	if enable {
		x |= RCTL_UPE
	} else {
		x &^= RCTL_UPE
		// XXX: ?
		x &^= RCTL_SBP
	}
	d.HW.RegWrite(RCTL, x)
}

// int eth_em_allmulticast_enable(struct rte_eth_dev *dev)
// int eth_em_allmulticast_disable(struct rte_eth_dev *dev)
func (d *Driver) SetAllMulticast(enable bool) {
	x := d.HW.RegRead(RCTL)
	if enable {
		x |= RCTL_MPE
	} else {
		x &^= RCTL_MPE
	}
	d.HW.RegWrite(RCTL, x)
}

func (d *Driver) GetMACAddr() ([6]byte, error) {
	return d.HW.MAC.Addr, nil
}

func (d *Driver) CounterGroup() *ethdev.CounterGroup {
	return nil
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

const ETHER_MAX_LEN = 1518
const FC_PAUSE_TIME = 0x0680

const FCSetting = FCModeFull

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

	err := d.HWInit()
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
	return nil
}

// void eth_em_tx_init(struct rte_eth_dev *dev)
func (d *Driver) TxInit() error {
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
		dev->data->nb_rx_queues = 0
		for _, txq := range d.txq {
			eth_em_tx_queue_release(dev, i)
		}
		dev->data->nb_tx_queues = 0
	*/
}
