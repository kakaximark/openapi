package service

import (
	"errors"
	"time"

	"openapi/internal/db"
	"openapi/internal/model"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserDisabled     = errors.New("user is disabled")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrNotAdmin         = errors.New("only admin users are allowed")
	ErrUserAlreadyExist = errors.New("username already exists")
)

// ValidateUser 验证用户登录
func ValidateUser(username, password string) (*model.User, error) {
	var user model.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}

	// 1. 验证密码
	if user.Password != password {
		return nil, ErrInvalidPassword
	}

	// 2. 验证用户状态
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	return &user, nil
}

// UpdateUserToken 更新用户的token和会话信息
func UpdateUserToken(userID uint, username, token, ip, userAgent string) error {
	// 先将该用户的所有活动会话设为失效
	if err := db.DB.Model(&model.UserSession{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false).Error; err != nil {
		return err
	}

	// 创建新会话
	session := &model.UserSession{
		UserID:    userID,
		Username:  username,
		Token:     token,
		IP:        ip,
		UserAgent: userAgent,
		ExpiredAt: time.Now().Add(24 * time.Hour),
		IsActive:  true,
	}

	return db.DB.Create(session).Error
}

// ValidateToken 验证token是否有效
func ValidateToken(userID uint, token string) error {
	var session model.UserSession
	err := db.DB.Where("user_id = ? AND token = ? AND is_active = ? AND expired_at > ?",
		userID, token, true, time.Now()).First(&session).Error

	if err != nil {
		return errors.New("invalid or expired session")
	}

	return nil
}

// InvalidateOtherSessions 使其他会话失效
func InvalidateOtherSessions(userID uint, currentToken string) error {
	return db.DB.Model(&model.UserSession{}).
		Where("user_id = ? AND token != ? AND is_active = ?", userID, currentToken, true).
		Update("is_active", false).Error
}

// LogoutSession 登出指定会话
func LogoutSession(token string) error {
	return db.DB.Model(&model.UserSession{}).
		Where("token = ?", token).
		Update("is_active", false).Error
}
