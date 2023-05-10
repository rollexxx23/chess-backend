package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rollexxx23/chess/database"
	"github.com/rollexxx23/chess/models"
)

type loginStruct struct {
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}
type loginResp struct {
	AccessToken string `json:"access_token"`
	Expires     string `json:"expires"`
}

type testStruct struct {
	Email string
}

var jwtKey = []byte("secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var loginCred loginStruct
	json.NewDecoder(r.Body).Decode(&loginCred)
	var users []models.User
	database.Instance.Where("email = ?", loginCred.Email).Find(&users)
	// no entry found
	if len(users) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		error := models.Error{Message: "No entry found"}
		json.NewEncoder(w).Encode(error)
		return
	}
	// check for password
	if users[0].CheckPassword(loginCred.Password) {
		// assign access token
		expirationTime := time.Now().Add(time.Minute * 120)
		claims := &Claims{
			Email: loginCred.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		resp := loginResp{
			AccessToken: tokenString,
			Expires:     expirationTime.String(),
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)
	} else {
		w.WriteHeader(http.StatusForbidden)
		error := models.Error{Message: "Wrong Pw"}
		json.NewEncoder(w).Encode(error)
	}

}

//refresh

func Refresh(w http.ResponseWriter, r *http.Request) {

	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		w.WriteHeader(http.StatusBadRequest)
		error := models.Error{Message: "Missing auth token"}
		json.NewEncoder(w).Encode(error)
		return
	}
	tokenStr := strings.Fields(bearer)[1]
	var req testStruct
	json.NewDecoder(r.Body).Decode(&req)
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

	expirationTime := time.Now().Add(time.Minute * 120)

	claims.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if claims.Email == req.Email {
		resp := loginResp{
			AccessToken: tokenString,
			Expires:     expirationTime.String(),
		}
		json.NewEncoder(w).Encode(resp)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Forbidden")))
	}

}

// for testing jwt
func Test(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(fmt.Sprintf("Hello,")))

}
