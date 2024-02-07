package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthHandler provides HTTP handlers for authentication operations
type AuthHandler struct {
	logger      *logrus.Logger
	authService *AuthService
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

// UpdateUserInfoRequest represents the update user info request payload
type UpdateUserInfoRequest struct {
	Nickname string `json:"nickname" binding:"required"`
}

// UpdatePasswordRequest represents the update password request payload
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(logger *logrus.Logger, authService *AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
	}
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid login request")
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "msg": "Invalid request format"})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.WithError(err).Warn("Login failed")

		statusCode := http.StatusInternalServerError
		errMsg := "Internal server error"
		errCode := 5000

		if err == ErrUserNotFound || err == ErrInvalidCredentials {
			statusCode = http.StatusUnauthorized
			errMsg = "Invalid username or password"
			errCode = 4001
		}

		c.JSON(statusCode, gin.H{"code": errCode, "msg": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 2000, "msg": "Login successful", "data": token})
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 4001, "msg": "Not authenticated"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(int64))
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5000, "msg": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 2000, "msg": "Success", "data": user})
}

// UpdateUserInfo updates the current user's information
func (h *AuthHandler) UpdateUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 4001, "msg": "未认证"})
		return
	}

	var req UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("无效的用户信息更新请求")
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "msg": "无效的请求格式"})
		return
	}

	err := h.authService.UpdateUserInfo(c.Request.Context(), userID.(int64), req.Nickname)
	if err != nil {
		h.logger.WithError(err).Error("更新用户信息失败")

		statusCode := http.StatusInternalServerError
		errMsg := "服务器内部错误"
		errCode := 5000

		if err == ErrUserNotFound {
			statusCode = http.StatusNotFound
			errMsg = "用户不存在"
			errCode = 4004
		}

		c.JSON(statusCode, gin.H{"code": errCode, "msg": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 2000, "msg": "用户信息更新成功"})
}

// UpdatePassword updates the current user's password
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 4001, "msg": "未认证"})
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("无效的密码更新请求")
		c.JSON(http.StatusBadRequest, gin.H{"code": 4000, "msg": "无效的请求格式"})
		return
	}

	err := h.authService.UpdateUserPassword(c.Request.Context(), userID.(int64), req.OldPassword, req.NewPassword)
	if err != nil {
		h.logger.WithError(err).Error("更新密码失败")

		statusCode := http.StatusInternalServerError
		errMsg := "服务器内部错误"
		errCode := 5000

		if err == ErrUserNotFound {
			statusCode = http.StatusNotFound
			errMsg = "用户不存在"
			errCode = 4004
		} else if err == ErrInvalidCredentials {
			statusCode = http.StatusUnauthorized
			errMsg = "原密码不正确"
			errCode = 4001
		}

		c.JSON(statusCode, gin.H{"code": errCode, "msg": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 2000, "msg": "密码更新成功"})
}
