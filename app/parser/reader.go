package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
)

var EmptyInputError = errors.New("empty input")

const BUFFER_SIZE = 512

type reader struct {
	conn io.Reader
}

type loggedReader struct {
	r io.Reader
}

func (l *loggedReader) Read(p []byte) (n int, err error) {
	n, err = l.r.Read(p)
	slog.Debug("[logged reader]", "data", string(p[:n]), "n", n, "err", err)
	return n, err
}

func NewReader(conn io.Reader) *reader {
	return &reader{
		conn: &loggedReader{conn},
	}
}

func (r *reader) ReadArray() ([]string, error) {
	rd := bufio.NewReader(r.conn)

	l, err := readHead(rd, '*')
	if err != nil {
		if l == 0 {
			return nil, EmptyInputError
		}
		return nil, fmt.Errorf("cannot read array len: %w", err)
	}

	res := make([]string, 0, l)

	for range l {
		size, err := readHead(rd, '$')
		if err != nil {
			return nil, fmt.Errorf("cannot parse bulk string len: %w", err)
		}
		if size < 0 {
			return nil, fmt.Errorf("zero bulk string len")
		}

		buf := make([]byte, size+2)
		if _, err := io.ReadFull(rd, buf); err != nil {
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
	if len(line) < 4 || line[0] != prefix || line[len(line)-2] != '\r' {
		return 0, fmt.Errorf("bad header %q", line)
	}
	n, err := strconv.Atoi(line[1 : len(line)-2])
	if err != nil {
		return 0, fmt.Errorf("bad length")
	}
	return n, nil
}
