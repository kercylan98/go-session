package session

import "errors"

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
}


// 构建一个新的会话
func newSession(sm Manager, sessionId string) *session {
	return &session{
		id: sessionId,
		sm: sm,
		storage: map[interface{}]interface{}{},
	}
}

// 会话结构
type session struct {
	id      string
	sm      Manager
	storage map[interface{}]interface{}
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


