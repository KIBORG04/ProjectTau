package repository

import "ssstatistics/internal/domain"

func RemoveAnnounces(ids []int32) {
	Database.Delete(&domain.CommunicationLogs{}, ids)
}
