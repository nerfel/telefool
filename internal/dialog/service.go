package dialog

import (
	"errors"
	"telefool/pkg/event"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DialogServiceDeps struct {
	EventBus         *event.Bus
	DialogRepository *DialogRepository
}

type DialogService struct {
	EventBus         *event.Bus
	DialogRepository *DialogRepository
}

func NewDialogService(deps *DialogServiceDeps) *DialogService {
	return &DialogService{
		EventBus:         deps.EventBus,
		DialogRepository: deps.DialogRepository,
	}
}

func (s *DialogService) IsExistingDialogEnabled(ChatId int64) bool {
	_, err := s.DialogRepository.GetEnabledDialog(ChatId)
	if err != nil {
		return false
	}
	return true
}

func (s *DialogService) GroupEventsListen() {
	for msg := range s.EventBus.Subscribe() {
		if msg.Type == event.EventAddToGroup {
			update := msg.Data.(tgbotapi.Update)
			s.DialogRepository.AddGroup(update.MyChatMember.Chat.ID, update.MyChatMember.Chat.Title)
		} else if msg.Type == event.EventRemoveFromGroup {
			update := msg.Data.(tgbotapi.Update)
			s.DialogRepository.RemoveFromGroup(update.MyChatMember.Chat.ID)
		}
	}
}

func (s *DialogService) GetChatPrompt(ChatId int64) (string, error) {
	dialog, err := s.DialogRepository.GetEnabledDialog(ChatId)
	if err != nil {
		return "", err
	}

	if dialog.ChatPrompt == "" {
		return "", errors.New("empty prompt for dialog")
	}

	return dialog.ChatPrompt, nil
}

func (s *DialogService) SetChatPrompt(ChatId int64, prompt string) error {
	dialog, err := s.DialogRepository.GetEnabledDialog(ChatId)
	if err != nil {
		return err
	}

	dialog.ChatPrompt = prompt

	dialog, err = s.DialogRepository.Update(dialog)
	if err != nil {
		return err
	}

	return nil
}
