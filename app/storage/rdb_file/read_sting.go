package rdbfile

import (
	"bufio"
	"encoding/binary"
	"io"
	"strconv"
)

func readString(br *bufio.Reader) (string, error) {
	l, isSpec, err := readLength(br)
	if err != nil {
		return "", err
	}

	if !isSpec {
		buf := make([]byte, l)
		if _, err := io.ReadFull(br, buf); err != nil {
			return "", noEOF(err)
		}
		return string(buf), nil
	}

	switch l {
	case 0: // C0: int8
		b, err := br.ReadByte()
		if err != nil {
			return "", noEOF(err)
		}
		return strconv.Itoa(int(int8(b))), nil
	case 1: // C1: int16 little-endian
		var b [2]byte
		if _, err := io.ReadFull(br, b[:]); err != nil {
			return "", noEOF(err)
		}
		return strconv.Itoa(int(int16(binary.LittleEndian.Uint16(b[:])))), nil
	case 2: // C2: int32 little-endian
		var b [4]byte
		if _, err := io.ReadFull(br, b[:]); err != nil {
			return "", noEOF(err)
		}
		return strconv.Itoa(int(int32(binary.LittleEndian.Uint32(b[:])))), nil
	case 3:
		return "", ErrLZF
	}

	return "", nil
}
