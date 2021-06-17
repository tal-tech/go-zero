package discovk8s

import "sort"
import v1 "k8s.io/api/core/v1"

func getReadyAddress(endpoints *v1.Endpoints) []*ServiceInstance {
	var readyAddesses []*ServiceInstance

	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			for _, port := range subset.Ports {
				si := ServiceInstance{
					Ip:   address.IP,
					Port: port.Port,
				}
				readyAddesses = append(readyAddesses, &si)
			}
		}
	}

	sort.Sort(ServiceInstanceSlice(readyAddesses))

	return readyAddesses
}
