package dao

import (
	"go-notification/config"
	"sync"
)

type Dao struct {
	rwMutex sync.RWMutex
	c       *config.Config
	Redis   *RedisPool
}

func New(c *config.Config) (d *Dao) {
	d = &Dao{
		c:     c,
		Redis: InitRedisPool(c),
	}
	return d
}

func (d *Dao) Ping() {

}

func (d *Dao) Close() {

}
