package em

import (
	"time"
)

func ResetHW82580(hw *HW) error {
	mac := &hw.MAC
	spec := hw.Spec.(*I82575DeviceSpec)

	global_device_reset := spec.GlobalDeviceReset
	spec.GlobalDeviceReset = false

	// 82580 does not reliably do global_device_reset due to hw errata
	if hw.MAC.Type == MACType82580 {
		global_device_reset = false
	}

	// Get current control state.
	ctrl := hw.RegRead(CTRL)

	// Prevent the PCI-E bus from sticking if there is no TLP connection
	// on the last TLP read/write transaction when MAC is reset.
	DisablePCIEMaster(hw)

	hw.RegWrite(IMC, ^uint32(0))
	hw.RegWrite(RCTL, 0)
	hw.RegWrite(TCTL, TCTL_PSP)
	hw.RegWriteFlush()

	time.Sleep(10 * time.Millisecond)

	// BH SW mailbox bit in SW_FW_SYNC
	swmbsw_mask := SW_SYNCH_MB
	// Determine whether or not a global dev reset is requested
	if global_device_reset && mac.Op.AcquireSWFWSync(swmbsw_mask) != nil {
		global_device_reset = false
	}

	if global_device_reset && hw.RegRead(STATUS)&STAT_DEV_RST_SET == 0 {
		ctrl |= CTRL_DEV_RST
	} else {
		ctrl |= CTRL_RST
	}
	hw.RegWrite(CTRL, ctrl)

	switch hw.DeviceID {
	case DEV_ID_DH89XXCC_SGMII:
	default:
		hw.RegWriteFlush()
	}

	// Add delay to insure DEV_RST or RST has time to complete
	time.Sleep(5 * time.Millisecond)

	// When auto config read does not complete, do not
	// return with an error. This can happen in situations
	// where there is no eeprom and prevents getting link.
	GetAutoRDDone(hw)

	// clear global device reset status bit
	hw.RegWrite(STATUS, STAT_DEV_RST_SET)

	// Clear any pending interrupt events.
	hw.RegWrite(IMC, ^uint32(0))
	hw.RegRead(ICR)

	ResetMDIConfig82580(hw)

	// Install any alternate MAC address into RAR0
	err := CheckAltMACAddr(hw)

	// Release semaphore
	if global_device_reset {
		mac.Op.ReleaseSWFWSync(swmbsw_mask)
	}

	return err
}

func ResetMDIConfig82580(hw *HW) error {
	if hw.MAC.Type != MACType82580 {
		return nil
	}
	if !SGMIIActive82575(hw) {
		return nil
	}

	var data [1]uint16
	err := hw.NVM.Op.Read(NVM_INIT_CONTROL3_PORT_A+NVM_82580_LAN_FUNC_OFFSET(hw.Bus.Func), data[:])
	if err != nil {
		return err
	}

	mdicnfg := hw.RegRead(MDICNFG)
	if data[0]&NVM_WORD24_EXT_MDIO != 0 {
		mdicnfg |= MDICNFG_EXT_MDIO
	}
	if data[0]&NVM_WORD24_COM_MDIO != 0 {
		mdicnfg |= MDICNFG_COM_MDIO
	}
	hw.RegWrite(MDICNFG, mdicnfg)
	return nil
}

func ReadPHYReg82580(hw *HW, offset uint32) (uint16, error) {
	phy := &hw.PHY
	err := phy.Op.Acquire()
	if err != nil {
		return 0, err
	}
	defer phy.Op.Release()
	return ReadPHYRegMDIC(hw, offset)
}

func WritePHYReg82580(hw *HW, offset uint32, data uint16) error {
	phy := &hw.PHY
	err := phy.Op.Acquire()
	if err != nil {
		return err
	}
	defer phy.Op.Release()
	return WritePHYRegMDIC(hw, offset, data)
}

func SetD0LpluState82580(hw *HW, active bool) error {
	phy := &hw.PHY
	data := hw.RegRead(I82580_PHY_POWER_MGMT)

	if active {
		data |= I82580_PM_D0_LPLU

		// When LPLU is enabled, we should disable SmartSpeed
		data &^= I82580_PM_SPD
	} else {
		data &^= I82580_PM_D0_LPLU

		// LPLU and SmartSpeed are mutually exclusive.  LPLU is used
		// during Dx states where the power conservation is most
		// important.  During driver activity we should enable
		// SmartSpeed, so performance is maintained.
		if phy.SmartSpeed == SmartSpeedOn {
			data |= I82580_PM_SPD
		} else if phy.SmartSpeed == SmartSpeedOff {
			data &^= I82580_PM_SPD
		}
	}

	hw.RegWrite(I82580_PHY_POWER_MGMT, data)
	return nil
}

func SetD3LpluState82580(hw *HW, active bool) error {
	phy := &hw.PHY

	data := hw.RegRead(I82580_PHY_POWER_MGMT)

	if !active {
		data &^= I82580_PM_D3_LPLU
		// LPLU and SmartSpeed are mutually exclusive.  LPLU is used
		// during Dx states where the power conservation is most
		// important.  During driver activity we should enable
		// SmartSpeed, so performance is maintained.
		if phy.SmartSpeed == SmartSpeedOn {
			data |= I82580_PM_SPD
		} else if phy.SmartSpeed == SmartSpeedOff {
			data &^= I82580_PM_SPD
		}
	} else if phy.AutonegAdvertised == ALL_SPEED_DUPLEX ||
		phy.AutonegAdvertised == ALL_NOT_GIG ||
		phy.AutonegAdvertised == ALL_10_SPEED {
		data |= I82580_PM_D3_LPLU
		// When LPLU is enabled, we should disable SmartSpeed
		data &^= I82580_PM_SPD
	}

	hw.RegWrite(I82580_PHY_POWER_MGMT, data)
	return nil
}
