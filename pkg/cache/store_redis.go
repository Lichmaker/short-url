package cache

import (
	"shorturl/pkg/config"
	"shorturl/pkg/redis"
	"sync"
	"time"
)

// RedisStore 实现 cache.Store interface
type RedisStore struct {
	RedisClient *redis.RedisClient
	KeyPrefix   string
}

var redisOnce sync.Once
var instance *RedisStore

func NewRedisStore(address string, username string, password string, db int) *RedisStore {

	redisOnce.Do(func() {
		instance = &RedisStore{}
		instance.RedisClient = redis.NewClient(address, username, password, db)
		instance.KeyPrefix = config.GetString("app.name") + ":cache:"
	})

	// RedisStoreInstance := &RedisStore{}
	// RedisStoreInstance.RedisClient = redis.NewClient(address, username, password, db)
	// RedisStoreInstance.KeyPrefix = config.GetString("app.name") + ":cache:"
	// logger.DebugString("查看值", "redis store", fmt.Sprintf("%v", instance))
	return instance
}

func GetRedisClient() *RedisStore {
	return instance
}

func (s *RedisStore) Set(key string, value string, expireTime time.Duration) {
	s.RedisClient.Set(s.KeyPrefix+key, value, expireTime)
}

func (s *RedisStore) Get(key string) string {
	return s.RedisClient.Get(s.KeyPrefix + key)
}

func (s *RedisStore) Has(key string) bool {
	return s.RedisClient.Has(s.KeyPrefix + key)
}

func (s *RedisStore) Forget(key string) {
	s.RedisClient.Del((s.KeyPrefix + key))
}

func (s *RedisStore) Forever(key string, value string) {
	s.RedisClient.Set(s.KeyPrefix+key, value, 0)
}

func (s *RedisStore) Flush() {
	s.RedisClient.FlushDB()
}

func (s *RedisStore) Increment(parameters ...interface{}) {
	s.RedisClient.Increment(parameters...)
}

func (s *RedisStore) Decrement(parameters ...interface{}) {
	s.RedisClient.Decrement(parameters...)
}

func (s *RedisStore) IsAlive() error {
	return s.RedisClient.Ping()
}
