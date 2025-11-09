package handlers

import (
	"errors"
	"log"
	"strings"
	"telefool/internal/dialog"
	"telefool/internal/reply"
	"telefool/pkg/di"
	"telefool/pkg/gpt"
	"telefool/pkg/memory"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func FallBackGPTHandle(ctx *di.UpdateContext, dialogService *dialog.DialogService) {
	if ctx.Update.Message.Chat.Type == "private" {
		return
	}
	if !dialogService.IsExistingDialogEnabled(ctx.Update.Message.Chat.ID) {
		return
	}

	if ctx.Config.YandexCloudConfig.IamToken == "" {
		err := gpt.GetYandexGPTOauthToken(ctx.Config)
		if err != nil {
			log.Println("Get YandexGPTOauthToken error ", err)
			return
		}
	}

	// store in memory last messages
	ctx.Memory.Add(memory.Message{
		UserID:         ctx.Update.Message.From.ID,
		FromCurrentBot: ctx.Update.Message.From.ID == ctx.Bot.Self.ID,
		ChatID:         ctx.Update.Message.Chat.ID,
		Text:           ctx.Update.Message.Text,
		Timestamp:      time.Now(),
	})

	var probability float64 = 0.8
	cooldown := 2 * time.Second
	if isMentioned(ctx.Update.Message.Text, ctx.Bot.Self.UserName) {
		probability = 1.0
	}

	if !reply.ShouldReply(ctx.Update.Message.Chat.ID, probability, cooldown) {
		return
	}

	history := ctx.Memory.ChatHistory(ctx.Update.Message.Chat.ID)
	chatPrompt, err := dialogService.GetChatPrompt(ctx.Update.Message.Chat.ID)
	if err != nil {
		log.Println("GetChatPrompt error", err)
		return
	}

	gptRequestPayload, err := gpt.BuildModelRequestPayload(history, chatPrompt, ctx.Config)

	if err != nil {
		log.Println("Gpt request payload error ", err)
		return
	}
	result, err := gpt.RequestModel(gptRequestPayload, ctx.Config)
	if errors.Is(err, gpt.ErrUnauthorized) {
		err = gpt.GetYandexGPTOauthToken(ctx.Config)
		if err != nil {
			log.Println("Get YandexGPTOauthToken error ", err)
			return
		}

		result, err = gpt.RequestModel(gptRequestPayload, ctx.Config)
		if err != nil {
			log.Println("Gpt request with new token failed")
			return
		}

		err = nil
	}
	if err != nil {
		log.Println("Request to model error", err)
		return
	}

	gptAnswer := result.Result.Alternatives[0].Message.Text
	ctx.Memory.Add(memory.Message{
		UserID:         ctx.Bot.Self.ID,
		FromCurrentBot: true,
		ChatID:         ctx.Update.Message.Chat.ID,
		Text:           gptAnswer,
		Timestamp:      time.Now(),
	})

	ctx.Bot.Send(tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, gptAnswer))
}

func isMentioned(textMessage string, botUserName string) bool {
	textMessage = strings.ToLower(textMessage)
	botName := strings.ToLower(botUserName)

	return strings.Contains(textMessage, botName)
}
