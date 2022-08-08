package ethdev

type TxOffloadCap uint64

const (
	TxOffloadVLANInsert TxOffloadCap = 1 << iota
	TxOffloadIPv4Checksum
	TxOffloadUDPChecksum
	TxOffloadTCPChecksum
	TxOffloadSCTPChecksum
	TxOffloadTCPTSO
	TxOffloadUDPTSO
	TxOffloadOuterIPv4Checksum
	TxOffloadQINQInsert
	TxOffloadVXLANTunnelTSO
	TxOffloadGRETunnelTSO
	TxOffloadIPIPTunnelTSO
	TxOffloadGENEVETunnelTSO
	TxOffloadMACSECInsert
	TxOffloadMTLockFree
	TxOffloadMultiSegs
	TxOffloadMbufFastFree
	TxOffloadSecurity
	TxOffloadUDPTunnelTSO
	TxOffloadIPTunnelTSO
	TxOffloadOuterUDPChecksum
	TxOffloadSendOnTimestamp
)
