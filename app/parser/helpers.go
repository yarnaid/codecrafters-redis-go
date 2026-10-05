package parser

func StringsToArray(s ...string) Array[BulkString] {
	res := make(Array[BulkString], len(s))
	for i := range len(s) {
		res[i] = BulkString(s[i])
	}
	return res
}

func StringsToBytes(s ...string) []byte {
	a := StringsToArray(s...)
	res, _ := a.Serialize()
	return res
}

func ToSerializableSlice[T any](vals []T) []Serializable {
	res := make([]Serializable, len(vals))
	for i, val := range vals {
		res[i] = ToSerializable(val)
	}
	return res
}

func Encode(s Serializable) []byte {
	r, err := s.Serialize()
	if err != nil {
		return ReplyError(err.Error())
	}
	return r
}
