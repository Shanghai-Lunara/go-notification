package api

import (
	"context"
	"errors"
	fb "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"go-notification/config"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"google.golang.org/api/option"
	"log"
	"strings"
	"sync"
	"time"
)

type FirebaseAPI struct {
	mu          sync.RWMutex
	c           *config.Config
	hub         map[int]*messaging.Client
	servicePath string
}

func NewFirebaseAPI(conf *config.Config) Push {
	servicePath := strings.Replace(conf.ConfigPath, "push.yml", "service_account.json", 1)
	var api Push = &FirebaseAPI{
		c:           conf,
		hub:         make(map[int]*messaging.Client, 0),
		servicePath: servicePath,
	}

	return api
}

func (f *FirebaseAPI) NewClient() (c *messaging.Client, err error) {
	app, err := fb.NewApp(context.Background(), nil, option.WithCredentialsFile(f.servicePath))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := app.Auth(ctx)
	cancel()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("app.Auth err:%v", err))
	}
	log.Println("client:", client)
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	m, err := app.Messaging(ctx)
	cancel()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Messaging() = (%v, %v); want (iid, nil)", m, err))
	}
	return m, nil
}

func (f *FirebaseAPI) GetClient(m *Message) (c *messaging.Client, err error) {
	if t, ok := f.hub[m.WorkerId]; ok {
		return t, nil
	}
	if c, err = f.NewClient(); err != nil {
		return nil, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.hub[m.WorkerId] = c
	return c, nil
}

func (f *FirebaseAPI) Send(m *Message) (result bool, err error) {
	var c *messaging.Client
	if c, err = f.GetClient(m); err != nil {
		return false, err
	}
	a := &messaging.Message{
		Notification: &messaging.Notification{
			Title: m.Title,
			Body:  m.Body,
		},
		Token: m.Token,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := c.Send(ctx, a)
	cancel()
	if err != nil {
		return false, errors.New(fmt.Sprintf("workerId:", m.WorkerId, " Send err:", err))
	}
	log.Println("workerId:", m.WorkerId, " pid:", m.Pid, " res:", res)
	return true, nil
}
