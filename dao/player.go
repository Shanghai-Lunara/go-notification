package dao

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
	"strings"
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
	if res, err := redis.Strings(redisConn.Do("lRange", key, 0, -1)); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (d *Dao) UpdateSinglePlayerList(pid int, m map[string]int) (err error) {
	redisConn := d.Redis.GetRedisPool(pid).Get()
	defer func() {
		if err := redisConn.Close(); err != nil {
			log.Println("GetSinglePlayerList redisConn.Close err:", err)
		}
	}()
	key := fmt.Sprintf("consumer:one:%d", pid)
	if _, err := redisConn.Do("multi"); err != nil {
		return err
	}
	for k, v := range m {
		if _, err := redisConn.Do("lRem", key, k, v); err != nil {
			return err
		}
	}
	if _, err := redisConn.Do("exec"); err != nil {
		return err
	}
	return nil
}

func (d *Dao) GetPlayerSettings(pid int) (cid string, close int, err error) {
	redisConn := d.Redis.GetRedisPool(pid).Get()
	defer func() {
		if err := redisConn.Close(); err != nil {
			log.Println("GetPlayerSettings redisConn.Close err:", err)
		}
	}()
	var (
		res map[string]string
	)
	settingKey := fmt.Sprintf("push_setting:%d", pid)
	if res, err = redis.StringMap(redisConn.Do("hGetAll", settingKey)); err != nil {
		return "", 0, err
	}
	if t, ok := res["cid"]; ok {
		cid = t
	}
	if settings, ok := res["settings"]; ok {
		tmp := strings.Split(settings, ",")
		if len(tmp) >= 7 {
			tmp2, err := strconv.Atoi(tmp[6])
			if err != nil {
				return cid, 0, err
			}
			close = tmp2
		}
	}
	return cid, close, nil
}

func (d *Dao) GetAllPlayer() {

}
