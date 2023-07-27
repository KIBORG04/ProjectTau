package repository

import (
	"fmt"
	"ssstatistics/internal/domain"
)

func FindSponsor(name string) (domain.ThetaSponsor, error) {
	var sponsor domain.ThetaSponsor
	Database.Raw(`select user_id from theta_sponsors where user_id = ?`, name).Scan(&sponsor)
	if sponsor.UserId == "" {
		return sponsor, fmt.Errorf("UserId not found")
	}
	return sponsor, nil
}
