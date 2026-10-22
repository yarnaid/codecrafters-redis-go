package bitmap

func (b BitMap) getInd(i int) byte {
	if i < len(b.data) {
		return b.data[i]
	}
	return 0
}

func applyOp(b1, b2 *BitMap, op func(byte, byte) byte) *BitMap {
	l := max(b1.DataLen(), b2.DataLen())
	dataRes := make([]byte, l)
	for i := range l {
		bb1 := b1.getInd(i)
		bb2 := b2.getInd(i)
		dataRes[i] = op(bb1, bb2)
	}
	return &BitMap{
		data: dataRes,
		len:  max(b1.len, b2.len),
	}
}

func And(b1, b2 *BitMap) *BitMap {
	f := func(a, b byte) byte { return a & b }
	return applyOp(b1, b2, f)
}
