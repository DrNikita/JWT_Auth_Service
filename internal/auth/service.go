package auth

import (
	"auth/config"
	"auth/internal/store"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
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
	userClaims := createUserClaims(user)
	accessToken, err := as.createAccessToken(&userClaims)
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

func (as *AuthService) createAccessToken(claims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(as.config.SecretKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (as *AuthService) createRefreshToken(accessToken string) (string, error) {
	sha256 := sha256.New()
	io.WriteString(sha256, as.config.SecretKey)

	salt := string(sha256.Sum(nil))[0:16]
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

func createUserClaims(user *store.User) jwt.MapClaims {
	return jwt.MapClaims{
		"email":    user.Email,
		"job_role": user.JobRole,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
	}
}


