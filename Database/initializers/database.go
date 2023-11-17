package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	dbConfig := postgres.Config{
		DSN:                  "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=5432",
		PreferSimpleProtocol: true, 
	}

	DB, err = gorm.Open(postgres.New(dbConfig), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
}
