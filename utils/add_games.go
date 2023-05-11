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

	blackPlayer := match.BlackPlayer
	whitePlayer := match.WhitePlayer

	if match.Result == 0 {

		var users []models.User
		// find black and add 1 to win counter
		database.Instance.Where("email = ?", blackPlayer).Find(&users)
		if len(users) != 0 {
			users[0].WinCnt = users[0].WinCnt + 1
		}
		database.Instance.Save(users)
		// find white and add 1 to loss counter

		database.Instance.Where("email = ?", whitePlayer).Find(&users)
		if len(users) != 0 {
			users[0].LossCnt = users[0].LossCnt + 1
		}
		database.Instance.Save(users)
	} else if match.Result == 1 {

		var users []models.User
		// find black and add 1 to win counter
		database.Instance.Where("email = ?", whitePlayer).Find(&users)
		if len(users) != 0 {
			users[0].WinCnt = users[0].WinCnt + 1
		}
		database.Instance.Save(users)
		// find white and add 1 to loss counter

		database.Instance.Where("email = ?", blackPlayer).Find(&users)
		if len(users) != 0 {
			users[0].LossCnt = users[0].LossCnt + 1
		}
		database.Instance.Save(users)
	} else {

		var users []models.User

		database.Instance.Where("email = ?", blackPlayer).Find(&users)
		if len(users) != 0 {
			users[0].DrawCnt = users[0].DrawCnt + 1
		}
		database.Instance.Save(users)

		database.Instance.Where("email = ?", whitePlayer).Find(&users)
		if len(users) != 0 {
			users[0].DrawCnt = users[0].DrawCnt + 1
		}
		database.Instance.Save(users)
	}

}
