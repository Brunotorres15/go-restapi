package database

import (
	"log"

	"github.com/Brunotorres15/go-restapi/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
	}

	log.Println("Connected to the dabase successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")
	// Add migrations
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
	Database = DbInstance{Db: db}
}
