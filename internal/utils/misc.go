package utils

import (
	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

func TrimPGDate(s string) string {
	return s[:10]
}
