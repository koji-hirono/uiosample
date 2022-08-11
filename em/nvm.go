package em

import (
	"errors"
	"time"
)

type NVMType int

const (
	NVMTypeUnknown NVMType = iota
	NVMTypeNone
	NVMTypeEepromSpi
	NVMTypeEepromMicrowire
	NVMTypeFlashHw
	NVMTypeInvm
	NVMTypeFlashSw
)

type NVMOverride int

const (
	NVMOverrideNone NVMOverride = iota
	NVMOverrideSpiSmall
	NVMOverrideSpiLarge
	NVMOverrideMicrowireSmall
	NVMOverrideMicrowireLarge
)

type NVMInfo struct {
	Op       NVMOp
	Type     NVMType
	Override NVMOverride

	FlashBankSize uint32
	FlashBaseAddr uint32

	WordSize    uint16
	DelayUsec   time.Duration
	AddressBits uint16
	OpcodeBits  uint16
	PageSize    uint16
}

type NVMOp interface {
	InitParams() error
	Acquire() error
	Read(uint16, []uint16) error
	Release()
	Reload()
	Update() error
	ValidLEDDefault() (uint16, error)
	Validate() error
	Write(uint16, []uint16) error
}

/*
   s32  (*init_params)(struct e1000_hw *);
   s32  (*acquire)(struct e1000_hw *);
   s32  (*read)(struct e1000_hw *, u16, u16, u16 *);
   void (*release)(struct e1000_hw *);
   void (*reload)(struct e1000_hw *);
   s32  (*update)(struct e1000_hw *);
   s32  (*valid_led_default)(struct e1000_hw *, u16 *);
   s32  (*validate)(struct e1000_hw *);
   s32  (*write)(struct e1000_hw *, u16, u16, u16 *);
*/

// NVM # attempts to gain grant
const NVM_GRANT_ATTEMPTS = 1000

func AcquireNVM(hw *HW) error {
	eecd := hw.RegRead(EECD)
	hw.RegWrite(EECD, eecd|EECD_REQ)
	eecd = hw.RegRead(EECD)
	tmo := NVM_GRANT_ATTEMPTS
	for ; tmo > 0; tmo-- {
		if eecd&EECD_GNT != 0 {
			break
		}
		time.Sleep(5 * time.Microsecond)
		eecd = hw.RegRead(EECD)
	}
	if tmo == 0 {
		eecd &^= EECD_REQ
		hw.RegWrite(EECD, eecd)
		return errors.New("could not acquire NVM grant")
	}
	return nil
}

func ReleaseNVM(hw *HW) {
	StopNVM(hw)

	eecd := hw.RegRead(EECD)
	eecd &^= EECD_REQ
	hw.RegWrite(EECD, eecd)
}

func StopNVM(hw *HW) {
	eecd := hw.RegRead(EECD)
	if hw.NVM.Type == NVMTypeEepromSpi {
		eecd |= EECD_CS
		LowerEECClk(hw, eecd)
	} else {
		eecd &^= EECD_CS | EECD_DI
		hw.RegWrite(EECD, eecd)
		eecd = RaiseEECClk(hw, eecd)
		LowerEECClk(hw, eecd)
	}
}

func RaiseEECClk(hw *HW, val uint32) uint32 {
	val |= EECD_SK
	hw.RegWrite(EECD, val)
	hw.RegWriteFlush()
	time.Sleep(hw.NVM.DelayUsec * time.Microsecond)
	return val
}

func LowerEECClk(hw *HW, val uint32) uint32 {
	val &^= EECD_SK
	hw.RegWrite(EECD, val)
	hw.RegWriteFlush()
	time.Sleep(hw.NVM.DelayUsec * time.Microsecond)
	return val
}

func NVMReload(hw *HW) {
	time.Sleep(10 * time.Microsecond)
	x := hw.RegRead(CTRL_EXT)
	x |= CTRL_EXT_EE_RST
	hw.RegWrite(CTRL_EXT, x)
	hw.RegWriteFlush()
}

func UpdateNVMChecksum(hw *HW) error {
	nvm := &hw.NVM
	var checksum uint16
	var x [1]uint16
	for i := 0; i < NVM_CHECKSUM_REG; i++ {
		err := nvm.Op.Read(uint16(i), x[:])
		if err != nil {
			return err
		}
		checksum += x[0]
	}
	checksum = uint16(NVM_SUM) - checksum
	x[0] = checksum
	return nvm.Op.Write(NVM_CHECKSUM_REG, x[:])
}

