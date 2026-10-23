package geo

const normMax = 1 << 26

func norm(v float64, mn float64, rng float64) float64 {
	return normMax * (v - mn) / rng
}

func normLong(v grad) float64 {
	return norm(float64(v), MinLong, LongRange)
}

func normLat(v grad) float64 {
	return norm(float64(v), MinLat, LatRange)
}

func round(v float64) uint32 {
	return uint32(v)
}

func spreadInt32ToInt64(v uint32) uint64 {
	result := uint64(v)
	result = (result | (result << 16)) & 0x0000FFFF0000FFFF
	result = (result | (result << 8)) & 0x00FF00FF00FF00FF
	result = (result | (result << 4)) & 0x0F0F0F0F0F0F0F0F
	result = (result | (result << 2)) & 0x3333333333333333
	result = (result | (result << 1)) & 0x5555555555555555
	return result
}

func interleave(x, y uint32) uint64 {
	xSpread := spreadInt32ToInt64(x)
	ySpread := spreadInt32ToInt64(y)
	yShifted := ySpread << 1
	return xSpread | yShifted
}

func compactInt64ToInt32(v uint64) uint32 {
	result := v & 0x5555555555555555
	result = (result | (result >> 1)) & 0x3333333333333333
	result = (result | (result >> 2)) & 0x0F0F0F0F0F0F0F0F
	result = (result | (result >> 4)) & 0x00FF00FF00FF00FF
	result = (result | (result >> 8)) & 0x0000FFFF0000FFFF
	result = (result | (result >> 16)) & 0x00000000FFFFFFFF
	return uint32(result)
}

func convertGridNumbersToCoordinates(gridLatitudeNumber, gridLongitudeNumber uint32) (float64, float64) {
	// Calculate the grid boundaries
	gridLatitudeMin := MinLat + LatRange*(float64(gridLatitudeNumber)/normMax)
	gridLatitudeMax := MinLat + LatRange*(float64(gridLatitudeNumber+1)/normMax)
	gridLongitudeMin := MinLong + LongRange*(float64(gridLongitudeNumber)/normMax)
	gridLongitudeMax := MinLong + LongRange*(float64(gridLongitudeNumber+1)/normMax)

	// Calculate the center point of the grid cell
	latitude := (gridLatitudeMin + gridLatitudeMax) / 2
	longitude := (gridLongitudeMin + gridLongitudeMax) / 2
	return latitude, longitude
}
