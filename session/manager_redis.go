package session

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

const (
	redisKeyPrefix = "github.com/kercylan98/session_"
	managerLock    = "manager_lock"
)

func NewManagerRedis(address string, password ...string) *managerRedis {
	var rs = &managerRedis{
		expire: 0,
	}
	var pwd = ""
	if len(password) > 0 {
		pwd = password[0]
	}
	rs.redis = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: pwd,
	})
	return rs
}

// 采用Redis存储的session管理器
type managerRedis struct {
	redis  *redis.Client // redis客户端
	expire time.Duration
}

func (slf *managerRedis) RegisterSession(sessionId string) (Session, error) {
	unlock, err := slf.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	var hashKey = slf.newKey(sessionId)
	_, err = slf.redis.Pipelined(func(pipeliner redis.Pipeliner) error {
		if _, err = pipeliner.HSet(hashKey, "session_id", sessionId).Result(); err != nil {
			return err
		}
		if slf.expire > 0 {
			if _, err = pipeliner.Expire(hashKey, slf.expire).Result(); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return newSessionRedis(slf, hashKey, sessionId, slf.redis), nil
}

func (slf *managerRedis) UnRegisterSession(session Session) error {
	unlock, err := slf.lock()
	if err != nil {
		return err
	}
	defer unlock()

	_, err = slf.redis.Del(slf.newKey(session.GetId())).Result()
	return err
}

func (slf *managerRedis) GetAllSession() ([]Session, error) {
	unlock, err := slf.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	keys, err := slf.redis.Keys(redisKeyPrefix + "*").Result()
	if err != nil {
		return nil, err
	}

	var sessions = make([]Session, len(keys))
	for i := 0; i < len(keys); i++ {
		sessions[i] = newSessionRedis(slf, keys[i], slf.formatKey(keys[i]), slf.redis)
	}
	return sessions, nil
}

func (slf *managerRedis) GetSession(sessionId string) (Session, error) {
	unlock, err := slf.lock()
	if err != nil {
		return nil, err
	}
	defer unlock()

	existCount, err := slf.redis.Exists(slf.newKey(sessionId)).Result()
	if err != nil {
		return nil, err
	}

	if existCount > 0 {
		return newSessionRedis(slf, slf.newKey(sessionId), sessionId, slf.redis), nil
	} else {
		return nil, errors.New("the session does not exist or has been closed")
	}

}

func (slf *managerRedis) SetExpire(expire time.Duration) error {
	unlock, err := slf.lock()
	if err != nil {
		return err
	}
	defer unlock()
	if expire < 0 {
		expire = 0
	}
	slf.expire = expire
	keys, err := slf.redis.Keys(redisKeyPrefix + "*").Result()
	if err != nil {
		return err
	}

	_, err = slf.redis.Pipelined(func(pipeliner redis.Pipeliner) error {
		for _, key := range keys {
			if _, err = pipeliner.Expire(key, slf.expire).Result(); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// 进行加锁
func (slf *managerRedis) lock() (unlock func() error, err error) {
	lockId := slf.newLookId()

tryLock:
	{
		success, err := slf.redis.SetNX(slf.newKey(managerLock), lockId, 5*time.Second).Result()
		if err != nil {
			return nil, err
		}

		if !success {
			goto tryLock
		}
	}
	// 返回解锁函数
	return func() error {
		_, err := slf.redis.Del(slf.newKey(managerLock)).Result()
		return err
	}, nil
}

// 生成redis key
func (slf *managerRedis) newKey(key string) string {
	return redisKeyPrefix + key
}

// 将特定key前缀移除后进行返回
func (slf *managerRedis) formatKey(key string) string {
	return strings.ReplaceAll(key, redisKeyPrefix, "")
}

// 生成锁id
func (slf *managerRedis) newLookId() string {
	return fmt.Sprintf(uuid.NewV4().String())
}
