package resolver

import (
	"fmt"
	"github.com/tal-tech/go-zero/core/discov/discovk8s"
	"github.com/tal-tech/go-zero/core/logx"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/grpc/resolver"
)

type discovK8sBuilder struct {
	registry discovk8s.Registry
	once     sync.Once
}

func (d *discovK8sBuilder) parseTarget(target resolver.Target) (*discovk8s.Service, error) {
	// k8s://default/service:port
	end := target.Endpoint
	snamespace := target.Authority
	// k8s://service.default:port/
	if end == "" {
		end = target.Authority
		snamespace = ""
	}
	ti := discovk8s.Service{}
	if end == "" {
		return nil, fmt.Errorf("target(%q) is empty", target)
	}
	var name string
	var port string
	if strings.LastIndex(end, ":") < 0 {
		name = end
	} else {
		var err error
		name, port, err = net.SplitHostPort(end)
		if err != nil {
			return nil, fmt.Errorf("target endpoint='%s' is invalid. grpc target is %#v, err=%v", end, target, err)
		}
	}

	namesplit := strings.SplitN(name, ".", 2)
	sname := name
	if len(namesplit) == 2 {
		sname = namesplit[0]
		snamespace = namesplit[1]
	}
	ti.Name = sname
	ti.Namespace = snamespace

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	ti.Port = int32(intPort)

	return &ti, nil
}

func (d *discovK8sBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	d.once.Do(func() {
		k8sClient, err := discovk8s.NewK8sClient()
		if err != nil {
			logx.Error("New k8s client error", err)
			os.Exit(1)
		}

		epController, err := discovk8s.NewEndpointController(k8sClient)
		if err != nil {
			logx.Error("New endpoint controller error", err)
			os.Exit(1)
		}

		d.registry = discovk8s.NewK8sRegistry(epController)
	})

	si, err := d.parseTarget(target)

	if err != nil {
		logx.Errorf("parse k8s service error: %v", err)
		os.Exit(1)
	}

	sub := d.registry.NewSubscriber(si)

	update := func() {
		var addrs []resolver.Address
		for _, val := range subset(sub.Values(), subsetSize) {
			addrs = append(addrs, resolver.Address{
				Addr: val,
			})
		}
		cc.UpdateState(resolver.State{
			Addresses: addrs,
		})
	}
	sub.SetUpdateFunc(update)

	update()

	return &nopResolver{cc: cc}, nil
}

func (d *discovK8sBuilder) Scheme() string {
	return DiscovK8sScheme
}
