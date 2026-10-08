// Package rdbfile ...
package rdbfile

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"my-redis/app/replication"
)

var logger = slog.Default().With("component", "rdbfile")

func Load(path, name string) (*Snapshot, error) {
	dbFilePath := filepath.Join(path, name)
	file, err := os.Open(dbFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			var err1 error
			file, err1 = createEmptyDB(path, name)
			if err1 != nil {
				logger.Error("failed creating a new rdb file", "err", err, "err1", err1)
			}
			return Parse(file)
		}
		return nil, fmt.Errorf("rdb loading: %w", err)
	}
	return Parse(file)
}

func createEmptyDB(path, name string) (*os.File, error) {
	data := replication.EmptyDBContent()
	if err := os.MkdirAll(path, 0o755); err != nil {
		return nil, err
	}
	dbFilePath := filepath.Join(path, name)
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
