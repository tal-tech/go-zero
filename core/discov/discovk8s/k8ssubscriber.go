package discovk8s

import "strconv"

type K8sSubscriber struct {
	Subscriber
	k8sRegistry *k8sRegistry
	instance    *Service
	updateFunc  func()
}

func (s *K8sSubscriber) SetUpdateFunc(f func()) {
	s.updateFunc = f
}

func (s *K8sSubscriber) OnUpdate() {
	if s.updateFunc != nil {
		s.updateFunc()
	}
}

func (s *K8sSubscriber) Values() []string {
	si := s.k8sRegistry.GetServices(s.instance)

	var ret []string
	for _, item := range si {
		ret = append(ret, item.Ip+":"+strconv.Itoa(int(item.Port)))
	}
	return ret
}
