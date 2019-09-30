package dao

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

func (d *Dao) LPopOne(addr string) (p string, err error) {
	redisConn := d.Redis.GetRedisClientByAddr(addr).Get()
	defer func() {
		if err := redisConn.Close(); err != nil {
			log.Println("LPopOne redisConn.Close err:", err)
		}
	}()
	key := "consumer_list"
	if res, err := redis.String(redisConn.Do("lPop", key)); err != nil {
		return "", err
	} else {
		return res, nil
	}
}

func (d *Dao) GetSinglePlayerList(pid int) (p []string, err error) {
	redisConn := d.Redis.GetRedisPool(pid).Get()
	defer func() {
		if err := redisConn.Close(); err != nil {
			log.Println("GetSinglePlayerList redisConn.Close err:", err)
		}
	}()
	key := fmt.Sprintf("consumer:one:%d", pid)
	if res, err := redis.Strings(redisConn.Do("lRange", key, 0, 200)); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (d *Dao) UpdateSinglePlayerList(pid, length int, m []string) (err error) {
	redisConn := d.Redis.GetRedisPool(pid).Get()
	defer func() {
		if err := redisConn.Close(); err != nil {
			log.Println("GetSinglePlayerList redisConn.Close err:", err)
		}
	}()
	key := fmt.Sprintf("consumer:one:%d", pid)
	if _, err := redisConn.Do("multi"); err != nil {
		return
	}
	if _, err := redisConn.Do("lTrim", key, length, -1); err != nil {
		return
	}
	for _, v := range m {
		if _, err := redisConn.Do("rPush", key, v); err != nil {
			return
		}
	}
	if _, err := redisConn.Do("exec"); err != nil {
		return
	}
	return nil
}

func (d *Dao) GetAllPlayer() {

}
