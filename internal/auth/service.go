package auth

import (
	"auth/config"
	"auth/internal/store"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	config *config.AuthConfig
	logger *slog.Logger
	ctx    *context.Context
}

func NewAuthService(config *config.AuthConfig, logger *slog.Logger, ctx *context.Context) *AuthService {
	return &AuthService{
		config: config,
		logger: logger,
		ctx:    ctx,
	}
}

func (as *AuthService) CreateToken(user *store.User) (*Token, error) {
	accessToken, _, err := as.createAccessToken(user.Id, user.Email, false, time.Minute*15)
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

func (as *AuthService) createAccessToken(id int64, email string, isAdmin bool, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, isAdmin, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(as.config.SecretKey))
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, claims, nil
}

func (as *AuthService) VerifyToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(as.config.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
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

func (as *AuthService) VerifyRefreshToken(token *Token) error {
	sha256 := sha256.New()
	io.WriteString(sha256, as.config.SecretKey)

	salt := string(sha256.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	data, err := base64.URLEncoding.DecodeString(token.Refresh)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	if string(plain) != token.Access {
		return errors.New("invalid token")
	}

	return nil
}
