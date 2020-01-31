package dao

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/nevercase/go-notification/config"
	"log"
	"time"
)

type RedisPool struct {
	RedisConf   [][]string
	Connections map[string]*redis.Pool
}

func InitRedisPool(c *config.Config) (r *RedisPool) {
	redisPool := &RedisPool{
		RedisConf:   c.GetRedisConfig(),
		Connections: make(map[string]*redis.Pool, len(c.GetRedisConfig())),
	}
	for _, v := range c.GetRedisConfig() {
		addr := fmt.Sprintf("%s:%s", v[0], v[1])
		redisPool.Connections[addr] = newPool(addr, v[2])
	}
	log.Println("redisPool:", redisPool)
	return redisPool
}

func newPool(addr string, pwd string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if pwd != "" {
				if _, err := c.Do("AUTH", pwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			//if _, err := c.Do("SELECT", db); err != nil {
			//	c.Close()
			//	return nil, err
			//}
			return c, nil
		},
	}
}

func (r *RedisPool) GetRedisPool(consumer int) *redis.Pool {
	remainder := consumer % len(r.RedisConf)
	t := r.RedisConf[remainder]
	addr := fmt.Sprintf("%s:%s", t[0], t[1])
	return r.Connections[addr]
}

func (r *RedisPool) GetRedisClientByAddr(addr string) *redis.Pool {
	return r.Connections[addr]
}
