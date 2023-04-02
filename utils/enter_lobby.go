package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var MatchReqs = struct {
	sync.RWMutex
	Lobby map[string]*websocket.Conn // email -> socket
}{Lobby: make(map[string]*websocket.Conn)}

type MessageStruct struct {
	Token       string `json:"token"`
	Email       string `json:"email"`
	MessageType string `json:"message_type"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) error {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return err
		}
		var req MessageStruct
		err = json.Unmarshal(p, &req)
		if err != nil {
			return err
		}

		MatchReqs.Lobby[req.Email] = conn

		LobbyConnect(conn, req)

	}

}

func LobbyEndpoint(w http.ResponseWriter, r *http.Request) {

	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = reader(ws)
	if err != nil {
		fmt.Println(err)
	}
}

func LobbyConnect(conn *websocket.Conn, m MessageStruct) {
	switch m.MessageType {
	case "fetch_matches":
		//send in array instead of sending individual
		opponent := ""
		for email, conn2 := range MatchReqs.Lobby {
			if email != m.Email {
				_ = conn2.WriteMessage(1, []byte("Your ooponent is "+m.Email))
				_ = conn.WriteMessage(1, []byte("Your ooponent is "+email))
			}
		}

		if opponent == "" {
			_ = conn.WriteMessage(1, []byte("Please wait"))

		} else {
			delete(MatchReqs.Lobby, opponent)
			delete(MatchReqs.Lobby, m.Email)
		}
	}
}
