package em

import (
	"uiosample/em/em82571"
)

type ManageMode uint32

const (
	ManageModeNone ManageMode = iota
	ManageModeASF
	ManageModePT
	ManageModeIPMI
	ManageModeHostIfOnly
)

type host_mng_dhcp_cookie struct {
	signature uint32
	status    uint8
	reserved0 uint8
	vlan_id   uint16
	reserved1 uint32
	reserved2 uint16
	reserved3 uint8
	checksum  uint8
}

// Host Interface "Rev 1"
type host_command_header struct {
	command_id      uint8
	command_length  uint8
	command_options uint8
	checksum        uint8
}

const HI_MAX_DATA_LENGTH = 252

type host_command_info struct {
	command_header host_command_header
	command_data   [HI_MAX_DATA_LENGTH]byte
}

type host_mng_command_header struct {
	command_id     uint8
	checksum       uint8
	reserved1      uint16
	reserved2      uint16
	command_length uint16
}

const HI_MAX_MNG_DATA_LENGTH = 0x6F8

type host_mng_command_info struct {
	command_header host_mng_command_header
	command_data   [HI_MAX_MNG_DATA_LENGTH]byte
}

// bool e1000_enable_mng_pass_thru(struct e1000_hw *hw)
func EnableManagePT(hw *HW) bool {
	mac := &hw.MAC

	if !mac.ASFFirmwarePresent {
		return false
	}

	manc := hw.RegRead(MANC)

	if manc&MANC_RCV_TCO_EN == 0 {
		return false
	}

	if mac.HasFWSM {
		fwsm := hw.RegRead(FWSM)
		factps := hw.RegRead(FACTPS)

		if factps&FACTPS_MNGCG == 0 &&
			fwsm&FWSM_MODE_MASK == uint32(ManageModePT)<<FWSM_MODE_SHIFT {
			return true
		}
	} else if mac.Type == MACType82574 || mac.Type == MACType82583 {
		nvm := &hw.NVM
		var data []uint16

		factps := hw.RegRead(FACTPS)
		err := nvm.Op.Read(NVM_INIT_CONTROL2_REG, data[:])
		if err != nil {
			return false
		}

		if factps&FACTPS_MNGCG == 0 &&
			data[0]&em82571.NVM_INIT_CTRL2_MNGM == uint16(ManageModePT)<<13 {
			return true
		}
	} else if manc&MANC_SMBUS_EN != 0 && manc&MANC_ASF_EN == 0 {
		return true
	}

	return false
}
