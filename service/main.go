package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "github.com/wartum/songbook-backend/db_access"
	jwt "github.com/wartum/songbook-backend/jwt"
)

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Token struct {
	Token string `json:"token"`
}

func verify_user(username string, password string) bool {
	user, err := db.GetUser(username)
	if err != nil {
		log.Printf("Authorization failed. %v\n", err)
		return false
	}

	hash_compute := sha256.New()
	hash_compute.Write([]byte(password + user.Salt))
	hash := fmt.Sprintf("%x", hash_compute.Sum(nil))

	if hash == user.Password {
		return true
	} else {
		log.Println("Password incorrect")
		return false
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	var creds LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Could not decode request")
		return
	}

	if creds.Token != "" && creds.Username != "" && creds.Password != "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Both token and password were provided")
		return
	}

	if creds.Token != "" {
		if !jwt.VerifyToken(creds.Token) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		if !verify_user(creds.Username, creds.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	token, err := jwt.GenerateToken(creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error generating token. %v\n", err)
		return
	}

	json.NewEncoder(w).Encode(Token{Token: token})
}

func main() {
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}
