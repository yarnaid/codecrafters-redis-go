// Package zset ...
package zset

type value struct {
	name  string
	score float64
}

type ZSet interface {
	Add(name string, score float64) bool
	// Remove(string) int
	// Score(string) float64
	// Count() int
	// Range(int, int) iter.Seq[value]
	Rank(string) (int, bool)
}
