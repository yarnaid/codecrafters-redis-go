package rdbfile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"my-redis/app/storage"
)

func Parse(r io.Reader) (*Snapshot, error) {
	br := bufio.NewReader(r)

	snap, err := readHeader(br)
	if err != nil {
		return nil, err
	}

	var expireAt time.Time
	for {
		op, err := br.ReadByte()
		if err != nil {
			if err == io.EOF {
				return snap, nil
			}
			return nil, parseErrorf(noEOF(err), "read op code")
		}

		switch op {
		case opAux:
			key, err := readString(br)
			if err != nil {
				return nil, parseErrorf(err, "invalid string")
			}
			val, err := readString(br)
			if err != nil {
				return nil, parseErrorf(err, "invalid string value")
			}
			snap.Metadata[key] = val
		case opSelectDB:
			if _, _, err := readLength(br); err != nil {
				return nil, parseErrorf(err, "invalid db index")
			}
		case opHashTableSize:
			if _, _, err := readLength(br); err != nil {
				return nil, parseErrorf(err, "invalid hash table size: kv")
			}
			if _, _, err := readLength(br); err != nil {
				return nil, parseErrorf(err, "invalid hash table size: expires keys")
			}
		case typeString:
			k, err := readString(br)
			if err != nil {
				return nil, parseErrorf(err, "invalid db key name")
			}
			val, err := readString(br)
			if err != nil {
				return nil, parseErrorf(err, "invalid db val")
			}
			snap.Data[k] = &storage.Value{Value: val, Kind: storage.KindString, ExpireAt: expireAt}
			expireAt = time.Time{}
		case opExpireSec:
			expireAt, err = readTimestampSec(br)
			if err != nil {
				return nil, err
			}
		case opExpireMs:
			expireAt, err = readTimestampMs(br)
			if err != nil {
				return nil, err
			}
		case opEOF:
			return snap, nil
		default:
			return snap, parseErrorf(ErrUnknownOpCode, "opcode: `%v`", op)
		}
	}

	return snap, nil
}

func readHeader(br *bufio.Reader) (*Snapshot, error) {
	snap := NewSnapshot()
	var hdr [9]byte

	if _, err := io.ReadFull(br, hdr[:]); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, &ParseError{Msg: "input too short for RDB header"}
		}
		return nil, fmt.Errorf("header reading error: %w", err)
	}

	if string(hdr[:5]) != "REDIS" {
		return nil, &ParseError{Msg: "`REDIS` keyword not found"}
	}

	v, err := strconv.Atoi(string(hdr[5:]))
	if err != nil {
		return nil, parseErrorf(err, "invalid version")
	}
	snap.Version = v
	return snap, nil
}
