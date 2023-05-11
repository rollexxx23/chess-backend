package utils

import (
	"encoding/json"
	"net/http"

	"github.com/rollexxx23/chess/database"
	"github.com/rollexxx23/chess/models"
)

type requestForAddmatches struct {
	WhitePlayerUserName string `json:"white_username"`
	BlackPlayerUserName string `json:"black_username"`
	WhitePlayer         string `json:"white_player"`
	BlackPlayer         string `json:"black_player"`
	GameMoves           string `json:"game_moves"`
	Moves               int    `json:"moves_cnt"`
	Result              int    `json:"result"`
	Comment             string `json:"comment"`
}

func AddMatch(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var match requestForAddmatches

	json.NewDecoder(r.Body).Decode(&match)

	game := models.Match{
		WhitePlayerUserName: match.WhitePlayerUserName,
		BlackPlayerUserName: match.BlackPlayerUserName,
		WhitePlayer:         match.WhitePlayer,
		BlackPlayer:         match.BlackPlayer,
		GameMoves:           match.GameMoves,
		Moves:               match.Moves,
		Result:              uint8(match.Result),
		Comment:             match.Comment,
	}

	database.Instance.Create(&game)

}
