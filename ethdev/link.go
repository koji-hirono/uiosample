package ethdev

type LinkSpeedCap uint32

const (
	LinkSpeedCapAutoneg LinkSpeedCap = 0
	LinkSpeedCap10MHalf LinkSpeedCap = 1 << iota
	LinkSpeedCap10M
	LinkSpeedCap100MHalf
	LinkSpeedCap100M
	LinkSpeedCap1G
	LinkSpeedCap2_5G
	LinkSpeedCap5G
	LinkSpeedCap10G
	LinkSpeedCap20G
	LinkSpeedCap25G
	LinkSpeedCap40G
	LinkSpeedCap50G
	LinkSpeedCap56G
	LinkSpeedCap100G
	LinkSpeedCap200G
)

type LinkDuplex uint8

const (
	LinkDuplexHalf LinkDuplex = iota
	LinkDuplexFull
)

func (d LinkDuplex) String() string {
	switch d {
	case LinkDuplexHalf:
		return "half"
	case LinkDuplexFull:
		return "full"
	default:
		return ""
	}
}

type LinkStatus struct {
	Speed   uint32
	Duplex  LinkDuplex
	Autoneg bool
	Up      bool
}

type Link interface {
	Up() error
	Down() error
	Status(bool) (*LinkStatus, error)
}
