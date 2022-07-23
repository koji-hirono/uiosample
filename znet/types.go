package znet

type Uint16 [2]byte

func (u Uint16) Get() uint16 {
	x := uint16(u[0]) << 8
	x |= uint16(u[1])
	return x
}

func (u *Uint16) Set(x uint16) {
	u[0] = byte(x >> 8)
	u[1] = byte(x)
}

type Uint32 [4]byte

func (u Uint32) Get() uint32 {
	x := uint32(u[0]) << 24
	x |= uint32(u[1]) << 16
	x |= uint32(u[2]) << 8
	x |= uint32(u[3])
	return x
}

func (u *Uint32) Set(x uint32) {
	u[0] = byte(x >> 24)
	u[1] = byte(x >> 16)
	u[2] = byte(x >> 8)
	u[3] = byte(x)
}

type Uint64 [8]byte

func (u Uint64) Get() uint64 {
	x := uint64(u[0]) << 56
	x |= uint64(u[1]) << 48
	x |= uint64(u[2]) << 40
	x |= uint64(u[3]) << 32
	x |= uint64(u[4]) << 24
	x |= uint64(u[5]) << 16
	x |= uint64(u[6]) << 8
	x |= uint64(u[7])
	return x
}

func (u *Uint64) Set(x uint64) {
	u[0] = byte(x >> 56)
	u[1] = byte(x >> 48)
	u[2] = byte(x >> 40)
	u[3] = byte(x >> 32)
	u[4] = byte(x >> 24)
	u[5] = byte(x >> 16)
	u[6] = byte(x >> 8)
	u[7] = byte(x)
}
