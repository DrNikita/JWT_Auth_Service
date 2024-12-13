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
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRepository struct {
	config *config.AuthConfig
	logger *slog.Logger
	ctx    *context.Context
}

func NewAuthService(config *config.AuthConfig, logger *slog.Logger, ctx *context.Context) *AuthRepository {
	return &AuthRepository{
		config: config,
		logger: logger,
		ctx:    ctx,
	}
}

func (as *AuthRepository) CreateToken(user *store.User) (*Token, error) {
	accessClaims, err := NewUserClaims(user.Id, user.Email, false, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	accessToken, err := as.CreateAccessToken(accessClaims)
	if err != nil {
		return nil, err
	}

	refreshClaims, err := NewUserClaims(user.Id, user.Email, false, 24*time.Hour)
	if err != nil {
		return nil, err
	}
	refreshToken, err := as.CreateRefreshToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	return &Token{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (as *AuthRepository) CreateAccessToken(claims *UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(as.config.SecretKey))
	if err != nil {
		as.logger.Error("failed to sign access token")
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (as *AuthRepository) CreateRefreshToken(claims *UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(as.config.SecretKey))
	if err != nil {
		as.logger.Error("failed to sign refresh token")
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (as *AuthRepository) VerifyAccessToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			as.logger.Error("failed to verify access token")
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(as.config.SecretKey), nil
	})
	if err != nil {
		as.logger.Error("failed to verify access token")
		return nil, fmt.Errorf("error parsing token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		as.logger.Error("failed to verify access token")
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (as *AuthRepository) VerifyRefreshToken(refreshToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			as.logger.Error("failed to verify refresh token")
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(as.config.SecretKey), nil
	})
	if err != nil {
		as.logger.Error("failed to verify refresh token")
		return nil, fmt.Errorf("error parsing token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		as.logger.Error("failed to verify refresh token")
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// Unused
// Deprecated
func (as *AuthRepository) DEPRECATED_createRefreshToken(accessToken string) (string, error) {
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
