package e1000

import (
	"uiosample/ethdev"
)

type Link struct {
}

func (l Link) Up() error {
	return nil
}

func (l Link) Down() error {
	return nil
}

func (l Link) Status(bool) (*ethdev.LinkStatus, error) {
	var status ethdev.LinkStatus
	return &status, nil
}
