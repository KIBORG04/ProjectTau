package repository

import (
	"fmt"
	"ssstatistics/internal/config"
	d "ssstatistics/internal/domain"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func CreateConnection() {
	c := &config.Config.DatabaseConfig
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.Dbname, c.Port,
	)

	Database, _ = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   logger.Default.LogMode(logger.Silent),
		})

	//tables, _ := Database.Migrator().GetTables()
	//for _, table := range tables {
	//	Database.Migrator().DropTable(table)
	//}
	AutoMigrate()

}

func AutoMigrate() {
	Database.AutoMigrate(d.Models...) // Not Fucking Auto
	// Manual migrate
	for _, model := range d.Models {
		switch t := model.(type) {
		case d.MyMigrator:
			t.ColumnsMigration(Database)
		}
	}
}

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

func Save(v interface{}) {
	Database.Save(v)
}

func FindByRoundId(id string) (*d.Root, error) {
	var root d.Root
	Database.Table("roots").Where("round_id = ?", id).First(&root)
	if root.RoundID == 0 {
		return nil, fmt.Errorf("not found %s id", id)
	}
	return &root, nil
}
