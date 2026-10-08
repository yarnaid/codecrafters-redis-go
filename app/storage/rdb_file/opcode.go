package rdbfile

const (
	opAux           = 0xFA
	opHashTableSize = 0xFB
	opExpireMs      = 0xFC
	opExpireSec     = 0xFD
	opSelectDB      = 0xFE
	opEOF           = 0xFF
	typeString      = 0x00

	maxStringLen = 512 << 20 // Redis' own proto-max-bulk-len default; guards against corrupt lengths
)
