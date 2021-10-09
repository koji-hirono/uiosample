package main

import (
	"fmt"
	"testing"
	"time"
)

func TestRdtsc_1msec(t *testing.T) {
	s := Rdtsc()
	time.Sleep(time.Millisecond)
	e := Rdtsc()
	fmt.Printf("s: %v\n", s)
	fmt.Printf("e: %v\n", e)
	fmt.Printf("d: %v\n", e-s)
}
