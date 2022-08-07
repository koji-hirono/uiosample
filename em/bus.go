package em

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
