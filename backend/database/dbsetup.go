package database

import (
	"log"
	"os"

	"github.com/TMP-The-Major-Project/Thrift-Store/backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func check(eText string, err error) {
	if err != nil {
		log.Fatalf("%v %v", eText, err)
	}
}

func Connect() *gorm.DB {
	err := godotenv.Load(".env")
	check("Could not load the .env file", err)

	db, err2 := gorm.Open(postgres.Open(os.Getenv("POSTGRES_URL")), &gorm.Config{})
	check("Could not connect to database!!", err2)

	db.AutoMigrate(&models.User{})

	return db
}
