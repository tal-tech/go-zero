package zrpc

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/load"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/netx"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/zrpc/internal"
	"github.com/tal-tech/go-zero/zrpc/internal/auth"
	"github.com/tal-tech/go-zero/zrpc/internal/serverinterceptors"
	"google.golang.org/grpc"
)

const (
	allEths  = "0.0.0.0"
	envPodIp = "POD_IP"
)

type RpcServer struct {
	server   internal.Server
	register internal.RegisterFn
}

func MustNewServer(c RpcServerConf, register internal.RegisterFn) *RpcServer {
	server, err := NewServer(c, register)
	if err != nil {
		log.Fatal(err)
	}

	return server
}

func NewServer(c RpcServerConf, register internal.RegisterFn) (*RpcServer, error) {
	var err error
	if err = c.Validate(); err != nil {
		return nil, err
	}

	var server internal.Server
	metrics := stat.NewMetrics(c.ListenOn)
	if c.HasEtcd() {
		listenOn := figureOutListenOn(c.ListenOn)

		if c.Etcd.EnableAuth(){
			server, err = internal.NewRpcPubServerWithEtcdAuth(c.Etcd.Hosts, c.Etcd.User, c.Etcd.Pass, c.Etcd.Key, listenOn, internal.WithMetrics(metrics))
		}else {
			server, err = internal.NewRpcPubServer(c.Etcd.Hosts, c.Etcd.Key, listenOn, internal.WithMetrics(metrics))
		}

		if err != nil {
			return nil, err
		}
	} else {
		server = internal.NewRpcServer(c.ListenOn, internal.WithMetrics(metrics))
	}

	server.SetName(c.Name)
	if err = setupInterceptors(server, c, metrics); err != nil {
		return nil, err
	}

	rpcServer := &RpcServer{
		server:   server,
		register: register,
	}
	if err = c.SetUp(); err != nil {
		return nil, err
	}

	return rpcServer, nil
}

func (rs *RpcServer) AddOptions(options ...grpc.ServerOption) {
	rs.server.AddOptions(options...)
}

func (rs *RpcServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	rs.server.AddStreamInterceptors(interceptors...)
}

func (rs *RpcServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	rs.server.AddUnaryInterceptors(interceptors...)
}

func (rs *RpcServer) Start() {
	if err := rs.server.Start(rs.register); err != nil {
		logx.Error(err)
		panic(err)
	}
}

func (rs *RpcServer) Stop() {
	logx.Close()
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	} else {
		return strings.Join(append([]string{ip}, fields[1:]...), ":")
	}
}

func setupInterceptors(server internal.Server, c RpcServerConf, metrics *stat.Metrics) error {
	if c.CpuThreshold > 0 {
		shedder := load.NewAdaptiveShedder(load.WithCpuThreshold(c.CpuThreshold))
		server.AddUnaryInterceptors(serverinterceptors.UnarySheddingInterceptor(shedder, metrics))
	}

	if c.Timeout > 0 {
		server.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(
			time.Duration(c.Timeout) * time.Millisecond))
	}

	if c.Auth {
		authenticator, err := auth.NewAuthenticator(c.Redis.NewRedis(), c.Redis.Key, c.StrictControl)
		if err != nil {
			return err
		}

		server.AddStreamInterceptors(serverinterceptors.StreamAuthorizeInterceptor(authenticator))
		server.AddUnaryInterceptors(serverinterceptors.UnaryAuthorizeInterceptor(authenticator))
	}

	return nil
}
