package em

import (
	"errors"
	"time"
)

// PHY ID
const (
	M88E1000_E_PHY_ID   uint32 = 0x01410C50
	M88E1000_I_PHY_ID          = 0x01410C30
	M88E1011_I_PHY_ID          = 0x01410C20
	IGP01E1000_I_PHY_ID        = 0x02A80380
	M88E1111_I_PHY_ID          = 0x01410CC0
	M88E1543_E_PHY_ID          = 0x01410EA0
	M88E1512_E_PHY_ID          = 0x01410DD0
	M88E1112_E_PHY_ID          = 0x01410C90
	I347AT4_E_PHY_ID           = 0x01410DC0
	M88E1340M_E_PHY_ID         = 0x01410DF0
	GG82563_E_PHY_ID           = 0x01410CA0
	IGP03E1000_E_PHY_ID        = 0x02A80390
	IFE_E_PHY_ID               = 0x02A80330
	IFE_PLUS_E_PHY_ID          = 0x02A80320
	IFE_C_E_PHY_ID             = 0x02A80310
	BME1000_E_PHY_ID           = 0x01410CB0
	BME1000_E_PHY_ID_R2        = 0x01410CB1
	I82577_E_PHY_ID            = 0x01540050
	I82578_E_PHY_ID            = 0x004DD040
	I82579_E_PHY_ID            = 0x01540090
	I217_E_PHY_ID              = 0x015400A0
	I82580_I_PHY_ID            = 0x015403A0
	I350_I_PHY_ID              = 0x015403B0
	I210_I_PHY_ID              = 0x01410C00
	IGP04E1000_E_PHY_ID        = 0x02A80391
	BCM54616_E_PHY_ID          = 0x03625D10
	M88_VENDOR                 = 0x0141
)

type PHYType int

const (
	PHYTypeUnknown PHYType = iota
	PHYTypeNone
	PHYTypeM88
	PHYTypeIgp
	PHYTypeIgp2
	PHYTypeGg82563
	PHYTypeIgp3
	PHYTypeIfe
	PHYTypeBm
	PHYType82578
	PHYType82577
	PHYType82579
	PHYTypeI217
	PHYType82580
	PHYTypeVf
	PHYTypeI210
)

// enum e1000_phy_type e1000_get_phy_type_from_id(u32 phy_id)
func PHYTypeGet(phyid uint32) PHYType {
	switch phyid {
	case M88E1000_I_PHY_ID:
		return PHYTypeM88
	case M88E1000_E_PHY_ID:
		return PHYTypeM88
	case M88E1111_I_PHY_ID:
		return PHYTypeM88
	case M88E1011_I_PHY_ID:
		return PHYTypeM88
	case M88E1543_E_PHY_ID:
		return PHYTypeM88
	case M88E1512_E_PHY_ID:
		return PHYTypeM88
	case I347AT4_E_PHY_ID:
		return PHYTypeM88
	case M88E1112_E_PHY_ID:
		return PHYTypeM88
	case M88E1340M_E_PHY_ID:
		return PHYTypeM88

	case IGP01E1000_I_PHY_ID:
		// IGP 1 & 2 share this
		return PHYTypeIgp2
	case GG82563_E_PHY_ID:
		return PHYTypeGg82563
	case IGP03E1000_E_PHY_ID:
		return PHYTypeIgp3

	case IFE_E_PHY_ID:
		return PHYTypeIfe
	case IFE_PLUS_E_PHY_ID:
		return PHYTypeIfe
	case IFE_C_E_PHY_ID:
		return PHYTypeIfe

	case BME1000_E_PHY_ID:
		return PHYTypeBm
	case BME1000_E_PHY_ID_R2:
		return PHYTypeBm

	case I82578_E_PHY_ID:
		return PHYType82578
	case I82577_E_PHY_ID:
		return PHYType82577
	case I82579_E_PHY_ID:
		return PHYType82579
	case I217_E_PHY_ID:
		return PHYTypeI217
	case I82580_I_PHY_ID:
		return PHYType82580
	case I210_I_PHY_ID:
		return PHYTypeI210
	default:
		return PHYTypeUnknown
	}
}

type E1000TRxStatus int

const (
	E1000TRxStatusNotOk E1000TRxStatus = iota
	E1000TRxStatusOk
	E1000TRxStatusUndefined = 0xff
)

