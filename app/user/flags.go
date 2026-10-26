package user

import "fmt"

type flag int16

const flagBits = 16

const (
	FlagEmpty  flag = 0
	FlagNoPass flag = 1 << iota
)

var flagNames = map[flag]string{
	FlagEmpty:  "empty",
	FlagNoPass: "nopass",
}

type flagSet flag

func (f flag) fs() flagSet {
	return flagSet(f)
}

func (f flagSet) flag() flag {
	return flag(f)
}

func (f flagSet) List() []string {
	if f == 0 {
		return []string{}
	}
	res := []string{}
	for i := range flagBits {
		var flg flagSet = 1 << i
		if f&flg == flg {
			name, ok := flagNames[flg.flag()]
			if !ok {
				name = fmt.Sprintf("invalid_flag_%d", flg)
			}
			res = append(res, name)
		}
	}
	return res
}

func (f flagSet) Has(flg flag) bool {
	return f&flg.fs() == flg.fs()
}
