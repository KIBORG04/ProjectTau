package cleaning

import (
	"fmt"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
)

type announcesCounter struct {
	contents string
	count    uint
	ids      []int32
}

func hasContents(slices []*announcesCounter, content string) (*announcesCounter, bool) {
	for _, announce := range slices {
		if announce.contents == content {
			return announce, true
		}
	}
	return nil, false
}

func CleanAnnounces() []string {
	query := r.Database.
		Preload("CommunicationLogs", r.PreloadSelect("RootID", "ID", "Content")).
		Omit("CompletionHTML")

	var roots []*domain.Root
	query.Find(&roots)

	duplicates := make([]*announcesCounter, 0)
	for _, root := range roots {
		for _, log := range root.CommunicationLogs {
			if v, ok := hasContents(duplicates, log.Content); ok {
				v.count++
				v.ids = append(v.ids, log.ID)
			} else {
				duplicates = append(duplicates, &announcesCounter{
					contents: log.Content,
					count:    1,
					ids:      []int32{log.ID},
				})
			}
		}
	}

	idsToRemove := make([]int32, 0)
	for _, duplicate := range duplicates {
		if duplicate.count > 3 {
			idsToRemove = append(idsToRemove, duplicate.ids...)
		}
	}

	r.RemoveAnnounces(idsToRemove)
	return []string{
		fmt.Sprintf("%d CommunicationLogs Removed", len(idsToRemove)),
	}
}
