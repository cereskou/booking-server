package utils

import (
	"errors"

	"github.com/go-redis/redis"
)

//RedisMultiDel -
func RedisMultiDel(client *redis.Client, pattern string) error {
	if client == nil {
		return errors.New("Redis client is null")
	}
	var cursor uint64
	var err error
	keys := make([]string, 0)
	for {
		var k []string
		if k, cursor, err = client.Scan(cursor, pattern, 100).Result(); err != nil {
			break
		}
		keys = append(keys, k...)

		if cursor == 0 {
			break
		}
	}

	pipe := client.Pipeline()
	for _, k := range keys {
		pipe.Del(k)
	}
	_, err = pipe.Exec()
	if err != nil {
		return err
	}

	return nil
}
