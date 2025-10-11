package message

import "telefool/pkg/db"

type MessageRepository struct {
	Database *db.Db
}

func NewMessageRepository(database *db.Db) *MessageRepository {
	return &MessageRepository{database}
}

func (repo *MessageRepository) Create(msg *Message) (*Message, error) {
	res := repo.Database.Create(msg)
	if res.Error != nil {
		return nil, res.Error
	}

	return msg, nil
}
