package reply

import (
	"log"
	"math/rand"
	"time"
)

var (
	lastReply = map[int64]time.Time{}
)

func ShouldReply(chatID int64, probability float64, cooldown time.Duration) bool {
	if probability == 1.0 {
		log.Println("Handled probability 100%")
		return true
	}

	if t, ok := lastReply[chatID]; ok && time.Since(t) < cooldown {
		log.Println("Skipped by cooldown")
		return false
	}

	if rand.Float64() > probability {
		log.Println("Skipped by probability")
		return false
	}

	lastReply[chatID] = time.Now()
	log.Println("Handling ...")
	return true
}