type MSType int

const (
	MSTypeHwDefault MSType = iota
	MSTypeForceMaster
	MSTypeForceSlave
	MSTypeAuto
)

type RevPolarity int

const (
	RevPolarityNormal RevPolarity = iota
	RevPolarityReversed
	RevPolarityUndefined = 0xff
)

type SmartSpeed int

const (
	SmartSpeedDefault SmartSpeed = iota
	SmartSpeedOn
	SmartSpeedOff
)

type MediaType int

const (
	MediaTypeUnknown MediaType = iota
	MediaTypeCopper
	MediaTypeFiber
	MediaTypeInternalSerdes
	NumMediaType
)

type PHYStats struct {
	IdleErrors    uint32
	ReceiveErrors uint32
}

const (
	ADVERTISE_10_HALF uint16 = 1 << iota
	ADVERTISE_10_FULL
	ADVERTISE_100_HALF
	ADVERTISE_100_FULL
	ADVERTISE_1000_HALF
	ADVERTISE_1000_FULL

	ALL_SPEED_DUPLEX = ADVERTISE_10_HALF | ADVERTISE_10_FULL | ADVERTISE_100_HALF | ADVERTISE_100_FULL | ADVERTISE_1000_FULL
	ALL_NOT_GIG      = ADVERTISE_10_HALF | ADVERTISE_10_FULL | ADVERTISE_100_HALF | ADVERTISE_100_FULL
	ALL_100_SPEED    = ADVERTISE_100_HALF | ADVERTISE_100_FULL
	ALL_10_SPEED     = ADVERTISE_10_HALF | ADVERTISE_10_FULL
	ALL_HALF_DUPLEX  = ADVERTISE_10_HALF | ADVERTISE_100_HALF

	AUTONEG_ADVERTISE_SPEED_DEFAULT = ALL_SPEED_DUPLEX
)

const CABLE_LENGTH_UNDEFINED uint16 = 0xff

const FIBER_LINK_UP_LIMIT = 50
const COPPER_LINK_UP_LIMIT = 10
const PHY_AUTO_NEG_LIMIT = 45
const PHY_FORCE_LIMIT = 20
const PHY_CFG_TIMEOUT = 100

const MAX_PHY_ADDR uint32 = 8

type PHYInfo struct {
	Op            PHYOp
	PHYType       PHYType
	LocalRx       E1000TRxStatus
	RemoteRx      E1000TRxStatus
	MSType        MSType
	OrigMSType    MSType
	CablePolarity RevPolarity
	SmartSpeed    SmartSpeed

	Addr         uint32
	ID           uint32
	ResetDelayUS time.Duration
	Revision     uint32

	MediaType MediaType

	AutonegAdvertised uint16
	AutonegMask       uint16
	CableLength       uint16
	MaxCableLength    uint16
	MinCableLength    uint16

	MDIX uint8

	DisablePolarityCorrection bool
	IsMDIX                    bool
	PolarityCorrection        bool
	SpeedDowngraded           bool
	AutonegWaitToComplete     bool
}

type PHYOp interface {
	InitParams() error
	Acquire() error
	CfgOnLinkUp() error
	CheckPolarity() error
	CheckResetBlock() error
	Commit() error
	ForceSpeedDuplex() error
	GetCfgDone() error
	GetCableLength() error
	GetInfo() error
	SetPage(uint16) error
	ReadReg(uint32) (uint16, error)
	ReadRegLocked(uint32) (uint16, error)
	ReadRegPage(uint32) (uint16, error)
	Release()
	Reset() error
	SetD0LpluState(bool) error
	SetD3LpluState(bool) error
	WriteReg(uint32, uint16) error
	WriteRegLocked(uint32, uint16) error
	WriteRegPage(uint32, uint16) error
	PowerUp()
	PowerDown()
	ReadI2CByte(uint8, uint8) (byte, error)
	WriteI2CByte(uint8, uint8, byte) error
}

