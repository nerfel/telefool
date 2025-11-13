package dialog

import (
	"time"

	"gorm.io/gorm"
)

type Dialog struct {
	gorm.Model
	ChatId      int64 `gorm:"uniqueIndex"`
	ChatTitle   string
	ChatPrompt  string
	Probability float64       `gorm:"default:0.3"`
	Cooldown    time.Duration `gorm:"default:2000000000"`
	IsEnabled   bool          `gorm:"default:false"`
}
