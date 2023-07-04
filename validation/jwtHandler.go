package validation

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateJwt(key string, payload *jwt.MapClaims) (string, error) {

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := t.SignedString([]byte(key))
	return token, err
}

func verifyJwt(jwtToken string) (jwt.MapClaims, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file!")
	}
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else {
			return []byte(os.Getenv("SECRET_KEY")), nil
		}

	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
