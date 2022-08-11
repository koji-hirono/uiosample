package em

import (
	"errors"
)

type BusType int

const (
	BusTypeUnknown BusType = iota
	BusTypePCI
	BusTypePCIX
	BusTypePCIExpress
	BusTypeReserved
)

type BusSpeed int

const (
	BusSpeedUnknown BusSpeed = iota
	BusSpeed33
	BusSpeed66
	BusSpeed100
	BusSpeed120
	BusSpeed133
	BusSpeed2500
	BusSpeed5000
	BusSpeedReserved
)

type BusWidth int

const (
	BusWidthUnknown  BusWidth = 0
	BusWidthPCIEx1            = 1
	BusWidthPCIEx2            = 2
	BusWidthPCIEx4            = 4
	BusWidthPCIEx8            = 8
	BusWidth32                = 9
	BusWidth64                = 10
	BusWidthReserved          = 11
)

const (
	BusFunc0 uint16 = iota
	BusFunc1
	BusFunc2
	BusFunc3
)

type BusInfo struct {
	Type  BusType
	Speed BusSpeed
	Width BusWidth

	Func       uint16
	PCICmdWord uint16
}

const (
	PCIE_LINK_STATUS      = 0x12
	PCIE_LINK_WIDTH_MASK  = 0x3f0
	PCIE_LINK_WIDTH_SHIFT = 4
	PCIE_LINK_SPEED_MASK  = 0x0f
	PCIE_LINK_SPEED_2500  = 0x01
	PCIE_LINK_SPEED_5000  = 0x02
)

func GetBusInfoPCI(hw *HW) error {
	mac := &hw.MAC
	bus := &hw.Bus
	status := hw.RegRead(STATUS)

	// PCI or PCI-X?
	if status&STATUS_PCIX_MODE != 0 {
		bus.Type = BusTypePCIX
	} else {
		bus.Type = BusTypePCI
	}

	// Bus speed
	if bus.Type == BusTypePCI {
		if status&STATUS_PCI66 != 0 {
			bus.Speed = BusSpeed66
		} else {
			bus.Speed = BusSpeed33
		}
	} else {
		switch status & STATUS_PCIX_SPEED {
		case STATUS_PCIX_SPEED_66:
			bus.Speed = BusSpeed66
		case STATUS_PCIX_SPEED_100:
			bus.Speed = BusSpeed100
		case STATUS_PCIX_SPEED_133:
			bus.Speed = BusSpeed133
		default:
			bus.Speed = BusSpeedReserved
		}
	}

	// Bus width
	if status&STATUS_BUS64 != 0 {
		bus.Width = BusWidth64
	} else {
		bus.Width = BusWidth32
	}

	// Which PCI(-X) function?
	mac.Op.SetLANID()

	return nil
}

func GetBusInfoPCIE(hw *HW) error {
	mac := &hw.MAC
	bus := &hw.Bus

	bus.Type = BusTypePCIExpress

	status, err := ReadPCIECapReg(hw, PCIE_LINK_STATUS)
	if err != nil {
		bus.Width = BusWidthUnknown
		bus.Speed = BusSpeedUnknown
	} else {
		switch status & PCIE_LINK_SPEED_MASK {
		case PCIE_LINK_SPEED_2500:
			bus.Speed = BusSpeed2500
		case PCIE_LINK_SPEED_5000:
			bus.Speed = BusSpeed5000
		default:
			bus.Speed = BusSpeedUnknown
		}
		bus.Width = BusWidth((status & PCIE_LINK_WIDTH_MASK) >> PCIE_LINK_WIDTH_SHIFT)
	}

	mac.Op.SetLANID()

	return nil
}

func SetLANIDMultiPortPCI(hw *HW) {
	bus := &hw.Bus
	bus.Func = 0
}

func ReadPCIECapReg(hw *HW, reg uint32) (uint16, error) {
	return 0, errors.New("not implemented")
}
