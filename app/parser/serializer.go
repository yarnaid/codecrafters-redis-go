package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Serializable interface {
	Serialize() ([]byte, error)
}

type SC interface {
	Serializable
	comparable
}

type (
	Int            int
	SimpleString   string
	BulkString     string
	SimpleError    string
	Array[t SC]    []t
	BulkNullString string
)

func Serialize(val Serializable) ([]byte, error) {
	return val.Serialize()
}

func (i Int) Serialize() ([]byte, error) {
	return []byte(":" + strconv.Itoa(int(i)) + "\r\n"), nil
}

func (s SimpleString) Serialize() ([]byte, error) {
	if strings.ContainsAny(string(s), "\r\n") {
		return nil, errors.New("SimpleString cannot contain CR or LF characters")
	}
	return []byte("+" + string(s) + "\r\n"), nil
}

func (b BulkString) Serialize() ([]byte, error) {
	return []byte("$" + strconv.Itoa(len(b)) + "\r\n" + string(b) + "\r\n"), nil
}

func (e SimpleError) Serialize() ([]byte, error) {
	return []byte("-" + string(e) + "\r\n"), nil
}

func (a Array[Serializable]) Serialize() ([]byte, error) {
	var serializedElements []byte
	for _, element := range a {
		serializedElement, err := element.Serialize()
		if err != nil {
			return nil, err
		}
		serializedElements = append(serializedElements, serializedElement...)
	}

	return []byte("*" + strconv.Itoa(len(a)) + "\r\n" + string(serializedElements)), nil
}

func (a *Array[T]) Equal(other *Array[T]) bool {
	if len(*a) != len(*other) {
		return false
	}
	for i := 0; i < len(*a); i++ {
		if (*a)[i] != (*other)[i] {
			return false
		}
	}
	return true
}

func (b BulkNullString) Serialize() ([]byte, error) {
	return []byte("$-1\r\n"), nil
}

func ToSerializable(val interface{}) Serializable {
	switch val := val.(type) {
	case int:
		return Int(val)
	case string:
		return BulkString(val)
	case []interface{}:
		var res Array[Serializable]
		for _, v := range val {
			res = append(res, ToSerializable(v))
		}
		return res
	default:
		panic(fmt.Sprintf("Not supported type for serializable: %T; %v", val, val))
	}
}
