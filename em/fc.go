package em

import (
	"errors"
)

type FCMode int

const (
	FCModeNone FCMode = iota
	FCModeRxPause
	FCModeTxPause
	FCModeFull
	FCModeDefault = 0xff
)

type FCInfo struct {
	HighWater     uint32
	LowWater      uint32
	PauseTime     uint16
	RefreshTime   uint16
	SendXON       bool
	StrictIEEE    bool
	CurrentMode   FCMode
	RequestedMode FCMode
}

// Flow Control Constants
const (
	FLOW_CONTROL_ADDRESS_LOW  = 0x00C28001
	FLOW_CONTROL_ADDRESS_HIGH = 0x00000100
	FLOW_CONTROL_TYPE         = 0x8808
)

func SetDefaultFC(hw *HW) error {
	// Read and store word 0x0F of the EEPROM. This word contains bits
	// that determine the hardware's default PAUSE (flow control) mode,
	// a bit that determines whether the HW defaults to enabling or
	// disabling auto-negotiation, and the direction of the
	// SW defined pins. If there is no SW over-ride of the flow
	// control setting, then the variable hw->fc will
	// be initialized based on a value in the EEPROM.
	var data [1]uint16
	if hw.MAC.Type == MACTypeI350 {
		offset := NVM_82580_LAN_FUNC_OFFSET(hw.Bus.Func)
		err := hw.NVM.Op.Read(NVM_INIT_CONTROL2_REG+offset, data[:])
		if err != nil {
			return err
		}
	} else {
		err := hw.NVM.Op.Read(NVM_INIT_CONTROL2_REG, data[:])
		if err != nil {
			return err
		}
	}

	if data[0]&NVM_WORD0F_PAUSE_MASK == 0 {
		hw.FC.RequestedMode = FCModeNone
	} else if data[0]&NVM_WORD0F_PAUSE_MASK == NVM_WORD0F_ASM_DIR {
		hw.FC.RequestedMode = FCModeTxPause
	} else {
		hw.FC.RequestedMode = FCModeFull
	}

	return nil
}

func SetFCWatermarks(hw *HW) error {
	var fcrtl uint32
	var fcrth uint32
	// Set the flow control receive threshold registers.  Normally,
	// these registers will be set to a default threshold that may be
	// adjusted later by the driver's runtime code.  However, if the
	// ability to transmit pause frames is not enabled, then these
	// registers will be set to 0.
	if hw.FC.CurrentMode&FCModeTxPause != 0 {
		// We need to set up the Receive Threshold high and low water
		// marks as well as (optionally) enabling the transmission of
		// XON frames.
		fcrtl = hw.FC.LowWater
		if hw.FC.SendXON {
			fcrtl |= FCRTL_XONE
		}
		fcrth = hw.FC.HighWater
	}
	hw.RegWrite(FCRTL, fcrtl)
	hw.RegWrite(FCRTH, fcrth)
	return nil
}

