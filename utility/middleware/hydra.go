package middleware

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
)

type authToken struct{}

func AuthForwarder() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, _ := ctx.Value(&authToken{}).(string)
		ctx = metautils.ExtractOutgoing(ctx).Add("Authorization", "Bearer "+token).ToOutgoing(ctx)
		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
