package db

import (
	"atur-dana/internal/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("DATABASE_URL")
	logLevel := os.Getenv("LogLevel")
	dbLog := logger.Silent

	switch logLevel {
	case "INFO":
		dbLog = logger.Info
	case "WARN":
		dbLog = logger.Warn
	case "ERROR":
		dbLog = logger.Error
	default:
		dbLog = logger.Silent
	}

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  dbLog,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         dbLogger,
		TranslateError: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Transaction{}, &models.Budget{})
}
