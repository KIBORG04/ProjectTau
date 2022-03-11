package repository

import (
	"fmt"
	"ssstatistics/internal/config"
	d "ssstatistics/internal/domain"

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

func Save(v interface{}) {
	Database.Save(v)
}


