package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sip-monitor/src/entity"
	"sip-monitor/src/model"
	"strings"
	"time"

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
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
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
	parts := strings.Split(tokenString, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidToken
	}

	// Extract payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Validate signature
	expectedSignature := s.generateSignature(parts[0])
	actualSignature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, ErrInvalidToken
	}

	// Parse claims
	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	// Check expiration
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrExpiredToken
	}

	return &claims, nil
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
	expirationTime := time.Now().Add(s.jwtExpiry)

	claims := &JWTClaims{
		UserID:    user.ID,
		Username:  user.Username,
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	// Create payload
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// Encode payload
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)

	// Generate signature
	signature := s.generateSignature(encodedPayload)

	// Encode signature
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	// Combine parts
	token := encodedPayload + "." + encodedSignature

	return token, nil
}

// generateSignature creates an HMAC-SHA256 signature for the given payload
func (s *AuthService) generateSignature(payload string) []byte {
	mac := hmac.New(sha256.New, s.jwtSecret)
	mac.Write([]byte(payload))
	return mac.Sum(nil)
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
