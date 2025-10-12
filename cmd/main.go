package main

import (
	"fmt"
	"log"
	"net/http"
	"telefool/configs"
	"telefool/internal/message"
	"telefool/pkg/db"
	"telefool/pkg/middleware"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func initHttpServer(port string) {
	httpServer := &http.Server{
		Addr: ":" + port,
	}

	fmt.Println("Listening on port " + port)
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}

func initBot(conf *configs.Config) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	bot.Request(&tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true})
	wh, _ := tgbotapi.NewWebhook(conf.BotWebhookUrl)

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return bot
}

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	fmt.Println("dbPointer:", database)

	go initHttpServer(conf.HttpPort)

	bot := initBot(conf)

	updates := bot.ListenForWebhook("/")

	stack := middleware.Chain(middleware.Logging)
	handler := stack(message.Handle)

	for update := range updates {
		handler(update)
	}
}
