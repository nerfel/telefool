package dialog

import "gorm.io/gorm"

type Dialog struct {
	gorm.Model
	ChatId     int64 `gorm:"uniqueIndex"`
	ChatTitle  string
	ChatPrompt string
	IsEnabled  bool `gorm:"default:false"`
}
