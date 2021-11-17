package session

import (
	"time"
)

// Manager 会话管理器实现
//
// 可供给各种用途使用session管理
type Manager interface {
	// RegisterSession 注册一个新的会话
	RegisterSession(sessionId string) (Session, error)
	// UnRegisterSession 注销一个会话
	UnRegisterSession(session Session) error
	// GetAllSession 获取所有会话
	GetAllSession() ([]Session, error)
	// GetSession 获取特定id的会话
	GetSession(sessionId string) (Session, error)
	// SetExpire 设置会话过期时间（默认 0, 永不过期）
	SetExpire(expire time.Duration) error
}
