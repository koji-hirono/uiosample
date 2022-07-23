#include "textflag.h"

// func Rdtsc() uint64
TEXT	·Rdtsc(SB), NOSPLIT, $0-8
	CPUID
	MFENCE
	RDTSC
	SHLQ	$32, DX
	ADDQ	DX, AX
	MOVQ	AX, ret+0(FP)
	RET
