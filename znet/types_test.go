package znet

import (
	"testing"
)

func TestUint8(t *testing.T) {
	a := Uint8(0x12)
	t.Run("get", func(t *testing.T) {
		want := uint8(0x12)
		got := a.Get()
		if got != want {
			t.Errorf("want %v; but got %v\n", want, got)
		}
	})
	t.Run("set", func(t *testing.T) {
		want := uint8(0xc7)
		a.Set(0xc7)
		got := a.Get()
		if got != want {
			t.Errorf("want %v; but got %v\n", want, got)
		}
	})
}

func TestUint16(t *testing.T) {
	a := Uint16([...]byte{0x12, 0x34})
	t.Run("get", func(t *testing.T) {
		want := uint16(0x1234)
		got := a.Get()
		if got != want {
			t.Errorf("want %v; but got %v\n", want, got)
		}
	})
	t.Run("set", func(t *testing.T) {
		want := uint16(0xc7a8)
		a.Set(0xc7a8)
		got := a.Get()
		if got != want {
			t.Errorf("want %v; but got %v\n", want, got)
		}
	})
}

func TestOptUint8(t *testing.T) {
	a := OptUint8([]byte{0x12})
	t.Run("get", func(t *testing.T) {
		want := uint8(0x12)
		got := a.Get()
		if got != want {
			t.Errorf("want %v; but got %v\n", want, got)
		}
	})
	t.Run("set", func(t *testing.T) {
		want := uint8(0xc7)
		a.Set(0xc7)
		got := a.Get()
		if got != want {
			t.Errorf("want %v; but got %v\n", want, got)
		}
	})
}

func TestOptUint16(t *testing.T) {
	a := OptUint16([]byte{0x12, 0x34})
	t.Run("get", func(t *testing.T) {
		want := uint16(0x1234)
		got := a.Get()
		if got != want {
			t.Errorf("want %v; but got %v\n", want, got)
		}
	})
	t.Run("set", func(t *testing.T) {
		want := uint16(0xc7a8)
		a.Set(0xc7a8)
		got := a.Get()
		if got != want {
			t.Errorf("want %v; but got %v\n", want, got)
		}
	})
}
