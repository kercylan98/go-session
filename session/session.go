package session

import (
	"time"
)

// 会话实现
type Session interface {
	// 获取会话ID
	GetId() string
	// 向会话中存储数据
	Store(key string, data interface{}) error
	// 加载会话中的数据
	Load(key string) (interface{}, error)
	// 删除会话中的数据
	Del(key string) error
	// 结束会话
	Close() error
	// 获取会话管理器
	GetSessionManager() Manager
	// 是否已过期
	IsExpire() bool
	// 设置过期时间
	SetExpire(expire time.Duration) error
}
