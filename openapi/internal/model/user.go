package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// User 用户表
type User struct {
	gorm.Model
	ID                 uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username           string    `gorm:"column:username;type:varchar(50);uniqueIndex;not null" json:"username"`
	Password           string    `gorm:"column:password;type:varchar(100);not null" json:"-"` // json:"-" 避免返回密码
	LastToken          string    `gorm:"column:last_token;type:varchar(500)" json:"-"`        // 最后一次登录的token
	IsAdmin            bool      `gorm:"column:is_admin;default:false" json:"is_admin"`       // 是否是管理员
	LastLoginAt        int64     `gorm:"column:last_login_at" json:"last_login_at"`           // 最后登录时间
	Status             int       `gorm:"column:status;default:1" json:"status"`               // 状态 1:正常 0:禁用
	LoginAttempts      int       `gorm:"column:login_attempts;default:0" json:"-"`
	LastLoginAttemptAt time.Time `gorm:"column:last_login_attempt_at" json:"-"`
}

// EncryptPassword 加密密码
func (u *User) EncryptPassword() error {
	if len(u.Password) == 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

// ValidatePassword 验证密码
func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserSession 用户会话表
type UserSession struct {
	gorm.Model
	UserID    uint      `gorm:"column:user_id;index" json:"user_id"`
	Username  string    `gorm:"column:username;type:varchar(50)" json:"username"`
	Token     string    `gorm:"column:token;type:varchar(500)" json:"token"`
	IP        string    `gorm:"column:ip;type:varchar(50)" json:"ip"`
	UserAgent string    `gorm:"column:user_agent;type:varchar(500)" json:"user_agent"`
	ExpiredAt time.Time `gorm:"column:expired_at" json:"expired_at"`
	IsActive  bool      `gorm:"column:is_active;default:true" json:"is_active"`
}

// TableName 指定表名
func (UserSession) TableName() string {
	return "user_sessions"
}