func ConfigFCAfterLinkUp(hw *HW) error {
	mac := &hw.MAC
	phy := &hw.PHY
	fc := &hw.FC

	// Check for the case where we have fiber media and auto-neg failed
	// so we had to force link.  In this case, we need to force the
	// configuration of the MAC to match the "fc" parameter.
	if mac.AutonegFailed {
		if phy.MediaType == MediaTypeFiber || phy.MediaType == MediaTypeInternalSerdes {
			err := ForceMACFC(hw)
			if err != nil {
				return err
			}
		}
	} else {
		if phy.MediaType == MediaTypeCopper {
			err := ForceMACFC(hw)
			if err != nil {
				return err
			}
		}
	}

	// Check for the case where we have copper media and auto-neg is
	// enabled.  In this case, we need to check and see if Auto-Neg
	// has completed, and if so, how the PHY and link partner has
	// flow control configured.
	if phy.MediaType == MediaTypeCopper && mac.Autoneg {
		// Read the MII Status Register and check to see if AutoNeg
		// has completed.  We read this twice because this reg has
		// some "sticky" (latched) bits.
		_, err := phy.Op.ReadReg(PHY_STATUS)
		if err != nil {
			return err
		}
		status, err := phy.Op.ReadReg(PHY_STATUS)
		if err != nil {
			return err
		}

		if status&MII_SR_AUTONEG_COMPLETE == 0 {
			return errors.New("not completed")
		}

		// The AutoNeg process has completed, so we now need to
		// read both the Auto Negotiation Advertisement
		// Register (Address 4) and the Auto_Negotiation Base
		// Page Ability Register (Address 5) to determine how
		// flow control was negotiated.
		adv, err := phy.Op.ReadReg(PHY_AUTONEG_ADV)
		if err != nil {
			return err
		}
		ability, err := phy.Op.ReadReg(PHY_LP_ABILITY)
		if err != nil {
			return err
		}
		// Two bits in the Auto Negotiation Advertisement Register
		// (Address 4) and two bits in the Auto Negotiation Base
		// Page Ability Register (Address 5) determine flow control
		// for both the PHY and the link partner.  The following
		// table, taken out of the IEEE 802.3ab/D6.0 dated March 25,
		// 1999, describes these PAUSE resolution bits and how flow
		// control is determined based upon these settings.
		// NOTE:  DC = Don't Care
		//
		//   LOCAL DEVICE  |   LINK PARTNER
		// PAUSE | ASM_DIR | PAUSE | ASM_DIR | NIC Resolution
		//-------|---------|-------|---------|--------------------
		//   0   |    0    |  DC   |   DC    | e1000_fc_none
		//   0   |    1    |   0   |   DC    | e1000_fc_none
		//   0   |    1    |   1   |    0    | e1000_fc_none
		//   0   |    1    |   1   |    1    | e1000_fc_tx_pause
		//   1   |    0    |   0   |   DC    | e1000_fc_none
		//   1   |   DC    |   1   |   DC    | e1000_fc_full
		//   1   |    1    |   0   |    0    | e1000_fc_none
		//   1   |    1    |   0   |    1    | e1000_fc_rx_pause
		//
		// Are both PAUSE bits set to 1?  If so, this implies
		// Symmetric Flow Control is enabled at both ends.  The
		// ASM_DIR bits are irrelevant per the spec.
		//
		// For Symmetric Flow Control:
		//
		//   LOCAL DEVICE  |   LINK PARTNER
		// PAUSE | ASM_DIR | PAUSE | ASM_DIR | Result
		//-------|---------|-------|---------|--------------------
		//   1   |   DC    |   1   |   DC    | E1000_fc_full
		//
		// For receiving PAUSE frames ONLY.
		//
		//   LOCAL DEVICE  |   LINK PARTNER
		// PAUSE | ASM_DIR | PAUSE | ASM_DIR | Result
		//-------|---------|-------|---------|--------------------
		//   0   |    1    |   1   |    1    | e1000_fc_tx_pause
		//
		// For transmitting PAUSE frames ONLY.
		//
		//   LOCAL DEVICE  |   LINK PARTNER
		// PAUSE | ASM_DIR | PAUSE | ASM_DIR | Result
		//-------|---------|-------|---------|--------------------
		//   1   |    1    |   0   |    1    | e1000_fc_rx_pause
		//
		if adv&NWAY_AR_PAUSE != 0 && ability&NWAY_LPAR_PAUSE != 0 {
			// Now we need to check if the user selected Rx ONLY
			// of pause frames.  In this case, we had to advertise
			// FULL flow control because we could not advertise Rx
			// ONLY. Hence, we must now check to see if we need to
			// turn OFF the TRANSMISSION of PAUSE frames.
			if fc.RequestedMode == FCModeFull {
				fc.CurrentMode = FCModeFull
			} else {
				fc.CurrentMode = FCModeRxPause
			}
		} else if adv&NWAY_AR_PAUSE == 0 && adv&NWAY_AR_ASM_DIR != 0 && ability&NWAY_LPAR_PAUSE != 0 && ability&NWAY_LPAR_ASM_DIR != 0 {
			fc.CurrentMode = FCModeTxPause
		} else if adv&NWAY_AR_PAUSE != 0 && adv&NWAY_AR_ASM_DIR != 0 && ability&NWAY_LPAR_PAUSE == 0 && ability&NWAY_LPAR_ASM_DIR != 0 {
			fc.CurrentMode = FCModeRxPause
		} else {
			// Per the IEEE spec, at this point flow control
			// should be disabled.
			fc.CurrentMode = FCModeNone
		}
		// Now we need to do one last check...  If we auto-
		// negotiated to HALF DUPLEX, flow control should not be
		// enabled per IEEE 802.3 spec.
		_, duplex, err := mac.Op.GetLinkUpInfo()
		if err != nil {
			return err
		}

		if duplex == HALF_DUPLEX {
			fc.CurrentMode = FCModeNone
		}

		// Now we call a subroutine to actually force the MAC
		// controller to use the correct flow control settings.
		err = ForceMACFC(hw)
		if err != nil {
			return err
		}
	}
	// Check for the case where we have SerDes media and auto-neg is
	// enabled.  In this case, we need to check and see if Auto-Neg
	// has completed, and if so, how the PHY and link partner has
	// flow control configured.
	if phy.MediaType == MediaTypeInternalSerdes && mac.Autoneg {
		// Read the PCS_LSTS and check to see if AutoNeg
		// has completed.
		status := hw.RegRead(PCS_LSTAT)
		if status&PCS_LSTS_AN_COMPLETE == 0 {
			return errors.New("not completed")
		}

		// The AutoNeg process has completed, so we now need to
		// read both the Auto Negotiation Advertisement
		// Register (PCS_ANADV) and the Auto_Negotiation Base
		// Page Ability Register (PCS_LPAB) to determine how
		// flow control was negotiated.
		adv := hw.RegRead(PCS_ANADV)
		ability := hw.RegRead(PCS_LPAB)

		// Two bits in the Auto Negotiation Advertisement Register
		// (PCS_ANADV) and two bits in the Auto Negotiation Base
		// Page Ability Register (PCS_LPAB) determine flow control
		// for both the PHY and the link partner.  The following
		// table, taken out of the IEEE 802.3ab/D6.0 dated March 25,
		// 1999, describes these PAUSE resolution bits and how flow
		// control is determined based upon these settings.
		// NOTE:  DC = Don't Care
		//
		//   LOCAL DEVICE  |   LINK PARTNER
		// PAUSE | ASM_DIR | PAUSE | ASM_DIR | NIC Resolution
		//-------|---------|-------|---------|--------------------
		//   0   |    0    |  DC   |   DC    | e1000_fc_none
		//   0   |    1    |   0   |   DC    | e1000_fc_none
		//   0   |    1    |   1   |    0    | e1000_fc_none
		//   0   |    1    |   1   |    1    | e1000_fc_tx_pause
		//   1   |    0    |   0   |   DC    | e1000_fc_none
		//   1   |   DC    |   1   |   DC    | e1000_fc_full
		//   1   |    1    |   0   |    0    | e1000_fc_none
		//   1   |    1    |   0   |    1    | e1000_fc_rx_pause
		//
		// Are both PAUSE bits set to 1?  If so, this implies
		// Symmetric Flow Control is enabled at both ends.  The
		// ASM_DIR bits are irrelevant per the spec.
		//
		// For Symmetric Flow Control:
		//
		//   LOCAL DEVICE  |   LINK PARTNER
		// PAUSE | ASM_DIR | PAUSE | ASM_DIR | Result
		//-------|---------|-------|---------|--------------------
		//   1   |   DC    |   1   |   DC    | e1000_fc_full
		//
		// For receiving PAUSE frames ONLY.
		//
		//   LOCAL DEVICE  |   LINK PARTNER
		// PAUSE | ASM_DIR | PAUSE | ASM_DIR | Result
		//-------|---------|-------|---------|--------------------
		//   0   |    1    |   1   |    1    | e1000_fc_tx_pause
		//
		// For transmitting PAUSE frames ONLY.
		//
		//   LOCAL DEVICE  |   LINK PARTNER
		// PAUSE | ASM_DIR | PAUSE | ASM_DIR | Result
		//-------|---------|-------|---------|--------------------
		//   1   |    1    |   0   |    1    | e1000_fc_rx_pause
		if adv&TXCW_PAUSE != 0 && ability&TXCW_PAUSE != 0 {
			// Now we need to check if the user selected Rx ONLY
			// of pause frames.  In this case, we had to advertise
			// FULL flow control because we could not advertise Rx
			// ONLY. Hence, we must now check to see if we need to
			// turn OFF the TRANSMISSION of PAUSE frames.
			if fc.RequestedMode == FCModeFull {
				fc.CurrentMode = FCModeFull
			} else {
				fc.CurrentMode = FCModeRxPause
			}
		} else if adv&TXCW_PAUSE == 0 && adv&TXCW_ASM_DIR != 0 && ability&TXCW_PAUSE != 0 && ability&TXCW_ASM_DIR != 0 {
			fc.CurrentMode = FCModeTxPause
		} else if adv&TXCW_PAUSE != 0 && adv&TXCW_ASM_DIR != 0 && ability&TXCW_PAUSE == 0 && ability&TXCW_ASM_DIR != 0 {
			fc.CurrentMode = FCModeRxPause
		} else {
			// Per the IEEE spec, at this point flow control
			fc.CurrentMode = FCModeNone
		}
		// Now we call a subroutine to actually force the MAC
		// controller to use the correct flow control settings.
		ctrl := hw.RegRead(PCS_LCTL)
		ctrl |= PCS_LCTL_FORCE_FCTRL
		hw.RegWrite(PCS_LCTL, ctrl)
		err := ForceMACFC(hw)
		if err != nil {
			return err
		}
	}
	return nil
}

