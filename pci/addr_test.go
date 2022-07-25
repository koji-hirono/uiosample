package pci

import (
	"testing"
)

func TestParseAddr(t *testing.T) {
	tests := []struct {
		name string
		text string
		addr Addr
	}{
		{
			name: "Simple BFD notation",
			text: "a2:3b.d",
			addr: Addr{
				Domain: 0,
				Bus:    0xa2,
				ID:     0x3b,
				Func:   0xd,
			},
		},
		{
			name: "BDF Notation Extension for PCI Domain",
			text: "95ec:a2:3b.d",
			addr: Addr{
				Domain: 0x95ec,
				Bus:    0xa2,
				ID:     0x3b,
				Func:   0xd,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := ParseAddr(tc.text)
			if err != nil {
				t.Fatal(err)
			}
			if addr.Domain != tc.addr.Domain {
				t.Errorf("domain want: %x; but got %x\n", tc.addr.Domain, addr.Domain)
			}
			if addr.Bus != tc.addr.Bus {
				t.Errorf("bus want: %x; but got %x\n", tc.addr.Bus, addr.Bus)
			}
			if addr.ID != tc.addr.ID {
				t.Errorf("id want: %x; but got %x\n", tc.addr.ID, addr.ID)
			}
			if addr.Func != tc.addr.Func {
				t.Errorf("func want: %x; but got %x\n", tc.addr.Func, addr.Func)
			}
		})
	}
}

func TestParseAddr_Abnormal(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr error
	}{
		{
			name:    "illegal Bus",
			text:    "ag:3b.d",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "illegal Domain",
			text:    "95e :a2:3b.d",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "illegal ID",
			text:    "95ed:a2:.d",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "illegal Func",
			text:    "95ed:a2:3b.y",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "illegal first colon",
			text:    "95ed.a2:3b.d",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "illegal second colon",
			text:    "95ed:a2,3b.d",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "too many fields",
			text:    "95ed:a2.3b.d",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "tail",
			text:    "95ed:a2:3b.dsss",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "head",
			text:    "   95ed:a2:3b.d",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "empty",
			text:    "",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "too few fields",
			text:    "95ed:a2",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "missing all fields",
			text:    "::.",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "missing domain field",
			text:    ":0:0.0",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "missing bus field",
			text:    "0::0.0",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "missing id field",
			text:    "0:0:.0",
			wantErr: ErrIllegalFormat,
		},
		{
			name:    "missing func field",
			text:    "0:0:0.",
			wantErr: ErrIllegalFormat,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseAddr(tc.text)
			if err != tc.wantErr {
				t.Errorf("wantErr: %v; but got %v\n", tc.wantErr, err)
			}
		})
	}
}
