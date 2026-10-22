package bitmap

import "slices"

func (b *BitMap) SetData(data []byte) {
	b.data = slices.Clone(data)
}

func (b *BitMap) SetLen(l int) {
	b.len = l
}

func (b *BitMap) Data() []byte {
	return b.data
}

func (b *BitMap) Len() int { return b.len }

func (b *BitMap) ExtendTo(l int) { b.extendTo(l) }
