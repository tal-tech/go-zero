package discovk8s

import (
	v1 "k8s.io/api/core/v1"
	"reflect"
	"testing"
)

func Test_getReadyAddress(t *testing.T) {
	type args struct {
		endpoints *v1.Endpoints
	}
	tests := []struct {
		name string
		args args
		want []*ServiceInstance
	}{
		{
			name: "case1",
			args: args{
				endpoints: &v1.Endpoints{
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP: "172.16.2.1",
								},
								{
									IP: "172.16.2.2",
								},
							},
							Ports: []v1.EndpointPort{
								{
									Port: 8080,
								},
								{
									Port: 8081,
								},
							},
						},
					},
				},
			},
			want: []*ServiceInstance{
				{
					Ip:   "172.16.2.1",
					Port: 8080,
				},
				{
					Ip:   "172.16.2.1",
					Port: 8081,
				},
				{
					Ip:   "172.16.2.2",
					Port: 8080,
				},
				{
					Ip:   "172.16.2.2",
					Port: 8081,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getReadyAddress(tt.args.endpoints)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getReadyAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
