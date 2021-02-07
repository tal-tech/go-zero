package internal

import "github.com/tal-tech/go-zero/core/discov"

func NewRpcPubServer(etcdEndpoints []string, etcdKey, listenOn string, opts ...ServerOption) (Server, error) {
	registerEtcd := func() error {
		pubClient := discov.NewPublisher(etcdEndpoints, etcdKey, listenOn)
		return pubClient.KeepAlive()
	}

	return newKeepAliveServer(registerEtcd,listenOn,opts...), nil
}

func NewRpcPubServerWithEtcdAuth(etcdEndpoints []string, user, pass, etcdKey, listenOn string, opts ...ServerOption) (Server, error) {
	registerEtcd := func() error {
		pubClient := discov.NewPublisherWithAuth(etcdEndpoints, user, pass, etcdKey, listenOn)
		return pubClient.KeepAliveWithAuth()
	}

	return newKeepAliveServer(registerEtcd,listenOn,opts...), nil
}

func newKeepAliveServer(registerEtcd func() error, listenOn string, opts ...ServerOption) keepAliveServer {
	return keepAliveServer{
		registerEtcd: registerEtcd,
		Server:       NewRpcServer(listenOn, opts...),
	}
}

type keepAliveServer struct {
	registerEtcd func() error
	Server
}

func (ags keepAliveServer) Start(fn RegisterFn) error {
	if err := ags.registerEtcd(); err != nil {
		return err
	}

	return ags.Server.Start(fn)
}