func ForceMACFC(hw *HW) error {
	ctrl := hw.RegRead(CTRL)

	// Because we didn't get link via the internal auto-negotiation
	// mechanism (we either forced link or we got link via PHY
	// auto-neg), we have to manually enable/disable transmit an
	// receive flow control.
	//
	// The "Case" statement below enables/disable flow control
	// according to the "hw->fc.current_mode" parameter.
	//
	// The possible values of the "fc" parameter are:
	//      0:  Flow control is completely disabled
	//      1:  Rx flow control is enabled (we can receive pause
	//          frames but not send pause frames).
	//      2:  Tx flow control is enabled (we can send pause frames
	//          frames but we do not receive pause frames).
	//      3:  Both Rx and Tx flow control (symmetric) is enabled.
	//  other:  No other values should be possible at this point.

	switch hw.FC.CurrentMode {
	case FCModeNone:
		ctrl &^= CTRL_TFCE | CTRL_RFCE
	case FCModeRxPause:
		ctrl &^= CTRL_TFCE
		ctrl |= CTRL_RFCE
	case FCModeTxPause:
		ctrl &^= CTRL_RFCE
		ctrl |= CTRL_TFCE
	case FCModeFull:
		ctrl |= CTRL_TFCE | CTRL_RFCE
	default:
		return errors.New("Flow control param set incorrectly")
	}

	hw.RegWrite(CTRL, ctrl)
	return nil
}

