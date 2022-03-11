package repository

import (
	d "ssstatistics/internal/domain"
	"time"
)

func FindAllByDate(date *time.Time) []string {
	var links []d.Link
	Database.Table("links").Where("date = ?", date.Format("2006-01-02")).Find(&links)

	strings := make([]string, 0, cap(links))
	for _, v := range links {
		strings = append(strings, v.Link)
	}

	return strings
}

func SaveDate(link *d.Link) {
	Database.Save(link)
}
