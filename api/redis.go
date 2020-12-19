package api

import (
	"ditto/booking/utils"
	"errors"
	"strings"
)

//HGetAll -
func (s *Service) HGetAll(key string) (map[string]string, error) {
	client := s.Redis()
	if client != nil {
		return client.HGetAll(key).Result()
	}

	return nil, errors.New("Not connect REDIS")
}

//HGet -
func (s *Service) HGet(key, field string) (string, error) {
	client := s.Redis()
	if client != nil {
		return client.HGet(key, field).Result()
	}
	return "", errors.New("Not connect REDIS")
}

//HSet -
func (s *Service) HSet(key, field string, value interface{}) error {
	client := s.Redis()
	if client != nil {
		err := client.HSet(key, field, value).Err()
		if err != nil {
			return err
		}

		return nil
	}
	return errors.New("Not connect REDIS")
}

//CacheGet -
func (s *Service) CacheGet(key string, val interface{}) error {
	client := s.Redis()
	if client != nil {
		p, err := client.Get(key).Result()
		if err != nil {
			return err
		}
		err = utils.JSON.NewDecoder(strings.NewReader(p)).Decode(val)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("Not connect REDIS")
}

//CacheSet -
func (s *Service) CacheSet(key string, val interface{}) error {
	client := s.Redis()
	if client != nil {
		p, err := utils.JSON.MarshalToString(val)
		if err != nil {
			return err
		}

		return client.Set(key, p, 0).Err()
	}

	return errors.New("Not connect REDIS")
}

//CacheDel -
func (s *Service) CacheDel(key string) error {
	client := s.Redis()
	if client != nil {
		return client.Del(key).Err()
	}

	return errors.New("Not connect REDIS")
}

//CacheMultiDel -
func (s *Service) CacheMultiDel(pattern string) error {
	client := s.Redis()
	if client != nil {
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

	return errors.New("Not connect REDIS")
}