/*
   s32  (*init_params)(struct e1000_hw *);
   s32  (*acquire)(struct e1000_hw *);
   s32  (*cfg_on_link_up)(struct e1000_hw *);
   s32  (*check_polarity)(struct e1000_hw *);
   s32  (*check_reset_block)(struct e1000_hw *);
   s32  (*commit)(struct e1000_hw *);
   s32  (*force_speed_duplex)(struct e1000_hw *);
   s32  (*get_cfg_done)(struct e1000_hw *hw);
   s32  (*get_cable_length)(struct e1000_hw *);
   s32  (*get_info)(struct e1000_hw *);
   s32  (*set_page)(struct e1000_hw *, u16);
   s32  (*read_reg)(struct e1000_hw *, u32, u16 *);
   s32  (*read_reg_locked)(struct e1000_hw *, u32, u16 *);
   s32  (*read_reg_page)(struct e1000_hw *, u32, u16 *);
   void (*release)(struct e1000_hw *);
   s32  (*reset)(struct e1000_hw *);
   s32  (*set_d0_lplu_state)(struct e1000_hw *, bool);
   s32  (*set_d3_lplu_state)(struct e1000_hw *, bool);
   s32  (*write_reg)(struct e1000_hw *, u32, u16);
   s32  (*write_reg_locked)(struct e1000_hw *, u32, u16);
   s32  (*write_reg_page)(struct e1000_hw *, u32, u16);
   void (*power_up)(struct e1000_hw *);
   void (*power_down)(struct e1000_hw *);
   s32 (*read_i2c_byte)(struct e1000_hw *, u8, u8, u8 *);
   s32 (*write_i2c_byte)(struct e1000_hw *, u8, u8, u8);
*/

func GetPHYID(hw *HW) error {
	phy := &hw.PHY
	for i := 0; i < 2; i++ {
		id1, err := phy.Op.ReadReg(PHY_ID1)
		if err != nil {
			return err
		}
		phy.ID = uint32(id1) << 16

		time.Sleep(20 * time.Microsecond)

		id2, err := phy.Op.ReadReg(PHY_ID2)
		if err != nil {
			return err
		}
		phy.ID |= uint32(id2) & PHY_REVISION_MASK
		phy.Revision = uint32(id2) &^ PHY_REVISION_MASK

		if phy.ID != 0 && phy.ID != PHY_REVISION_MASK {
			return nil
		}
	}
	return nil
}

func GetCfgDone(hw *HW) error {
	time.Sleep(10 * time.Millisecond)
	return nil
}

func PHYSWReset(hw *HW) error {
	phy := &hw.PHY

	ctrl, err := phy.Op.ReadReg(PHY_CONTROL)
	if err != nil {
		return err
	}

	ctrl |= MII_CR_RESET
	err = phy.Op.WriteReg(PHY_CONTROL, ctrl)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Microsecond)
	return nil
}

func PHYHWReset(hw *HW) error {
	phy := &hw.PHY
	err := phy.Op.CheckResetBlock()
	if err != nil {
		return err
	}

	err = phy.Op.Acquire()
	if err != nil {
		return err
	}

	ctrl := hw.RegRead(CTRL)
	hw.RegWrite(CTRL, ctrl|CTRL_PHY_RST)
	hw.RegWriteFlush()

	time.Sleep(phy.ResetDelayUS * time.Microsecond)

	hw.RegWrite(CTRL, ctrl)
	hw.RegWriteFlush()

	phy.Op.Release()

	time.Sleep(150 * time.Microsecond)

	return phy.Op.GetCfgDone()
}

func PHYHasLink(hw *HW, n int, interval time.Duration) (bool, error) {
	for i := 0; i < n; i++ {
		// Some PHYs require the PHY_STATUS register to be read
		// twice due to the link bit being sticky.  No harm doing
		// it across the board.
		_, err := hw.PHY.Op.ReadReg(PHY_STATUS)
		if err != nil {
			// If the first read fails, another entity may have
			// ownership of the resources, wait and try again to
			// see if they have relinquished the resources yet.
			time.Sleep(interval * time.Microsecond)
		}
		status, err := hw.PHY.Op.ReadReg(PHY_STATUS)
		if err != nil {
			return false, err
		}
		if status&MII_SR_LINK_STATUS != 0 {
			return true, nil
		}
		time.Sleep(interval * time.Microsecond)
	}
	return false, nil
}

func PHYResetDSP(hw *HW) error {
	err := hw.PHY.Op.WriteReg(M88E1000_PHY_GEN_CONTROL, 0xc1)
	if err != nil {
		return err
	}

	return hw.PHY.Op.WriteReg(M88E1000_PHY_GEN_CONTROL, 0)
}

