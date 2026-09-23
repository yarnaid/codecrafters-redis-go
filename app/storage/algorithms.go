package storage

import (
	"my-redis/app/parser"
)

func rangeList(list parser.Array[parser.Serializable], start, end int) parser.Array[parser.Serializable] {
	l := len(list)
	if end >= 0 {
		end = min(l-1, end)
	}
	start, end = invert_index(start, l), invert_index(end, l)
	// slog.Debug("[Storage][LRange] returning", "list", list, "start", start, "end", end)
	return list[start : end+1]
}

func invert_index(i, ln int) int {
	if i >= 0 {
		return i
	} else {
		return max(ln+i, 0)
	}
}

func popListLeft(list []parser.Serializable, n int) (popped, remaining []parser.Serializable) {
	if n == 0 {
		return []parser.Serializable{}, list
	}
	if n >= 0 {
		n = min(len(list), n)
	}

	return list[:n], list[n:]
}

func PrependReversed[T any](s []T, values ...T) []T {
	result := make([]T, 0, len(values)+len(s))
	for i := len(values) - 1; i >= 0; i-- {
		result = append(result, values[i])
	}
	return append(result, s...)
}
