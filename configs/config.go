package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db            DbConfig
	BotToken      string
	HttpPort      string
	BotWebhookUrl string
}

type DbConfig struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     string
}

func (conf *DbConfig) GetDsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		conf.Host,
		conf.User,
		conf.Password,
		conf.DbName,
		conf.Port,
	)
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Db: DbConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			DbName:   os.Getenv("DB_DATABASE"),
			Port:     os.Getenv("DB_PORT"),
		},
		BotToken:      os.Getenv("BOT_TOKEN"),
		HttpPort:      os.Getenv("HTTP_PORT"),
		BotWebhookUrl: os.Getenv("BOT_WEBHOOK_URL"),
	}
}
