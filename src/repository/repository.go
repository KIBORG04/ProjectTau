package repository

import (
	"fmt"
	"ssstatistics/domain"
	"ssstatistics/config"

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

}

func AutoMigrate() {
	Database.AutoMigrate(domain.Models...)
}
