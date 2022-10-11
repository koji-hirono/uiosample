package e1000

type LED struct{}

func (LED) On() error {
	return nil
}

func (LED) Off() error {
	return nil
}
