package session

import (
	"errors"
	"sync"
	"time"
)

// 会话管理器实现
//
// 可供给各种用途使用session管理
type Manager interface {
	// 注册一个新的会话
	RegisterSession(sessionId string) Session
	// 注销一个会话
	UnRegisterSession(session Session)
	// 获取所有会话
	GetAllSession() []Session
	// 获取特定id的会话
	GetSession(sessionId string) (Session, error)
	// 设置会话过期时间（默认 0, 永不过期）
	SetExpire(expire time.Duration)
}

func NewManager() *manager {
	return &manager{
		expire:   0,
		sessions: map[string]Session{},
	}
}

type manager struct {
	sync.Mutex
	expire   time.Duration
	sessions map[string]Session
}

func (slf *manager) SetExpire(expire time.Duration) {
	slf.clearExpire()
	slf.Lock()
	slf.expire = expire
	if slf.expire < 0 {
		slf.expire = 0
	}
	for _, session := range slf.sessions {
		session.setExpire(slf.expire)
	}
	slf.Unlock()
}

func (slf *manager) RegisterSession(sessionId string) Session {
	slf.Lock()
	defer slf.Unlock()
	if session, exist := slf.sessions[sessionId]; exist {
		return session
	} else {
		session = newSession(slf, sessionId, slf.expire)
		slf.sessions[sessionId] = session
		return session
	}
}

func (slf *manager) UnRegisterSession(session Session) {
	slf.Lock()
	delete(slf.sessions, session.GetId())
	slf.Unlock()
}

func (slf *manager) GetAllSession() []Session {
	slf.clearExpire()
	slf.Lock()
	defer slf.Unlock()
	var sessions []Session
	for _, session := range slf.sessions {
		if !session.IsExpire() {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

func (slf *manager) GetSession(sessionId string) (Session, error) {
	slf.clearExpire()
	slf.Lock()
	defer slf.Unlock()
	if session, exist := slf.sessions[sessionId]; exist {
		if session.IsExpire() {
			return nil, errors.New("the session has expired")
		}
		return session, nil
	}
	return nil, errors.New("not found session with session id: " + sessionId)
}

// 清理过期
func (slf *manager) clearExpire() {
	for _, session := range slf.sessions {
		if session.IsExpire() {
			slf.UnRegisterSession(session)
		}
	}
}
