package storage

// deref returns *p, or the string "<nil>" for logging purposes.
func deref[T any](p *T) any {
	if p == nil {
		return "<nil>"
	}
	return *p
}
