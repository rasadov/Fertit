package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func GetPostgresDB(PostgresUrl string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(PostgresUrl))

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	err = sqlDB.Ping()

	if err != nil {
		log.Fatal("Failed to ping:", err)
	}

	return db
}
