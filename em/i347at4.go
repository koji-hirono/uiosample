package em

// Intel I347AT4 Registers
const (
	I347AT4_PCDL        = 0x10 // PHY Cable Diagnostics Length
	I347AT4_PCDC        = 0x15 // PHY Cable Diagnostics Control
	I347AT4_PAGE_SELECT = 0x16
)

// I347AT4 Extended PHY Specific Control Register

// Number of times we will attempt to autonegotiate before downshifting if we
// are the master
const (
	I347AT4_PSCR_DOWNSHIFT_ENABLE = 0x0800
	I347AT4_PSCR_DOWNSHIFT_MASK   = 0x7000
	I347AT4_PSCR_DOWNSHIFT_1X     = 0x0000
	I347AT4_PSCR_DOWNSHIFT_2X     = 0x1000
	I347AT4_PSCR_DOWNSHIFT_3X     = 0x2000
	I347AT4_PSCR_DOWNSHIFT_4X     = 0x3000
	I347AT4_PSCR_DOWNSHIFT_5X     = 0x4000
	I347AT4_PSCR_DOWNSHIFT_6X     = 0x5000
	I347AT4_PSCR_DOWNSHIFT_7X     = 0x6000
	I347AT4_PSCR_DOWNSHIFT_8X     = 0x7000
)

// I347AT4 PHY Cable Diagnostics Control
const (
	I347AT4_PCDC_CABLE_LENGTH_UNIT = 0x0400 // 0=cm 1=meters
)
