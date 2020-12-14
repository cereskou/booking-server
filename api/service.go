package api

import (
	"ditto/booking/config"
	"ditto/booking/db"
	"ditto/booking/logger"
	"ditto/booking/rsa"

	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis"
)

//Service -
type Service struct {
	_db       *db.Database
	_rsa      *rsa.RSA
	_client   *redis.Client
	_enforcer *casbin.Enforcer
}

//New -
func New(db *db.Database, r *rsa.RSA, client *redis.Client) (*Service, error) {
	logger.Debug("Service New")
	s := &Service{
		_db:     db,
		_rsa:    r,
		_client: client,
	}
	e, err := casbin.NewEnforcer("auth.conf")
	if err != nil {
		return nil, err
	}
	s._enforcer = e

	//Load casbin policies
	policies, err := db.GetCasbinPolicies()
	if err != nil {
		return nil, err
	}

	conf := config.Load()
	//add policy to casbin enforcer
	for _, p := range policies {
		path := p.Path
		if path[0] == '/' {
			path = conf.BaseURL + p.Path
		} else {
			path = conf.BaseURL + "/" + p.Path
		}

		e.AddPolicy(p.Role, path, p.Method)
		logger.Tracef("AddPolicy role:%v, path:%v, method:%v", p.Role, path, p.Method)
	}

	// go func() {
	// 	if s._client == nil {
	// 		return
	// 	}

	// 	count := 0
	// 	ticker := time.Tick(3 * time.Second)
	// 	for now := range ticker {
	// 		//check the server is alive
	// 		if err := client.Ping().Err(); err != nil {
	// 			logger.Error(now, err)
	// 			count++
	// 			if count > 3 {
	// 				s._client = nil
	// 				break
	// 			}
	// 		}
	// 	}
	// }()

	return s, nil
}

//Close -
func (s *Service) Close() {
	logger.Debug("Service Close()")
}

//DB -
func (s *Service) DB() *db.Database {
	return s._db
}

//Redis -
func (s *Service) Redis() *redis.Client {
	return s._client
}
