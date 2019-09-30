package common

import (
	"context"
	"errors"
	"fmt"
	"go-notification/config"
	"go-notification/dao"
	"log"
	"sync"
	"time"

	pb "go-notification/service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

type RpcClient struct {
	mutex         sync.Mutex
	conn          *grpc.ClientConn
	gatewayClient pb.GatewayClient
	id            int32
	ctx           context.Context
	cancel        context.CancelFunc
	closeChan     chan int
	status        int
}

const (
	RpcClosed = iota
	RpcAlive
)

func (s *Service) initRpcClient(conf *config.Config) (err error) {
	if s.rpcClient.conn != nil {
		if err := s.rpcClient.conn.Close(); err != nil {
			log.Printf("maintainRpcClient initRpcClient conn.Close err:")
		}
	}
	addr := fmt.Sprintf("%s:%d", conf.Dispatch.InternalIP, conf.Dispatch.Port)
	s.rpcClient.conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		return errors.New(fmt.Sprintf("grpc.Dial err: %v", err))
	}
	c := pb.NewGatewayClient(s.rpcClient.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	res, err := c.Register(ctx, &pb.RegisterRequest{Message: time.Now().String(), ExternalAddr: ""})
	if err != nil {
		return errors.New(fmt.Sprintf("unexpected error from Register: %v", err))
	}
	cancel()
	s.rpcClient.status = RpcAlive
	s.rpcClient.gatewayClient = c
	s.rpcClient.id = res.Id
	return nil
}

func (s *Service) maintainRpcClient() {
	s.rpcClient = &RpcClient{
		closeChan: make(chan int),
		status:    RpcClosed,
	}
	s.rpcClient.ctx, s.rpcClient.cancel = context.WithCancel(s.ctx)
	go func() {
		s.rpcClient.closeChan <- 1
	}()
	go s.rpcPing()
	for {
		select {
		case <-s.rpcClient.ctx.Done():
			return
		case <-s.rpcClient.closeChan:
			if s.rpcClient.status == RpcAlive {
				continue
			}
			if err := s.initRpcClient(s.c); err != nil {
				log.Printf("maintainRpcClient initRpcClient err: %v", err)
				continue
			}
		}
	}
}

func (s *Service) rpcPing() {
	tick := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-s.rpcClient.ctx.Done():
			return
		case <-tick.C:
			if s.rpcClient.status == RpcClosed {
				log.Printf("rpcPing RpcClosed")
				s.rpcClose()
				continue
			}
			_, err := s.rpcClient.gatewayClient.Ping(s.ctx, &pb.PingRequest{Id: s.rpcClient.id})
			if err != nil {
				log.Printf("rpcPing unexpected error from ping: %v", err)
				s.rpcClose()
				continue
			}
		}
	}
}

func (s *Service) rpcClose() {
	s.rpcClient.status = RpcClosed
	s.rpcClient.closeChan <- 1
}

func (s *Service) getAllocatedNode() (addr string, err error) {
	if res, err := s.rpcClient.gatewayClient.GetAllocatedNode(s.ctx, &pb.CommonRequest{Id: s.rpcClient.id, Addr: ""}); err != nil {
		return "", err
	} else {
		return res.Addr, err
	}
}

type Service struct {
	c         *config.Config
	dao       *dao.Dao
	rpcClient *RpcClient
	workers   *Workers
	ctx       context.Context
	cancel    context.CancelFunc
}

func New(conf *config.Config) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		c:      conf,
		dao:    dao.New(conf),
		ctx:    ctx,
		cancel: cancel,
	}
	s.workers = s.NewWorkers()
	return s
}

func (s *Service) Close() {
	s.dao.Close()
	s.cancel()
}
