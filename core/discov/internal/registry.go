package internal

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/contextx"
	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/syncx"
	"github.com/tal-tech/go-zero/core/threading"
	"go.etcd.io/etcd/clientv3"
)

var (
	registryInstance = Registry{
		clusters: make(map[string]*cluster),
	}
	connManager = syncx.NewResourceManager()
)

type Registry struct {
	clusters map[string]*cluster
	lock     sync.Mutex
}

func GetRegistry() *Registry {
	return &registryInstance
}

func (r *Registry) getCluster(endpoints []string) *cluster {
	clusterKey := getClusterKey(endpoints)
	r.lock.Lock()
	defer r.lock.Unlock()
	c, ok := r.clusters[clusterKey]
	if !ok {
		c = newCluster(endpoints)
		r.clusters[clusterKey] = c
	}

	return c
}

func (r *Registry) getClusterWithAuth(endpoints []string, user, pass string) *cluster {
	clusterKey := getClusterKey(endpoints)
	r.lock.Lock()
	defer r.lock.Unlock()
	c, ok := r.clusters[clusterKey]
	if !ok {
		c = newClusterWithAuth(endpoints, user, pass)
		r.clusters[clusterKey] = c
	}

	return c
}

func (r *Registry) GetConn(endpoints []string) (EtcdClient, error) {
	return r.getCluster(endpoints).getClient()
}

func (r *Registry) GetConnWithAuth(endpoints []string, user, pass string) (EtcdClient, error) {
	return r.getClusterWithAuth(endpoints, user, pass).getClient()
}

func (r *Registry) Monitor(endpoints []string, key string, l UpdateListener) error {
	return r.getCluster(endpoints).monitor(key, l)
}

func (r *Registry) MonitorWithAuth(endpoints []string, user, pass, key string, l UpdateListener) error {
	return r.getClusterWithAuth(endpoints, user, pass).monitor(key, l)
}

type cluster struct {
	endpoints  []string
	user       string
	pass       string
	key        string
	values     map[string]map[string]string
	listeners  map[string][]UpdateListener
	watchGroup *threading.RoutineGroup
	done       chan lang.PlaceholderType
	lock       sync.Mutex
}

func newCluster(endpoints []string) *cluster {
	return &cluster{
		endpoints:  endpoints,
		key:        getClusterKey(endpoints),
		values:     make(map[string]map[string]string),
		listeners:  make(map[string][]UpdateListener),
		watchGroup: threading.NewRoutineGroup(),
		done:       make(chan lang.PlaceholderType),
	}
}

func newClusterWithAuth(endpoints []string, user, pass string) *cluster {
	return &cluster{
		endpoints:  endpoints,
		user:       user,
		pass:       pass,
		key:        getClusterKey(endpoints),
		values:     make(map[string]map[string]string),
		listeners:  make(map[string][]UpdateListener),
		watchGroup: threading.NewRoutineGroup(),
		done:       make(chan lang.PlaceholderType),
	}
}

func (c *cluster) context(cli EtcdClient) context.Context {
	return contextx.ValueOnlyFrom(cli.Ctx())
}

func (c *cluster) getClient() (EtcdClient, error) {
	val, err := connManager.GetResource(c.key, func() (io.Closer, error) {
		return c.newClient()
	})
	if err != nil {
		return nil, err
	}

	return val.(EtcdClient), nil
}

func (c *cluster) handleChanges(key string, kvs []KV) {
	var add []KV
	var remove []KV
	c.lock.Lock()
	listeners := append([]UpdateListener(nil), c.listeners[key]...)
	vals, ok := c.values[key]
	if !ok {
		add = kvs
		vals = make(map[string]string)
		for _, kv := range kvs {
			vals[kv.Key] = kv.Val
		}
		c.values[key] = vals
	} else {
		m := make(map[string]string)
		for _, kv := range kvs {
			m[kv.Key] = kv.Val
		}
		for k, v := range vals {
			if val, ok := m[k]; !ok || v != val {
				remove = append(remove, KV{
					Key: k,
					Val: v,
				})
			}
		}
		for k, v := range m {
			if val, ok := vals[k]; !ok || v != val {
				add = append(add, KV{
					Key: k,
					Val: v,
				})
			}
		}
		c.values[key] = m
	}
	c.lock.Unlock()

	for _, kv := range add {
		for _, l := range listeners {
			l.OnAdd(kv)
		}
	}
	for _, kv := range remove {
		for _, l := range listeners {
			l.OnDelete(kv)
		}
	}
}

