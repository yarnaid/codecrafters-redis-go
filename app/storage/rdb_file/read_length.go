package rdbfile

import (
	"bufio"
	"encoding/binary"
	"io"
)

func readLength(br *bufio.Reader) (uint64, bool, error) {
	b, err := br.ReadByte()
	if err != nil {
		return 0, false, noEOF(err)
	}

	switch b >> 6 {
	case 0b00: // simple 6-bit value
		return uint64(b & 0x3F), false, nil
	case 0b01: // 14-bit value
		next, err := br.ReadByte()
		if err != nil {
			return 0, false, noEOF(err)
		}
		return uint64(b&0x3F)<<8 | uint64(next), false, nil
	case 0b10: // ignore this and read next 4 bytes in big-endian
		var buf [4]byte
		if _, err := io.ReadFull(br, buf[:]); err != nil {
			return 0, false, noEOF(err)
		}
		return uint64(binary.BigEndian.Uint32(buf[:])), false, nil
	default: // 0b11 special
		return uint64(b & 0x3F), true, nil
	}
}
