package ethdev

type DeviceInfo struct {
	DriverName        string
	MinMTU            uint16
	MaxMTU            uint16
	MaxRxQueue        int
	MaxTxQueue        int
	MaxMACAddrs       int
	MaxVFs            int
	RxOffloadCap      RxOffloadCap
	TxOffloadCap      TxOffloadCap
	RxQueueOffloadCap RxOffloadCap
	TxQueueOffloadCap RxOffloadCap
	RxConfig          RxConfig
	TxConfig          TxConfig
	LinkSpeedCap      LinkSpeedCap
	RxQueueNums       int
	TxQueueNums       int
}
