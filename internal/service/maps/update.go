package maps

import (
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
)

func FixMaxShit() []string {
	r.Database.Model(new(domain.Root)).Where("map = ?", "Prometheus").Update("map", "Prometheus Station")
	return []string{"Map Shit Fixed"}
}
