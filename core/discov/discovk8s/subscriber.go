package discovk8s

type Subscriber interface {
	SetUpdateFunc(func())
	OnUpdate()
	Values() []string
}
