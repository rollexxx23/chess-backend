package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/rollexxx23/chess/models"
	"github.com/rollexxx23/chess/utils"
)

type mwStruct struct {
	Email string
}

var jwtKey = []byte("secret_key")

func AuthMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("middleware", r.URL)
		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			w.WriteHeader(http.StatusBadRequest)
			error := models.Error{Message: "Missing auth token"}
			json.NewEncoder(w).Encode(error)
			return
		}
		tokenStr := strings.Fields(bearer)[1]
		var req mwStruct
		json.NewDecoder(r.Body).Decode(&req)
		claims := &utils.Claims{}

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
		if claims.Email == req.Email {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Forbidden")))
		}

	}
}
