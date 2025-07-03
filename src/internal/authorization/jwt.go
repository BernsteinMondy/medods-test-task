package authorization

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/medods-test-task/src/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	refreshTokenTTL = time.Hour * 24 * 30
	accessTokenTTL  = time.Minute * 30
)

var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrTokenExpired            = errors.New("token expired")
	ErrInvalidClaims           = errors.New("invalid token claims")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidTokenFormat      = errors.New("invalid token format")
)

type TokenService struct {
	secretKey []byte
}

var _ service.TokenService = (*TokenService)(nil)

func NewTokenService(secretKey string) *TokenService {
	return &TokenService{
		secretKey: []byte(secretKey),
	}
}

func (ts *TokenService) GenerateAccessToken(userID, tokenID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      userID.String(),
		"token_id": tokenID.String(),
		"exp":      now.Add(accessTokenTTL).Unix(),
		"iat":      now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(ts.secretKey)
}

func (ts *TokenService) GenerateRefreshToken(userID uuid.UUID, userAgent, ip string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":        userID.String(),
		"exp":        now.Add(refreshTokenTTL).Unix(),
		"iat":        now.Unix(),
		"user_agent": userAgent,
		"ip":         ip,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(ts.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	encodedToken := base64.StdEncoding.EncodeToString([]byte(tokenStr))

	return encodedToken, nil
}

func (ts *TokenService) VerifyRefreshToken(token, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify refresh token: %w", err)
	}
	return true, nil
}

func (ts *TokenService) ParseToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return ts.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, ErrInvalidTokenFormat
	}
	if time.Now().After(time.Unix(int64(exp), 0)) {
		return nil, ErrTokenExpired
	}

	if iat, ok := claims["iat"].(float64); ok {
		if time.Now().Before(time.Unix(int64(iat), 0)) {
			return nil, ErrInvalidTokenFormat
		}
	}

	result := make(map[string]interface{})
	for k, v := range claims {
		result[k] = v
	}

	return result, nil
}
