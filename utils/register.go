package utils

import (
	"encoding/json"
	"net/http"

	"github.com/rollexxx23/chess/database"
	"github.com/rollexxx23/chess/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	var users []models.User
	json.NewDecoder(r.Body).Decode(&user)
	user.HashPassword(user.Password)
	// check for duplicate entry
	database.Instance.Where("email = ? OR username = ?", user.Email, user.Username).Find(&users)
	if len(users) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		error := models.Error{Message: "username/ email already exists"}
		json.NewEncoder(w).Encode(error)
		return
	}
	database.Instance.Create(&user)
	json.NewEncoder(w).Encode(user)

}
