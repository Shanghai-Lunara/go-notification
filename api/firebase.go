package api

import (
	"firebase.google.com/go/messaging"
	"go-notification/config"
	"sync"
)

type FirebaseAPI struct {
	mu  sync.RWMutex
	c   *config.Config
	hub map[int]*messaging.Client
}

func NewFirebaseAPI(conf *config.Config) Push {
	var api Push = &FirebaseAPI{
		c:   conf,
		hub: make(map[int]*messaging.Client, 0),
	}
	return api
}

func (f *FirebaseAPI) NewClient() *messaging.Client {

}

func (f *FirebaseAPI) GetClient(workerId int) *messaging.Client {

}

func (f *FirebaseAPI) Send(workerId, pid, infoType int, token string) (result bool, err error) {
	return true, nil
}
