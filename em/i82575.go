package em

import (
	"errors"
	"time"
)

// SFP modules ID memory locations
const (
	SFF_IDENTIFIER_OFFSET = 0x00
	SFF_IDENTIFIER_SFF    = 0x02
	SFF_IDENTIFIER_SFP    = 0x03

	SFF_ETH_FLAGS_OFFSET = 0x06
)

// Flags for SFP modules compatible with ETH up to 1Gb
type SFPFlags uint8

const (
	SFPFlags1000BaseSX SFPFlags = 1 << iota
	SFPFlags1000BaseLX
	SFPFlags1000BaseCX
	SFPFlags1000BaseT
	SFPFlags100BaseLX
	SFPFlags100BaseFX
	SFPFlags10BaseBX10
	SFPFlags10BasePX
)

type I82575DeviceSpec struct {
	SGMIIActive        bool
	GlobalDeviceReset  bool
	EEEDisable         bool
	ModulePlugged      bool
	ClearSemaphoreOnce bool
	MTU                uint32
	Ethflags           SFPFlags
	MediaPort          uint8
	MediaChanged       bool
}

func I82575Init(hw *HW) {
	hw.Spec = &I82575DeviceSpec{}
	nvm := NewI82575NVM(hw)
	phy := NewI82575PHY(hw)
	mac := NewI82575MAC(hw, nvm, phy)
	hw.MAC.Op = mac
	hw.PHY.Op = phy
	hw.NVM.Op = nvm
}

func AcquireSWFWSync82575(hw *HW, mask uint16) error {
	swmask := uint32(mask)
	fwmask := uint32(mask) << 16
	timeout := 200
	for i := 0; i < timeout; i++ {
		err := GetHWSemaphore(hw)
		if err != nil {
			return err
		}
		swfw_sync := hw.RegRead(SW_FW_SYNC)
		if swfw_sync&(fwmask|swmask) == 0 {
			swfw_sync |= swmask
			hw.RegWrite(SW_FW_SYNC, swfw_sync)
			PutHWSemaphore(hw)
			return nil
		}
		// Firmware currently using resource (fwmask)
		// or other software thread using resource (swmask)
		PutHWSemaphore(hw)
		time.Sleep(5 * time.Millisecond)
	}
	return errors.New("SW_FW_SYNC timeout")
}

func ReleaseSWFWSync82575(hw *HW, mask uint16) {
	for {
		err := GetHWSemaphore(hw)
		if err == nil {
			break
		}
	}

	swfw_sync := hw.RegRead(SW_FW_SYNC)
	swfw_sync &^= uint32(mask)
	hw.RegWrite(SW_FW_SYNC, swfw_sync)

	PutHWSemaphore(hw)
}

func SGMIIActive82575(hw *HW) bool {
	spec := hw.Spec.(*I82575DeviceSpec)
	return spec.SGMIIActive
}

func SGMIIUsesMDIO82575(hw *HW) bool {
	switch hw.MAC.Type {
	case MACType82575, MACType82576:
		x := hw.RegRead(MDIC)
		return x&MDIC_DEST != 0
	case MACType82580, MACTypeI350, MACTypeI354, MACTypeI210, MACTypeI211:
		x := hw.RegRead(MDICNFG)
		return x&MDICNFG_EXT_MDIO != 0
	default:
		return false
	}
}

func PHYHWResetSGMII82575(hw *HW) error {
	phy := &hw.PHY
	// SFP documentation requires the following to configure the SPF module
	// to work on SGMII.  No further documentation is given.
	err := phy.Op.WriteReg(0x1B, 0x8084)
	if err != nil {
		return err
	}
	err = phy.Op.Commit()
	if err != nil {
		return err
	}
	if phy.ID == M88E1512_E_PHY_ID {
		return InitM88E1512PHY(hw)
	}
	return nil
}

func SetD0LpluState82575(hw *HW, active bool) error {
	phy := &hw.PHY
	data, err := phy.Op.ReadReg(IGP02E1000_PHY_POWER_MGMT)
	if err != nil {
		return err
	}

	if active {
		data |= IGP02E1000_PM_D0_LPLU
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
		err = phy.Op.WriteReg(IGP01E1000_PHY_PORT_CONFIG, data)
		if err != nil {
			return err
		}
	} else {
		data &^= IGP02E1000_PM_D0_LPLU
		phy.Op.WriteReg(IGP02E1000_PHY_POWER_MGMT, data)
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
	}
	return nil
}
