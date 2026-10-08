package server

func GetNextAOFIncr(prev []string, fname string) (int, error) {
	return getNextAOFIncr(prev, fname)
}

func NewAOFFileName(aofPath, name string) (string, error) {
	return newAOFFileName(aofPath, name)
}

func ParseManifest(s string) (Manifest, error) {
	return parseManifest(s)
}
