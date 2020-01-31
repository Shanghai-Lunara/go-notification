package main

import (
	"flag"
	"fmt"
	"github.com/nevercase/go-notification/config"
	"github.com/nevercase/go-notification/dao"
	"log"
	"time"
)

var (
	d *dao.Dao
)

func main() {
	flag.Parse()

	if err := config.Init(); err != nil {
		log.Fatal(err)
	}
	d = dao.New(config.GetConfig())

	ConsumerSettings(1)
	time.Sleep(time.Second * 2)

	ConsumerInfo()
	ConsumerList()

	ConsumerSettings(0)
	time.Sleep(time.Second * 2)
	//ConsumerSettings2(1)
	time.Sleep(time.Second * 2)
	ConsumerList()
	ConsumerInfo()
	//time.Sleep(time.Second * 2)
	//ConsumerSettings(0)
	//time.Sleep(time.Second * 10)
	//ConsumerSettings(1)
}

const (
	//start = 1000
	//end   = 1100
	start = 1
	end   = 100
)

func ConsumerList() {

	for i := start; i <= end; i++ {
		appendOne(i)
	}

}

func appendOne(consumerId int) {

	redisInstance := d.Redis.GetRedisPool(consumerId).Get()
	defer redisInstance.Close()

	key := "consumer_list"

	_, err := redisInstance.Do("rpush", key, consumerId)
	if err != nil {
		log.Println("rpush err:", err)
	}
}

func ConsumerInfo() {
	for i := start; i <= end; i++ {
		consumerOne(i)
	}
}

func consumerOne(consumerId int) {
	redisInstance := d.Redis.GetRedisPool(consumerId).Get()
	defer redisInstance.Close()

	key := fmt.Sprintf("%s:%d", "consumer:one", consumerId)

	now := int(time.Now().Unix())

	log.Println("now:", now)

	for j := 1; j <= 6; j++ {
		for i := 0; i < 3; i++ {

			t := fmt.Sprintf("%d:%d:%d:%d:%d", consumerId, j, now+i*10+j*3, 0, 0)

			_, err := redisInstance.Do("rpush", key, t)
			if err != nil {
				log.Println("Redis rpush err:", err)
			}
		}
	}
}

func ConsumerSettings(close int) {
	for i := start; i <= end; i++ {
		changeOne(i, close)
	}
}

func changeOne(consumerId, close int) {
	redisInstance := d.Redis.GetRedisPool(consumerId).Get()
	defer redisInstance.Close()

	key := fmt.Sprintf("%s:%d", "push_setting", consumerId)
	settings := fmt.Sprintf("1,1,1,1,1,1,%d", close)
	_, err := redisInstance.Do("hSet", key, "settings", settings)
	if err != nil {
		log.Println("rpush err:", err)
	}

	key2 := "consumer_list"
	_, err = redisInstance.Do("rpush", key2, fmt.Sprintf("%d:%d", consumerId, 0))
	if err != nil {
		log.Println("rpush err:", err)
	}
}

func ConsumerSettings2(close int) {
	for i := start; i <= end; i++ {
		changeOne2(i, close)
	}
}

func changeOne2(consumerId, close int) {
	redisInstance := d.Redis.GetRedisPool(consumerId).Get()
	defer redisInstance.Close()

	key := fmt.Sprintf("%s:%d", "push_setting", consumerId)
	settings := fmt.Sprintf("1,1,1,1,1,1,%d", close)
	_, err := redisInstance.Do("hSet", key, "settings", settings)
	if err != nil {
		log.Println("rpush err:", err)
	}
}
