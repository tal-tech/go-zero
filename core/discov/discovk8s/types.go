package discovk8s

import (
	"strconv"
	"strings"
)

type Service struct {
	Name      string
	Namespace string
	Port      int32
}

func (s *Service) EpName() string {
	return s.Name + "." + s.Namespace
}

type ServiceInstance struct {
	Ip   string
	Port int32
}

func (s *ServiceInstance) FullName() string {
	return s.Ip + ":" + strconv.Itoa(int(s.Port))
}

type ServiceInstanceSlice []*ServiceInstance

func (s ServiceInstanceSlice) Len() int {
	return len(s)
}

func (s ServiceInstanceSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ServiceInstanceSlice) Less(i, j int) bool {
	res := strings.Compare(s[i].Ip, s[j].Ip)
	if res < 0 {
		return true
	}

	if res == 0 && s[i].Port < s[j].Port {
		return true
	}

	return false
}
