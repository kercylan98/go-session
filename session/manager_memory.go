package session

import (
	"errors"
	"sync"
	"time"
)

func NewManagerMemory() *managerMemory {
	return &managerMemory{
		expire:   0,
		sessions: map[string]Session{},
	}
}

// 采用内存存储的Session管理器实现
type managerMemory struct {
	sync.Mutex
	expire   time.Duration
	sessions map[string]Session
}

func (slf *managerMemory) SetExpire(expire time.Duration) error {
	if err := slf.clearExpire(); err != nil {
		return err
	}
	slf.Lock()
	slf.expire = expire
	if slf.expire < 0 {
		slf.expire = 0
	}
	for _, session := range slf.sessions {
		if err := session.SetExpire(slf.expire); err != nil {
			return err
		}
	}
	slf.Unlock()
	return nil
}

func (slf *managerMemory) RegisterSession(sessionId string) (Session, error) {
	slf.Lock()
	defer slf.Unlock()
	if session, exist := slf.sessions[sessionId]; exist {
		return session, nil
	} else {
		session = newSessionMemory(slf, sessionId, slf.expire)
		slf.sessions[sessionId] = session
		return session, nil
	}
}

func (slf *managerMemory) UnRegisterSession(session Session) error {
	slf.Lock()
	delete(slf.sessions, session.GetId())
	slf.Unlock()
	return nil
}

func (slf *managerMemory) GetAllSession() ([]Session, error) {
	if err := slf.clearExpire(); err != nil {
		return nil, err
	}
	slf.Lock()
	defer slf.Unlock()
	var sessions []Session
	for _, session := range slf.sessions {
		if !session.IsExpire() {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (slf *managerMemory) GetSession(sessionId string) (Session, error) {
	if err := slf.clearExpire(); err != nil {
		return nil, err
	}
	slf.Lock()
	defer slf.Unlock()
	if session, exist := slf.sessions[sessionId]; exist {
		if session.IsExpire() {
			return nil, errors.New("the sessionMemory has expired")
		}
		return session, nil
	}
	return nil, errors.New("not found sessionMemory with sessionMemory id: " + sessionId)
}

// 清理过期
func (slf *managerMemory) clearExpire() error {
	for _, session := range slf.sessions {
		if session.IsExpire() {
			if err := slf.UnRegisterSession(session); err != nil {
				return err
			}
		}
	}
	return nil
}
