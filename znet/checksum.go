package znet

func CalcChecksum(d []byte) uint16 {
	var sum uint32
	n := len(d)
	for i := 0; i < n; i += 2 {
		x := uint32(d[i]) << 8
		if i+1 < n {
			x |= uint32(d[i+1])
		}
		sum += x
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return uint16(^sum)
}

type Calc struct {
	sum uint32
}

func NewCalc() *Calc {
	return &Calc{}
}

func (c *Calc) Append(p []byte) {
	n := len(p)
	for i := 0; i < n; i += 2 {
		x := uint32(p[i]) << 8
		if i+1 < n {
			x |= uint32(p[i+1])
		}
		c.sum += x
		c.sum = (c.sum & 0xffff) + (c.sum >> 16)
	}
}

func (c *Calc) Checksum() uint16 {
	return uint16(^c.sum)
}
