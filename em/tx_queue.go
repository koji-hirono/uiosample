package em

import (
	"uiosample/ethdev"
)

type TxQueue struct {
	ID      int
	NumDesc int
}

const TxQueueOffloadCap = ethdev.TxOffloadCapMultiSegs |
	ethdev.TxOffloadCapVLANInsert |
	ethdev.TxOffloadCapIPv4Checksum |
	ethdev.TxOffloadCapUDPChecksum |
	ethdev.TxOffloadCapTCPChecksum

func (q *TxQueue) Start() error {
	return nil
}

func (q *TxQueue) Stop() error {
	return nil
}

// int eth_em_tx_descriptor_status(void *tx_queue, uint16_t offset)
func (q *TxQueue) Status(offset uint16) int {
	return 0
}

// uint16_t eth_em_xmit_pkts(void *tx_queue, struct rte_mbuf **tx_pkts,
//
//	uint16_t nb_pkts)
func (q *TxQueue) Do([][]byte) int {
	return 0
}

// uint16_t eth_em_prep_pkts(void *tx_queue, struct rte_mbuf **tx_pkts,
//
//	uint16_t nb_pkts)
func (q *TxQueue) Prepare([][]byte) int {
	return 0
}
