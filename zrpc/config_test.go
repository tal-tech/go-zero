package zrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

func TestRpcClientConf(t *testing.T) {
	conf := NewDirectClientConf([]string{"localhost:1234"}, "foo", "bar")
	assert.True(t, conf.HasCredential())
	conf = NewEtcdClientConf([]string{"localhost:1234", "localhost:5678"}, "key", "foo", "bar")
	assert.True(t, conf.HasCredential())
	// ssl on
	conf = NewDirectClientConf([]string{"localhost:1234", "localhost:5678"}, "foo", "bar")
	assert.False(t, conf.HasSslVerify())
	conf.InsecureVerify = true
	assert.True(t, conf.HasSslVerify())
}

func TestRpcServerConf(t *testing.T) {
	conf := RpcServerConf{
		ServiceConf: service.ServiceConf{},
		ListenOn:    "",
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:1234"},
			Key:   "key",
		},
		Auth: true,
		Redis: redis.RedisKeyConf{
			RedisConf: redis.RedisConf{
				Type: redis.NodeType,
			},
			Key: "foo",
		},
		StrictControl: false,
		Timeout:       0,
		CpuThreshold:  0,
	}
	assert.True(t, conf.HasEtcd())
	assert.NotNil(t, conf.Validate())
	conf.Redis.Host = "localhost:5678"
	assert.Nil(t, conf.Validate())
}
