package persistance

import (
	"github.com/garyburd/redigo/redis"
	"os"
)

type Store interface {
	Get(key string) (value int, ok bool)
	Set(key string, value int)
}

type RealStore struct{
	Pool *redis.Pool
}

func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 10,
		MaxActive: 50,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("WB_DB_HOST"), redis.DialPassword(os.Getenv("WB_DB_PASSWORD")))
		},
	}
}

func (store *RealStore) Get(key string) (value int, ok bool) {
	conn := store.Pool.Get()
	defer conn.Close()

	value, err := redis.Int(conn.Do("GET", key))
	ok = err == nil
	return
}

func (store *RealStore) Set(key string, value int) {
	conn := store.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return
	}
}
