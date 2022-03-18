package parser

import (
	"fmt"
	"strings"
	"time"
)

const (
	urlBody     = "https://stat.taucetistation.org/json"
	filePostfix = "stat.json"
)

// first formats date to yyyy/dd/mm then formats date to yyyy/mm/dd
func normalFormattedDate(date *time.Time) string {
	return strings.ReplaceAll(date.Format("2006-01-02"), "-", "/")
}

func dateUrl(date *time.Time) string {
	return fmt.Sprintf("%s/%s", urlBody, normalFormattedDate(date))
}

func statUrl(date *time.Time, roundPostfix string) string {
	return fmt.Sprintf("%s/%s/%s/%s",
		urlBody, normalFormattedDate(date), roundPostfix, filePostfix)
}
