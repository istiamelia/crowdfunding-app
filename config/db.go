package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// global variable for db connection
var DB *gorm.DB

func InitDB(){
	// Getting data from the environment variable
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	if host == "" || user == "" || password == "" || dbName == "" || port == "" {
		log.Fatal("One or more required environment variables are not set")
	}
	// Save it to variable dsn
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta", host, user, password, dbName, port)

	// Connect to PostgreSQL Database
	var err error
	// Open initialize db section
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Fatalf("Failed to connect database %v", err)
	}

	log.Println("Connected to database")
}