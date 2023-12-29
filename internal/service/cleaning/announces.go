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

type duplicatedAnnounces []*announcesCounter

func (d duplicatedAnnounces) HasAttributes(args []string) (any, bool) {
	for _, announce := range d {
		if announce.contents == args[0] {
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

	duplicates := make(duplicatedAnnounces, 0)
	for _, root := range roots {
		for _, log := range root.CommunicationLogs {
			if v, ok := duplicates.HasAttributes([]string{log.Content}); ok {
				communicationLog := v.(*announcesCounter)
				communicationLog.count++
				communicationLog.ids = append(communicationLog.ids, log.ID)
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
