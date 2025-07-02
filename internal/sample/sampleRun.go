package sample

import (
	"lesta-battleship/server-core/internal/match"
	"log"
	"time"
)

func RunSampleMatch(handler *match.MatchHandler) {

	// какая то логика игры, что то происходит

	// и чел юзает какой то предмет
	if err := handler.ItemUsed(match.InventoryEvent{
		PlayerID: 1,
		ItemID:   42,
	}); err != nil {
		log.Printf("Error processing item used: %v", err)
	}

	// в конце игры/после игры
	if err := handler.MatchEnded(match.MatchResult{
		WinnerID:           1,
		LoserID:            2,
		MatchID:            1001,
		MatchDurationInSec: 300,
		MatchDate:          time.Now(),
		MatchType:          "ranked",
		Experience: &match.Experience{
			WinnerGain: 30,
			LoserGain:  -15,
		},
	}); err != nil {
		log.Printf("Error processing match ended: %v", err)
	}

	time.Sleep(1 * time.Second)
}
