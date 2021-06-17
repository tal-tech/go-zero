package resolver

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

const (
	// DirectScheme stands for direct scheme.
	DirectScheme = "direct"
	// DiscovScheme stands for discov scheme.
	DiscovScheme = "discov"
	// DiscovK8sScheme stands for k8s schema.
	DiscovK8sScheme = "k8s"
	// EndpointSepChar is the separator cha in endpoints.
	EndpointSepChar = ','

	subsetSize = 32
)

var (
	// EndpointSep is the separator string in endpoints.
	EndpointSep = fmt.Sprintf("%c", EndpointSepChar)

	dirBuilder    directBuilder
	disBuilder    discovBuilder
	disK8sBuilder discovK8sBuilder
)

// RegisterResolver registers the direct and discov schemes to the resolver.
func RegisterResolver() {
	resolver.Register(&dirBuilder)
	resolver.Register(&disBuilder)
	resolver.Register(&disK8sBuilder)
}

type nopResolver struct {
	cc resolver.ClientConn
}

func (r *nopResolver) Close() {
}

func (r *nopResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
