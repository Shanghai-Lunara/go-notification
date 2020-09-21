package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Shanghai-Lunara/go-notification/config"
)

type InternalAPI struct {
	c *config.Config
}

func NewInternalAPI(conf *config.Config) Push {
	var api Push = &InternalAPI{
		c: conf,
	}
	return api
}

func (c *InternalAPI) Send(m *Message) (result bool, err error) {
	requestUrl := fmt.Sprintf("%s?pid=%d&cid=%s&ntype=%d", c.c.HttpRequestApi, m.Pid, m.Token, m.InfoType)
	start := time.Now()
	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Println("http.get err:", err)
		return false, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Send http Body.Close err:", err)
		}
	}()
	end := time.Now()
	log.Println(requestUrl, " request_time: ", end.Sub(start))
	return true, nil
}