func ReadPHYRegMDIC(hw *HW, offset uint32) (uint16, error) {
	phy := &hw.PHY

	if offset > MAX_PHY_REG_ADDRESS {
		return 0, errors.New("out of range")
	}

	// Set up Op-code, Phy Address, and register offset in the MDI
	// Control register.  The MAC will take care of interfacing with the
	// PHY to retrieve the desired data.
	mdic := offset << MDIC_REG_SHIFT
	mdic |= phy.Addr << MDIC_PHY_SHIFT
	mdic |= MDIC_OP_READ
	hw.RegWrite(MDIC, mdic)

	// Poll the ready bit to see if the MDI read completed
	// Increasing the time out as testing showed failures with
	// the lower time out
	for i := 0; i < GEN_POLL_TIMEOUT*3; i++ {
		time.Sleep(50 * time.Microsecond)
		mdic = hw.RegRead(MDIC)
		if mdic&MDIC_READY != 0 {
			break
		}
	}
	if mdic&MDIC_READY == 0 {
		return 0, errors.New("MDI Read did not complete")
	}
	if mdic&MDIC_ERROR != 0 {
		return 0, errors.New("MDI Error")
	}
	if (mdic&MDIC_REG_MASK)>>MDIC_REG_SHIFT != offset {
		return 0, errors.New("MDI Read offset error")
	}

	// Allow some time after each MDIC transaction to avoid
	// reading duplicate data in the next MDIC transaction.
	if hw.MAC.Type == MACTypePch2lan {
		time.Sleep(100 * time.Microsecond)
	}

	return uint16(mdic), nil
}

func WritePHYRegMDIC(hw *HW, offset uint32, val uint16) error {
	phy := &hw.PHY

	if offset > MAX_PHY_REG_ADDRESS {
		return errors.New("out of range")
	}

	// Set up Op-code, Phy Address, and register offset in the MDI
	// Control register.  The MAC will take care of interfacing with the
	// PHY to retrieve the desired data.
	mdic := uint32(val)
	mdic |= offset << MDIC_REG_SHIFT
	mdic |= phy.Addr << MDIC_PHY_SHIFT
	mdic |= MDIC_OP_WRITE
	hw.RegWrite(MDIC, mdic)

	// Poll the ready bit to see if the MDI read completed
	// Increasing the time out as testing showed failures with
	// the lower time out
	for i := 0; i < GEN_POLL_TIMEOUT*3; i++ {
		time.Sleep(50 * time.Microsecond)
		mdic = hw.RegRead(MDIC)
		if mdic&MDIC_READY != 0 {
			break
		}
	}
	if mdic&MDIC_READY == 0 {
		return errors.New("MDI Write did not complete")
	}
	if mdic&MDIC_ERROR != 0 {
		return errors.New("MDI Error")
	}
	if (mdic&MDIC_REG_MASK)>>MDIC_REG_SHIFT != offset {
		return errors.New("MDI Write offset error")
	}
	// Allow some time after each MDIC transaction to avoid
	// reading duplicate data in the next MDIC transaction.
	if hw.MAC.Type == MACTypePch2lan {
		time.Sleep(100 * time.Microsecond)
	}

	return nil
}

func DeterminePHYAddress(hw *HW) error {
	hw.PHY.PHYType = PHYTypeUnknown
	for addr := uint32(0); addr < MAX_PHY_ADDR; addr++ {
		hw.PHY.Addr = addr
		for i := 0; i < 10; i++ {
			GetPHYID(hw)
			t := PHYTypeGet(hw.PHY.ID)
			// If phy_type is valid, break - we found our
			// PHY address
			if t != PHYTypeUnknown {
				return nil
			}
			time.Sleep(1 * time.Millisecond)
		}
	}
	return errors.New("not found")
}