func CommitFCSettings(hw *HW) error {
	mac := &hw.MAC
	fc := &hw.FC

	// Check for a software override of the flow control settings, and
	// setup the device accordingly.  If auto-negotiation is enabled, then
	// software will have to set the "PAUSE" bits to the correct value in
	// the Transmit Config Word Register (TXCW) and re-start auto-
	// negotiation.  However, if auto-negotiation is disabled, then
	// software will have to manually configure the two flow control enable
	// bits in the CTRL register.
	//
	// The possible values of the "fc" parameter are:
	//      0:  Flow control is completely disabled
	//      1:  Rx flow control is enabled (we can receive pause frames,
	//          but not send pause frames).
	//      2:  Tx flow control is enabled (we can send pause frames but we
	//          do not support receiving pause frames).
	//      3:  Both Rx and Tx flow control (symmetric) are enabled.
	var txcw uint32
	switch fc.CurrentMode {
	case FCModeNone:
		// Flow control completely disabled by a software over-ride.
		txcw = TXCW_ANE | TXCW_FD
	case FCModeRxPause:
		// Rx Flow control is enabled and Tx Flow control is disabled
		// by a software over-ride. Since there really isn't a way to
		// advertise that we are capable of Rx Pause ONLY, we will
		// advertise that we support both symmetric and asymmetric Rx
		// PAUSE.  Later, we will disable the adapter's ability to send
		// PAUSE frames.
		txcw = TXCW_ANE | TXCW_FD | TXCW_PAUSE_MASK
	case FCModeTxPause:
		// Tx Flow control is enabled, and Rx Flow control is disabled,
		// by a software over-ride.
		txcw = TXCW_ANE | TXCW_FD | TXCW_ASM_DIR
	case FCModeFull:
		// Flow control (both Rx and Tx) is enabled by a software
		// over-ride.
		txcw = TXCW_ANE | TXCW_FD | TXCW_PAUSE_MASK
	default:
		return errors.New("Flow control param set incorrectly")
	}

	hw.RegWrite(TXCW, txcw)
	mac.TxCW = txcw

	return nil
}
