package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramId  int64
	UserName    string
	FirstName   string
	LastName    string
	InWhitelist bool
}
