package session

import (
	"time"
)

// Session 会话实现
type Session interface {
	// GetId 获取会话ID
	GetId() string
	// Store 向会话中存储数据
	Store(key string, data interface{}) error
	// Load 加载会话中的数据
	Load(key string) (interface{}, error)
	// Del 删除会话中的数据
	Del(key string) error
	// Close 结束会话
	Close() error
	// GetSessionManager 获取会话管理器
	GetSessionManager() Manager
	// IsExpire 是否已过期
	IsExpire() bool
	// SetExpire 设置过期时间
	SetExpire(expire time.Duration) error
}
