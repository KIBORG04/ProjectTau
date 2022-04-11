package repository

import (
	"fmt"
	"os"
	"ssstatistics/internal/config"
	d "ssstatistics/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func CreateConnection() {
	dsn := ""

	_, exists := os.LookupEnv("POSTGRES_HOST")
	if exists {
		e := func(v string) string {
			r, _ := os.LookupEnv(v)
			return r
		}
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			e("POSTGRES_HOST"), e("POSTGRES_USER"), e("POSTGRES_PASSWORD"), e("POSTGRES_DB"), e("POSTGRES_PORT"),
		)
		fmt.Println("MY DSN ENV:", dsn)
	} else {
		c := &config.Config.DatabaseConfig
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			c.Host, c.User, c.Password, c.Dbname, c.Port,
		)
		fmt.Println("MY DSN ENV:", dsn)
	}

	if dsn == "" {
		panic("Database configuration not created!!!")
	}

	Database, _ = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			//Logger:                                   logger.Default.LogMode(logger.Silent),
		})

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

func Save(v any) {
	Database.Save(v)
}
