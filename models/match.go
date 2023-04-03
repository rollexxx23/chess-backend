package models

import (
	"gorm.io/gorm"
)

type Match struct {
	gorm.Model
	WhitePlayerUserName string `json:"white_username"`
	BlackPlayerUserName string `json:"black_username"`
	WhitePlayer         string `json:"white_name"`
	BlackPlayer         string `json:"black_name"`
	GameMoves           string `json:"game_moves"`
	Moves               int    `json:"moves"`
	Result              uint8  `json:"result"` // 0 -> black, 1 -> white, 2 -> draw
	Comment             string `json:"comment"`
}
