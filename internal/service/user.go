package service

import (
	"errors"
	"time"

	"openapi/internal/db"
	"openapi/internal/model"
)

const (
	MaxLoginAttempts = 5                // 最大登录尝试次数
	LockoutDuration  = 15 * time.Minute // 锁定时间
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserDisabled     = errors.New("user is disabled")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrNotAdmin         = errors.New("only admin users are allowed")
	ErrUserAlreadyExist = errors.New("username already exists")
	ErrTooManyAttempts  = errors.New("too many login attempts, please try again later")
)

// ValidateUser 验证用户登录
func ValidateUser(username, password string) (*model.User, error) {
	var user model.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}

	// 验证用户状态
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 检查用户是否被锁定
	if user.LoginAttempts >= MaxLoginAttempts {
		if time.Since(user.LastLoginAttemptAt) < LockoutDuration {
			return nil, ErrTooManyAttempts
		}
		db.DB.Model(&user).Updates(map[string]interface{}{
			"login_attempts": 0,
		})
	}

	// 1. 验证密码
	if err := user.ValidatePassword(password); err != nil {
		db.DB.Model(&user).Updates(map[string]interface{}{
			"login_attempts":        user.LoginAttempts + 1,
			"last_login_attempt_at": time.Now(),
		})
		return nil, ErrInvalidPassword
	}

	// 登陆成功，重制尝试次数
	db.DB.Model(&user).Updates(map[string]interface{}{
		"login_attempts":        0,
		"last_login_attempt_at": time.Now(),
		"last_login_at":         time.Now().Unix(),
	})

	return &user, nil
}

// GetRemainingAttempts 获取剩余尝试次数和锁定剩余时间
func GetRemainingAttempts(username string) (attempts int, lockoutRemaining time.Duration, err error) {
	var user model.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, 0, ErrUserNotFound
	}

	attempts = MaxLoginAttempts - user.LoginAttempts
	if attempts < 0 {
		attempts = 0
	}

	if user.LoginAttempts >= MaxLoginAttempts {
		lockoutRemaining = LockoutDuration - time.Since(user.LastLoginAttemptAt)
		if lockoutRemaining < 0 {
			lockoutRemaining = 0
		}
	}

	return attempts, lockoutRemaining, nil
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
