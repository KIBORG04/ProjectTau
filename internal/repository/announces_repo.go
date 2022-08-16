package repository

import "ssstatistics/internal/domain"

func RemoveAnnounces(ids []int32) {
	if len(ids) > 0 {
		Database.Delete(&domain.CommunicationLogs{}, ids)
	}
}
