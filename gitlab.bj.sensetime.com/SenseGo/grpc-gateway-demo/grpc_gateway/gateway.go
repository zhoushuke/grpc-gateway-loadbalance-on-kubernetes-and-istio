package grpc_gateway

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type GatewayConfig struct {
	Fn       func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	GrpcPort string
	HttpPort string
	Mux      *runtime.ServeMux
	Opts     []grpc.DialOption
}

func RegisterGRPCGateway(cf *GatewayConfig) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cf.GrpcPort == "" {
		return errors.New("grpc port invalid")
	}

	if cf.HttpPort == "" {
		cf.HttpPort = ":8088"
	}

	if cf.Mux == nil {
		cf.Mux = runtime.NewServeMux()
	}

	if len(cf.Opts) == 0 {
		cf.Opts = DialOptionGRPC()
	}

	err := cf.Fn(ctx, cf.Mux, cf.GrpcPort, cf.Opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(cf.HttpPort, cf.Mux)
}

func DialOptionGRPC() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    time.Minute,
			Timeout: time.Second * 30,
		}),
	}
}
