package em

import (
	"errors"
	"time"
)

func InitHWI210(hw *HW) error {
	if hw.MAC.Type >= MACTypeI210 && !GetFlashPresenceI210(hw) {
		err := PllWorkaroundI210(hw)
		if err != nil {
			return err
		}
	}
	// TODO
	// hw.PHY.Op.GetCfgDone = e1000_get_cfg_done_i210

	// Initialize identification LED
	hw.MAC.Op.IDLEDInit()

	return InitHWBase(hw)
}

func GetFlashPresenceI210(hw *HW) bool {
	eecd := hw.RegRead(EECD)
	return eecd&EECD_FLASH_DETECTED_I210 != 0
}

func PllWorkaroundI210(hw *HW) error {
	// TODO:
	return nil
}

func AcquireSWFWSyncI210(hw *HW, mask uint16) error {
	swmask := uint32(mask)
	fwmask := uint32(mask) << 16
	timeout := 200
	for i := 0; i < timeout; i++ {
		err := GetHWSemaphoreI210(hw)
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

func ReleaseSWFWSyncI210(hw *HW, mask uint16) {
	for {
		err := GetHWSemaphoreI210(hw)
		if err == nil {
			break
		}
	}

	swfw_sync := hw.RegRead(SW_FW_SYNC)
	swfw_sync &^= uint32(mask)
	hw.RegWrite(SW_FW_SYNC, swfw_sync)

	PutHWSemaphore(hw)
}

func GetHWSemaphoreI210(hw *HW) error {
	// Get the SW semaphore
	timeout := int(hw.NVM.WordSize) + 1
	var i int
	for i < timeout {
		swsm := hw.RegRead(SWSM)
		if swsm&SWSM_SMBI == 0 {
			break
		}
		time.Sleep(50 * time.Microsecond)
		i++
	}
	if i == timeout {
		// In rare circumstances, the SW semaphore may already be held
		// unintentionally. Clear the semaphore once before giving up.
		spec := hw.Spec.(*I82575DeviceSpec)
		if spec.ClearSemaphoreOnce {
			spec.ClearSemaphoreOnce = false
			PutHWSemaphore(hw)
			for i = 0; i < timeout; i++ {
				swsm := hw.RegRead(SWSM)
				if swsm&SWSM_SMBI == 0 {
					break
				}
				time.Sleep(50 * time.Microsecond)
			}
		}
		// If we do not have the semaphore here, we have to give up.
		if i == timeout {
			return errors.New("SMBI bit is set")
		}
	}

	// Get the FW semaphore.
	for i = 0; i < timeout; i++ {
		swsm := hw.RegRead(SWSM)
		hw.RegWrite(SWSM, swsm|SWSM_SWESMBI)
		// Semaphore acquired if bit latched
		if hw.RegRead(SWSM)&SWSM_SWESMBI != 0 {
			break
		}
		time.Sleep(50 * time.Microsecond)
	}
	if i == timeout {
		// Release semaphores
		PutHWSemaphore(hw)
		return errors.New("timeout")
	}
	return nil
}
