package persistance

import (
	"github.com/garyburd/redigo/redis"
	"os"
)

type Store interface {
	Get(key string) (value int, ok bool)
	Set(key string, value int)
}

type RealStore struct{}

func (store RealStore) Get(key string) (value int, ok bool) {

	client, err := redis.Dial("tcp", os.Getenv("WB_DB_HOST"), redis.DialPassword(os.Getenv("WB_DB_PASSWORD")))
	if err != nil {
		ok = false
	}
	defer client.Close()

	value, err = redis.Int(client.Do("GET", key))
	ok = err == nil
	return
}

func (store RealStore) Set(key string, value int) {
	client, err := redis.Dial("tcp", os.Getenv("WB_DB_HOST"), redis.DialPassword(os.Getenv("WB_DB_PASSWORD")))
	if err != nil {
		return
	}
	defer client.Close()

	_, err = client.Do("SET", key, value)
	if err != nil {
		return
	}
}
