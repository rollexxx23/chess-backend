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
	ID                  int        `json:"id"`
	WhitePlayerUserName string     `json:"white_username"` // username
	BlackPlayerUserName string     `json:"black_username"`
	WhitePlayer         string     `json:"white_email"`
	BlackPlayer         string     `json:"black_email"`
	gameMoves           []GameMove //list of moves

	Result      uint8 `json:"result"`       // 0 -> black, 1 -> white, 2 -> draw
	PendingDraw bool  `json:"pendind_draw"` //draw offer cnt
	validator   *chess.Game

	WhiteTurn   bool   `json:"white_move"` // white -> true
	MessageType string `json:"message_type"`
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
	GameId      int    `json:"game_id"`
	Src         string `json:"src"`
	Dest        string `json:"des"`
	Prom        string `json:"prom"`
	MessageType string `json:"message_type"`
	Fen         string `json:"fen"`
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
	game.MessageType = "GAME_INIT"
	game.ID = rand.Intn(999999)
	if rand.Intn(2) == 0 {
		game.WhitePlayerUserName = player1UserName
		game.WhitePlayer = player1
		game.BlackPlayerUserName = player2UserName
		game.BlackPlayer = player2
	} else {
		game.BlackPlayerUserName = player1UserName
		game.BlackPlayer = player1
		game.WhitePlayerUserName = player2UserName
		game.WhitePlayer = player2
	}
	// init struct

	game.gameMoves = nil
	game.validator = chess.NewGame(chess.UseNotation(chess.UCINotation{}))
	game.WhiteTurn = true
	return game.ID, game
}

func playMove(game *ChessGame, move GameMove) error {

	if game.WhiteTurn {
		if move.Email != game.WhitePlayer {
			return errors.New("Not Your Turn")
		}
	} else {
		if move.Email != game.BlackPlayer {
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
		GameId:      game.ID,
		Src:         move.Src,
		Dest:        move.Dest,
		Prom:        move.Prom,
		MessageType: "GAME_MOVE",
		Fen:         game.validator.FEN(),
	}

	text, _ := json.Marshal(moveSnd)
	if game.WhiteTurn {
		Directory.EmailToSocketMap[game.BlackPlayer].WriteMessage(1, []byte(text))
	} else {
		Directory.EmailToSocketMap[game.WhitePlayer].WriteMessage(1, []byte(text))
	}
	if game.validator.Outcome() != chess.NoOutcome {
		// game end logic
		gameEndLogic(game.ID)
		Directory.EmailToSocketMap[game.BlackPlayer].WriteMessage(1, []byte("game over"))
		Directory.EmailToSocketMap[game.WhitePlayer].WriteMessage(1, []byte("game over"))
	}

	game.WhiteTurn = !game.WhiteTurn
	game.gameMoves = append(game.gameMoves, move)

	return nil
}

func gameEndLogic(id int) {
	fmt.Println("called...")
	chessgame := ActiveMatches.Match[id]
	//now store game in MySQL database
	allMoves, err := json.Marshal(chessgame.gameMoves)
	if err != nil {
		fmt.Println("Error marshalling data to store in MySQL")
	}
	updateCnt(chessgame)
	//gets length of all the moves in the game
	totalMoves := (len(chessgame.gameMoves) + 1) / 2
	storeGame(totalMoves, string(allMoves), chessgame)

	//now delete game from memory
	delete(ActiveMatches.Match, id)
}

func updateCnt(chessGame *ChessGame) {
	if chessGame.validator.Outcome() == chess.BlackWon {
		chessGame.Result = 0
		var users []models.User
		// find black and add 1 to win counter
		database.Instance.Where("email = ?", chessGame.BlackPlayer).Find(&users)
		if len(users) != 0 {
			users[0].WinCnt = users[0].WinCnt + 1
		}
		database.Instance.Save(users)
		// find white and add 1 to loss counter

		database.Instance.Where("email = ?", chessGame.WhitePlayer).Find(&users)
		if len(users) != 0 {
			users[0].LossCnt = users[0].LossCnt + 1
		}
		database.Instance.Save(users)
	} else if chessGame.validator.Outcome() == chess.WhiteWon {
		chessGame.Result = 1
		var users []models.User
		// find black and add 1 to win counter
		database.Instance.Where("email = ?", chessGame.WhitePlayer).Find(&users)
		if len(users) != 0 {
			users[0].WinCnt = users[0].WinCnt + 1
		}
		database.Instance.Save(users)
		// find white and add 1 to loss counter

		database.Instance.Where("email = ?", chessGame.BlackPlayer).Find(&users)
		if len(users) != 0 {
			users[0].LossCnt = users[0].LossCnt + 1
		}
		database.Instance.Save(users)
	} else {
		chessGame.Result = 2
		var users []models.User

		database.Instance.Where("email = ?", chessGame.BlackPlayer).Find(&users)
		if len(users) != 0 {
			users[0].DrawCnt = users[0].DrawCnt + 1
		}
		database.Instance.Save(users)

		database.Instance.Where("email = ?", chessGame.WhitePlayer).Find(&users)
		if len(users) != 0 {
			users[0].DrawCnt = users[0].DrawCnt + 1
		}
		database.Instance.Save(users)
	}
}

func storeGame(moveCnt int, moves string, chessgame *ChessGame) error {

	game := models.Match{
		WhitePlayerUserName: chessgame.WhitePlayerUserName,
		BlackPlayerUserName: chessgame.BlackPlayerUserName,
		WhitePlayer:         chessgame.WhitePlayer,
		BlackPlayer:         chessgame.BlackPlayer,
		GameMoves:           moves,
		Moves:               moveCnt,
		Result:              chessgame.Result,
		Comment:             chessgame.validator.Outcome().String(),
	}
	database.Instance.Create(&game)
	return nil
}

/*

{"token": "1234","email": "arin2@gmail.com","game_id": 499379,"src": "d2","des": "d4","prom": ""}


*/
