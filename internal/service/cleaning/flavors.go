package cleaning

import (
	"fmt"
	"gorm.io/gorm"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
)

type StatisticsSliceCleaner interface {
	HasAttributes(args []string) (any, bool)
}

type manifestEntriesCounter struct {
	flavor string
	name   string
	ids    []int32
}

type duplicatedFlavors []*manifestEntriesCounter

func (d duplicatedFlavors) HasAttributes(args []string) (any, bool) {
	for _, flavor := range d {
		if flavor.name == args[0] && flavor.flavor == args[1] {
			return flavor, true
		}
	}
	return nil, false
}

func CleanDuplicatedFlavors() []string {
	query := r.Database.
		Preload("ManifestEntries", r.PreloadSelect("RootID", "ID", "Name", "Flavor"),
			func(tx *gorm.DB) *gorm.DB {
				return tx.Where("flavor IS NOT null and flavor <> ''")
			})

	var roots []*domain.Root
	query.Find(&roots)

	duplicates := make(duplicatedFlavors, 0)
	for _, root := range roots {
		for _, log := range root.ManifestEntries {
			if v, ok := duplicates.HasAttributes([]string{log.Name, log.Flavor}); ok {
				manifest := v.(*manifestEntriesCounter)
				manifest.ids = append(manifest.ids, log.ID)
			} else {
				duplicates = append(duplicates, &manifestEntriesCounter{
					flavor: log.Flavor,
					name:   log.Name,
					ids:    []int32{log.ID},
				})
			}
		}
	}

	idsToNull := make([]int32, 0)
	for _, duplicate := range duplicates {
		if len(duplicate.ids) > 1 {
			// add all except last element
			for i := 0; i < len(duplicate.ids)-1; i++ {
				idsToNull = append(idsToNull, duplicate.ids[i])
			}
		}
	}

	for _, id := range idsToNull {
		r.SetToNullFlavorById(id)
	}

	return []string{
		fmt.Sprintf("%d Flavors Nulled", len(idsToNull)),
	}
}
