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
