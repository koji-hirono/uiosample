package hugetlb

import (
	"log"
	"testing"
	"unsafe"
)

func TestAlloc(t *testing.T) {
	err := SetPages(8)
	if err != nil {
		t.Fatal(err)
	}
	Init()

	for i := 0; i < 2; i++ {
		func() {
			buf, phys, err := Alloc(4 * 1024)
			if err != nil {
				t.Fatal(err)
			}

			log.Printf("buf: %p\n", &buf[0])

			virt := uintptr(unsafe.Pointer(&buf[0]))
			log.Printf("virt: %x\n", virt)

			log.Printf("phys: %x\n", phys)
		}()
	}
}
