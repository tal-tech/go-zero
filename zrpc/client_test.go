package zrpc

import (
	"context"
	"fmt"
	"github.com/tal-tech/go-zero/zrpc/internal"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/zrpc/internal/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func init() {
	logx.Disable()
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	mock.RegisterDepositServiceServer(server, &mock.DepositServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestDepositServer_Deposit(t *testing.T) {
	tests := []struct {
		name    string
		amount  float32
		res     *mock.DepositResponse
		errCode codes.Code
		errMsg  string
	}{
		{
			"invalid request with negative amount",
			-1.11,
			nil,
			codes.InvalidArgument,
			fmt.Sprintf("cannot deposit %v", -1.11),
		},
		{
			"valid request with non negative amount",
			0.00,
			&mock.DepositResponse{Ok: true},
			codes.OK,
			"",
		},
	}

	directClient := MustNewClient(
		RpcClientConf{
			Endpoints: []string{"foo"},
			App:       "foo",
			Token:     "bar",
			Timeout:   1000,
		},
		WithDialOption(grpc.WithInsecure()),
		WithDialOption(grpc.WithContextDialer(dialer())),
		WithUnaryClientInterceptor(func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	targetClient, err := NewClientWithTarget("foo", WithDialOption(grpc.WithInsecure()),
		WithDialOption(grpc.WithContextDialer(dialer())), WithUnaryClientInterceptor(
			func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
				invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				return invoker(ctx, method, req, reply, cc, opts...)
			}))
	assert.Nil(t, err)
	clients := []Client{
		directClient,
		targetClient,
	}
	for _, tt := range tests {
		for _, client := range clients {
			t.Run(tt.name, func(t *testing.T) {
				cli := mock.NewDepositServiceClient(client.Conn())
				request := &mock.DepositRequest{Amount: tt.amount}
				response, err := cli.Deposit(context.Background(), request)
				if response != nil {
					assert.True(t, len(response.String()) > 0)
					if response.GetOk() != tt.res.GetOk() {
						t.Error("response: expected", tt.res.GetOk(), "received", response.GetOk())
					}
				}
				if err != nil {
					if e, ok := status.FromError(err); ok {
						if e.Code() != tt.errCode {
							t.Error("error code: expected", codes.InvalidArgument, "received", e.Code())
						}
						if e.Message() != tt.errMsg {
							t.Error("error message: expected", tt.errMsg, "received", e.Message())
						}
					}
				}
			})
		}
	}
}

func TestNewClientWithTarget(t *testing.T) {
	//client, err := NewClientWithTarget("kubernetes://nginx-service.default:80")
	//client, err := NewClientWithTarget("discov://nginx-service.default:80")

	client, err := internal.NewClient(internal.BuildDiscovK8sTarget("k8s:///transform-svc.default:8081"))
	fmt.Println(client, err)
	select {}
	//assert.Nil(t, err)

	//assert.Equal(t, client.Conn().GetState(), connectivity.Ready)
}
