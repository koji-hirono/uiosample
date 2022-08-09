package em

import (
	"uiosample/ethdev"
)

type RxQueue struct {
	ID      int
	NumDesc int
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
func (q *RxQueue) Do([][]byte) int {
	return 0
}
