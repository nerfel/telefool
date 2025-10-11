package message

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ChatId        int
	ChatType      string
	Text          string
	FromId        int
	FromFirstName string
	FromLastName  string
	FromUsername  string
}
