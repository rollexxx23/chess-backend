package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rollexxx23/chess/database"
	"github.com/rollexxx23/chess/middleware"
	"github.com/rollexxx23/chess/utils"
)

func initializeRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/register", utils.Register).Methods("POST")
	r.HandleFunc("/login", utils.Login).Methods("GET")
	r.HandleFunc("/test", middleware.AuthMW(utils.Test)).Methods("GET")
	r.HandleFunc("/forgot-password", utils.RequestForgotPassword)
	r.HandleFunc("/refresh", utils.Refresh).Methods("GET").Methods("GET")
	r.HandleFunc("/forgot-password/{token}", utils.ForgotPassword).Methods("PUT")
	log.Fatal(http.ListenAndServe(":5000", r))
}

func main() {
	// Initialize Database
	database.Connect("root:qazplm45@tcp(localhost:3306)/chess?parseTime=true")
	database.Migrate()
	initializeRouter()

}