func ValidateNVMChecksum(hw *HW) error {
	nvm := &hw.NVM
	var checksum uint16
	var x [1]uint16
	for i := 0; i < NVM_CHECKSUM_REG+1; i++ {
		err := nvm.Op.Read(uint16(i), x[:])
		if err != nil {
			return err
		}
		checksum += x[0]
	}
	if checksum != uint16(NVM_SUM) {
		return errors.New("NVM Checksum Invalid")
	}
	return nil
}

func StandbyNVM(hw *HW) {
	nvm := &hw.NVM
	eecd := hw.RegRead(EECD)
	switch nvm.Type {
	case NVMTypeEepromMicrowire:
		eecd &^= EECD_CS | EECD_SK
		hw.RegWrite(EECD, eecd)
		hw.RegWriteFlush()
		eecd = RaiseEECClk(hw, eecd)
		eecd |= EECD_CS
		hw.RegWrite(EECD, eecd)
		hw.RegWriteFlush()
		time.Sleep(nvm.DelayUsec * time.Microsecond)
		LowerEECClk(hw, eecd)
	case NVMTypeEepromSpi:
		eecd |= EECD_CS
		hw.RegWrite(EECD, eecd)
		hw.RegWriteFlush()
		time.Sleep(nvm.DelayUsec * time.Microsecond)
		eecd &^= EECD_CS
		hw.RegWrite(EECD, eecd)
		hw.RegWriteFlush()
		time.Sleep(nvm.DelayUsec * time.Microsecond)
	}
}

func ReadyNVMEeprom(hw *HW) error {
	nvm := &hw.NVM
	eecd := hw.RegRead(EECD)
	switch nvm.Type {
	case NVMTypeEepromMicrowire:
		// Clear SK and DI
		eecd &^= EECD_DI | EECD_SK
		hw.RegWrite(EECD, eecd)
		// Set CS
		eecd |= EECD_CS
		hw.RegWrite(EECD, eecd)
	case NVMTypeEepromSpi:
		// Clear SK and CS
		eecd &^= EECD_CS | EECD_SK
		hw.RegWrite(EECD, eecd)
		hw.RegWriteFlush()
		time.Sleep(1 * time.Microsecond)
		tmo := NVM_MAX_RETRY_SPI
		for tmo > 0 {
			shiftOutEECbits(hw, NVM_RDSR_OPCODE_SPI, nvm.OpcodeBits)
			x := shiftInEECbits(hw, 8)
			if x&NVM_STATUS_RDY_SPI == 0 {
				break
			}
			time.Sleep(5 * time.Microsecond)
			StandbyNVM(hw)
			tmo--
		}
		if tmo == 0 {
			return errors.New("SPI NVM Status error")
		}
	}
	return nil
}

func shiftOutEECbits(hw *HW, data uint16, count uint16) {
	nvm := &hw.NVM
	eecd := hw.RegRead(EECD)
	mask := uint16(0x01) << (count - 1)
	switch nvm.Type {
	case NVMTypeEepromMicrowire:
		eecd &^= EECD_DO
	case NVMTypeEepromSpi:
		eecd |= EECD_DO
	}

	for mask != 0 {
		eecd &^= EECD_DI

		if data&mask != 0 {
			eecd |= EECD_DI
		}

		hw.RegWrite(EECD, eecd)
		hw.RegWriteFlush()

		time.Sleep(nvm.DelayUsec * time.Microsecond)

		eecd = RaiseEECClk(hw, eecd)
		eecd = LowerEECClk(hw, eecd)

		mask >>= 1
	}

	eecd &^= EECD_DI
	hw.RegWrite(EECD, eecd)
}

func shiftInEECbits(hw *HW, count uint16) uint16 {
	eecd := hw.RegRead(EECD)
	eecd &^= EECD_DO | EECD_DI
	var data uint16
	for i := uint16(0); i < count; i++ {
		data <<= 1
		RaiseEECClk(hw, eecd)
		eecd = hw.RegRead(EECD)
		eecd &^= EECD_DI
		if eecd&EECD_DO != 0 {
			data |= 1
		}
		eecd = LowerEECClk(hw, eecd)
	}
	return data
}

func ReadMACAddr(hw *HW) error {
	high := hw.RegRead(RAH(0))
	low := hw.RegRead(RAL(0))

	hw.MAC.PermAddr[0] = byte(low)
	hw.MAC.PermAddr[1] = byte(low >> 8)
	hw.MAC.PermAddr[2] = byte(low >> 16)
	hw.MAC.PermAddr[3] = byte(low >> 24)
	hw.MAC.PermAddr[4] = byte(high)
	hw.MAC.PermAddr[5] = byte(high >> 8)

	copy(hw.MAC.Addr[:], hw.MAC.PermAddr[:])
	return nil
}
