// Package rdbfile ...
package rdbfile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"my-redis/app/replication"
	"my-redis/app/storage"
)

var logger = slog.Default().With("component", "rdbfile")

func Load(path, name string) (*Snapshot, error) {
	dbFilePath := filepath.Join(path, name)
	file, err := os.Open(dbFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			var err1 error
			file, err1 = createEmptyDB(dbFilePath)
			if err1 != nil {
				logger.Error("failed creating a new rdb file", "err", err, "err1", err1)
			}
			return Parse(file)
		}
		return nil, fmt.Errorf("rdb loading: %w", err)
	}
	return Parse(file)
}

func Parse(r io.Reader) (*Snapshot, error) {
	br := bufio.NewReader(r)

	snap, err := readHeader(br)
	if err != nil {
		return nil, err
	}

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
			snap.Data[k] = &storage.Value{Value: val, Kind: storage.KindString}
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

func createEmptyDB(dbFilePath string) (*os.File, error) {
	data := replication.EmptyDBContent()
	os.WriteFile(dbFilePath, data, 0o666)
	return os.Open(dbFilePath)
	// return nil, fmt.Errorf("empty db creation is not implemented")
}

func noEOF(err error) error {
	if errors.Is(err, io.EOF) {
		return io.ErrUnexpectedEOF
	}
	return err
}
