package vf

type Device struct {
	vf_number   uint32
	v2p_mailbox uint32
}

type Stats struct {
	base_gprc   uint64
	base_gptc   uint64
	base_gorc   uint64
	base_gotc   uint64
	base_mprc   uint64
	base_gotlbc uint64
	base_gptlbc uint64
	base_gorlbc uint64
	base_gprlbc uint64

	last_gprc   uint32
	last_gptc   uint32
	last_gorc   uint32
	last_gotc   uint32
	last_mprc   uint32
	last_gotlbc uint32
	last_gptlbc uint32
	last_gorlbc uint32
	last_gprlbc uint32

	gprc   uint64
	gptc   uint64
	gorc   uint64
	gotc   uint64
	mprc   uint64
	gotlbc uint64
	gptlbc uint64
	gorlbc uint64
	gprlbc uint64
}
