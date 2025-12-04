package config

import (
	"log"
	"github.com/AllanKariuki/fiber-go-rest-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(AppConfig.DBUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the models
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to auto-migrate models:", err)
	}

	log.Println("Database connected and migrated!")
}

func GetDB() *gorm.DB {
	return DB
}