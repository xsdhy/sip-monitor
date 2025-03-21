package services

import (
	"context"
	"errors"
	"sip-monitor/src/entity"
	"sip-monitor/src/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token expired")
)

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthService provides authentication functionality
type AuthService struct {
	logger     *logrus.Logger
	repository model.Repository
	jwtSecret  []byte
	jwtExpiry  time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(logger *logrus.Logger, repository model.Repository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		logger:     logger,
		repository: repository,
		jwtSecret:  []byte(jwtSecret),
		jwtExpiry:  jwtExpiry,
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user")
		return "", err
	}

	if user == nil {
		return "", ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate token")
		return "", err
	}

	return token, nil
}

// ValidateToken validates the provided token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user by ID")
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	// Don't return password
	user.Password = ""
	return user, nil
}

// generateToken creates a new token for the given user
func (s *AuthService) generateToken(user *entity.User) (string, error) {
	now := time.Now()
	expirationTime := now.Add(s.jwtExpiry)

	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// UpdateUserInfo updates a user's nickname
func (s *AuthService) UpdateUserInfo(ctx context.Context, userID int64, nickname string) error {
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user")
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	user.Nickname = nickname
	user.UpdateAt = time.Now()

	if err := s.repository.UpdateUser(ctx, user); err != nil {
		s.logger.WithError(err).Error("Failed to update user")
		return err
	}

	return nil
}

// UpdateUserPassword updates a user's password
func (s *AuthService) UpdateUserPassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user")
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash password")
		return err
	}

	user.Password = string(hashedPassword)
	user.UpdateAt = time.Now()

	if err := s.repository.UpdateUser(ctx, user); err != nil {
		s.logger.WithError(err).Error("Failed to update user password")
		return err
	}

	return nil
}
