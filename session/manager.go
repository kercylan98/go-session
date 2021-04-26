package session

import (
	"errors"
)

// 会话管理器实现
//
// 可供给各种用途使用session管理
// todo: 过期时间
type Manager interface {
	// 注册一个新的会话
	RegisterSession(sessionId string) Session
	// 注销一个会话
	UnRegisterSession(session Session)
	// 获取所有会话
	GetAllSession() []Session
	// 获取特定id的会话
	GetSession(sessionId string) (Session, error)
}

func NewManager() *manager {
	return &manager{
		sessions: map[string]Session{},
	}
}

type manager struct {
	sessions 		map[string]Session
}

func (slf *manager) RegisterSession(sessionId string) Session {
	if session, exist := slf.sessions[sessionId]; exist {
		return session
	}else {
		session = newSession(slf, sessionId)
		slf.sessions[sessionId] = session
		return session
	}
}

func (slf *manager) UnRegisterSession(session Session) {
	delete(slf.sessions, session.GetId())
}

func (slf *manager) GetAllSession() []Session {
	var sessions = make([]Session, len(slf.sessions))
	var index = 0
	for _, session := range slf.sessions {
		sessions[index] = session
		index++
	}
	return sessions
}

func (slf *manager) GetSession(sessionId string) (Session, error) {
	if session, exist := slf.sessions[sessionId]; exist {
		return session, nil
	}
	return nil, errors.New("not found session with session id: " + sessionId)
}
