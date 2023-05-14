package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

var jwt_secret []byte

func init() {
	var err error
	jwt_secret, err = os.ReadFile("jwt_secret")
	if err != nil {
		log.Fatal("Could not read jwt secret from file")
	}
}

func VerifyToken(tokenStr string) bool {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwt_secret, nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v\n", err)
		return false
	}

	if token.Valid {
		log.Println("Token verified")
		return true
	} else {
		log.Println("Invalid token")
		return false
	}
}

func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenStr, err := token.SignedString(jwt_secret)
	if err != nil {
		return "", fmt.Errorf("Could not sign token. %v", err)
	}

	return tokenStr, nil
}

func IfAuthorized(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("No authorization header in request")
		return
	}

	tokenStr := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwt_secret, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Error parsing token: %v\n", err)
		return
	}

	if token.Valid {
		w.WriteHeader(http.StatusOK)
		log.Println("Token verified")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Invalid token")
	}
}
