package internal

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/zrpc/internal/resolver"
)

// BuildDirectTarget returns a string that represents the given endpoints with direct schema.
func BuildDirectTarget(endpoints []string) string {
	return fmt.Sprintf("%s:///%s", resolver.DirectScheme,
		strings.Join(endpoints, resolver.EndpointSep))
}

// BuildDiscovTarget returns a string that represents the given endpoints with discov schema.
func BuildDiscovTarget(endpoints []string, key string) string {
	return fmt.Sprintf("%s://%s/%s", resolver.DiscovScheme,
		strings.Join(endpoints, resolver.EndpointSep), key)
}

// The target format is: k8s:///service-name.namespace:8080
func BuildDiscovK8sTarget(name string) string {
	return name
}
