package dialog

import (
	"telefool/pkg/db"

	"gorm.io/gorm/clause"
)

type DialogRepository struct {
	*db.Db
}

func NewDialogRepository(db *db.Db) *DialogRepository {
	return &DialogRepository{db}
}

func (repo *DialogRepository) AddGroup(ChatId int64, ChatTitle string) {
	repo.Db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "chat_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"deleted_at": nil,
			"chat_title": ChatTitle,
		}),
	}).Create(&Dialog{
		ChatId:    ChatId,
		ChatTitle: ChatTitle,
	})
}

func (repo *DialogRepository) RemoveFromGroup(ChatId int64) {
	repo.Db.Where(Dialog{ChatId: ChatId}).Delete(&Dialog{})
}

func (repo *DialogRepository) GetEnabledDialog(ChatId int64) (*Dialog, error) {
	var dialog Dialog
	result := repo.Db.
		Where(Dialog{ChatId: ChatId}).
		Where(Dialog{IsEnabled: true}).
		First(&dialog)

	if result.Error != nil {
		return nil, result.Error
	}

	return &dialog, nil
}
