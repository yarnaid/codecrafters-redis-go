package replication

import (
	"encoding/base64"
	"strconv"
)

func EmptyDB() []byte {
	data, _ := base64.StdEncoding.DecodeString("UkVESVMwMDEx+glyZWRpcy12ZXIFNy4yLjD6CnJlZGlzLWJpdHPAQPoFY3RpbWXCbQi8ZfoIdXNlZC1tZW3CsMQQAPoIYW9mLWJhc2XAAP/wbjv+wP9aog==")
	header := strconv.AppendInt([]byte{'$'}, int64(len(data)), 10)
	header = append(header, '\r', '\n')

	out := make([]byte, 0, len(header)+len(data))
	out = append(out, header...)
	out = append(out, data...) // raw bytes, no CRLF after
	return out
}
