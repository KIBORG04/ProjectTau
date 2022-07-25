package utils

import "regexp"

var (
	RoundId   = regexp.MustCompile(`\d\d\d\d\d`)
	IsDrone   = regexp.MustCompile(`maintenance drone \(\d+\)`)
	IsMobName = regexp.MustCompile(`\w+ \(\d+\)`)
)
