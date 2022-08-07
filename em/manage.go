package em

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
