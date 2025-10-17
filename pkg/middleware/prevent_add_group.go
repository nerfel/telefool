package middleware

import (
	"telefool/configs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PreventAddGroup(next Handler) Handler {
	return func(update tgbotapi.Update, config *configs.Config, bot *tgbotapi.BotAPI) {
		if update.MyChatMember == nil {
			next(update, config, bot)
			return
		}

		old := update.MyChatMember.OldChatMember
		newm := update.MyChatMember.NewChatMember

		if newm.User == nil || newm.User.ID != bot.Self.ID {
			next(update, config, bot)
			return
		}

		if old.Status == "left" && newm.Status == "member" {
			if update.MyChatMember.From.UserName != config.AdminUserName {
				bot.Request(tgbotapi.LeaveChatConfig{ChatID: update.MyChatMember.Chat.ID})
			}

			return
		}

		if newm.Status == "administrator" {
			return
		}

		if newm.Status == "left" {
			return
		}

		next(update, config, bot)
	}
}
