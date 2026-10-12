// Package zset ...
package zset

import "iter"

type value struct {
	name  string
	score float64
}

type ZSet interface {
	Add(name string, score float64) bool
	Remove(string) bool
	Score(string) (float64, bool)
	Count() int
	Range(int, int) iter.Seq[string]
	Rank(string) (int, bool)
}
