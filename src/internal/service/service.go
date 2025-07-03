package service

import (
	"context"
	"fmt"
	"github.com/BernsteinMondy/medods-test-task/src/internal/entities"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	UserExistsByUsername(ctx context.Context, username string) (bool, error)

	GetRefreshToken(ctx context.Context, id uuid.UUID) (string, error)
	UpdateRefreshToken(ctx context.Context, token string) (string, error)
}

type TokenService interface {
	GenerateAccessToken(userID, tokenID uuid.UUID) (string, error)
	GenerateRefreshToken(userID uuid.UUID, userAgent, ip string) (string, error)

	VerifyRefreshToken(token, hash string) (bool, error)
	ParseToken(tokenStr string) (map[string]interface{}, error)
}
type Hasher interface {
	Hash(src string) (string, error)
	CompareHashAndPassword(hash, password string) (bool, error)
}

type TokenEncoder interface {
	Encode(token []byte) string
}

type Service struct {
	TokenService TokenService
	repo         Repository
	hasher       Hasher
}

func NewService(
	repo Repository,
	hasher Hasher,
	tokenService TokenService,
) *Service {
	return &Service{
		repo:         repo,
		hasher:       hasher,
		TokenService: tokenService,
	}
}

type RegisterUserData struct {
	Username  string
	Password  string
	UserAgent string
	IP        string
}

func (s *Service) RegisterUser(ctx context.Context, data *RegisterUserData) (string, string, error) {
	exists, err := s.repo.UserExistsByUsername(ctx, data.Username)
	if err != nil {
		return "", "", fmt.Errorf("repository: user exists by username: %w", err)
	}
	if exists {
		return "", "", ErrAlreadyExists
	}

	hash, err := s.hasher.Hash(data.Password)
	if err != nil {
		return "", "", fmt.Errorf("password service: hash password: %w", err)
	}

	userID := uuid.New()
	refreshToken, err := s.TokenService.GenerateRefreshToken(userID, data.UserAgent, data.IP)
	if err != nil {
		return "", "", fmt.Errorf("token service: generate token: %w", err)
	}

	hashedToken, err := s.hasher.Hash(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("hasher: hash token: %w", err)
	}

	var (
		accessTokenID = uuid.New()

		user = &entities.User{
			ID:           userID,
			Username:     data.Username,
			Hash:         hash,
			RefreshToken: hashedToken,
		}
	)

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("repo: create user: %w", err)
	}

	accessToken, err := s.TokenService.GenerateAccessToken(userID, accessTokenID)
	if err != nil {
		return "", "", fmt.Errorf("token service: generate token: %w", err)
	}

	return refreshToken, accessToken, nil
}

func (s *Service) GetUserID(ctx context.Context, token string) (uuid.UUID, error) {
	// TODO: Finish
	return uuid.Nil, nil
}

func (s *Service) GetTokenPair(ctx context.Context, id uuid.UUID) (string, string, error) {
	// TODO: Finish
	return "", "", nil
}

func (s *Service) RefreshTokenPair(ctx context.Context, refresh, userAgent, ip string) (string, string, error) {
	// TODO: Finish
	return "", "", nil
}
