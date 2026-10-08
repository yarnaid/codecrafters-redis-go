package rdbfile

import (
	"bufio"
	"encoding/binary"
	"io"
	"time"
)

func readTimestampMs(br *bufio.Reader) (time.Time, error) {
	var buf [8]byte
	if _, err := io.ReadFull(br, buf[:]); err != nil {
		return time.Time{}, parseErrorf(err, "incorrect timestamp")
	}
	val := int64(binary.LittleEndian.Uint64(buf[:]))
	return time.UnixMilli(val), nil
}

func readTimestampSec(br *bufio.Reader) (time.Time, error) {
	var buf [4]byte
	if _, err := io.ReadFull(br, buf[:]); err != nil {
		return time.Time{}, parseErrorf(err, "incorrect timestamp")
	}
	val := int64(binary.LittleEndian.Uint32(buf[:]))
	return time.Unix(val, 0), nil
}
