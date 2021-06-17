package discovk8s

import (
	"sync"
)

type Registry interface {
	NewSubscriber(service *Service) Subscriber
}

type k8sRegistry struct {
	Registry

	k8sEpController EndpointController

	services map[string][]*ServiceInstance

	lock sync.Mutex
}

func NewK8sRegistry(k8sEpController EndpointController) Registry {

	registry := k8sRegistry{
		k8sEpController: k8sEpController,
		services:        make(map[string][]*ServiceInstance),
	}

	return &registry
}

func (r *k8sRegistry) NewSubscriber(service *Service) Subscriber {

	serviceFullName := service.EpName()

	k8sSubscriber := K8sSubscriber{
		k8sRegistry: r,
		instance:    service,
	}

	r.k8sEpController.AddOnUpdateFunc(serviceFullName, func(addresses []*ServiceInstance) {
		r.lock.Lock()
		r.services[serviceFullName] = addresses
		r.lock.Unlock()

		k8sSubscriber.OnUpdate()
	})

	return &k8sSubscriber
}

func (r *k8sRegistry) GetServices(service *Service) []*ServiceInstance {
	r.lock.Lock()
	defer r.lock.Unlock()

	if value, ok := r.services[service.EpName()]; ok {
		var ret []*ServiceInstance

		for _, item := range value {
			if item.Port == service.Port {
				ret = append(ret, item)
			}
		}
		return ret
	} else {
		fullService, _ := r.k8sEpController.GetEndpoints(service.Name, service.Namespace)
		var ret []*ServiceInstance

		for _, item := range fullService {
			if item.Port == service.Port {
				ret = append(ret, item)
			}
		}
		return ret
	}

}
