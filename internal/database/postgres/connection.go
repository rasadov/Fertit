package postgres

import (
	"github.com/rasadov/MailManagerApp/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.AppConfig.PostgresUrl))

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	err = sqlDB.Ping()

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	return db
}
