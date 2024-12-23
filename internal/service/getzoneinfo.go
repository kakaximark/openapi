package service

import (
	"errors"
	"openapi/internal/db"
	"openapi/internal/model"
	"strconv"
)

// GetZoneInfo 获取账号设置信息
func GetZoneInfo(token string) (string, string, error) {
	// 从user_sessions表获取用户信息
	var session model.UserSession
	if err := db.DB.Where("token = ? AND is_active = true", token).First(&session).Error; err != nil {
		return "", "", errors.New("invalid session")
	}

	return strconv.FormatUint(uint64(session.UserID), 10), session.Username, nil
}
