package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/timebetov/readerblog/internals/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Variable for database
var DB *gorm.DB

func ConnectDB() {
	var err error
	p := os.Getenv("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		log.Println("Error Parsing Database Port")
	}

	// Getting Database env values
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Checking if all env variables are set
	if host == "" || user == "" || password == "" || dbName == "" {
		log.Fatalf("Missing required environment variables for database connection")
	}

	// Connection URL to connect to Postgres Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbName, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failde to connect database")
	}

	fmt.Println("Connection Opened to Database")

	// Migrate the schema
	if err = DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}
	fmt.Println("Schema was successfully migrated to database!")
}
