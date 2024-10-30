package auth

import (
	"auth/internal/store"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(user *store.User) (string, error) {
	key := "secret"

	t := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"iss":        "my-auth-server",
			"user_id":    user.Id,
			"sub":        user.Name,
			"user_email": user.Email,
			"foo":        2,
		})

	signedJWT, err := t.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}
