package match

import "time"

type MatchResult struct {
	WinnerID   int         `json:"winner_id"`
	LoserID    int         `json:"loser_id"`
	MatchID    int         `json:"match_id"`
	MatchDate  time.Time   `json:"match_date"`
	MatchType  string      `json:"match_type"`
	Experience *Experience `json:"experience,omitempty"`
}

type Experience struct {
	WinnerGain int `json:"winner_gain,omitempty"`
	LoserGain  int `json:"loser_gain,omitempty"`
}

type Item struct {
	PlayerID int `json:"player_id"`
	ItemID   int `json:"item_id"`
}
