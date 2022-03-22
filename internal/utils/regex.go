package utils

import "regexp"

var (
	RoundId = regexp.MustCompile(`[\d][\d][\d][\d][\d]`)
)
