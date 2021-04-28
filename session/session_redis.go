package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"time"
)

func newSessionRedis(rm *managerRedis, hashKey string, sessionId string, redis *redis.Client) *sessionRedis {
	return &sessionRedis{
		hashKey:   hashKey,
		sessionId: sessionId,
		redis:     redis,
		rm:        rm,
	}
}

// 基于Redis存储的session
type sessionRedis struct {
	hashKey   string
	sessionId string
	redis     *redis.Client
	rm        *managerRedis
}

func (slf *sessionRedis) GetId() string {
	return slf.sessionId
}

func (slf *sessionRedis) Store(key string, data interface{}) error {
	unlock, err := slf.lock()
	if err != nil {
		return err
	}
	defer unlock()

	if slf.IsExpire() {
		return errors.New("session has expired")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = slf.redis.HSet(slf.hashKey, key, string(jsonData)).Result()
	return err
}

func (slf *sessionRedis) Load(key string) (interface{}, error) {
	unlock, err := slf.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	if slf.IsExpire() {
		return nil, errors.New("session has expired")
	}

	r, err := slf.redis.HGet(slf.hashKey, key).Result()
	if err != nil {
		return nil, err
	}

	var a interface{}
	err = json.Unmarshal([]byte(r), &a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (slf *sessionRedis) Del(key string) error {
	unlock, err := slf.lock()
	if err != nil {
		return err
	}
	defer unlock()

	if slf.IsExpire() {
		return errors.New("session has expired")
	}

	_, err = slf.redis.HDel(slf.hashKey, key).Result()
	return err
}

func (slf *sessionRedis) Close() error {
	return slf.rm.UnRegisterSession(slf)
}

func (slf *sessionRedis) GetSessionManager() Manager {
	return slf.rm
}

func (slf *sessionRedis) IsExpire() bool {
	session, _ := slf.rm.GetSession(slf.sessionId)
	return session == nil
}

func (slf *sessionRedis) SetExpire(expire time.Duration) error {
	if _, err := slf.redis.Expire(slf.hashKey, expire).Result(); err != nil {
		return err
	}
	return nil
}

// 进行加锁
func (slf *sessionRedis) lock() (unlock func() error, err error) {

tryLock:
	{
		success, err := slf.redis.SetNX(slf.hashKey+"lock", slf.sessionId, 5*time.Second).Result()
		if err != nil {
			return nil, err
		}

		if !success {
			goto tryLock
		}
	}
	// 返回解锁函数
	return func() error {
		_, err := slf.redis.Del(slf.hashKey + "lock").Result()
		return err
	}, nil
}

// 生成锁id
func (slf *sessionRedis) newLookId() string {
	return fmt.Sprintf(uuid.NewV4().String())
}
