package zrpc

import (
	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

type (
	// A RpcServerConf is a rpc server config.
	RpcServerConf struct {
		service.ServiceConf
		ListenOn      string
		Etcd          discov.EtcdConf    `json:",optional"`
		Auth          bool               `json:",optional"`
		Redis         redis.RedisKeyConf `json:",optional"`
		StrictControl bool               `json:",optional"`
		// setting 0 means no timeout
		Timeout      int64 `json:",default=2000"`
		CpuThreshold int64 `json:",default=900,range=[0:1000]"`
	}

	// A RpcClientConf is a rpc client config.
	RpcClientConf struct {
		Etcd           discov.EtcdConf `json:",optional"`
		Endpoints      []string        `json:",optional=!Etcd"`
		App            string          `json:",optional"`
		Token          string          `json:",optional"`
		Timeout        int64           `json:",default=2000"`
		InsecureVerify bool            `json:",default=false"`
	}
)

// NewDirectClientConf returns a RpcClientConf.
func NewDirectClientConf(endpoints []string, app, token string) RpcClientConf {
	return RpcClientConf{
		Endpoints: endpoints,
		App:       app,
		Token:     token,
	}
}

// NewEtcdClientConf returns a RpcClientConf.
func NewEtcdClientConf(hosts []string, key, app, token string) RpcClientConf {
	return RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: hosts,
			Key:   key,
		},
		App:   app,
		Token: token,
	}
}

// HasEtcd checks if there is etcd settings in config.
func (sc RpcServerConf) HasEtcd() bool {
	return len(sc.Etcd.Hosts) > 0 && len(sc.Etcd.Key) > 0
}

// Validate validates the config.
func (sc RpcServerConf) Validate() error {
	if !sc.Auth {
		return nil
	}

	return sc.Redis.Validate()
}

// HasCredential checks if there is a credential in config.
func (cc RpcClientConf) HasCredential() bool {
	return len(cc.App) > 0 && len(cc.Token) > 0
}

//HasTls checks if there is a SSL in config.
func (cc RpcClientConf) HasSslVerify() bool {
	return cc.InsecureVerify
}
