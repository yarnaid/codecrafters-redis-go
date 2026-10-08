package rdbfile

import "bufio"

func ReadLength(br *bufio.Reader) (uint64, bool, error) {
	return readLength(br)
}

func ReadString(br *bufio.Reader) (string, error) {
	return readString(br)
}
