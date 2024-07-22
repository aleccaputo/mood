package api

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

func withJWTAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		id := r.PathValue("id")
		// remove the word bearer
		sanitzedTokenString := strings.Split(tokenString, " ")[1]
		_, err := validateJWT(sanitzedTokenString, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

func validateJWT(tokenString string, id string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	// not asking for data by userId
	if id == "" {
		return token, nil
	}
	if err != nil {
		return nil, err
	}

	aud, err := token.Claims.GetAudience()
	if err != nil {
		return nil, err
	}

	for _, curAud := range aud {
		if curAud == id {
			return token, nil
		}
	}
	return nil, errors.New("unauthorized")

}

func createJwt(id string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	claims := &jwt.RegisteredClaims{
		Audience: []string{id},
		Issuer:   "Alec LLC",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
