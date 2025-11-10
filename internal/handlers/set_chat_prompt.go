package handlers

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"telefool/pkg/di"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SetChatPromptData struct {
	ChatId int64
	Prompt string
}

const setChatPromptRegexpString = `^set-chat-prompt\s+(-?\d+)\s+(.+)$`

func SetChatPromptRoute(ctx *di.UpdateContext) bool {
	setChatPromptRegexp := regexp.MustCompile(setChatPromptRegexpString)
	textMessage := strings.TrimSpace(ctx.Update.Message.Text)

	if !setChatPromptRegexp.MatchString(textMessage) {
		return false
	}

	matches := setChatPromptRegexp.FindStringSubmatch(textMessage)

	if len(matches) < 3 {
		return false
	}

	num, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		log.Println("Failed to parse chat prompt number")
		return false
	}

	ctx.RoutePayload = &SetChatPromptData{
		ChatId: num,
		Prompt: matches[2],
	}

	return true
}

func SetChatPromptHandler(ctx *di.UpdateContext, container *di.Container) {
	chatData := ctx.RoutePayload.(*SetChatPromptData)

	err := container.DialogService.SetChatPrompt(chatData.ChatId, chatData.Prompt)
	if err == nil {
		ctx.Bot.Send(tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Новый промпт установлен!"))
	} else {
		ctx.Bot.Send(tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "Не удалось установить промпт"))
	}
}
