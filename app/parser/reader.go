package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

var EmptyInputError = errors.New("empty input")

const BUFFER_SIZE = 512

type Reader struct {
	conn io.Reader
	rd   *bufio.Reader
}

type loggedReader struct {
	r io.Reader
}

func (l *loggedReader) Read(p []byte) (n int, err error) {
	n, err = l.r.Read(p)
	// slog.Debug("[logged reader]", "data", string(p[:n]), "n", n, "err", err)
	return n, err
}

func NewReader(conn io.Reader) *Reader {
	return &Reader{
		conn: &loggedReader{conn},
		rd:   bufio.NewReader(conn),
	}
}

func (r *Reader) Empty() bool {
	// slog.Debug("reading buff size", "size", r.rd.Buffered())
	return r.rd.Buffered() == 0
}

func (r *Reader) ReadArrays() ([][]string, error) {
	res := make([][]string, 1)
	first, err := r.ReadArray()
	if err != nil {
		return nil, err
	}
	res[0] = first
	for !r.Empty() {
		arr, err := r.ReadArray()
		if err != nil {
			return nil, err
		}
		res = append(res, arr)
	}
	return res, nil
}

func (r *Reader) ReadArray() ([]string, error) {
	l, err := readHead(r.rd, '*')
	if err != nil {
		if l == 0 {
			if errors.Is(err, io.EOF) {
				return nil, err
			}
			return nil, EmptyInputError
		}
		return nil, fmt.Errorf("cannot read array len: %w", err)
	}

	res := make([]string, 0, l)

	for range l {
		size, err := readHead(r.rd, '$')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, err
			}
			return nil, fmt.Errorf("cannot parse bulk string len: %w", err)
		}
		if size < 0 {
			return nil, fmt.Errorf("zero bulk string len")
		}

		buf := make([]byte, size+2)
		if _, err := io.ReadFull(r.rd, buf); err != nil {
			if errors.Is(err, io.EOF) {
				return nil, err
			}
			return nil, fmt.Errorf("invalid bulk string: %w", err)
		}
		if buf[size] != '\r' && buf[size+1] != '\n' {
			return nil, fmt.Errorf("malformed string: %s", buf[:])
		}
		res = append(res, string(buf[:size]))
	}
	return res, nil
}

func readHead(r *bufio.Reader, prefix byte) (int, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return 0, err
	}
	if len(line) == 0 {
		return 0, fmt.Errorf("got empty string from connection, left %d", r.Buffered())
	}
	if len(line) < 4 || line[0] != prefix || line[len(line)-2] != '\r' {
		return 0, fmt.Errorf("bad header %q", line)
	}
	n, err := strconv.Atoi(line[1 : len(line)-2])
	if err != nil {
		return 0, fmt.Errorf("bad length")
	}
	return n, nil
}
