package repository

import (
	"fmt"
	"ssstatistics/internal/config"
	d "ssstatistics/internal/domain"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		})

	//Database.Migrator().DropTable("links")
	//tables, _ := Database.Migrator().GetTables()
	//Database.Migrator().DropTable(tables)
	AutoMigrate()

}

func AutoMigrate() {
	Database.AutoMigrate(d.Models...)
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
