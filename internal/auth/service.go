package auth

import (
	"auth/internal/store"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("very-secret-key")

func CreateJWT(user *store.User) (string, error) {

	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return t, nil
}
