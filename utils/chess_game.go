package utils

import (
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
	Status              uint8      `json:"status"`       // 0 -> white to move, 1 -> black to move, 2 -> white won, 3 -> black won, 4 -> draw
	Result              uint8      `json:"result"`       // 0 -> black, 1 -> white, 2 -> draw
	PendingDraw         bool       `json:"pendind_draw"` //draw offer cnt
	validator           *chess.Game

	WhiteTurn bool `json:"white_move"` // white -> true
}

type GameMove struct {
	Type string
	ID   int
	Src  string
	Dest string
	Prom string
	Fen  string // current FEN
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
	game.Status = 0
	game.gameMoves = nil
	game.validator = chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	game.WhiteTurn = true
	return game.ID, game
}
