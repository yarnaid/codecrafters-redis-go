package rdbfile

import (
	"bufio"
	"time"
)

func ReadLength(br *bufio.Reader) (uint64, bool, error) {
	return readLength(br)
}

func ReadString(br *bufio.Reader) (string, error) {
	return readString(br)
}

func ReadTimestampSec(br *bufio.Reader) (time.Time, error) {
	return readTimestampSec(br)
}

func ReadTimestampMs(br *bufio.Reader) (time.Time, error) {
	return readTimestampMs(br)
}
