package dbConnect

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "log"
    "os"
    "fmt"
    
)

func InitDB() (*gorm.DB, error) {
    var db *gorm.DB
    host := os.Getenv("DB_HOST")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    databaseName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", host, user, password, databaseName)
    
    // DB , err = gorm.Open(postgres.New(postgres.Config{
	// 	DSN: os.Getenv("DB_URL"),
	// 	PreferSimpleProtocol: true, // disables implicit prepared statement usage
	//   }), &gorm.Config{})

    db, err := gorm.Open(postgres.New(postgres.Config{ DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
        return nil, err
    }

    log.Println("Successfully connected to the database.")
    return db, nil
}
