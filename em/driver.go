package em

import (
	"log"

	"uiosample/pci"
)

type Driver struct {
	Dev    *pci.Device
	Logger *log.Logger
	HW     *HW
	MAC    [][6]byte
}

// int eth_em_dev_init(struct rte_eth_dev *eth_dev)
func OpenDriver(dev *pci.Device, logger *log.Logger) (*Driver, error) {
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
	d.HW = hw

	d.MAC = [][6]byte{d.HW.MAC.Addr}

	return d, nil
}

// int eth_em_dev_uninit(struct rte_eth_dev *)
func (d *Driver) Close() {
	d.Stop()

	// em_dev_free_queues(dev)

	d.HW.ResetPHY()

	// em_release_manageability(hw)
	// em_hw_control_release(hw)
}

// uint32_t eth_em_rx_queue_count(void *rx_queue)
func (d *Driver) RxQueueCount() int {
	return 0
}

// int eth_em_rx_descriptor_status(void *rx_queue, uint16_t offset)
func (d *Driver) RxDescStatus(offset uint16) int {
	return 0
}

// int eth_em_tx_descriptor_status(void *tx_queue, uint16_t offset)
func (d *Driver) TxDescStatus(offset uint16) int {
	return 0
}

// uint16_t eth_em_recv_pkts(void *rx_queue, struct rte_mbuf **rx_pkts, uint16_t nb_pkts)
func (d *Driver) RxPktBurst([][]byte) int {
	return 0
}

// uint16_t eth_em_xmit_pkts(void *tx_queue, struct rte_mbuf **tx_pkts, uint16_t nb_pkts)
func (d *Driver) TxPktBurst([][]byte) int {
	return 0
}

// uint16_t eth_em_prep_pkts(void *tx_queue, struct rte_mbuf **tx_pkts, uint16_t nb_pkts)
func (d *Driver) TxPktPrep([][]byte) int {
	return 0
}

// int eth_em_start(struct rte_eth_dev *dev)
func (d *Driver) Start() error {
	d.Stop()

	d.HW.PowerUpPHY()

	d.HW.SetPBA()

	d.HW.SetRAR(d.HW.MAC.Addr, 0)

	if d.HW.MAC.Type == MACType82571 {
		// e1000_set_laa_state_82571(hw, TRUE)
		// e1000_rar_set(hw, hw->mac.addr, E1000_RAR_ENTRIES - 1)
	}

	// em_hardware_init(hw)

	// VETレジスタにRTE_ETHER_TYPE_VLANを設定する。

	// em_init_manageability(hw)

	// eth_em_tx_init(dev)
	// eth_em_rx_init(dev)

	// e1000_clear_hw_cntrs_base_generic(hw)

	// VLANのオフロードを設定する。

	// ITRレジスタにuint16_max(0xffff)を設定
	d.HW.RegWrite(ITR, 0xffff)

	// autoneg

	d.HW.SetupLink()

	// eth_em_rxtx_control(dev, true)

	d.UpdateLink(false)
	return nil
}

// int eth_em_stop(struct rte_eth_dev *dev)
func (d *Driver) Stop() {
	// eth_em_rxtx_control(dev, false)

	d.HW.ResetHW()

	switch d.HW.MAC.Type {
	case MACTypePch_spt, MACTypePch_cnp:
		// em_flush_desc_rings(dev)
	default:
	}

	if d.HW.MAC.Type >= MACType82544 {
		d.HW.RegWrite(WUC, 0)
	}

	d.HW.PowerDownPHY()

	// em_dev_clear_queues(dev)
}

// int eth_em_promiscuous_enable(struct rte_eth_dev *dev)
// int eth_em_promiscuous_disable(struct rte_eth_dev *dev)
func (d *Driver) SetPromisc(enable bool) {
	x := d.HW.RegRead(RCTL)
	if enable {
		x |= RCTL_UPE
		x |= RCTL_MPE
	} else {
		x &^= RCTL_UPE
		x &^= RCTL_SBP
		// TODO: multicast
		x &^= RCTL_MPE
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

// int eth_em_link_update(struct rte_eth_dev *dev, int wait_to_complete)
func (d *Driver) UpdateLink(block bool) {

}
