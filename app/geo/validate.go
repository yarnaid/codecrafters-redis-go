package geo

func ValidateCoord(long, lat float64) bool {
	return long <= MaxLong && long >= MinLong && lat >= MinLat && lat <= MaxLat
}
