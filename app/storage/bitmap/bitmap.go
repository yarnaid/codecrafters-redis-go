package bitmap

import (
	"fmt"
	"math/bits"
	"strings"
)

type BitMap struct {
	data []byte
	len  int
}

func (b *BitMap) Bit(offset int) int {
	if offset > b.len {
		return 0
	}
	i, ofs := idx(offset)
	bt := b.data[i]
	return int((bt >> ofs) & 1)
}

func (b *BitMap) Set(offset int) {
	i, ofs := idx(offset)
	b.extendTo(i + 1)
	b.data[i] |= 1 << uint(ofs)
	b.len = max(b.len, offset+1)
}

func (b *BitMap) Clear(offset int) {
	if offset > b.len {
		return
	}
	i, ofs := idx(offset)
	var mask byte = 1 << uint(ofs)
	b.extendTo(i)
	b.data[i] &^= mask
	b.len = max(b.len, offset+1)
}

func (b *BitMap) extendTo(l int) {
	if len(b.data) < l {
		addLen := l - len(b.data)
		b.data = append(b.data, make([]byte, addLen)...)
	}
}

func idx(offset int) (int, int) {
	i := offset / 8
	ofs := offset - i*8
	return i, ofs
}

func (b *BitMap) String() string {
	builder := strings.Builder{}
	builder.Grow(b.len)
	for _, bt := range b.data {
		s := fmt.Sprintf("%08b", uint64(bits.Reverse8(bt)))
		builder.WriteString(s)
	}
	if ex := b.len - builder.Len(); ex > 0 {
		s := strings.Repeat("0", ex)
		builder.WriteString(s)
	}
	return builder.String()[:b.len]
}

func (b *BitMap) AsString() string {
	bld := strings.Builder{}
	for _, bt := range b.data {
		bld.WriteRune(rune(bits.Reverse8(uint8(bt))))
	}
	return bld.String()
}
