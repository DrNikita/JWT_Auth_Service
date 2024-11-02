package auth

import (
	"auth/config"
	"auth/internal/store"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	config *config.AuthConfig
	logger *slog.Logger
}

func NewAuthService(config *config.AuthConfig, logger *slog.Logger) *AuthService {
	return &AuthService{
		config: config,
		logger: logger,
	}
}

func (as *AuthService) CreateToken(user *store.User) (*Token, error) {

	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	accessToken, err := token.SignedString([]byte(as.config.SecretKey))
	if err != nil {
		return nil, err
	}

	refreshToken, err := as.createRefreshToken(accessToken)
	if err != nil {
		return nil, err
	}

	return &Token{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (as *AuthService) createRefreshToken(accessToken string) (string, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, as.config.SecretKey)

	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		fmt.Println(err.Error())

		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	refreshToken := base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(accessToken), nil))

	return refreshToken, nil
}
