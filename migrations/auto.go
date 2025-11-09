package main

import (
	"telefool/configs"
	"telefool/internal/dialog"
	"telefool/internal/message"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	config := configs.LoadConfig()

	db, err := gorm.Open(postgres.Open(config.Db.GetDsn()), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&message.Message{}, &dialog.Dialog{})
}
