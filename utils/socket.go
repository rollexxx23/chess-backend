package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var MatchReqs = struct {
	sync.RWMutex
	Lobby map[string]*websocket.Conn // currently searching for matches
}{Lobby: make(map[string]*websocket.Conn)}

var ActiveMatches = struct {
	sync.RWMutex
	Match map[int]*ChessGame
}{Match: make(map[int]*ChessGame)}

var Directory = struct {
	EmailToSocketMap map[string]*websocket.Conn
}{EmailToSocketMap: make(map[string]*websocket.Conn)}

type findMatchMsgStruct struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type movesMsgStruct struct {
	Token  string `json:"token"`
	Email  string `json:"email"`
	GameId int    `json:"game_id"`
	Src    string `json:"src"`
	Dest   string `json:"des"`
	Prom   string `json:"prom"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func findMatch(conn *websocket.Conn) error {

	// read in a message
	_, p, err := conn.ReadMessage()

	if err != nil {
		return err
	}

	var req findMatchMsgStruct
	err = json.Unmarshal(p, &req)
	fmt.Println(req)
	if err != nil {
		return err
	}

	MatchReqs.Lobby[req.Email] = conn
	Directory.EmailToSocketMap[req.Email] = conn

	var found bool
	found = false

	for !found {
		_, found = MatchReqs.Lobby[req.Email]
		found = !found
		text := LobbyConnect(conn, req)
		if text != nil {
			break
		}
		time.Sleep(time.Second)

	}

	return nil

}

func getMoves(conn *websocket.Conn) error {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var req movesMsgStruct
		err = json.Unmarshal(p, &req)
		if err != nil {
			return err
		}

		move := GameMove{
			ID:    req.GameId,
			Src:   req.Src,
			Email: req.Email,
			Dest:  req.Dest,
			Prom:  req.Prom,
		}

		game := ActiveMatches.Match[req.GameId]

		err = playMove(game, move)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func LobbyEndpoint(w http.ResponseWriter, r *http.Request) {

	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	defer ws.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Coonected...")
	err = findMatch(ws)
	if err != nil {
		fmt.Println(err)
	}

	err = getMoves(ws)
	if err != nil {
		fmt.Println(err)
	}
}

func LobbyConnect(conn *websocket.Conn, m findMatchMsgStruct) []byte {

	opponent := ""
	for email, conn2 := range MatchReqs.Lobby {
		if email != m.Email {
			opponent = email
			delete(MatchReqs.Lobby, opponent)
			delete(MatchReqs.Lobby, m.Email)

			id, game := initGame(m.Email, opponent)
			ActiveMatches.Match[id] = &game

			text, _ := json.Marshal(game)

			_ = conn.WriteMessage(1, []byte(text))
			_ = conn2.WriteMessage(1, []byte(text))

		}
	}

	if opponent != "" {
		return []byte("found")
	}
	return nil

}

/*
{"message_type":"find_match", "token":"abc", "email": "arin2@gmail.com"}
*/
