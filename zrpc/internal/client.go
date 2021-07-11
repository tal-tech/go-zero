package internal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/zrpc/internal/balancer/p2c"
	"github.com/tal-tech/go-zero/zrpc/internal/clientinterceptors"
	"github.com/tal-tech/go-zero/zrpc/internal/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	dialTimeout = time.Second * 3
	separator   = '/'
)

func init() {
	resolver.RegisterResolver()
}

type (
	// Client interface wraps the Conn method.
	Client interface {
		Conn() *grpc.ClientConn
	}

	// A ClientOptions is a client options.
	ClientOptions struct {
		Timeout     time.Duration
		DialOptions []grpc.DialOption
	}

	// ClientOption defines the method to customize a ClientOptions.
	ClientOption func(options *ClientOptions)

	client struct {
		conn *grpc.ClientConn
	}
)

// NewClient returns a Client.
func NewClient(target string, opts ...ClientOption) (Client, error) {
	var cli client
	opts = append([]ClientOption{WithDialOption(grpc.WithBalancerName(p2c.Name))}, opts...)
	if err := cli.dial(target, opts...); err != nil {
		return nil, err
	}

	return &cli, nil
}

func (c *client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *client) buildDialOptions(opts ...ClientOption) []grpc.DialOption {
	var cliOpts ClientOptions
	for _, opt := range opts {
		opt(&cliOpts)
	}

	options := []grpc.DialOption{
		grpc.WithBlock(),
		WithUnaryClientInterceptors(
			clientinterceptors.TracingInterceptor,
			clientinterceptors.DurationInterceptor,
			clientinterceptors.BreakerInterceptor,
			clientinterceptors.PrometheusInterceptor,
			clientinterceptors.TimeoutInterceptor(cliOpts.Timeout),
		),
	}

	return append(options, cliOpts.DialOptions...)
}

func (c *client) dial(server string, opts ...ClientOption) error {
	options := c.buildDialOptions(opts...)
	timeCtx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	conn, err := grpc.DialContext(timeCtx, server, options...)
	if err != nil {
		service := server
		if errors.Is(err, context.DeadlineExceeded) {
			pos := strings.LastIndexByte(server, separator)
			// len(server) - 1 is the index of last char
			if 0 < pos && pos < len(server)-1 {
				service = server[pos+1:]
			}
		}
		return fmt.Errorf("rpc dial: %s, error: %s, make sure rpc service %q is already started",
			server, err.Error(), service)
	}

	c.conn = conn
	return nil
}

// WithDialOption returns a func to customize a ClientOptions with given dial option.
func WithDialOption(opt grpc.DialOption) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions, opt)
	}
}

// WithInsecure returns a func to customize a ClientOptions with secure option.
func WithInsecure() ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions, grpc.WithInsecure())
	}
}

// WithTimeout returns a func to customize a ClientOptions with given timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(options *ClientOptions) {
		options.Timeout = timeout
	}
}

// WithUnaryClientInterceptor returns a func to customize a ClientOptions with given interceptor.
func WithUnaryClientInterceptor(interceptor grpc.UnaryClientInterceptor) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions, WithUnaryClientInterceptors(interceptor))
	}
}

// WithTlsClientFromUnilateralism return a func to customize a ClientOptions Verify with Unilateralism authentication.
func WithTlsClientFromUnilateralism(crt, domainName string) ClientOption {
	return func(options *ClientOptions) {
		c, err := credentials.NewClientTLSFromFile(crt, domainName)
		if err != nil {
			log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
		}
		options.DialOptions = append(options.DialOptions, grpc.WithTransportCredentials(c))
	}
}

// WithTlsClientFromMutual return a func to customize a ClientOptions Verify with mutual authentication.
func WithTlsClientFromMutual(crtFile, keyFile, caFile string) ClientOption {
	return func(options *ClientOptions) {
		cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			log.Fatalf("tls.LoadX509KeyPair err: %v", err)
		}
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatalf("credentials: failed to ReadFile CA certificates err: %v", err)
		}

		if !certPool.AppendCertsFromPEM(ca) {
			log.Fatalf("credentials: failed to append certificates err: %v", err)
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		}

		options.DialOptions = append(options.DialOptions, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	}
}
