package parser

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

func Deserialize(data []byte) (interface{}, error) {
	// slog.Info("start deserialization", "data", data)
	if len(data) == 0 {
		return nil, errors.New("Empty input data")
	}

	s := string(data)
	res, rest, err := deserializeString(s)
	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("Extra data after deserialization: %s", rest)
	}
	slog.Debug("Deserialize", "output", res)
	return res, nil
}

func deserializeString(s string) (interface{}, string, error) {
	switch s[0] {
	case '*': // array
		n, s, err := parseNumber(s[1:])
		if err != nil {
			return nil, s, err
		}
		if n < 0 {
			return nil, s, errors.New("Invalid array length")
		}
		array := make([]Serializable, n)
		for i := 0; i < n; i++ {
			element, remaining, err := deserializeString(s)
			if err != nil {
				return nil, remaining, err
			}
			array[i] = element.(Serializable)
			s = remaining
		}
		return array, s, nil
	case '+': // simple string
		var val string
		val, s, _ = strings.Cut(s[1:], "\r\n")
		return SimpleString(val), s, nil
	case '-': // simple error
		var val string
		val, s, _ = strings.Cut(s[1:], "\r\n")
		return SimpleError(val), s, nil
	case ':': // integer
		var val string
		val, s, _ = strings.Cut(s[1:], "\r\n")
		val_int, err := strconv.Atoi(val)
		if err != nil {
			return nil, s, errors.New("Invalid integer value")
		}
		return Int(val_int), s, nil
	case '$': // bulk string
		n, s, err := parseNumber(s[1:])
		if err != nil {
			return nil, s, err
		}
		if n < 0 {
			return nil, s, errors.New("Invalid bulk string length")
		}
		if len(s) < n+2 {
			return nil, s, errors.New("Incomplete bulk string data")
		}
		val := s[:n]
		s = s[n+2:] // Skip the CRLF after the bulk string
		return BulkString(val), s, nil
	default:
		return nil, s, fmt.Errorf("Invalid input data type ``%s''", s)
	}
}

func parseNumber(input string) (int, string, error) {
	num, left, exec := strings.Cut(input, "\r\n")
	if !exec {
		return 0, "", errors.New("Invalid number format")
	}
	n, err := strconv.Atoi(num)
	if err != nil {
		return 0, "", errors.New("Invalid number value")
	}
	return n, left, nil
}
