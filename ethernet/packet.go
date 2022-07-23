package ethernet

type Encoder interface {
	Len() int
	Sum() uint32
	Encode([]byte) error
}

type Packet []Encoder

func (p Packet) Len() int {
	n := 0
	for _, e := range p {
		n += e.Len()
	}
	return n
}

func (p Packet) Sum() uint32 {
	sum := uint32(0)
	for _, e := range p {
		sum += e.Sum()
	}
	return sum
}

func (p Packet) Encode(b []byte) error {
	off := 0
	for _, e := range p {
		err := e.Encode(b[off:])
		if err != nil {
			return err
		}
		off += e.Len()
	}
	return nil
}

type Data []byte

func (d Data) Len() int {
	return len(d)
}

func (d Data) Sum() uint32 {
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
	return sum
}

func (d Data) Encode(b []byte) error {
	copy(b, d)
	return nil
}