func WaitAutoneg(hw *HW) error {
	phy := &hw.PHY
	// Break after autoneg completes or PHY_AUTO_NEG_LIMIT expires.
	for i := PHY_AUTO_NEG_LIMIT; i > 0; i-- {
		_, err := phy.Op.ReadReg(PHY_STATUS)
		if err != nil {
			return err
		}
		status, err := phy.Op.ReadReg(PHY_STATUS)
		if err != nil {
			return err
		}
		if status&MII_SR_AUTONEG_COMPLETE != 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// PHY_AUTO_NEG_TIME expiration doesn't guarantee auto-negotiation
	// has completed.
	return nil
}

func PHYSetupAutoneg(hw *HW) error {
	phy := &hw.PHY

	phy.AutonegAdvertised &= phy.AutonegMask

	// Read the MII Auto-Neg Advertisement Register (Address 4).
	adv, err := phy.Op.ReadReg(PHY_AUTONEG_ADV)
	if err != nil {
		return err
	}

	var ctrl uint16
	if phy.AutonegMask&ADVERTISE_1000_FULL != 0 {
		// Read the MII 1000Base-T Control Register (Address 9).
		x, err := phy.Op.ReadReg(PHY_1000T_CTRL)
		if err != nil {
			return err
		}
		ctrl = x
	}

	// Need to parse both autoneg_advertised and fc and set up
	// the appropriate PHY registers.  First we will parse for
	// autoneg_advertised software override.  Since we can advertise
	// a plethora of combinations, we need to check each bit
	// individually.

	// First we clear all the 10/100 mb speed bits in the Auto-Neg
	// Advertisement Register (Address 4) and the 1000 mb speed bits in
	// the  1000Base-T Control Register (Address 9).
	adv &^= NWAY_AR_100TX_FD_CAPS | NWAY_AR_100TX_HD_CAPS | NWAY_AR_10T_FD_CAPS | NWAY_AR_10T_HD_CAPS
	ctrl &^= CR_1000T_HD_CAPS | CR_1000T_FD_CAPS

	// Do we want to advertise 10 Mb Half Duplex?
	if phy.AutonegAdvertised&ADVERTISE_10_HALF != 0 {
		adv |= NWAY_AR_10T_HD_CAPS
	}

	// Do we want to advertise 10 Mb Full Duplex?
	if phy.AutonegAdvertised&ADVERTISE_10_FULL != 0 {
		adv |= NWAY_AR_10T_FD_CAPS
	}

	// Do we want to advertise 100 Mb Half Duplex?
	if phy.AutonegAdvertised&ADVERTISE_100_HALF != 0 {
		adv |= NWAY_AR_100TX_HD_CAPS
	}

	// Do we want to advertise 100 Mb Full Duplex?
	if phy.AutonegAdvertised&ADVERTISE_100_FULL != 0 {
		adv |= NWAY_AR_100TX_FD_CAPS
	}

	// We do not allow the Phy to advertise 1000 Mb Half Duplex
	if phy.AutonegAdvertised&ADVERTISE_1000_HALF != 0 {
	}

	// Do we want to advertise 1000 Mb Full Duplex?
	if phy.AutonegAdvertised&ADVERTISE_1000_FULL != 0 {
		ctrl |= CR_1000T_FD_CAPS
	}

	// Check for a software override of the flow control settings, and
	// setup the PHY advertisement registers accordingly.  If
	// auto-negotiation is enabled, then software will have to set the
	// "PAUSE" bits to the correct value in the Auto-Negotiation
	// Advertisement Register (PHY_AUTONEG_ADV) and re-start auto-
	// negotiation.
	//
	// The possible values of the "fc" parameter are:
	//      0:  Flow control is completely disabled
	//      1:  Rx flow control is enabled (we can receive pause frames
	//          but not send pause frames).
	//      2:  Tx flow control is enabled (we can send pause frames
	//          but we do not support receiving pause frames).
	//      3:  Both Rx and Tx flow control (symmetric) are enabled.
	//  other:  No software override.  The flow control configuration
	//          in the EEPROM is used.
	switch hw.FC.CurrentMode {
	case FCModeNone:
		// Flow control (Rx & Tx) is completely disabled by a
		// software over-ride.
		adv &^= NWAY_AR_ASM_DIR | NWAY_AR_PAUSE
	case FCModeRxPause:
		// Rx Flow control is enabled, and Tx Flow control is
		// disabled, by a software over-ride.
		//
		// Since there really isn't a way to advertise that we are
		// capable of Rx Pause ONLY, we will advertise that we
		// support both symmetric and asymmetric Rx PAUSE.  Later
		// (in e1000_config_fc_after_link_up) we will disable the
		// hw's ability to send PAUSE frames.
		adv |= NWAY_AR_ASM_DIR | NWAY_AR_PAUSE
	case FCModeTxPause:
		// Tx Flow control is enabled, and Rx Flow control is
		// disabled, by a software over-ride.
		adv |= NWAY_AR_ASM_DIR
		adv &^= NWAY_AR_PAUSE
	case FCModeFull:
		// Flow control (both Rx and Tx) is enabled by a software
		// over-ride.
		adv |= NWAY_AR_ASM_DIR | NWAY_AR_PAUSE
	default:
		return errors.New("Flow control param set incorrectly")
	}

	err = phy.Op.WriteReg(PHY_AUTONEG_ADV, adv)
	if err != nil {
		return err
	}

	if phy.AutonegMask&ADVERTISE_1000_FULL != 0 {
		err := phy.Op.WriteReg(PHY_1000T_CTRL, ctrl)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetupFiberSerdesLink(hw *HW) error {
	ctrl := hw.RegRead(CTRL)

	// Take the link out of reset
	ctrl &^= CTRL_LRST

	hw.MAC.Op.ConfigCollisionDist()

	err := CommitFCSettings(hw)
	if err != nil {
		return err
	}

	// Since auto-negotiation is enabled, take the link out of reset (the
	// link will be in reset, because we previously reset the chip). This
	// will restart auto-negotiation.  If auto-negotiation is successful
	// then the link-up status bit will be set and the flow control enable
	// bits (RFCE and TFCE) will be set according to their negotiated value.

	hw.RegWrite(CTRL, ctrl)
	hw.RegWriteFlush()
	time.Sleep(1 * time.Millisecond)

	// For these adapters, the SW definable pin 1 is set when the optics
	// detect a signal.  If we have a signal, then poll for a "Link-Up"
	// indication.
	if hw.PHY.MediaType != MediaTypeInternalSerdes {
		return nil
	}
	ctrl = hw.RegRead(CTRL)
	if ctrl&CTRL_SWDPIN1 == 0 {
		return nil
	}
	return PollFiberSerdesLink(hw)
}

func PollFiberSerdesLink(hw *HW) error {
	mac := &hw.MAC

	// If we have a signal (the cable is plugged in, or assumed true for
	// serdes media) then poll for a "Link-Up" indication in the Device
	// Status Register.  Time-out if a link isn't seen in 500 milliseconds
	// seconds (Auto-negotiation should complete in less than 500
	// milliseconds even if the other end is doing it in SW).
	var i int
	for ; i < FIBER_LINK_UP_LIMIT; i++ {
		time.Sleep(10 * time.Millisecond)
		status := hw.RegRead(STATUS)
		if status&STATUS_LU != 0 {
			break
		}
	}
	if i == FIBER_LINK_UP_LIMIT {
		mac.AutonegFailed = true
		// AutoNeg failed to achieve a link, so we'll call
		// mac->check_for_link. This routine will force the
		// link up if we detect a signal. This will allow us to
		// communicate with non-autonegotiating link partners.
		err := mac.Op.CheckForLink()
		if err != nil {
			return err
		}
		mac.AutonegFailed = false
	} else {
		mac.AutonegFailed = false
	}

	return nil
}

func CheckDownshift(hw *HW) error {
	phy := &hw.PHY
	var offset uint32
	var mask uint16
	switch phy.PHYType {
	case PHYTypeI210, PHYTypeM88, PHYTypeGg82563, PHYTypeBm, PHYType82578:
		offset = M88E1000_PHY_SPEC_STATUS
		mask = M88E1000_PSSR_DOWNSHIFT
		break
	case PHYTypeIgp, PHYTypeIgp2, PHYTypeIgp3:
		offset = IGP01E1000_PHY_LINK_HEALTH
		mask = IGP01E1000_PLHR_SS_DOWNGRADE
		break
	default:
		// speed downshift not supported
		phy.SpeedDowngraded = false
		return nil
	}

	data, err := phy.Op.ReadReg(offset)
	if err != nil {
		return err
	}
	phy.SpeedDowngraded = data&mask != 0

	return nil
}

func SetMasterSlaveMode(hw *HW) error {
	phy := &hw.PHY

	// Resolve Master/Slave mode
	data, err := phy.Op.ReadReg(PHY_1000T_CTRL)
	if err != nil {
		return err
	}

	// load defaults for future use
	if data&CR_1000T_MS_ENABLE != 0 {
		if data&CR_1000T_MS_VALUE != 0 {
			phy.OrigMSType = MSTypeForceMaster
		} else {
			phy.OrigMSType = MSTypeForceSlave
		}
	} else {
		phy.OrigMSType = MSTypeAuto
	}

	switch phy.MSType {
	case MSTypeForceMaster:
		data |= CR_1000T_MS_ENABLE | CR_1000T_MS_VALUE
	case MSTypeForceSlave:
		data |= CR_1000T_MS_ENABLE
		data &^= CR_1000T_MS_VALUE
	case MSTypeAuto:
		data &^= CR_1000T_MS_ENABLE
	}

	return phy.Op.WriteReg(PHY_1000T_CTRL, data)
}

func CheckResetBlock(hw *HW) error {
	manc := hw.RegRead(MANC)
	if manc&MANC_BLK_PHY_RST_ON_IDE != 0 {
		return errors.New("block phy reset")
	}
	return nil
}

func AcquirePHYBase(hw *HW) error {
	var mask uint16
	switch hw.Bus.Func {
	case 1:
		mask = SWFW_PHY1_SM
	case 2:
		mask = SWFW_PHY2_SM
	case 3:
		mask = SWFW_PHY3_SM
	default:
		mask = SWFW_PHY0_SM
	}
	return hw.MAC.Op.AcquireSWFWSync(mask)
}

func ReleasePHYBase(hw *HW) {
	var mask uint16
	switch hw.Bus.Func {
	case 1:
		mask = SWFW_PHY1_SM
	case 2:
		mask = SWFW_PHY2_SM
	case 3:
		mask = SWFW_PHY3_SM
	default:
		mask = SWFW_PHY0_SM
	}
	hw.MAC.Op.ReleaseSWFWSync(mask)
}

func SetD3LpluState(hw *HW, active bool) error {
	phy := &hw.PHY

	data, err := phy.Op.ReadReg(IGP02E1000_PHY_POWER_MGMT)
	if err != nil {
		return err
	}

	if !active {
		data &^= IGP02E1000_PM_D3_LPLU
		err := phy.Op.WriteReg(IGP02E1000_PHY_POWER_MGMT, data)
		if err != nil {
			return err
		}
		// LPLU and SmartSpeed are mutually exclusive.  LPLU is used
		// during Dx states where the power conservation is most
		// important.  During driver activity we should enable
		// SmartSpeed, so performance is maintained.
		if phy.SmartSpeed == SmartSpeedOn {
			data, err := phy.Op.ReadReg(IGP01E1000_PHY_PORT_CONFIG)
			if err != nil {
				return err
			}
			data |= IGP01E1000_PSCFR_SMART_SPEED
			err = phy.Op.WriteReg(IGP01E1000_PHY_PORT_CONFIG, data)
			if err != nil {
				return err
			}
		} else if phy.SmartSpeed == SmartSpeedOff {
			data, err := phy.Op.ReadReg(IGP01E1000_PHY_PORT_CONFIG)
			if err != nil {
				return err
			}
			data &^= IGP01E1000_PSCFR_SMART_SPEED
			err = phy.Op.WriteReg(IGP01E1000_PHY_PORT_CONFIG, data)
			if err != nil {
				return err
			}
		}
	} else if phy.AutonegAdvertised == ALL_SPEED_DUPLEX ||
		phy.AutonegAdvertised == ALL_NOT_GIG ||
		phy.AutonegAdvertised == ALL_10_SPEED {
		data |= IGP02E1000_PM_D3_LPLU
		err := phy.Op.WriteReg(IGP02E1000_PHY_POWER_MGMT, data)
		if err != nil {
			return err
		}
		// When LPLU is enabled, we should disable SmartSpeed
		data, err = phy.Op.ReadReg(IGP01E1000_PHY_PORT_CONFIG)
		if err != nil {
			return err
		}
		data &^= IGP01E1000_PSCFR_SMART_SPEED
		return phy.Op.WriteReg(IGP01E1000_PHY_PORT_CONFIG, data)
	}
	return nil
}
