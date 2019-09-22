package common

import (
	"context"
	"fmt"
	"go-notification/config"
	"go-notification/dao"
	"log"
	"net"
	"time"

	pb "go-notification/service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
)

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}

type server struct {
	service *Service
}

func newRegisterServer(s *Service) *server {
	return &server{service: s}
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		log.Println("addr:", p.Addr)
		log.Println("AuthInfo:", p.AuthInfo)
	} else {
		log.Println("peer.FromContext !ok")
	}
	id := s.service.handleInit(p.Addr, req.ExternalAddr)
	return &pb.RegisterResponse{Id: id, Message: req.ExternalAddr}, nil
}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	s.service.handlePing(req.Id)
	return &pb.PongResponse{Message: "pong"}, nil
}

func (s *Service) initRpcService() {
	addr := fmt.Sprintf("%s:%d", s.c.Dispatch.IP, s.c.Dispatch.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	registerServer := newRegisterServer(s)
	pb.RegisterGatewayServer(server, registerServer)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type Service struct {
	c       *config.Config
	dao     *dao.Dao
	hub     *Hub
	nodeHub *NodeHub
	ctx     context.Context
	cancel  context.CancelFunc
}

func New(conf *config.Config) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		c:      conf,
		dao:    dao.New(conf),
		hub:    NewHub(conf, ctx),
		ctx:    ctx,
		cancel: cancel,
	}
	s.nodeHub = s.NewNodeHub()
	go s.initRpcService()
	return s
}

func (s *Service) Close() {
	s.dao.Close()
	s.cancel()
}
