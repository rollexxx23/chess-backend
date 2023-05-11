package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rollexxx23/chess/database"
	"github.com/rollexxx23/chess/models"
)

func GetMatches(w http.ResponseWriter, r *http.Request) {
	var matches []models.Match
	database.Instance.Find(&matches)

	json.NewEncoder(w).Encode(matches)
}

func GetMatchesByUserName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	var matches []models.Match
	database.Instance.Where("white_player_user_name = ? OR black_player_user_name = ?", username, username).Find(&matches)

	json.NewEncoder(w).Encode(matches)
}

func GetMatchesByEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]
	var matches []models.Match
	database.Instance.Where("white_player = ? OR black_player = ?", email, email).Find(&matches)

	json.NewEncoder(w).Encode(matches)
}

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]
	var users []models.User
	database.Instance.Where("email = ?", email).Find(&users)
	json.NewEncoder(w).Encode(users[0])

}
