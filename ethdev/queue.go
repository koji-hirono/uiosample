package ethdev

type RxQueue interface {
	Start() error
	Stop() error
	Burst([][]byte) int
	Count() int
}

type TxQueue interface {
	Start() error
	Stop() error
	Burst([][]byte) int
	Prepare([][]byte) int
}
