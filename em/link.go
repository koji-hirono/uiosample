package em

import (
	"errors"
	"time"

	"uiosample/ethdev"
)

const LINK_UPDATE_CHECK_TIMEOUT = 90   // 9s
const LINK_UPDATE_CHECK_INTERVAL = 100 // ms

type Link struct {
	hw   *HW
	conf *ethdev.Config
}

func NewLink(hw *HW) *Link {
	return &Link{hw: hw}
}

func (l *Link) Up() error {
	return errors.New("not support")
}

func (l *Link) Down() error {
	return errors.New("not support")
}

func (l *Link) Status(block bool) (*ethdev.LinkStatus, error) {
	return l.UpdateLink(block)
}

// int eth_em_link_update(struct rte_eth_dev *dev, int wait_to_complete)
func (l *Link) UpdateLink(block bool) (*ethdev.LinkStatus, error) {
	mac := &l.hw.MAC
	phy := &l.hw.PHY
	conf := l.conf

	mac.GetLinkStatus = true

	// possible wait-to-complete in up to 9 seconds
	var link_up bool
	for count := 0; count < LINK_UPDATE_CHECK_TIMEOUT; count++ {
		// Read the real link status
		switch phy.MediaType {
		case MediaTypeCopper:
			// Do the work to read phy
			mac.Op.CheckForLink()
			link_up = !mac.GetLinkStatus
		case MediaTypeFiber:
			mac.Op.CheckForLink()
			link_up = l.hw.RegRead(STATUS)&STATUS_LU != 0
		case MediaTypeInternalSerdes:
			mac.Op.CheckForLink()
			link_up = mac.SerdesHasLink
		}
		if link_up || !block {
			break
		}
		time.Sleep(LINK_UPDATE_CHECK_INTERVAL * time.Millisecond)
	}

	link := &ethdev.LinkStatus{}

	// Now we check if a transition has happened
	if link_up {
		speed, duplex, err := mac.Op.GetLinkUpInfo()
		if err != nil {
			return nil, err
		}
		if duplex == FULL_DUPLEX {
			link.Duplex = ethdev.LinkDuplexFull
		} else {
			link.Duplex = ethdev.LinkDuplexHalf
		}
		link.Speed = uint32(speed)
		link.Up = true
		if conf.LinkSpeedCap == ethdev.LinkSpeedCapAutoneg {
			link.Autoneg = true
		} else {
			link.Autoneg = false
		}
	} else {
		link.Speed = 0
		link.Duplex = ethdev.LinkDuplexHalf
		link.Up = false
		link.Autoneg = false
	}

	return link, nil
}
