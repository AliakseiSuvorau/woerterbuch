package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"woerterbuch-backend/src/global"
	"woerterbuch-backend/src/model/dtos"
)

func Init() {
	var (
		user   = os.Getenv("POSTGRES_USER")
		pass   = os.Getenv("POSTGRES_PASSWORD")
		host   = os.Getenv("POSTGRES_HOST")
		port   = os.Getenv("POSTGRES_PORT")
		dbname = os.Getenv("POSTGRES_DB_NAME")
	)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, dbname)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err.Error()))
	}

	global.DB = db
	if err := global.DB.AutoMigrate(&dtos.Word{}); err != nil {
		panic(fmt.Sprintf("failed to auto migrate database: %v", err.Error()))
	}
}
