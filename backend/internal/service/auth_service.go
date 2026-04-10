package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users      domain.UserRepository
	tokens     domain.RefreshTokenRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(
	users domain.UserRepository,
	tokens domain.RefreshTokenRepository,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		users:      users,
		tokens:     tokens,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (a *AuthService) Register(ctx context.Context, email, password string, orgID uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	_, err := a.users.GetByEmail(ctx, email)
	if err == nil {
		return uuid.Nil, domain.ErrAlreadyExists
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.Nil, fmt.Errorf("authService.Register: %w", err)
	}
	newUser := domain.User{
		ID:           id,
		OrgID:        orgID,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	id, err = a.users.Create(ctx, &newUser)
	if err != nil {
		return uuid.Nil, fmt.Errorf("authService.Register: %w", err)
	}
	return id, nil
}
func (a *AuthService) Login(ctx context.Context, email, password string, orgID uuid.UUID) (TokenPair, error) {
	gotUser, err := a.users.GetByEmail(ctx, email)
	if err != nil {
		return TokenPair{}, domain.ErrUnauthorized
	}
	err = bcrypt.CompareHashAndPassword([]byte(gotUser.PasswordHash), []byte(password))
	if err != nil {
		return TokenPair{}, domain.ErrUnauthorized
	}
	accessToken, err := a.generateAccessToken(gotUser)
	if err != nil {
		return TokenPair{}, fmt.Errorf("authService.Login: %w", err)
	}
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return TokenPair{}, fmt.Errorf("authService.Login: %w", err)
	}
	hash := hashToken(refreshToken)
	err = a.tokens.Create(ctx, &domain.RefreshToken{
		UserID:    gotUser.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add((a.refreshTTL)),
	})
	if err != nil {
		return TokenPair{}, fmt.Errorf("authService.Login: %w", err)
	}
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthService) generateAccessToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		OrgID:  user.OrgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
func (a *AuthService) Refresh(ctx context.Context, hashToken string) (TokenPair, error) {
	token, err := a.tokens.GetByTokenHash(ctx, hashToken)
	if err != nil {
		return TokenPair{}, fmt.Errorf("authService.Refresh: %w", err)
	}
	if token.ExpiresAt.Before(time.Now()) {
		return TokenPair{}, domain.ErrUnauthorized
	}
	gotUser, err := a.users.GetByID(ctx, token.UserID)
	if err != nil {
		return TokenPair{}, domain.ErrNotFound
	}
	accessToken, err := a.generateAccessToken(gotUser)
	if err != nil {
		return TokenPair{}, fmt.Errorf("authService.Refresh: %w", err)
	}
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return TokenPair{}, fmt.Errorf("authService.Refresh: %w", err)
	}
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (a *AuthService) Logout(ctx context.Context, refreshToken string) error {
	hash := hashToken(refreshToken)
	return a.tokens.DeleteByTokenHash(ctx, string(hash))
}
func (a *AuthService) ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return a.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func hashToken(refreshToken string) string {
	h := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(h[:])
}
