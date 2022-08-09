package ethdev

type RxQueue interface {
	Start() error
	Stop() error
	Do([][]byte) int
	Count() int
}

type TxQueue interface {
	Start() error
	Stop() error
	Do([][]byte) int
	Prepare([][]byte) int
}
