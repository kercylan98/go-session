package session

import (
	"errors"
	"time"
)

// 会话实现
type Session interface {
	// 获取会话ID
	GetId() string
	// 向会话中存储数据
	Store(key interface{}, data interface{})
	// 加载会话中的数据
	Load(key interface{}) (interface{}, error)
	// 删除会话中的数据
	Del(key interface{})
	// 结束会话
	Close()
	// 获取会话管理器
	GetSessionManager() Manager
	// 是否已过期
	IsExpire() bool
	// 设置过期时间
	setExpire(expire time.Duration)
}

// 构建一个新的会话
func newSession(sm Manager, sessionId string, expireTime time.Duration) *session {
	return &session{
		id:          sessionId,
		sm:          sm,
		storage:     map[interface{}]interface{}{},
		createdTime: time.Now(),
		expireTime:  expireTime,
	}
}

// 会话结构
type session struct {
	id          string
	sm          Manager
	storage     map[interface{}]interface{}
	createdTime time.Time
	expireTime  time.Duration
}

func (slf *session) setExpire(expire time.Duration) {
	slf.expireTime = expire
}

func (slf *session) IsExpire() bool {
	return slf.expireTime != 0 && time.Now().Sub(slf.createdTime).Milliseconds() > slf.expireTime.Milliseconds()
}

func (slf *session) Store(key interface{}, data interface{}) {
	slf.storage[key] = data
}

func (slf *session) Load(key interface{}) (interface{}, error) {
	if val, exist := slf.storage[key]; exist {
		return val, nil
	}
	return nil, errors.New("can not found session value with key")
}

func (slf *session) Del(key interface{}) {
	delete(slf.storage, key)
}

func (slf *session) Close() {
	slf.sm.UnRegisterSession(slf)
}

func (slf *session) GetSessionManager() Manager {
	return slf.sm
}

func (slf *session) GetId() string {
	return slf.id
}
