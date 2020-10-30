package models

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// main gorm instance
var DB *gorm.DB

// InitDBConnection – Initialize database connection
func InitDBConnection(connDSN string) (err error) {
	// todo: replace db connections settings to env vars / flags
	// dsn := "host=0.0.0.0 user=postgres password=postgres dbname=hydropony port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(connDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if tableExists := DB.Migrator().HasTable(&User{}); tableExists == false {
		if err := DB.Migrator().CreateTable(&User{}); err != nil {
			log.Fatal("Cannot create table", err)
		}
	}

	return err
}
