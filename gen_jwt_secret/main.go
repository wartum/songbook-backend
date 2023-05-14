package main

import (
	"crypto/rand"
	"log"
	"os"
)

func main() {
	jwtSecret := make([]byte, 32)
	_, err := rand.Read(jwtSecret)
	if err != nil {
		log.Fatal("Could not generate secret")
	}
	os.WriteFile("jwt_secret", jwtSecret, 0644)
}
