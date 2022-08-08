package ethdev

type RxOffloadCap uint64

const (
	RxOffloadCapVLANStrip RxOffloadCap = 1 << iota
	RxOffloadCapIPv4Checksum
	RxOffloadCapUDPChecksum
	RxOffloadCapTCPChecksum
	RxOffloadCapTCPLRO
	RxOffloadCapQINQStrip
	RxOffloadCapOuterIPv4Checksum
	RxOffloadCapMACSECStrip
	RxOffloadCapHeaderSplit
	RxOffloadCapVLANFilter
	RxOffloadCapVLANExtend
	RxOffloadCapScatter
	RxOffloadCapTimestamp
	RxOffloadCapSecurity
	RxOffloadCapKeepCRC
	RxOffloadCapSCTPChecksum
	RxOffloadCapOuterUDPChecksum
	RxOffloadCapRSSHash
	RxOffloadCapBufferSplit

	RxOffloadCapChecksum = RxOffloadCapIPv4Checksum |
		RxOffloadCapUDPChecksum |
		RxOffloadCapTCPChecksum

	RxOffloadCapVLAN = RxOffloadCapVLANStrip |
		RxOffloadCapVLANFilter |
		RxOffloadCapVLANExtend |
		RxOffloadCapQINQStrip
)
