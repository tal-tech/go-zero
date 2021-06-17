package discovk8s

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"sync"
	"testing"
)

var (
	services []*ServiceInstance
)

func TestMain(m *testing.M) {
	services = []*ServiceInstance{
		{
			Ip:   "172.16.1.2",
			Port: 8080,
		},
		{
			Ip:   "172.16.1.3",
			Port: 8080,
		},
	}
	m.Run()
}

func TestK8sRegistry_GetServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	epMock := NewMockEndpointController(ctrl)
	epMock.EXPECT().GetEndpoints("foo", "default").Return(
		services, nil)

	registry := NewK8sRegistry(epMock)

	type fields struct {
		Registry        Registry
		k8sEpController EndpointController
		services        map[string][]*ServiceInstance
		lock            sync.Mutex
	}
	type args struct {
		service *Service
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*ServiceInstance
	}{
		{
			name: "services key doesn't exists in map",
			fields: fields{
				Registry:        registry,
				k8sEpController: epMock,
				services: map[string][]*ServiceInstance{
					"foo": services,
				},
			},
			args: args{
				service: &Service{
					Name:      "foo",
					Namespace: "default",
					Port:      8080,
				},
			},
			want: services,
		},
		{
			name: "services key exists in map",
			fields: fields{
				Registry:        registry,
				k8sEpController: epMock,
				services: map[string][]*ServiceInstance{
					"foo.default": services,
				},
			},
			args: args{
				service: &Service{
					Name:      "foo",
					Namespace: "default",
					Port:      8080,
				},
			},
			want: services,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &k8sRegistry{
				Registry:        tt.fields.Registry,
				k8sEpController: tt.fields.k8sEpController,
				services:        tt.fields.services,
				lock:            tt.fields.lock,
			}
			if got := r.GetServices(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestK8sRegistry_NewSubscriber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	epMock := NewMockEndpointController(ctrl)
	epMock.EXPECT().GetEndpoints("foo", "default").Return(
		services, nil)
	epMock.EXPECT().AddOnUpdateFunc("foo.default", gomock.Any()).Times(2)

	registry := NewK8sRegistry(epMock)

	type fields struct {
		Registry        Registry
		k8sEpController EndpointController
		services        map[string][]*ServiceInstance
		lock            sync.Mutex
	}
	type args struct {
		service *Service
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "services key error",
			fields: fields{
				Registry:        registry,
				k8sEpController: epMock,
				services: map[string][]*ServiceInstance{
					"foo": services,
				},
			},
			args: args{
				service: &Service{
					Name:      "foo",
					Namespace: "default",
					Port:      8080,
				},
			},
			want: []string{"172.16.1.2:8080", "172.16.1.3:8080"},
		},
		{
			name: "services key right",
			fields: fields{
				Registry:        registry,
				k8sEpController: epMock,
				services: map[string][]*ServiceInstance{
					"foo.default": services,
				},
			},
			args: args{
				service: &Service{
					Name:      "foo",
					Namespace: "default",
					Port:      8080,
				},
			},
			want: []string{"172.16.1.2:8080", "172.16.1.3:8080"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &k8sRegistry{
				Registry:        tt.fields.Registry,
				k8sEpController: tt.fields.k8sEpController,
				services:        tt.fields.services,
				lock:            tt.fields.lock,
			}
			if got := r.NewSubscriber(tt.args.service); !reflect.DeepEqual(got.Values(), tt.want) {
				t.Errorf("NewSubscriber() = %v, want %v", got, tt.want)
			}
		})
	}
}
