package ethdev

type LED interface {
	On() error
	Off() error
}
