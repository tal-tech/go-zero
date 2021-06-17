//go:generate mockgen -package discovk8s -destination k8sendpointcontroller_mock.go -source k8sendpointcontroller.go EndpointController
package discovk8s

import (
	"errors"
	"github.com/tal-tech/go-zero/core/logx"
	v1 "k8s.io/api/core/v1"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"reflect"
	"sync"
)

type OnUpdateFunc func(addresses []*ServiceInstance)

type EndpointController interface {
	AddOnUpdateFunc(serviceName string, updateFunc OnUpdateFunc)
	GetEndpoints(serviceName string, namespace string) ([]*ServiceInstance, error)
}

type endpointController struct {
	endpointsInformer cache.SharedIndexInformer
	endpointsLister   listerv1.EndpointsLister
	updatefuncs       map[string][]OnUpdateFunc
	lock              sync.Mutex
}

func NewEndpointController(k8sClient *K8sClient) (EndpointController, error) {

	informer := k8sClient.InformerFactory.Core().V1().Endpoints().Informer()
	lister := k8sClient.InformerFactory.Core().V1().Endpoints().Lister()

	c := endpointController{
		endpointsInformer: informer,
		endpointsLister:   lister,
		updatefuncs:       make(map[string][]OnUpdateFunc),
	}

	informer.AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.endpointAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.endpointUpdate,
			// Called on resource deletion.
			DeleteFunc: c.endpointDelete,
		},
	)

	k8sClient.InformerFactory.Start(k8sClient.stop)

	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(k8sClient.stop, informer.HasSynced) {
		return nil, errors.New("Failed to sync k8s")
	}

	return &c, nil
}

func (e *endpointController) AddOnUpdateFunc(key string, updateFunc OnUpdateFunc) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, ok := e.updatefuncs[key]; !ok {
		e.updatefuncs[key] = make([]OnUpdateFunc, 1)
	}

	e.updatefuncs[key] = append(e.updatefuncs[key], updateFunc)
}

func (e *endpointController) GetEndpoints(name string, namespace string) ([]*ServiceInstance, error) {
	ep, err := e.endpointsLister.Endpoints(namespace).Get(name)
	if err != nil {
		logx.Errorf("get endpoints error, %v", err)
		return nil, err
	}

	return getReadyAddress(ep), nil
}

func (e *endpointController) endpointAdd(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)

	epName := buildEpName(endpoints.Name, endpoints.Namespace)
	if _, ok := e.updatefuncs[epName]; !ok {
		return
	}

	e.lock.Lock()
	funcs := append([]OnUpdateFunc(nil), e.updatefuncs[epName]...)
	e.lock.Unlock()

	for _, f := range funcs {
		if f != nil {
			f(getReadyAddress(endpoints))
		}
	}
}

func (e *endpointController) endpointUpdate(old, new interface{}) {
	oldEndpoint := old.(*v1.Endpoints)
	newEndpoint := new.(*v1.Endpoints)

	epName := buildEpName(oldEndpoint.Name, newEndpoint.Namespace)

	if _, ok := e.updatefuncs[epName]; !ok {
		return
	}

	oldAddress := getReadyAddress(oldEndpoint)
	newAddress := getReadyAddress(newEndpoint)

	if reflect.DeepEqual(oldAddress, newAddress) {
		return
	}

	e.lock.Lock()
	var funcs []OnUpdateFunc
	funcs = append(funcs, e.updatefuncs[epName]...)
	e.lock.Unlock()

	for _, f := range funcs {
		if f != nil {
			f(newAddress)
		}
	}

}

func (e *endpointController) endpointDelete(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)

	epName := buildEpName(endpoints.Name, endpoints.Namespace)
	if _, ok := e.updatefuncs[epName]; !ok {
		return
	}
	e.lock.Lock()
	funcs := append([]OnUpdateFunc(nil), e.updatefuncs[epName]...)
	e.lock.Unlock()

	for _, f := range funcs {
		if f != nil {
			f(nil)
		}
	}
}

func buildEpName(name string, namespace string) string {
	return name + "." + namespace
}
