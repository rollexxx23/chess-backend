package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/notnil/chess"
	"github.com/rollexxx23/chess/database"
	"github.com/rollexxx23/chess/models"
)

type ChessGame struct {
	ID                  int    `json:"id"`
	WhitePlayerUserName string `json:"white_username"` // username
	BlackPlayerUserName string `json:"black_username"`
	whitePlayer         string // email
	blackPlayer         string
	gameMoves           []GameMove //list of moves

	Result      uint8 `json:"result"`       // 0 -> black, 1 -> white, 2 -> draw
	PendingDraw bool  `json:"pendind_draw"` //draw offer cnt
	validator   *chess.Game

	WhiteTurn bool `json:"white_move"` // white -> true
}

type GameMove struct {
	Email string
	ID    int
	Src   string
	Dest  string
	Prom  string
	Fen   string // current FEN
}

type sndMoveStruct struct {
	GameId int    `json:"game_id"`
	Src    string `json:"src"`
	Dest   string `json:"des"`
	Prom   string `json:"prom"`
}

func initGame(player1, player2 string) (int, ChessGame) {
	var users []models.User
	var player1UserName, player2UserName string
	// get username of players
	database.Instance.Where("email = ?", player1).Find(&users)
	player1UserName = users[0].Username
	database.Instance.Where("email = ?", player2).Find(&users)
	player2UserName = users[0].Username
	// chess game struct
	var game ChessGame
	game.ID = rand.Intn(999999)
	if rand.Intn(2) == 0 {
		game.WhitePlayerUserName = player1UserName
		game.whitePlayer = player1
		game.BlackPlayerUserName = player2UserName
		game.blackPlayer = player2
	} else {
		game.BlackPlayerUserName = player1UserName
		game.blackPlayer = player1
		game.WhitePlayerUserName = player2UserName
		game.whitePlayer = player2
	}
	// init struct

	game.gameMoves = nil
	game.validator = chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	game.WhiteTurn = true
	return game.ID, game
}

func playMove(game *ChessGame, move GameMove) error {

	if game.WhiteTurn {
		if move.Email != game.whitePlayer {
			return errors.New("Not Your Turn")
		}
	} else {
		if move.Email != game.blackPlayer {
			return errors.New("Not Your Turn")
		}
	}

	err := game.validator.MoveStr(move.Src + move.Dest + move.Prom)

	if err != nil {
		fmt.Println("Could not decode", err, move.Src+move.Dest+move.Prom)
		return err
	}

	// send move to opponent
	moveSnd := sndMoveStruct{
		GameId: game.ID,
		Src:    move.Src,
		Dest:   move.Dest,
		Prom:   move.Prom,
	}

	text, _ := json.Marshal(moveSnd)
	if game.WhiteTurn {
		fmt.Println(game.blackPlayer)
		Directory.EmailToSocketMap[game.blackPlayer].WriteMessage(1, []byte(text))
	} else {
		Directory.EmailToSocketMap[game.whitePlayer].WriteMessage(1, []byte(text))
	}
	if game.validator.Outcome() != chess.NoOutcome {
		// game end logic
	}

	game.WhiteTurn = !game.WhiteTurn
	game.gameMoves = append(game.gameMoves, move)

	return nil
}

/*

{"token": "1234","email": "arin2@gmail.com","game_id": 499379,"src": "d2","des": "d4","prom": ""}


*/
