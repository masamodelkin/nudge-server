package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/masamodelkin/nudge-server/internal/auth"
	"github.com/masamodelkin/nudge-server/internal/model"
	"github.com/masamodelkin/nudge-server/internal/store"
)

var (
	ErrUsernameTaken      = errors.New("username already taken")
	ErrPasswordHashing    = errors.New("failed to hash password")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrTokenGeneration    = errors.New("failed to generate tokens")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type AuthService struct {
	store  *store.Store
	tokens *auth.TokenService
}

func NewAuthService(s *store.Store, tokens *auth.TokenService) *AuthService {
	return &AuthService{
		store:  s,
		tokens: tokens,
	}
}

type RegisterRequest struct {
	Username string
	Password string
	Email    *string
}

type RegisterResponse struct {
	Username string  `json:"username"`
	Email    *string `json:"email"`
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	Username     string  `json:"username"`
	Email        *string `json:"email"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Register(req RegisterRequest) (*RegisterResponse, error) {
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, ErrPasswordHashing
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: hash,
		Email:        req.Email,
		AuthProvider: "local",
	}

	if err := s.store.CreateUser(user); err != nil {
		return nil, ErrUsernameTaken
	}

	response := &RegisterResponse{
		Username: user.Username,
		Email:    user.Email,
	}

	return response, nil
}

// Login verifies credentials and returns JWT tokens.
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.store.GetUserByUsername(req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.generateAndStoreTokens(user.ID)
	if err != nil {
		return nil, ErrTokenGeneration
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Username:     user.Username,
		Email:        user.Email,
	}, nil
}

func (s *AuthService) Refresh(rawRefreshToken string) (*RefreshResponse, error) {
	userID, tokenType, err := s.tokens.ValidateToken(rawRefreshToken)
	if err != nil || tokenType != "refresh" {
		return nil, ErrInvalidToken
	}

	tokenHash := auth.HashToken(rawRefreshToken)
	stored, err := s.store.GetRefreshTokenByHash(tokenHash)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if stored.UserID != userID {
		return nil, ErrInvalidToken
	}

	s.store.DeleteRefreshToken(stored.ID)

	accessToken, refreshToken, err := s.generateAndStoreTokens(userID)
	if err != nil {
		return nil, ErrTokenGeneration
	}

	return &RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateAndStoreTokens(userID string) (string, string, error) {
	accessToken, err := s.tokens.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, expiresAt, err := s.tokens.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshRecord := &model.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		TokenHash: auth.HashToken(refreshToken),
		ExpiresAt: expiresAt.Unix(),
	}

	if err := s.store.CreateRefreshToken(refreshRecord); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	tokenHash := auth.HashToken(refreshToken)
	return s.store.DeleteRefreshTokenByHash(tokenHash)
}
