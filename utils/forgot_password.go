package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rollexxx23/chess/database"
	"github.com/rollexxx23/chess/misc"
	"github.com/rollexxx23/chess/models"
)

type requestForgotPasswordReq struct {
	Email string
}

type forgotPwStruct struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func RequestForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req requestForgotPasswordReq
	json.NewDecoder(r.Body).Decode(&req)

	expirationTime := time.Now().Add(time.Minute * 120)
	claims := &Claims{
		Email: req.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	url := "http://chess.com/forgotpassword/" + tokenString
	mailer := misc.New("sandbox.smtp.mailtrap.io", 2525, "7656ea062cdd81", "08e44d0b8badb1", "arin@example.com")
	err = mailer.Send(req.Email, url)
	if err == nil {
		w.Write([]byte(fmt.Sprintf("Email sent to %s", req.Email)))
	} else {
		w.Write([]byte(fmt.Sprintf(err.Error())))
	}

}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req forgotPwStruct
	json.NewDecoder(r.Body).Decode(&req)
	if req.ConfirmPassword != req.Password {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)
	tokenStr := params["token"]
	if tokenStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user models.User
	var users []models.User
	database.Instance.Where("email = ?", claims.Email).Find(&users)
	if len(users) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	user = users[0]

	user.HashPassword(req.Password)
	database.Instance.Save(user)

}
