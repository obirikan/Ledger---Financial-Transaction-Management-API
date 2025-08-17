package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DbConfig struct {
	DBHost     string 
	DBPort     string 
	DBUser     string 
	DBPassword string 
	DBName     string 
	ServerPort string 
}

func SetupDB() *sql.DB {

 err := godotenv.Load("./.env")

	if err != nil {
		log.Println("Couldn't find the secret note, using defaults")
	}

	 data := &DbConfig{
		DBHost:     getEnv("DB_HOST", "default"),  
		DBPort:     getEnv("DB_PORT", "default"),     
		DBUser:     getEnv("DB_USER", "default"),  
		DBPassword: getEnv("DB_PASSWORD", "default"), 
		DBName:     getEnv("DB_NAME", "default"),     
		ServerPort: getEnv("SERVER_PORT", "default"),
	}
	
	dsn := fmt.Sprintf("host=localhost port=%v user=%v password=%v dbname=%v sslmode=disable", data.DBPort, data.DBUser, data.DBPassword, data.DBName)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	fmt.Println("Connected to PostgreSQL!")
	return db
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)  
	if value == "" {        
		return defaultValue 
	}
	return value
}