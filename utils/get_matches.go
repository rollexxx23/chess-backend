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
