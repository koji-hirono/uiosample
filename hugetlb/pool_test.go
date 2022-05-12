package hugetlb

import (
	"testing"
)

func TestPool(t *testing.T) {
	buf := make([]byte, 128)
	p := NewPool(buf, 32)
	if p.free != nil {
		t.Fatal("free is not nil")
	}
	if p.used != 0 {
		t.Fatal("used is not 0")
	}
	t.Run("1st get 1 and put 1", func(t *testing.T) {
		e, ok := p.Get()
		if !ok {
			t.Fatalf("want %v; but got %v", true, ok)
		}
		if &e[0] != &buf[0] {
			t.Errorf("want %v; but got %v", &buf[0], &e[0])
		}
		if p.free != nil {
			t.Fatal("free is not nil")
		}
		if p.used != 32 {
			t.Fatal("used is not 32")
		}
		p.Put(e)
		if p.free == nil {
			t.Fatal("free is nil")
		}
		if p.used != 32 {
			t.Fatal("used is not 32")
		}
	})
	t.Run("2nd get 1 and put 1", func(t *testing.T) {
		e, ok := p.Get()
		if !ok {
			t.Fatalf("want %v; but got %v", true, ok)
		}
		if &e[0] != &buf[0] {
			t.Errorf("want %v; but got %v", &buf[0], &e[0])
		}
		if p.free != nil {
			t.Fatal("free is not nil")
		}
		if p.used != 32 {
			t.Fatal("used is not 32")
		}
		p.Put(e)
		if p.free == nil {
			t.Fatal("free is nil")
		}
		if p.used != 32 {
			t.Fatal("used is not 32")
		}
	})
	t.Run("get 2 and put 2", func(t *testing.T) {
		e1, ok := p.Get()
		if !ok {
			t.Fatalf("want %v; but got %v", true, ok)
		}
		if &e1[0] != &buf[0] {
			t.Errorf("want %v; but got %v", &buf[0], &e1[0])
		}
		if p.free != nil {
			t.Fatal("free is not nil")
		}
		if p.used != 32 {
			t.Fatal("used is not 32")
		}

		e2, ok := p.Get()
		if !ok {
			t.Fatalf("want %v; but got %v", true, ok)
		}
		if &e2[0] != &buf[32] {
			t.Errorf("want %v; but got %v", &buf[32], &e2[0])
		}
		if p.free != nil {
			t.Fatal("free is not nil")
		}
		if p.used != 64 {
			t.Fatal("used is not 64")
		}

		p.Put(e1)
		if p.free == nil {
			t.Fatal("free is nil")
		}
		if p.used != 64 {
			t.Fatal("used is not 64")
		}

		p.Put(e2)
		if p.free == nil {
			t.Fatal("free is nil")
		}
		if p.used != 64 {
			t.Fatal("used is not 64")
		}
	})
	t.Run("get 2 and put 1 and get 1 and put 2", func(t *testing.T) {
		e1, ok := p.Get()
		if !ok {
			t.Fatalf("want %v; but got %v", true, ok)
		}
		if &e1[0] != &buf[32] {
			t.Errorf("want %v; but got %v", &buf[0], &e1[32])
		}
		if p.free == nil {
			t.Fatal("free is nil")
		}
		if p.used != 64 {
			t.Fatal("used is not 64")
		}

		e2, ok := p.Get()
		if !ok {
			t.Fatalf("want %v; but got %v", true, ok)
		}
		if &e2[0] != &buf[0] {
			t.Errorf("want %v; but got %v", &buf[0], &e2[0])
		}
		if p.free != nil {
			t.Fatalf("free is not nil %p", p.free)
		}
		if p.used != 64 {
			t.Fatal("used is not 64")
		}

		p.Put(e1)
		if p.free == nil {
			t.Fatal("free is nil")
		}
		if p.used != 64 {
			t.Fatal("used is not 64")
		}

		e3, ok := p.Get()
		if !ok {
			t.Fatalf("want %v; but got %v", true, ok)
		}
		if &e3[0] != &buf[32] {
			t.Errorf("want %v; but got %v", &buf[32], &e3[0])
		}
		if p.free != nil {
			t.Fatal("free is not nil")
		}
		if p.used != 64 {
			t.Fatal("used is not 64")
		}

		p.Put(e2)
		if p.free == nil {
			t.Fatal("free is nil")
		}
		if p.used != 64 {
			t.Fatal("used is not 64")
		}

		p.Put(e3)
		if p.free == nil {
			t.Fatal("free is nil")
		}
		if p.used != 64 {
			t.Fatal("used is not 64")
		}
	})
}

func TestPool_Get(t *testing.T) {
	buf := make([]byte, 8)
	p := NewPool(buf, 8)
	if p.free != nil {
		t.Fatal("free is not nil")
	}
	if p.used != 0 {
		t.Fatal("used is not 0")
	}
	t.Run("get 2", func(t *testing.T) {
		e1, ok := p.Get()
		if !ok {
			t.Fatalf("want %v; but got %v", true, ok)
		}
		defer p.Put(e1)
		e2, ok := p.Get()
		if ok {
			p.Put(e2)
			t.Fatalf("want %v; but got %v", false, ok)
		}
	})
}
