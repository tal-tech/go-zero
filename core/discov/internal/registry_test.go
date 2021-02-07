package internal

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/contextx"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stringx"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

var mockLock sync.Mutex

func init() {
	logx.Disable()
}

func setMockClient(cli EtcdClient) func() {
	mockLock.Lock()
	NewClient = func([]string, string, string) (EtcdClient, error) {
		return cli, nil
	}
	return func() {
		NewClient = DialClient
		mockLock.Unlock()
	}
}

func TestGetCluster(t *testing.T) {
	c1 := GetRegistry().getCluster([]string{"first"})
	c2 := GetRegistry().getCluster([]string{"second"})
	c3 := GetRegistry().getCluster([]string{"first"})
	assert.Equal(t, c1, c3)
	assert.NotEqual(t, c1, c2)
}

func TestGetClusterKey(t *testing.T) {
	assert.Equal(t, getClusterKey([]string{"localhost:1234", "remotehost:5678"}),
		getClusterKey([]string{"remotehost:5678", "localhost:1234"}))
}

func TestCluster_HandleChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	l := NewMockUpdateListener(ctrl)
	l.EXPECT().OnAdd(KV{
		Key: "first",
		Val: "1",
	})
	l.EXPECT().OnAdd(KV{
		Key: "second",
		Val: "2",
	})
	l.EXPECT().OnDelete(KV{
		Key: "first",
		Val: "1",
	})
	l.EXPECT().OnDelete(KV{
		Key: "second",
		Val: "2",
	})
	l.EXPECT().OnAdd(KV{
		Key: "third",
		Val: "3",
	})
	l.EXPECT().OnAdd(KV{
		Key: "fourth",
		Val: "4",
	})
	c := newCluster([]string{"any"})
	c.listeners["any"] = []UpdateListener{l}
	c.handleChanges("any", []KV{
		{
			Key: "first",
			Val: "1",
		},
		{
			Key: "second",
			Val: "2",
		},
	})
	assert.EqualValues(t, map[string]string{
		"first":  "1",
		"second": "2",
	}, c.values["any"])
	c.handleChanges("any", []KV{
		{
			Key: "third",
			Val: "3",
		},
		{
			Key: "fourth",
			Val: "4",
		},
	})
	assert.EqualValues(t, map[string]string{
		"third":  "3",
		"fourth": "4",
	}, c.values["any"])
}

func TestCluster_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	cli.EXPECT().Get(gomock.Any(), "any/", gomock.Any()).Return(&clientv3.GetResponse{
		Kvs: []*mvccpb.KeyValue{
			{
				Key:   []byte("hello"),
				Value: []byte("world"),
			},
		},
	}, nil)
	cli.EXPECT().Ctx().Return(context.Background())
	c := &cluster{
		values: make(map[string]map[string]string),
	}
	c.load(cli, "any")
}

func TestCluster_Watch(t *testing.T) {
	tests := []struct {
		name      string
		method    int
		eventType mvccpb.Event_EventType
	}{
		{
			name:      "add",
			eventType: clientv3.EventTypePut,
		},
		{
			name:      "delete",
			eventType: clientv3.EventTypeDelete,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cli := NewMockEtcdClient(ctrl)
			restore := setMockClient(cli)
			defer restore()
			ch := make(chan clientv3.WatchResponse)
			cli.EXPECT().Watch(gomock.Any(), "any/", gomock.Any()).Return(ch)
			cli.EXPECT().Ctx().Return(context.Background())
			var wg sync.WaitGroup
			wg.Add(1)
			c := &cluster{
				listeners: make(map[string][]UpdateListener),
				values:    make(map[string]map[string]string),
			}
			listener := NewMockUpdateListener(ctrl)
			c.listeners["any"] = []UpdateListener{listener}
			listener.EXPECT().OnAdd(gomock.Any()).Do(func(kv KV) {
				assert.Equal(t, "hello", kv.Key)
				assert.Equal(t, "world", kv.Val)
				wg.Done()
			}).MaxTimes(1)
			listener.EXPECT().OnDelete(gomock.Any()).Do(func(_ interface{}) {
				wg.Done()
			}).MaxTimes(1)
			go c.watch(cli, "any")
			ch <- clientv3.WatchResponse{
				Events: []*clientv3.Event{
					{
						Type: test.eventType,
						Kv: &mvccpb.KeyValue{
							Key:   []byte("hello"),
							Value: []byte("world"),
						},
					},
				},
			}
			wg.Wait()
		})
	}
}

func TestClusterWatch_RespFailures(t *testing.T) {
	resps := []clientv3.WatchResponse{
		{
			Canceled: true,
		},
		{
			// cause resp.Err() != nil
			CompactRevision: 1,
		},
	}

	for _, resp := range resps {
		t.Run(stringx.Rand(), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cli := NewMockEtcdClient(ctrl)
			restore := setMockClient(cli)
			defer restore()
			ch := make(chan clientv3.WatchResponse)
			cli.EXPECT().Watch(gomock.Any(), "any/", gomock.Any()).Return(ch)
			cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
			c := new(cluster)
			go func() {
				ch <- resp
			}()
			c.watch(cli, "any")
		})
	}
}

func TestClusterWatch_CloseChan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := NewMockEtcdClient(ctrl)
	restore := setMockClient(cli)
	defer restore()
	ch := make(chan clientv3.WatchResponse)
	cli.EXPECT().Watch(gomock.Any(), "any/", gomock.Any()).Return(ch)
	cli.EXPECT().Ctx().Return(context.Background()).AnyTimes()
	c := new(cluster)
	go func() {
		close(ch)
	}()
	c.watch(cli, "any")
}

func TestValueOnlyContext(t *testing.T) {
	ctx := contextx.ValueOnlyFrom(context.Background())
	ctx.Done()
	assert.Nil(t, ctx.Err())
}
