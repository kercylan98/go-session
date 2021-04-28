package session

import (
	"errors"
	"time"
)

// 构建一个新的会话
func newSessionMemory(sm Manager, sessionId string, expireTime time.Duration) *sessionMemory {
	return &sessionMemory{
		id:          sessionId,
		sm:          sm,
		storage:     map[interface{}]interface{}{},
		createdTime: time.Now(),
		expireTime:  expireTime,
	}
}

// 单机器内存存储session
type sessionMemory struct {
	id          string
	sm          Manager
	storage     map[interface{}]interface{}
	createdTime time.Time
	expireTime  time.Duration
}

func (slf *sessionMemory) SetExpire(expire time.Duration) error {
	slf.expireTime = expire
	return nil
}

func (slf *sessionMemory) IsExpire() bool {
	return slf.expireTime != 0 && time.Now().Sub(slf.createdTime).Milliseconds() > slf.expireTime.Milliseconds()
}

func (slf *sessionMemory) Store(key string, data interface{}) error {
	slf.storage[key] = data
	return nil
}

func (slf *sessionMemory) Load(key string) (interface{}, error) {
	if val, exist := slf.storage[key]; exist {
		return val, nil
	}
	return nil, errors.New("can not found sessionMemory value with key")
}

func (slf *sessionMemory) Del(key string) error {
	delete(slf.storage, key)
	return nil
}

func (slf *sessionMemory) Close() error {
	return slf.sm.UnRegisterSession(slf)
}

func (slf *sessionMemory) GetSessionManager() Manager {
	return slf.sm
}

func (slf *sessionMemory) GetId() string {
	return slf.id
}
