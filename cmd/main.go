package main

import (
	"fmt"
	"log"
	"net/http"
	"telefool/configs"
	"telefool/internal/dialog"
	"telefool/internal/handlers"
	"telefool/internal/user"
	"telefool/pkg/db"
	"telefool/pkg/di"
	"telefool/pkg/event"
	"telefool/pkg/memory"
	"telefool/pkg/router"

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
	eventBus := event.NewEventBus()

	go initHttpServer(conf.HttpPort)

	bot := initBot(conf)

	// Repositories
	userRepository := user.NewUserRepository(database)
	dialogRepository := dialog.NewDialogRepository(database)

	// Services
	userService := user.NewUserService(userRepository)
	dialogService := dialog.NewDialogService(&dialog.DialogServiceDeps{
		EventBus:         eventBus,
		DialogRepository: dialogRepository,
	})

	go dialogService.GroupEventsListen()

	// GlobalHandler
	gmh := handlers.NewUpdateHandler(&handlers.UpdateHandlerDeps{
		Config:   conf,
		EventBus: eventBus,
		Bot:      bot,
		Router:   router.NewUpdateRouter(),
		Container: &di.Container{
			DialogService: dialogService,
			UserService:   userService,
		},
		Memory: memory.NewShortTermMemory(100),
	})

	gmh.Handle()
}
