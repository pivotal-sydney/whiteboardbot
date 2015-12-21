package persistance

import (
	"github.com/garyburd/redigo/redis"
	"os"
	"fmt"
)

type Store interface {
	Get(key string) (value string, ok bool)
	Set(key string, value string)
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

func (store *RealStore) Get(key string) (value string, ok bool) {
	conn := store.Pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))
	ok = err == nil
	if !ok {
		fmt.Printf("Error occurred GETing from Redis: %v", err)
	}
	return
}

func (store *RealStore) Set(key string, value string) {
	conn := store.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return
	}
}