func (c *cluster) handleWatchEvents(key string, events []*clientv3.Event) {
	c.lock.Lock()
	listeners := append([]UpdateListener(nil), c.listeners[key]...)
	c.lock.Unlock()

	for _, ev := range events {
		switch ev.Type {
		case clientv3.EventTypePut:
			c.lock.Lock()
			if vals, ok := c.values[key]; ok {
				vals[string(ev.Kv.Key)] = string(ev.Kv.Value)
			} else {
				c.values[key] = map[string]string{string(ev.Kv.Key): string(ev.Kv.Value)}
			}
			c.lock.Unlock()
			for _, l := range listeners {
				l.OnAdd(KV{
					Key: string(ev.Kv.Key),
					Val: string(ev.Kv.Value),
				})
			}
		case clientv3.EventTypeDelete:
			if vals, ok := c.values[key]; ok {
				delete(vals, string(ev.Kv.Key))
			}
			for _, l := range listeners {
				l.OnDelete(KV{
					Key: string(ev.Kv.Key),
					Val: string(ev.Kv.Value),
				})
			}
		default:
			logx.Errorf("Unknown event type: %v", ev.Type)
		}
	}
}

func (c *cluster) load(cli EtcdClient, key string) {
	var resp *clientv3.GetResponse
	for {
		var err error
		ctx, cancel := context.WithTimeout(c.context(cli), RequestTimeout)
		resp, err = cli.Get(ctx, makeKeyPrefix(key), clientv3.WithPrefix())
		cancel()
		if err == nil {
			break
		}

		logx.Error(err)
		time.Sleep(coolDownInterval)
	}

	var kvs []KV
	c.lock.Lock()
	for _, ev := range resp.Kvs {
		kvs = append(kvs, KV{
			Key: string(ev.Key),
			Val: string(ev.Value),
		})
	}
	c.lock.Unlock()

	c.handleChanges(key, kvs)
}

func (c *cluster) monitor(key string, l UpdateListener) error {
	c.lock.Lock()
	c.listeners[key] = append(c.listeners[key], l)
	c.lock.Unlock()

	cli, err := c.getClient()
	if err != nil {
		return err
	}

	c.load(cli, key)
	c.watchGroup.Run(func() {
		c.watch(cli, key)
	})

	return nil
}

func (c *cluster) newClient() (EtcdClient, error) {
	cli, err := NewClient(c.endpoints, c.user, c.pass)
	if err != nil {
		return nil, err
	}

	go c.watchConnState(cli)

	return cli, nil
}

func (c *cluster) reload(cli EtcdClient) {
	c.lock.Lock()
	close(c.done)
	c.watchGroup.Wait()
	c.done = make(chan lang.PlaceholderType)
	c.watchGroup = threading.NewRoutineGroup()
	var keys []string
	for k := range c.listeners {
		keys = append(keys, k)
	}
	c.lock.Unlock()

	for _, key := range keys {
		k := key
		c.watchGroup.Run(func() {
			c.load(cli, k)
			c.watch(cli, k)
		})
	}
}

func (c *cluster) watch(cli EtcdClient, key string) {
	rch := cli.Watch(clientv3.WithRequireLeader(c.context(cli)), makeKeyPrefix(key), clientv3.WithPrefix())
	for {
		select {
		case wresp, ok := <-rch:
			if !ok {
				logx.Error("etcd monitor chan has been closed")
				return
			}
			if wresp.Canceled {
				logx.Error("etcd monitor chan has been canceled")
				return
			}
			if wresp.Err() != nil {
				logx.Error(fmt.Sprintf("etcd monitor chan error: %v", wresp.Err()))
				return
			}

			c.handleWatchEvents(key, wresp.Events)
		case <-c.done:
			return
		}
	}
}

func (c *cluster) watchConnState(cli EtcdClient) {
	watcher := newStateWatcher()
	watcher.addListener(func() {
		go c.reload(cli)
	})
	watcher.watch(cli.ActiveConnection())
}

func DialClient(endpoints []string, user, pass string) (EtcdClient, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:            endpoints,
		AutoSyncInterval:     autoSyncInterval,
		DialTimeout:          DialTimeout,
		DialKeepAliveTime:    dialKeepAliveTime,
		DialKeepAliveTimeout: DialTimeout,
		RejectOldCluster:     true,
		Username:             user,
		Password:             pass,
	})
}

func getClusterKey(endpoints []string) string {
	sort.Strings(endpoints)
	return strings.Join(endpoints, endpointsSeparator)
}

func makeKeyPrefix(key string) string {
	return fmt.Sprintf("%s%c", key, Delimiter)
}
