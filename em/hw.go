package em

import (
	"errors"

	"uiosample/pci"
)

// Statistics counters collected by the MAC
type HWStats struct {
	crcerrs  uint64
	algnerrc uint64
	symerrs  uint64
	rxerrc   uint64
	mpc      uint64
	scc      uint64
	ecol     uint64
	mcc      uint64
	latecol  uint64
	colc     uint64
	dc       uint64
	tncrs    uint64
	sec      uint64
	cexterr  uint64
	rlec     uint64
	xonrxc   uint64
	xontxc   uint64
	xoffrxc  uint64
	xofftxc  uint64
	fcruc    uint64
	prc64    uint64
	prc127   uint64
	prc255   uint64
	prc511   uint64
	prc1023  uint64
	prc1522  uint64
	gprc     uint64
	bprc     uint64
	mprc     uint64
	gptc     uint64
	gorc     uint64
	gotc     uint64
	rnbc     uint64
	ruc      uint64
	rfc      uint64
	roc      uint64
	rjc      uint64
	mgprc    uint64
	mgpdc    uint64
	mgptc    uint64
	tor      uint64
	tot      uint64
	tpr      uint64
	tpt      uint64
	ptc64    uint64
	ptc127   uint64
	ptc255   uint64
	ptc511   uint64
	ptc1023  uint64
	ptc1522  uint64
	mptc     uint64
	bptc     uint64
	tsctc    uint64
	tsctfc   uint64
	iac      uint64
	icrxptc  uint64
	icrxatc  uint64
	ictxptc  uint64
	ictxatc  uint64
	ictxqec  uint64
	ictxqmtc uint64
	icrxdmtc uint64
	icrxoc   uint64
	cbtmpc   uint64
	htdpmc   uint64
	cbrdpc   uint64
	cbrmpc   uint64
	rpthc    uint64
	hgptc    uint64
	htcbdpc  uint64
	hgorc    uint64
	hgotc    uint64
	lenerrs  uint64
	scvpc    uint64
	hrmpc    uint64
	doosync  uint64
	o2bgptc  uint64
	o2bspc   uint64
	b2ospc   uint64
	b2ogprc  uint64
}

// RevisionID
const (
	Revision0 uint8 = iota
	Revision1
	Revision2
	Revision3
	Revision4
)

type HW struct {
	MAC MACInfo
	FC  FCInfo
	PHY PHYInfo
	NVM NVMInfo
	Bus BusInfo
	MBX MBXInfo

	DeviceID          DeviceID
	SubsystemVendorID uint16
	SubsystemDeviceID uint16
	VendorID          uint16
	RevisionID        uint8

	BAR0 pci.Resource
	BAR1 pci.Resource

	VNIC bool
	Spec interface{}
}

func NewHW(id DeviceID, bar0 pci.Resource, bar1 pci.Resource) (*HW, error) {
	hw := new(HW)
	hw.DeviceID = id
	hw.BAR0 = bar0
	hw.BAR1 = bar1
	hw.MAC.Type = MACTypeGet(id)
	return hw, nil
}

func (hw *HW) RegRead(reg int) uint32 {
	return hw.BAR0.Read32(reg)
}

func (hw *HW) RegWrite(reg int, val uint32) {
	hw.BAR0.Write32(reg, val)
}

func (hw *HW) RegMaskWrite(reg int, val, mask uint32) {
	hw.BAR0.MaskWrite32(reg, val, mask)
}

func (hw *HW) RegWriteFlush() {
	hw.RegRead(STATUS)
}

// e1000_setup_init_funcs
func SetupInitFuncs(hw *HW, initdev bool) error {
	// e1000_init_mac_ops_generic(hw)
	// e1000_init_phy_ops_generic(hw)
	// e1000_init_nvm_ops_generic(hw)
	// e1000_init_mbx_ops_generic(hw)

	switch hw.MAC.Type {
	case MACType82542:
		//I82542Init(hw)
	case MACType82543, MACType82544:
		//I82543Init(hw)
	case MACType82540, MACType82545, MACType82545Rev3, MACType82546, MACType82546Rev3:
		I82540Init(hw)
	case MACType82541, MACType82541Rev2, MACType82547, MACType82547Rev2:
		//I82541Init(hw)
	case MACType82571, MACType82572, MACType82573, MACType82574, MACType82583:
		//I82571Init(hw)
	case MACType80003es2lan:
		//I80003es2lanInit(hw)
	case MACTypeIch8lan, MACTypeIch9lan, MACTypeIch10lan, MACTypePchlan, MACTypePch2lan, MACTypePch_lpt, MACTypePch_spt, MACTypePch_cnp, MACTypePch_adp:
		//ICH8lanInit(hw)
	case MACType82575, MACType82576, MACType82580, MACTypeI350, MACTypeI354:
		//I82575Init(hw)
	case MACTypeI210, MACTypeI211:
		//I210Init(hw)
	case MACTypeVfadapt:
		//VFInit(hw)
	case MACTypeVfadaptI350:
		//VFInit(hw)
	default:
		return errors.New("not supported")
	}

	if initdev {
		err := hw.MAC.Op.InitParams()
		if err != nil {
			return err
		}
		err = hw.NVM.Op.InitParams()
		if err != nil {
			return err
		}
		err = hw.PHY.Op.InitParams()
		if err != nil {
			return err
		}
		//err = hw.MBX.Op.InitParams()
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

func InitHWBase(hw *HW) error {
	mac := &hw.MAC

	// Setup the receive address
	InitRxAddrs(hw, int(mac.RAREntryCount))

	// Zero out the Multicast HASH table
	for i := 0; i < int(mac.MTARegCount); i++ {
		hw.RegWrite(MTA+(i<<2), 0)
	}

	// Zero out the Unicast HASH table
	for i := 0; i < int(mac.UTARegCount); i++ {
		hw.RegWrite(UTA+(i<<2), 0)
	}

	// Setup link and flow control
	err := mac.Op.SetupLink()

	// Clear all of the statistics registers (clear on read).  It is
	// important that we do this after we have tried to establish link
	// because the symbol error count will increment wildly if there
	// is no link.
	ClearHWCounters(hw)

	return err
}
