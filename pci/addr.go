package pci

import (
	"errors"
	"fmt"
)

type Addr struct {
	Domain uint32
	Bus    uint8
	ID     uint8
	Func   uint8
}

var (
	ErrIllegalFormat = errors.New("illegal format")
)

// BDF notation:
// expr -> (domain ':')? bus ':' device '.' func
// domain -> hexdigit{1,4}
// bus -> hexdigit{1,2}
// device -> hexdigit{1,2}
// func -> hexdigit{1}
// hexdigit -> [0-9a-fA-F]
func ParseAddr(s string) (*Addr, error) {
	addr := new(Addr)
	var i int
	n := len(s)

	x1, m, ok := hextoi(s)
	if !ok {
		return nil, ErrIllegalFormat
	}
	i += m
	if i >= n {
		return nil, ErrIllegalFormat
	}

	if s[i] != ':' {
		return nil, ErrIllegalFormat
	}
	i++
	if i >= n {
		return nil, ErrIllegalFormat
	}

	x2, m, ok := hextoi(s[i:])
	if !ok {
		return nil, ErrIllegalFormat
	}
	i += m
	if i >= n {
		return nil, ErrIllegalFormat
	}

	if s[i] == ':' {
		i++
		if i >= n {
			return nil, ErrIllegalFormat
		}
		x3, m, ok := hextoi(s[i:])
		if !ok {
			return nil, ErrIllegalFormat
		}
		i += m
		if i >= n {
			return nil, ErrIllegalFormat
		}
		addr.Domain = uint32(x1)
		addr.Bus = uint8(x2)
		addr.ID = uint8(x3)
	} else {
		addr.Domain = 0
		addr.Bus = uint8(x1)
		addr.ID = uint8(x2)
	}

	if s[i] != '.' {
		return nil, ErrIllegalFormat
	}
	i++
	if i >= n {
		return nil, ErrIllegalFormat
	}

	x4, m, ok := hextoi(s[i:])
	if !ok {
		return nil, ErrIllegalFormat
	}
	i += m
	if i != n {
		return nil, ErrIllegalFormat
	}
	addr.Func = uint8(x4)
	return addr, nil
}

func hextoi(s string) (uint64, int, bool) {
	var x uint64
	var match bool
	n := len(s)
	for i := 0; i < n; i++ {
		switch s[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			x <<= 4
			x |= uint64(s[i] - '0')
			match = true
		case 'a', 'b', 'c', 'd', 'e', 'f':
			x <<= 4
			x |= uint64(s[i]-'a') + 10
			match = true
		case 'A', 'B', 'C', 'D', 'E', 'F':
			x <<= 4
			x |= uint64(s[i]-'A') + 10
			match = true
		default:
			return x, i, match
		}
	}
	return x, n, match
}

func (a *Addr) String() string {
	return fmt.Sprintf("%04x:%02x:%02x.%01x", a.Domain, a.Bus, a.ID, a.Func)
}
